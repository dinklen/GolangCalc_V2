package application

type Application struct {
}

func NewApplication() *Application {
	return &Application{}
}

func (app *Application) Run() error {
	return nil
}
