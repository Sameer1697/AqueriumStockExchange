package matchingengine

import (
	"goStockExchange/order"
	"goStockExchange/orderbook"
)

type CommandType int

const (
	CmdNewOrder CommandType = iota
	CmdCancel
)

type Command struct {
	Type    CommandType
	Order   *order.Order
	OrderID uint64
}

type MatchingEngine struct {
	book *orderbook.OrderBook
	cmds chan Command
}

func NewMatchingEngine(symbol string, buffer int) *MatchingEngine {
	return &MatchingEngine{
		book: orderbook.NewOrderBook(symbol),
		cmds: make(chan Command, buffer),
	}
}

func (me *MatchingEngine) SubmitOrder(o *order.Order) {
	me.cmds <- Command{
		Type:  CmdNewOrder,
		Order: o,
	}
}

func (me *MatchingEngine) CancelOrder(orderID uint64) {
	me.cmds <- Command{
		Type:    CmdCancel,
		OrderID: orderID,
	}
}

func (me *MatchingEngine) Run() {
	for cmd := range me.cmds {
		switch cmd.Type {

		case CmdNewOrder:
			me.book.Add(cmd.Order)

		case CmdCancel:
			me.book.Cancel(cmd.OrderID)
		}
	}
}

func (me *MatchingEngine) Stop() {
	close(me.cmds)
}
