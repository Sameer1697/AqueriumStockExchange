package orderbook

import (
	"container/list" //doubly linked list to maintain fifo order queue
	"goStockExchange/order"
)

type PriceLevel struct {
	Price  float64
	orders *list.List // FIFO queue of *order.Order

	//orderMap map[uint64]*list.Element
	// we will user ordermap later to optimize order removal
	//as it will give us o(1) time complexity instead of o(n)
	//But it will consume extra memory
}

func NewPriceLevel(price float64) *PriceLevel {
	return &PriceLevel{
		Price:  price,
		orders: list.New(),
	}
}

// Add appends an order to the FIFO queue
func (pl *PriceLevel) Add(o *order.Order) *list.Element {
	return pl.orders.PushBack(o)
}

// Front returns the oldest resting order without removing it
func (pl *PriceLevel) Front() *order.Order {
	if e := pl.orders.Front(); e != nil {
		return e.Value.(*order.Order)
	}
	return nil
}

// Remove deletes a specific order from this price level (cancel path)
func (pl *PriceLevel) Remove(e *list.Element) {
	pl.orders.Remove(e)
}

// IsEmpty indicates whether this price level has no resting orders
func (pl *PriceLevel) IsEmpty() bool {
	return pl.orders.Len() == 0
}

// Len returns the number of resting orders at this price
func (pl *PriceLevel) Len() int {
	return pl.orders.Len()
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

// Match executes FIFO matching against the incoming order
func (pl *PriceLevel) Match(incoming *order.Order) {
	for e := pl.orders.Front(); e != nil && incoming.Remaining() > 0; {
		resting := e.Value.(*order.Order)

		qty := min(incoming.Remaining(), resting.Remaining())

		incoming.Fill(qty)
		resting.Fill(qty)
		if resting.IsFilled() {
			next := e.Next()
			pl.orders.Remove(e)
			e = next
		} else {

			break
		}
	}
}
