package main

import (
	"birthdays/Db"
	"birthdays/calendar"
	"birthdays/server"
	"github.com/gorilla/sessions"
	"net/http"
)

type Application struct {
	config       *AppConfig
	server       *server.ApiServer
	storage      *Db.Storage
	sessionStore sessions.Store
	Calendar     *calendar.Calendar
}

type AppConfig struct {
	ServerPort      string `toml:"server_port" `
	DbConnectionStr string `toml:"DbConnectionStr" `
	SessionKey      string `toml:"sessionKey" `
}

func (a *Application) StartServer() error {
	a.spinUpRoutes()
	return http.ListenAndServe(a.server.Port, a.server.Router)

}

func (a *Application) spinUpRoutes() {
	a.server.Router.Handle("/hello", a.sayHello()).Methods(http.MethodGet)
	a.server.Router.Handle("/users", a.InsertUser()).Methods(http.MethodPost)
	a.server.Router.Handle("/users/{id}", a.showUser()).Methods(http.MethodGet)
	a.server.Router.Handle("/sessions", a.CreateSession()).Methods(http.MethodPost)

	private := a.server.Router.PathPrefix("/private").Subrouter()
	private.Use(a.authUser)
	private.HandleFunc("/whoami", a.WhoAmI()).Methods(http.MethodGet)
	private.HandleFunc("/logout", a.LogOut()).Methods(http.MethodGet)
	private.HandleFunc("/subscribe", a.subscribe()).Methods(http.MethodPost)
	private.HandleFunc("/subscriptions", a.getUserSubscriptions()).Methods(http.MethodGet)

}
