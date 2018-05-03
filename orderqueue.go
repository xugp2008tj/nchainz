package main

import (
	"bytes"
	"container/heap"
	"errors"
	"fmt"
)

type OrderQueue struct {
	Items []OrderQueueItem
	Side  OrderSide
}

type OrderQueueItem struct {
	price float64
	order *Order
	index int
}

type OrderSide bool

const (
	BASE  OrderSide = true
	QUOTE OrderSide = false
)

func NewOrderQueue(side OrderSide) *OrderQueue {
	items := make([]OrderQueueItem, 0)
	return &OrderQueue{items, side}
}

func (oq *OrderQueue) Enq(order *Order) {
	item := OrderQueueItem{}
	item.order = order

	if oq.Side == QUOTE {
		item.price = float64(order.AmountToSell) / float64(order.AmountToBuy)
	} else { // oq.Side == BASE
		item.price = float64(order.AmountToBuy) / float64(order.AmountToSell)
	}

	heap.Push(oq, item)
}

func (oq *OrderQueue) Deq() (order *Order, price float64, err error) {
	if oq.Len() == 0 {
		return nil, 0.0, errors.New("empty queue")
	}

	item := heap.Pop(oq).(OrderQueueItem)
	return item.order, item.price, nil
}

func (oq *OrderQueue) Peek() (order *Order, price float64, err error) {
	if oq.Len() == 0 {
		return nil, 0.0, errors.New("empty queue")
	}

	item := oq.Items[0]
	return item.order, item.price, nil
}

func (oq *OrderQueue) String() string {
	var buffer bytes.Buffer

	for _, oqi := range oq.Items {
		buffer.WriteString(oq.oqiString(&oqi))
	}

	if oq.Side == BASE {
		return fmt.Sprintf("BASE %s", buffer.String())
	} else {
		return fmt.Sprintf("QUOTE %s", buffer.String())
	}
}

func (oq *OrderQueue) oqiString(oqi *OrderQueueItem) string {
	if oq.Side == BASE {
		return fmt.Sprintf("<%v @ %f #%v>", oqi.order.AmountToSell, oqi.price, oqi.order.ID)
	} else {
		return fmt.Sprintf("<%v @ %f #%v>", oqi.order.AmountToBuy, oqi.price, oqi.order.ID)
	}
}

/////////////////////
// Utils used by heap

func (oq *OrderQueue) Len() int {
	return len(oq.Items)
}

func (oq *OrderQueue) Less(i, j int) bool {
	if oq.Side == BASE {
		return oq.Items[i].price < oq.Items[j].price
	} else { // oq.Side == QUOTE
		return oq.Items[i].price > oq.Items[j].price
	}
}

func (oq *OrderQueue) Swap(i, j int) {
	oq.Items[i], oq.Items[j] = oq.Items[j], oq.Items[i]
	oq.Items[i].index = i
	oq.Items[j].index = j
}

func (oq *OrderQueue) Push(x interface{}) {
	item := x.(OrderQueueItem)
	item.index = len(oq.Items)
	oq.Items = append(oq.Items, item)
}

func (oq *OrderQueue) Pop() interface{} {
	item := oq.Items[len(oq.Items)-1]
	item.index = -1
	oq.Items = oq.Items[0 : len(oq.Items)-1]
	return item
}
