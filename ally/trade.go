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

func ValidateOrder(order *Order) error {
	if order.Transmit && order.LmtPrice >= MINIMUM_STOCK_PRICE {
		return nil
	}

	data, _ := xml.Marshal(Render(order))
	return fmt.Errorf("not a valid order: %s", string(data))
}

var OverrideMap = map[WarningCode]bool{
	DuplicateOrder:       false,
	UnsettledFunds:       false,
	HigherMarginReq:      false,
	NotMarginable:        false,
	MktOrderWhileClosed:  false,
	ExchangeClosed:       false,
	ForeignSettlementFee: false,
}

func handleWarning(order *Order, resp IPostOrderWarn) error {
	override := OverrideMap[resp.Warning.Code]

	if !override {
		order.Transmit = false
		return fmt.Errorf("error code %d; not overridden", resp.Warning.Code)
	} else {
		switch resp.Warning.Code {
		case MktOrderWhileClosed:
			order.OrderType = "LMT"
			order.LmtPrice = resp.FinData.Quote.LastPx
		case ExchangeClosed:
			order.OrderType = "LMT"
			order.LmtPrice = resp.FinData.Quote.LastPx
		default:
			log.Printf("error code %d; overridden\n", resp.Warning.Code)
		}
	}
	return nil
}

// POST accounts/:id/orders
func PostOrder(order *Order) (IPostOrder, error) {
	err := ValidateOrder(order)
	if err != nil {
		return IPostOrder{}, err
	}

	data, err := xml.Marshal(Render(order))
	if err != nil {
		return IPostOrder{}, fmt.Errorf("got order that couldn't be rendered to xml: %s (%s)", Render(order), err.Error())
	}

	headers := RequestHeaders()
	resp, body, err := post[IPostOrder](fmt.Sprintf("accounts/%s/orders", order.Account), data, headers, url.Values{}, tradesRL)

	if err == nil {
		return resp, err
	}

	warn := IPostOrderWarn{}
	err = xml.Unmarshal(body, &warn)
	log.Printf("order had warnings: %s", warn.Warning.Text)

	if err != nil {
		return IPostOrder{}, fmt.Errorf("got invalid order response: %s (%s)", Render(order), err.Error())
	}

	err = handleWarning(order, warn)

	if err == nil {
		err = ValidateOrder(order)
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
