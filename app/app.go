package app

type app struct {
	errChan chan error
}

type App interface {
}

func New() (App, error) {
	a := &app{}
	a.errChan = make(chan error)

	return a, nil
}

func (a *app) Start() {
	// run controller
	// run go to listen a.errChan and store it into file
}
