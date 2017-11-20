package common

import (
	"github.com/emirpasic/gods/sets"
	"github.com/emirpasic/gods/sets/hashset"
)

type (
	/*
		Observable that uses channels to signal it's Observers.
	*/
	Observable interface {
		Subscribe() chan struct{}
		Unsubscribe(subscription chan struct{})
		Notify()
	}

	/*
		Observable that signals empty messages.
		The meaning of the signal depends on the use and must be known by the Observer.
		Usually means that the Observable has changed and the Observer must fetch it's current state.
	*/
	EmptySignalObservable struct {
		Subscriptions sets.Set
	}
)

func NewEmptySignalObservable() *EmptySignalObservable {
	return &EmptySignalObservable{
		Subscriptions: hashset.New(),
	}
}

func (o *EmptySignalObservable) Subscribe() chan struct{} {
	subscription := make(chan struct{}, 1)
	o.Subscriptions.Add(subscription)

	return subscription
}

func (o *EmptySignalObservable) Unsubscribe(subscription chan struct{}) {
	o.Subscriptions.Remove(subscription)
}

func (o *EmptySignalObservable) Notify() {
	for _, s := range o.Subscriptions.Values() {
		subscription, okToCast := s.(chan struct{})
		if okToCast {
			go func() { subscription <- struct{}{} }()
		}
	}
}
