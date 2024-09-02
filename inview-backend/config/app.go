package config

import "insense-local/database"

type Application struct {
	Env   *Env
	Mongo database.Client
}

func App() Application {
	app := &Application{}
	app.Env = NewEnv()
	app.Mongo = NewMongoDatabase(app.Env)
	return *app
}

func (app *Application) CloseDbConnection() {
	CloseMongoDBConnection(app.Mongo)
}
