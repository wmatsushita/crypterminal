package common

import (
	"testing"

	"sync"

	"fmt"
)

func TestObserverReceivesSignal(t *testing.T) {
	observable1 := NewEmptySignalObservable()
	observable2 := NewEmptySignalObservable()
	observer := NewEmptySignalObserver()

	signalsReceived := 0
	wg := sync.WaitGroup{}
	signals := 10
	wg.Add(signals)
	c1 := observer.Watch(observable1, func() {
		defer wg.Done()
		signalsReceived += 1
		fmt.Println("Signal received from observer1!")
	})

	c2 := observer.Watch(observable2, func() {
		defer wg.Done()
		signalsReceived += 1
		fmt.Println("Signal received from observer2!")
	})

	for i := 0; i < signals/2; i++ {
		observable1.Notify()
		observable2.Notify()
	}

	wg.Wait()

	if signalsReceived != signals {
		t.Fatalf("Expected %v signals received, got %v", signals, signalsReceived)
	}

	observer.Ignore(observable1, c1)
	observer.Ignore(observable2, c2)

}
