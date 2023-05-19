package ally

import (
	"encoding/xml"
	"fmt"
	"reflect"
	"testing"
)

func TestRender(t *testing.T) {
	mkt_order := Order{
		Symbol:        "AAPL",
		SecType:       "CS",
		Account:       "TEST",
		Action:        "BUY",
		OrderType:     "MKT",
		Tif:           "DAY",
		TotalQuantity: 5,
		LmtPrice:      1,
		StopPrice:     0,
		Transmit:      false,
	}

	mkt_fixml := "<FIXML xmlns=\"http://www.fixprotocol.org/FIXML-5-0-SP2\"><Order Acct=\"TEST\" Side=\"1\" TmInForce=\"0\" Typ=\"1\"><Instrmt SecTyp=\"CS\" Sym=\"AAPL\"></Instrmt><OrdQty Qty=\"5\"></OrdQty></Order></FIXML>"

	lmt_order := Order{
		Symbol:        "AAPL",
		SecType:       "CS",
		Account:       "TEST",
		Action:        "SELL",
		OrderType:     "LMT",
		Tif:           "GTC",
		TotalQuantity: 1.25,
		LmtPrice:      12.00,
	}

	lmt_fixml := "<FIXML xmlns=\"http://www.fixprotocol.org/FIXML-5-0-SP2\"><Order Acct=\"TEST\" Px=\"12.00\" Side=\"2\" TmInForce=\"1\" Typ=\"2\"><Instrmt SecTyp=\"CS\" Sym=\"AAPL\"></Instrmt><OrdQty Qty=\"1.25\"></OrdQty></Order></FIXML>"

	data, _ := xml.Marshal(lmt_fixml)
	fmt.Printf("%s", data)

	type args struct {
		o *Order
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "market",
			args: args{&mkt_order},
			want: mkt_fixml,
		},
		{
			name: "limit",
			args: args{&lmt_order},
			want: lmt_fixml,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := xml.Marshal(Render(tt.args.o)); !reflect.DeepEqual(string(got), tt.want) {
				t.Errorf("Render() = %v, want %v", string(got), tt.want)
			}
		})
	}
}

func TestRenderChange(t *testing.T) {
	lmt_order := Order{
		Symbol:        "AAPL",
		SecType:       "CS",
		Account:       "TEST",
		Action:        "SELL",
		OrderType:     "LMT",
		Tif:           "GTC",
		TotalQuantity: 1.25,
		LmtPrice:      12.00,
		Orig_ID:       "SVI-TEST",
	}
	lmt_fixml := "<FIXML xmlns=\"http://www.fixprotocol.org/FIXML-5-0-SP2\"><OrderCxlRplcReq Acct=\"TEST\" OrigID=\"SVI-TEST\" Px=\"12.00\" Side=\"2\" TmInForce=\"1\" Typ=\"2\"><Instrmt SecTyp=\"CS\" Sym=\"AAPL\"></Instrmt><OrdQty Qty=\"1.25\"></OrdQty></OrderCxlRplcReq></FIXML>"
	type args struct {
		o *Order
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "limit",
			args: args{&lmt_order},
			want: lmt_fixml,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := xml.Marshal(RenderChange(tt.args.o)); !reflect.DeepEqual(string(got), tt.want) {
				t.Errorf("RenderChange() = %v, want %v", string(got), tt.want)
			}
		})
	}
}

func TestRenderCancel(t *testing.T) {
	lmt_order := Order{
		Symbol:        "AAPL",
		SecType:       "CS",
		Account:       "TEST",
		Action:        "SELL",
		OrderType:     "LMT",
		Tif:           "GTC",
		TotalQuantity: 1.25,
		LmtPrice:      12.00,
		Orig_ID:       "SVI-TEST",
	}
	lmt_fixml := "<FIXML xmlns=\"http://www.fixprotocol.org/FIXML-5-0-SP2\"><OrderCxlReq Acct=\"TEST\" OrigID=\"SVI-TEST\" Px=\"12.00\" Side=\"2\" TmInForce=\"1\" Typ=\"2\"><Instrmt SecTyp=\"CS\" Sym=\"AAPL\"></Instrmt><OrdQty Qty=\"1.25\"></OrdQty></OrderCxlReq></FIXML>"
	type args struct {
		o *Order
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "limit",
			args: args{&lmt_order},
			want: lmt_fixml,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := xml.Marshal(RenderCancel(tt.args.o)); !reflect.DeepEqual(string(got), tt.want) {
				t.Errorf("RenderCancel() = %v, want %v", string(got), tt.want)
			}
		})
	}
}

func TestGetOrders(t *testing.T) {
	t.Error("not imp")
	type args struct {
		id string
	}
	tests := []struct {
		name string
		args args
		want IGetOrders
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetOrders(tt.args.id); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetOrders() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPostOrder(t *testing.T) {
	t.Error("not imp")
	type args struct {
		order *Order
	}
	tests := []struct {
		name    string
		args    args
		want    IPostOrder
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := PostOrder(tt.args.order)
			if (err != nil) != tt.wantErr {
				t.Errorf("PostOrder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PostOrder() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPostCancel(t *testing.T) {
	t.Error("not imp")
	type args struct {
		order *Order
	}
	tests := []struct {
		name    string
		args    args
		want    IPostOrder
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := PostCancel(tt.args.order)
			if (err != nil) != tt.wantErr {
				t.Errorf("PostCancel() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PostCancel() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPostChange(t *testing.T) {
	t.Error("not imp")
	type args struct {
		order *Order
	}
	tests := []struct {
		name    string
		args    args
		want    IPostOrder
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := PostChange(tt.args.order)
			if (err != nil) != tt.wantErr {
				t.Errorf("PostChange() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PostChange() = %v, want %v", got, tt.want)
			}
		})
	}
}
