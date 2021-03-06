package components

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
)


var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRST0123456789")

type OrderEvent struct {
	ID string
	OrderDetail string
	CurrentState StateHandlerInterface
	Done chan struct{}
	Err error
}

func NewOrderEvent(or string) *OrderEvent {
	state := GenerateStates()
	o := &OrderEvent{
		ID: generateOrderID(),
		OrderDetail: or,
		CurrentState: state,
		Done: make(chan struct{}),
	}
	return o
}
func (oe *OrderEvent)GetStatus() string  {
	if oe.CurrentState != nil {
		return oe.CurrentState.Name()
	}
	return "all done"
}

func (o *OrderEvent)HandleEvent(op , event string) *OrderEvent {
	if o.CurrentState == nil {
		o.Err = errors.New("no procedure is ongoing")
		return o
	}
	err :=o.CurrentState.StateHandler(op, event)
	if err != nil {
		o.Err = err
	}
	if o.CurrentState.IsFinished() {
		next := o.CurrentState.Next()
		if next == nil {
			fmt.Println("Waiting for confirmation")
		}
		o.CurrentState = o.CurrentState.Next()
	}
	return o
}

func (o *OrderEvent)Close()  {
	o.CurrentState = nil
	close(o.Done)
}

func (o *OrderEvent)IsFinished() bool {
	return o.CurrentState == nil
}

// help functions
func generateOrderID() string {
	return fmt.Sprintf("order-%s", randn(10))
}

func randn(n int) string {
	res := make([]rune, n)
	seed := time.Now().UnixNano()
	rand.Seed(seed)
	for i := range res {
		r := rand.Intn(len(letters))
		res[i] = letters[r]
	}
	return string(res)
}
