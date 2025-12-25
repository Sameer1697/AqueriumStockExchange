package matchingengine

import (
	"fmt"
	"testing"
	"time"

	"goStockExchange/order"
)

/* =========================
   HELPERS
========================= */

func startEngine(me *MatchingEngine) {
	go me.Run()
}

func stopEngine(me *MatchingEngine) {
	me.Stop()
	time.Sleep(10 * time.Millisecond) // allow goroutine to exit
}

func wait() {
	time.Sleep(10 * time.Millisecond)
}

/* =========================
   CORE MATCHING TESTS
========================= */

func TestEngineExactMatch(t *testing.T) {
	me := NewMatchingEngine("TEST", 10)
	startEngine(me)
	defer stopEngine(me)

	sell := order.NewOrder(1, 1, order.Sell, 100, 10, order.Limit)
	buy := order.NewOrder(2, 2, order.Buy, 100, 10, order.Limit)

	me.SubmitOrder(sell)
	me.SubmitOrder(buy)

	wait()

	if sell.Remaining() != 0 || buy.Remaining() != 0 {
		t.Fatalf("exact match failed sell=%f buy=%f",
			sell.Remaining(), buy.Remaining())
	}
}

func TestEnginePartialFill(t *testing.T) {
	me := NewMatchingEngine("TEST", 10)
	startEngine(me)
	defer stopEngine(me)

	sell := order.NewOrder(1, 1, order.Sell, 100, 10, order.Limit)
	buy := order.NewOrder(2, 2, order.Buy, 100, 4, order.Limit)

	me.SubmitOrder(sell)
	me.SubmitOrder(buy)

	wait()

	if sell.Remaining() != 6 {
		t.Fatalf("expected sell remaining 6, got %f", sell.Remaining())
	}
	if buy.Remaining() != 0 {
		t.Fatalf("buy should be filled")
	}
}

/* =========================
   FIFO TEST
========================= */

func TestEngineFIFO(t *testing.T) {
	me := NewMatchingEngine("TEST", 10)
	startEngine(me)
	defer stopEngine(me)

	a1 := order.NewOrder(1, 1, order.Sell, 100, 5, order.Limit)
	a2 := order.NewOrder(2, 1, order.Sell, 100, 5, order.Limit)
	b := order.NewOrder(3, 2, order.Buy, 100, 6, order.Limit)

	me.SubmitOrder(a1)
	me.SubmitOrder(a2)
	me.SubmitOrder(b)

	wait()

	if a1.Remaining() != 0 {
		t.Fatal("first order should be fully filled")
	}
	if a2.Remaining() != 4 {
		t.Fatal("second order should have 4 remaining")
	}
}

/* =========================
   NO-CROSS TEST
========================= */

func TestEngineNoCross(t *testing.T) {
	me := NewMatchingEngine("TEST", 10)
	startEngine(me)
	defer stopEngine(me)

	sell := order.NewOrder(1, 1, order.Sell, 105, 10, order.Limit)
	buy := order.NewOrder(2, 2, order.Buy, 100, 10, order.Limit)

	me.SubmitOrder(sell)
	me.SubmitOrder(buy)

	wait()

	if sell.Remaining() != 10 || buy.Remaining() != 10 {
		t.Fatal("orders should not match")
	}
}

/* =========================
   SWEEP MULTIPLE LEVELS
========================= */

func TestEngineSweep(t *testing.T) {
	me := NewMatchingEngine("TEST", 20)
	startEngine(me)
	defer stopEngine(me)

	me.SubmitOrder(order.NewOrder(1, 1, order.Sell, 100, 5, order.Limit))
	me.SubmitOrder(order.NewOrder(2, 1, order.Sell, 101, 5, order.Limit))
	me.SubmitOrder(order.NewOrder(3, 1, order.Sell, 102, 5, order.Limit))

	buy := order.NewOrder(4, 2, order.Buy, 105, 12, order.Limit)
	me.SubmitOrder(buy)

	wait()
	fmt.Printf("Final buy: filled=%f remaining=%f\n", buy.Filled, buy.Remaining())

	if buy.Remaining() != 0 {
		t.Fatal("buy should be fully filled")
	}
}

/* =========================
   CANCEL ORDER TEST
========================= */

func TestEngineCancelOrder(t *testing.T) {
	me := NewMatchingEngine("TEST", 10)
	startEngine(me)
	defer stopEngine(me)

	sell := order.NewOrder(1, 1, order.Sell, 100, 10, order.Limit)
	me.SubmitOrder(sell)

	wait()

	me.CancelOrder(1)
	wait()

	if sell.Remaining() != 10 {
		t.Fatal("cancel should not change remaining size")
	}
	if sell.Status != order.Cancelled {
		t.Fatal("order should be cancelled")
	}
}

/* =========================
   BURST TEST (BUFFERED CHANNEL)
========================= */

func TestEngineBurst(t *testing.T) {
	me := NewMatchingEngine("TEST", 1000)
	startEngine(me)
	defer stopEngine(me)

	for i := 0; i < 100; i++ {
		me.SubmitOrder(order.NewOrder(
			uint64(i+1), 1, order.Sell, 100, 1, order.Limit,
		))
	}

	buy := order.NewOrder(999, 2, order.Buy, 100, 100, order.Limit)
	me.SubmitOrder(buy)

	wait()

	if buy.Remaining() != 0 {
		t.Fatal("burst orders not fully matched")
	}
}
