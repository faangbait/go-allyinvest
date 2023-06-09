package ally

import (
	"encoding/xml"
	"fmt"
	"log"
	"net/url"
)

// GET accounts/:id/orders
func GetOrders(id string) IGetOrders {
	return get[IGetOrders](fmt.Sprintf("accounts/%s/orders", id), url.Values{}, authedRL)
}

// POST accounts/:id/orders/preview
// Attempt to preview the order. By default, we don't use this in the order flow.
func PostPreview(order *Order) (IPostPreview, error) {
	data, err := xml.Marshal(Render(order))

	if err != nil {
		return IPostPreview{}, fmt.Errorf("got order that couldn't be rendered to xml: %s (%s)", Render(order), err.Error())
	}

	headers := RequestHeaders()
	resp, _, err := post[IPostPreview](fmt.Sprintf("accounts/%s/orders/preview", order.Account), data, headers, url.Values{}, tradesRL)
	return resp, err
}

// POST accounts/:id/orders
// Post the order, overriding warnings.
func postOverride(order *Order) (IPostOrder, error) {
	err := ValidateOrder(order)
	if err != nil {
		return IPostOrder{}, err
	}

	data, err := xml.Marshal(Render(order))
	if err != nil {
		return IPostOrder{}, fmt.Errorf("got order that couldn't be rendered to xml: %s (%s)", Render(order), err.Error())
	}

	headers := RequestHeaders()
	headers["TKI_OVERRIDE"] = []string{"true"}
	resp, _, err := post[IPostOrder](fmt.Sprintf("accounts/%s/orders", order.Account), data, headers, url.Values{}, tradesRL)
	return resp, err

}

// Attempt to validate the order. Will coax the order into a valid order by default.
// To prevent this behavior, set order.DoNotCoax = true
func ValidateOrder(o *Order) error {
	if o.Transmit && o.LmtPrice >= MINIMUM_STOCK_PX {
		return nil
	}

	// Blank out OrderType for MOC contracts
	if o.Tif == "MOC" {
		if o.SecType != "CS" {
			return fmt.Errorf("market on close orders valid only for Stock contracts: %v", o)
		}
		o.OrderType = ""
	}

	// Set SecType to CS if undefined
	if o.SecType == "" {
		if o.DoNotCoax {
			return fmt.Errorf("no security type defined: %v", o)
		}

		o.SecType = "CS"
	}

	// Set Tif to 'Day' for market orders
	if o.Tif == "GTC" && o.OrderType == "MKT" {
		if o.DoNotCoax {
			return fmt.Errorf("gtc order not valid for market ordertype")
		}
		o.Tif = "DAY"
	}

	data, _ := xml.Marshal(Render(o))
	return fmt.Errorf("not a valid order: %s", string(data))
}

// Attempts to handle a warning in the manner determined by the OverrideMap
// In other words, if making this a Limit order solves our problem, we do that.
// Returns an error if can't or won't override the warning; this should abort the trade.
func handleWarning(order *Order, resp *IPostOrderWarn) error {
	if resp.Warning.Code == 0 {
		return nil
	}

	log.Printf("order had warnings: %s", resp.Warning.Text)
	override := OverrideMap[resp.Warning.Code]

	if !override {
		order.Transmit = false
		return fmt.Errorf("error code %d; not overridden", resp.Warning.Code)
	} else {
		switch resp.Warning.Code {
		case MktOrderWhileClosed, ExchangeClosed, NoMarketOrder1, NoMarketOrder2:
			if order.OrderType == "MKT" || order.OrderType == "MOC" {
				order.OrderType = "LMT"
				order.LmtPrice = resp.FinData.Quote.LastPx
			} else if order.OrderType == "STP" {
				order.OrderType = "STP LMT"
				order.LmtPrice = resp.FinData.Quote.LastPx
			} else {
				return fmt.Errorf("error code %d; couldn't override", resp.Warning.Code)
			}
		default:
			log.Printf("error code %d; overridden\n", resp.Warning.Code)
		}
	}
	resp.Warning.Text = ""
	resp.Warning.Code = 0
	return nil
}

