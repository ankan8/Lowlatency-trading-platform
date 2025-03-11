package service

import (
    "container/heap"
    "log"
    "time"

    
)

// Basic order struct
type OrderSide string

const (
    BUY  OrderSide = "BUY"
    SELL OrderSide = "SELL"
)

type OrderType string

const (
    MARKET OrderType = "MARKET"
    LIMIT  OrderType = "LIMIT"
)

type InMemoryOrder struct {
    OrderID    string
    UserID     string
    Symbol     string
    Side       OrderSide
    OrderType  OrderType
    Quantity   float64
    Price      float64
    Timestamp  time.Time
}

// We store two heaps: one for BUY (max-heap by price), one for SELL (min-heap by price).
type BuyHeap []InMemoryOrder
func (h BuyHeap) Len() int           { return len(h) }
func (h BuyHeap) Less(i, j int) bool { return h[i].Price > h[j].Price } // highest price first
func (h BuyHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *BuyHeap) Push(x interface{}) {
    *h = append(*h, x.(InMemoryOrder))
}
func (h *BuyHeap) Pop() interface{} {
    old := *h
    n := len(old)
    item := old[n-1]
    *h = old[0 : n-1]
    return item
}

type SellHeap []InMemoryOrder
func (h SellHeap) Len() int           { return len(h) }
func (h SellHeap) Less(i, j int) bool { return h[i].Price < h[j].Price } // lowest price first
func (h SellHeap) Swap(i, j int) { h[i], h[j] = h[j], h[i] }
func (h *SellHeap) Push(x interface{}) {
    *h = append(*h, x.(InMemoryOrder))
}
func (h *SellHeap) Pop() interface{} {
    old := *h
    n := len(old)
    item := old[n-1]
    *h = old[0 : n-1]
    return item
}

// OrderBook for a single symbol
type OrderBook struct {
    Symbol string
    Buys   *BuyHeap
    Sells  *SellHeap
}

func NewOrderBook(symbol string) *OrderBook {
    bh := &BuyHeap{}
    sh := &SellHeap{}
    heap.Init(bh)
    heap.Init(sh)
    return &OrderBook{
        Symbol: symbol,
        Buys:   bh,
        Sells:  sh,
    }
}

// PlaceMarketOrder matches immediately with the opposite side
func (ob *OrderBook) PlaceMarketOrder(o InMemoryOrder) {
    if o.Side == BUY {
        // Fill against the Sells
        for o.Quantity > 0 && ob.Sells.Len() > 0 {
            bestSell := heap.Pop(ob.Sells).(InMemoryOrder)
            // For a real matching engine, you'd determine the execution price.
            // Here, we just use bestSell.Price
            if bestSell.Price > 0 {
                fillQty := min(o.Quantity, bestSell.Quantity)
                // record a trade (fill) for fillQty at bestSell.Price
                log.Printf("[MARKET BUY FILL] %s buys %.2f of %s at %.2f\n", o.UserID, fillQty, o.Symbol, bestSell.Price)

                o.Quantity -= fillQty
                bestSell.Quantity -= fillQty
                if bestSell.Quantity > 0 {
                    // put the partially unfilled order back
                    heap.Push(ob.Sells, bestSell)
                }
            }
        }
        if o.Quantity > 0 {
            // leftover, but market orders typically fill as much as possible
            log.Printf("[MARKET BUY PARTIAL] leftover=%.2f not filled", o.Quantity)
        }
    } else {
        // SELL market => fill against Buys
        for o.Quantity > 0 && ob.Buys.Len() > 0 {
            bestBuy := heap.Pop(ob.Buys).(InMemoryOrder)
            if bestBuy.Price > 0 {
                fillQty := min(o.Quantity, bestBuy.Quantity)
                // record a trade (fill) for fillQty at bestBuy.Price
                log.Printf("[MARKET SELL FILL] %s sells %.2f of %s at %.2f\n", o.UserID, fillQty, o.Symbol, bestBuy.Price)

                o.Quantity -= fillQty
                bestBuy.Quantity -= fillQty
                if bestBuy.Quantity > 0 {
                    heap.Push(ob.Buys, bestBuy)
                }
            }
        }
        if o.Quantity > 0 {
            log.Printf("[MARKET SELL PARTIAL] leftover=%.2f not filled", o.Quantity)
        }
    }
}

// PlaceLimitOrder tries to match if it crosses the opposite side, otherwise rests
func (ob *OrderBook) PlaceLimitOrder(o InMemoryOrder) {
    if o.Side == BUY {
        // Check if we can match with best SELL
        for o.Quantity > 0 && ob.Sells.Len() > 0 {
            bestSell := heap.Pop(ob.Sells).(InMemoryOrder)
            // If bestSell.Price <= o.Price => cross
            if bestSell.Price <= o.Price {
                fillQty := min(o.Quantity, bestSell.Quantity)
                log.Printf("[LIMIT BUY FILL] %s buys %.2f of %s at %.2f\n", o.UserID, fillQty, o.Symbol, bestSell.Price)

                o.Quantity -= fillQty
                bestSell.Quantity -= fillQty
                if bestSell.Quantity > 0 {
                    heap.Push(ob.Sells, bestSell)
                    break
                }
            } else {
                // no cross, push bestSell back
                heap.Push(ob.Sells, bestSell)
                break
            }
        }
        // leftover rests in the BuyHeap
        if o.Quantity > 0 {
            heap.Push(ob.Buys, o)
            log.Printf("[LIMIT BUY REST] user=%s leftover=%.2f at price=%.2f\n", o.UserID, o.Quantity, o.Price)
        }
    } else {
        // SELL side
        for o.Quantity > 0 && ob.Buys.Len() > 0 {
            bestBuy := heap.Pop(ob.Buys).(InMemoryOrder)
            if bestBuy.Price >= o.Price {
                fillQty := min(o.Quantity, bestBuy.Quantity)
                log.Printf("[LIMIT SELL FILL] %s sells %.2f of %s at %.2f\n", o.UserID, fillQty, o.Symbol, bestBuy.Price)

                o.Quantity -= fillQty
                bestBuy.Quantity -= fillQty
                if bestBuy.Quantity > 0 {
                    heap.Push(ob.Buys, bestBuy)
                    break
                }
            } else {
                heap.Push(ob.Buys, bestBuy)
                break
            }
        }
        if o.Quantity > 0 {
            heap.Push(ob.Sells, o)
            log.Printf("[LIMIT SELL REST] user=%s leftover=%.2f at price=%.2f\n", o.UserID, o.Quantity, o.Price)
        }
    }
}

// Helper
func min(a, b float64) float64 {
    if a < b {
        return a
    }
    return b
}
