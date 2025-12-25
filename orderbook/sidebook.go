package orderbook

import "github.com/huandu/skiplist"

// SideBook maintains price-priority ordering of PriceLevels
// Buy side  -> highest price first
// Sell side -> lowest price first
type SideBook struct {
	isBuy  bool
	levels *skiplist.SkipList
}

// NewSideBook creates a price-ordered book for one side
func NewSideBook(isBuy bool) *SideBook {
	var comparator skiplist.GreaterThanFunc

	if isBuy {
		// BUY: higher prices first (descending)
		comparator = func(a, b interface{}) int {
			priceA := a.(float64)
			priceB := b.(float64)
			if priceA > priceB {
				return -1
			} else if priceA < priceB {
				return 1
			}
			return 0
		}
	} else {
		// SELL: lower prices first (ascending)
		comparator = func(a, b interface{}) int {
			priceA := a.(float64)
			priceB := b.(float64)
			if priceA < priceB {
				return -1
			} else if priceA > priceB {
				return 1
			}
			return 0
		}
	}

	return &SideBook{
		isBuy:  isBuy,
		levels: skiplist.New(comparator),
	}
}

func (sb *SideBook) Best() *PriceLevel {
	if e := sb.levels.Front(); e != nil {
		if pl, ok := e.Value.(*PriceLevel); ok {
			return pl
		}
	}
	return nil
}

func (sb *SideBook) GetOrCreate(price float64) *PriceLevel {
	if e := sb.levels.Get(price); e != nil {
		if pl, ok := e.Value.(*PriceLevel); ok {
			return pl
		}
	}

	pl := NewPriceLevel(price)
	sb.levels.Set(price, pl)
	return pl
}

func (sb *SideBook) Remove(price float64) {
	sb.levels.Remove(price)
}

func (sb *SideBook) RemoveIfEmpty(pl *PriceLevel) {
	if pl.IsEmpty() {
		sb.levels.Remove(pl.Price)
	}
}

func (sb *SideBook) IsEmpty() bool {
	return sb.levels.Len() == 0
}

func (sb *SideBook) Levels() int {
	return sb.levels.Len()
}

func (sb *SideBook) ForEachFromBest(fn func(*PriceLevel) bool) {
	for e := sb.levels.Front(); e != nil; e = e.Next() {
		if pl, ok := e.Value.(*PriceLevel); ok {
			if !fn(pl) {
				return
			}
		}
	}
}
