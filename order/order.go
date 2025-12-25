package order

import "time"

type Orderside int

const (
	Buy Orderside = iota
	Sell
)

type OrderStatus int

const (
	New OrderStatus = iota
	PartiallyFilled
	Filled
	Cancelled
	Rejected
)

type OrderType int

const (
	Limit OrderType = iota
	Market
)

type Order struct {
	ID        uint64
	UserID    uint64
	Symbol    string
	Side      Orderside
	Size      float64
	Type      OrderType
	Status    OrderStatus
	Price     float64
	Filled    float64
	Timestamp int64
}

func NewOrder(id uint64, userID uint64, side Orderside, price float64, size float64, ordertype OrderType) *Order {
	return &Order{
		UserID:    userID,
		ID:        id,
		Size:      size,
		Side:      side,
		Price:     price,
		Type:      ordertype,
		Timestamp: time.Now().UnixNano(),
	}
}

func (o *Order) Remaining() float64 {
	return o.Size - o.Filled
}

func (o *Order) IsFilled() bool {
	return o.Filled >= o.Size
}

func (o *Order) Fill(qty float64) {
	o.Filled += qty
	if o.Filled >= o.Size {
		o.Status = Filled
	} else {
		o.Status = PartiallyFilled
	}
}

func (o *Order) IsActive() bool {
	return o.Status == New || o.Status == PartiallyFilled
}

func (o *Order) IsBuy() bool  { return o.Side == Buy }
func (o *Order) IsSell() bool { return o.Side == Sell }
func (o *Order) IsMarket() bool {
	return o.Type == Market
}
func (o *Order) IsLimit() bool {
	return o.Type == Limit
}

func (o *Order) Cancel() bool {
	if o.Status == Filled || o.Status == Cancelled {
		return false
	}
	o.Status = Cancelled
	return true
}

func (o *Order) CanMatch(price float64) bool {
	if o.Type == Market {
		return true
	}
	if o.Side == Buy {
		return o.Price >= price
	}
	return o.Price <= price
}
