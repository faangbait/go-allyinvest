package ally

import (
	"encoding/xml"
	"fmt"
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

// Attempts to handle a warning in the manner determined by the OverrideMap
// In other words, if making this a Limit order solves our problem, we do that.
// Returns an error if can't or won't override the warning; this should abort the trade.
func handleWarning(order *Order, preview *IPostPreview) (IPostOrder, error) {
	switch preview.Error {
	case "Success":
		return PostOverride(order, true)
	case "Please check your Order Status and Holdings pages.  You may have an existing open order or you may be entering a quantity greater than the position you hold in your account.  Please review and replace this order with the revised quantity as needed.":
		order.Transmit = false
		return IPostOrder{}, fmt.Errorf("error code %d; couldn't override", DuplicateOrder)

	case "UnsettledFunds":
		if !OverrideMap[UnsettledFunds] {
			order.Transmit = false
			return IPostOrder{}, fmt.Errorf("error code %d; couldn't override", UnsettledFunds)
		}

	case "HigherMarginReq":
		if !OverrideMap[HigherMarginReq] {
			order.Transmit = false
			return IPostOrder{}, fmt.Errorf("error code %d; couldn't override", HigherMarginReq)
		}

	case "NotMarginable":
		if !OverrideMap[NotMarginable] {
			order.Transmit = false
			return IPostOrder{}, fmt.Errorf("error code %d; couldn't override", NotMarginable)
		}

	case "MktOrderWhileClosed":
		if !OverrideMap[MktOrderWhileClosed] {
			order.Transmit = false
			return IPostOrder{}, fmt.Errorf("error code %d; couldn't override", MktOrderWhileClosed)
		}
		switch order.OrderType {
		case "MKT", "MOC":
			order.OrderType = "LMT"
		case "STP":
			order.OrderType = "STP LMT"
		default:
			return IPostOrder{}, fmt.Errorf("error code %d; couldn't override", MktOrderWhileClosed)
		}

	case "ExchangeClosed":
		order.Transmit = false
		return IPostOrder{}, fmt.Errorf("error code %d; couldn't override", ExchangeClosed)

	case "ForeignSettlementFee":
		if !OverrideMap[ForeignSettlementFee] {
			order.Transmit = false
			return IPostOrder{}, fmt.Errorf("error code %d; couldn't override", ForeignSettlementFee)
		}

	case "We are not currently accepting Market orders for this security. Please change your order to a Limit order.":
		if !OverrideMap[NoMarketOrder1] {
			order.Transmit = false
			return IPostOrder{}, fmt.Errorf("error code %d; couldn't override", NoMarketOrder1)
		}
		switch order.OrderType {
		case "MKT", "MOC":
			order.OrderType = "LMT"
		case "STP":
			order.OrderType = "STP LMT"
		default:
			return IPostOrder{}, fmt.Errorf("error code %d; couldn't override", NoMarketOrder1)
		}

	case "We are not currently accepting Market or Stop orders. Please place a Limit order.":
		if !OverrideMap[NoMarketOrder2] {
			order.Transmit = false
			return IPostOrder{}, fmt.Errorf("error code %d; couldn't override", NoMarketOrder2)
		}
		switch order.OrderType {
		case "MKT", "MOC":
			order.OrderType = "LMT"
		case "STP":
			order.OrderType = "STP LMT"
		default:
			return IPostOrder{}, fmt.Errorf("error code %d; couldn't override", NoMarketOrder2)
		}

	case "Due to nightly processing we are unable to accept orders between 11:30 PM and 12:00 AM EST. Please replace your order after 12:00 AM .":
		return IPostOrder{}, fmt.Errorf("tried to place order during maintenance window")

	case "This is a Pinksheet or over-the-counter (OTC) security.  Stop and Stop Limit orders are not available for this security.  Please use a Limit order.":
		if !OverrideMap[NoStopLimitOrders] || order.OrderType == "STP" {
			order.Transmit = false
			return IPostOrder{}, fmt.Errorf("error code %d; couldn't override", NoStopLimitOrders)
		}
		order.OrderType = "LMT"

	case "You are placing a stop order on the wrong side of the market by either (1) purchasing at a price below the current Ask price or (2) selling at a price above the current Bid price.  Please adjust your Stop price based on the current quote.", "For stop-limit order, limit price should be equal or above stop price for buys, and equal or below stop price for sells.", "The limit price you have entered is too aggressive either buying at a limit above the Ask, or selling at a limit below the Bid.  Please adjust your price.":
		if !OverrideMap[WrongDirectionTrade] {
			order.Transmit = false
			return IPostOrder{}, fmt.Errorf("error code %d; couldn't override", WrongDirectionTrade)
		}

	default:
		return IPostOrder{}, fmt.Errorf("error text %s; no override implemented", preview.WarningText)
	}

	return PostOverride(order, true)
}

// POST accounts/:id/orders
// Attempt to post the order. Validate and handle errors according to the rules defined by OverrideMap.
func PostOrder(order *Order) (IPostOrder, error) {
	// Validate the order; if we can't, return an error
	err := ValidateOrder(order)
	if err != nil {
		return IPostOrder{}, err
	}

	// Render the order to fixml; if we can't, return an error
	_, err = xml.Marshal(Render(order))
	if err != nil {
		return IPostOrder{}, fmt.Errorf("got order that couldn't be rendered to xml: %s (%s)", Render(order), err.Error())
	}

	// Preview the order, check for warnings
	preview, err := PostPreview(order)
	if err != nil {
		return IPostOrder{}, fmt.Errorf("got order that couldn't be previewed")
	}

	// Handle any warnings that occurred
	resp, err := handleWarning(order, &preview)

	if err != nil {
		return resp, fmt.Errorf("got warning we couldn't handle")
	}

	return resp, err
}

// POST accounts/:id/orders
// Post the order, overriding warnings.
func PostOverride(order *Order, validate bool) (IPostOrder, error) {
	if validate {
		err := ValidateOrder(order)
		if err != nil {
			return IPostOrder{}, err
		}
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

	if o.LmtPrice < MINIMUM_STOCK_PX {
		if o.DoNotCoax {
			return fmt.Errorf("lmtprice below minimum")
		}
		o.LmtPrice = MINIMUM_STOCK_PX
	}

	if !o.Transmit {
		return fmt.Errorf("order marked no-transmit")
	}

	data, err := xml.Marshal(Render(o))
	if err != nil {
		return fmt.Errorf("not a valid order: %s", string(data))
	}
	return nil
}

// // POST accounts/:id/orders
// // Attempt to post the order. Validate and handle errors according to the rules defined by OverrideMap.
// func PostOrder(order *Order) (IPostOrder, error) {
// 	warn := IPostOrderWarn{}
// 	headers := RequestHeaders()

// 	// Validate the order; if we can't, return an error
// 	err := ValidateOrder(order)
// 	if err != nil {
// 		return IPostOrder{}, err
// 	}

// 	// Render the order to fixml; if we can't, return an error
// 	data, err := xml.Marshal(Render(order))
// 	if err != nil {
// 		return IPostOrder{}, fmt.Errorf("got order that couldn't be rendered to xml: %s (%s)", Render(order), err.Error())
// 	}

// 	// Attempt to post the order.
// 	resp, body, err := post[IPostOrder](fmt.Sprintf("accounts/%s/orders", order.Account), data, headers, url.Values{}, tradesRL)

// 	if err != nil {
// 		return IPostOrder{}, fmt.Errorf("got invalid order response: %s (%s)", Render(order), err.Error())
// 	}

// 	// Check if order contained an unhandled error and handle the custom warning code.
// 	switch resp.Error {
// 	case "Success":
// 		// Order successfully posted as-is. Woohoo!
// 		return resp, nil
// 	case "We are not currently accepting Market orders for this security. Please change your order to a Limit order.":
// 		warn.Warning.Text = resp.Error
// 		warn.Warning.Code = NoMarketOrder1
// 	case "We are not currently accepting Market or Stop orders. Please place a Limit order.":
// 		warn.Warning.Text = resp.Error
// 		warn.Warning.Code = NoMarketOrder2
// 	case "Due to nightly processing we are unable to accept orders between 11:30 PM and 12:00 AM EST. Please replace your order after 12:00 AM .":
// 		warn.Warning.Text = resp.Error
// 		warn.Warning.Code = MaintenanceWindow
// 	}

// 	// Now that we know the order could not post, attempt to unmarshal the response into an IPostOrderWarn.
// 	err = xml.Unmarshal(body, &warn)
// 	if err != nil {
// 		return IPostOrder{}, fmt.Errorf("got invalid order response: %s (%s)", Render(order), err.Error())
// 	}

// 	// Handle any additional warnings; perhaps overriding them.
// 	err = handleWarning(order, &warn)

// 	if err == nil {
// 		if order.Transmit {
// 			return postOverride(order)
// 		}
// 	}

// 	return resp, err
// }

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
	NoStopLimitOrders:    false,
	WrongDirectionTrade:  false,
}
