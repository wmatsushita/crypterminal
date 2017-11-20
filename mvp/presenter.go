package mvp

type (
	Presenter interface {
		ProcessUiEvent(event Event)
	}

	PortfolioPresenter struct {
	}
)