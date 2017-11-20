package mvp

type EventType int

type Event struct {
	Type EventType
}

const (
	portfolioRefresh EventType = iota
	programQuit
)
