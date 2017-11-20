package common

type (

	Observable interface {
		Subscribe() chan struct{}
		Unsubscribe(chan struct{})
		Notify()
	}

)