package orderbook

import (
	"goStockExchange/order"
)

type OrderBook struct {
	Symbol string

	Bids *SideBook // BUY side
	Asks *SideBook // SELL side

	index map[uint64]*OrderRef
}

func NewOrderBook(symbol string) *OrderBook {
	return &OrderBook{
		Symbol: symbol,
		Bids:   NewSideBook(true),
		Asks:   NewSideBook(false),
		index:  make(map[uint64]*OrderRef),
	}
}

func (ob *OrderBook) Add(o *order.Order) {
	if o.IsBuy() {
		ob.matchBuy(o)
	} else {
		ob.matchSell(o)
	}

	if o.Remaining() > 0 && o.IsLimit() {
		ob.addResting(o)
	}
}

func (ob *OrderBook) matchBuy(o *order.Order) {
	var empty []*PriceLevel
	ob.Asks.ForEachFromBest(func(pl *PriceLevel) bool {
		if o.IsLimit() && pl.Price > o.Price {
			return false // price crossed
		}

		pl.Match(o)
		if pl.IsEmpty() {
			empty = append(empty, pl)
		}

		return o.Remaining() > 0

	})

	for _, pl := range empty {
		ob.Asks.Remove(pl.Price)
	}
}

func (ob *OrderBook) matchSell(o *order.Order) {
	ob.Bids.ForEachFromBest(func(pl *PriceLevel) bool {
		if o.IsLimit() && pl.Price < o.Price {
			return false
		}

		pl.Match(o)
		ob.Bids.RemoveIfEmpty(pl)

		return o.Remaining() > 0
	})
}

func (ob *OrderBook) addResting(o *order.Order) {
	var side *SideBook
	if o.IsBuy() {
		side = ob.Bids
	} else {
		side = ob.Asks
	}

	pl := side.GetOrCreate(o.Price)
	node := pl.Add(o)

	ob.index[o.ID] = &OrderRef{
		Order: o,
		Level: pl,
		Node:  node,
	}
}

func (ob *OrderBook) BestBid() *PriceLevel {
	return ob.Bids.Best()
}

func (ob *OrderBook) BestAsk() *PriceLevel {
	return ob.Asks.Best()
}

func (ob *OrderBook) IsEmpty() bool {
	return ob.Bids.IsEmpty() && ob.Asks.IsEmpty()
}

func (ob *OrderBook) Cancel(orderID uint64) bool {
	ref, ok := ob.index[orderID]
	if !ok {
		return false
	}

	ref.Level.Remove(ref.Node)
	ref.Order.Cancel()

	delete(ob.index, orderID)

	if ref.Level.IsEmpty() {
		if ref.Order.IsBuy() {
			ob.Bids.Remove(ref.Level.Price)
		} else {
			ob.Asks.Remove(ref.Level.Price)
		}
	}

	return true
}
