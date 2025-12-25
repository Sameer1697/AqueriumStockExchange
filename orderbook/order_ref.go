package orderbook

import (
	"container/list"
	"goStockExchange/order"
)

type OrderRef struct {
	Order *order.Order
	Level *PriceLevel
	Node  *list.Element
}