// POST accounts/:id/orders
// Attempt to post the order. Validate and handle errors according to the rules defined by OverrideMap.
func PostOrder(order *Order) (IPostOrder, error) {
	warn := IPostOrderWarn{}
	headers := RequestHeaders()

	// Validate the order; if we can't, return an error
	err := ValidateOrder(order)
	if err != nil {
		return IPostOrder{}, err
	}

	// Render the order to fixml; if we can't, return an error
	data, err := xml.Marshal(Render(order))
	if err != nil {
		return IPostOrder{}, fmt.Errorf("got order that couldn't be rendered to xml: %s (%s)", Render(order), err.Error())
	}

	// Attempt to post the order.
	resp, body, err := post[IPostOrder](fmt.Sprintf("accounts/%s/orders", order.Account), data, headers, url.Values{}, tradesRL)

	if err != nil {
		return IPostOrder{}, fmt.Errorf("got invalid order response: %s (%s)", Render(order), err.Error())
	}

	// Check if order contained an unhandled error and handle the custom warning code.
	switch resp.Error {
	case "Success":
		// Order successfully posted as-is. Woohoo!
		return resp, nil
	case "We are not currently accepting Market orders for this security. Please change your order to a Limit order.":
		warn.Warning.Text = resp.Error
		warn.Warning.Code = NoMarketOrder1
	case "We are not currently accepting Market or Stop orders. Please place a Limit order.":
		warn.Warning.Text = resp.Error
		warn.Warning.Code = NoMarketOrder2
	case "Due to nightly processing we are unable to accept orders between 11:30 PM and 12:00 AM EST. Please replace your order after 12:00 AM .":
		warn.Warning.Text = resp.Error
		warn.Warning.Code = MaintenanceWindow
	}

	if warn.Warning.Code >= 2000 {
		err = handleWarning(order, &warn)
		if err != nil {
			return resp, err
		}
	}

	// Now that we know the order could not post, attempt to unmarshal the response into an IPostOrderWarn.
	err = xml.Unmarshal(body, &warn)
	if err != nil {
		return IPostOrder{}, fmt.Errorf("got invalid order response: %s (%s)", Render(order), err.Error())
	}

	// Handle any additional warnings; perhaps overriding them.
	err = handleWarning(order, &warn)

	if err == nil {
		if order.Transmit {
			return postOverride(order)
		}
	}

	return resp, err
}

// POST accounts/:id/orders
func PostCancel(order *Order) (IPostOrder, error) {
	data, err := xml.Marshal(RenderCancel(order))

	if err != nil {
		return IPostOrder{}, fmt.Errorf("got cancel order that couldn't be rendered to xml: %s (%s)", RenderCancel(order), err.Error())
	}

	headers := RequestHeaders()
	resp, _, err := post[IPostOrder](fmt.Sprintf("accounts/%s/orders", order.Account), data, headers, url.Values{}, tradesRL)
	return resp, err
}

// POST accounts/:id/orders
func PostChange(order *Order) (IPostOrder, error) {
	data, err := xml.Marshal(RenderChange(order))

	if err != nil {
		return IPostOrder{}, fmt.Errorf("got change order that couldn't be rendered to xml: %s (%s)", RenderChange(order), err.Error())
	}

	headers := RequestHeaders()
	resp, _, err := post[IPostOrder](fmt.Sprintf("accounts/%s/orders", order.Account), data, headers, url.Values{}, tradesRL)
	return resp, err
}

// By default, override no values. Strategies can update this map to change behavior.
var OverrideMap = map[WarningCode]bool{
	DuplicateOrder:       false,
	UnsettledFunds:       false,
	HigherMarginReq:      false,
	NotMarginable:        false,
	MktOrderWhileClosed:  false,
	ExchangeClosed:       false,
	ForeignSettlementFee: false,
	NoMarketOrder1:       false,
	NoMarketOrder2:       false,
	MaintenanceWindow:    false,
}
