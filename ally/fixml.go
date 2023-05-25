package ally

import (
	"encoding/xml"
	"fmt"
	"log"
)

func prepRender(o *Order) fixmlOrder {
	if o.SecType == "" {
		o.SecType = "CS"
	}

	if o.Tif == "MOC" {
		o.OrderType = ""
		if o.SecType != "CS" {
			log.Panicln("Market on Close valid only for Stock contracts")
		}
	}

	order_types := map[string]int{
		"MKT":     1,
		"LMT":     2,
		"STP":     3,
		"STP LMT": 4,
	}

	tifs := map[string]int{
		"DAY": 0,
		"GTC": 1,
		"MOC": 7,
	}

	actions := map[string]int{
		"BUY":   1,
		"SELL":  2,
		"SHORT": 5,
		"COVER": 1,
	}

	var qty string
	if o.Orig_ID != "" {
		if o.Action == "BUY" {
			qty = fmt.Sprintf("%.0f", o.TotalQuantity)
		} else {
			qty = fmt.Sprintf("%.2f", o.TotalQuantity)
		}
	} else {
		qty = fmt.Sprint(o.TotalQuantity)
	}

	var px string
	if o.OrderType == "LMT" || o.OrderType == "STP LMT" {
		px = fmt.Sprintf("%.2f", o.LmtPrice)
	} else {
		px = ""
	}

	var stoppx string
	if o.OrderType == "STP" || o.OrderType == "STP LMT" {
		stoppx = fmt.Sprintf("%.2f", o.StopPrice)
	} else {
		stoppx = ""
	}

	var accttyp string
	if o.Action == "COVER" {
		accttyp = "5"
	}

	return fixmlOrder{
		Acct:        o.Account,
		AcctTyp:     accttyp,
		OrigID:      o.Orig_ID,
		OrigClOrdID: "",
		PosEfct:     "",
		Px:          px,
		Side:        fmt.Sprint(actions[o.Action]),
		StopPx:      stoppx,
		TmInForce:   fmt.Sprint(tifs[o.Tif]),
		Typ:         fmt.Sprint(order_types[o.OrderType]),
		Instrmt: fixmlInstrmt{
			CFI:    "",
			MatDt:  "",
			StrkPx: "",
			SecTyp: o.SecType,
			Sym:    o.Symbol,
		},
		OrdQty: fixmlOrdQty{
			Qty: qty,
		},
	}

}

func Render(o *Order) FIXMLOrder {
	order := &o
	renderedOrder := prepRender(*order)
	renderedOrder.OrigID = ""
	return FIXMLOrder{
		XMLName: xml.Name{},
		Xmlns:   "http://www.fixprotocol.org/FIXML-5-0-SP2",
		Order:   renderedOrder,
	}
}

func RenderCancel(o *Order) FIXMLCancel {
	order := &o
	renderedOrder := prepRender(*order)

	return FIXMLCancel{
		XMLName: xml.Name{},
		Xmlns:   "http://www.fixprotocol.org/FIXML-5-0-SP2",
		Cancel:  renderedOrder,
	}
}

func RenderChange(o *Order) FIXMLChange {
	order := &o
	renderedOrder := prepRender(*order)

	return FIXMLChange{
		XMLName: xml.Name{},
		Xmlns:   "http://www.fixprotocol.org/FIXML-5-0-SP2",
		Change:  renderedOrder,
	}
}

func FetchOrder(fml *FIXMLFetch) (Order, error) {
	var err error

	order_types := map[int]string{
		1: "MKT",
		2: "LMT",
		3: "STP",
		4: "STP LMT",
	}

	tifs := map[int]string{
		0: "DAY",
		1: "GTC",
		7: "MOC",
	}

	actions := map[int]string{
		1: "BUY",
		2: "SELL",
		5: "SHORT",
	}

	o := Order{
		Symbol:        fml.Rpt.Instrmt.Sym,
		SecType:       fml.Rpt.Instrmt.SecTyp,
		Account:       fml.Rpt.Acct,
		Action:        actions[fml.Rpt.Side],
		OrderType:     order_types[fml.Rpt.Typ],
		Tif:           tifs[fml.Rpt.TmInForce],
		TotalQuantity: fml.Rpt.LeavesQty,
		LmtPrice:      fml.Rpt.Px,
		StopPrice:     fml.Rpt.StopPx,
		Orig_ID:       fml.Rpt.OrdID,
	}

	return o, err
}
