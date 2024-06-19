package main

import (
	"birthdays/Db"
	"birthdays/calendar"
	"birthdays/server"
	"flag"
	"github.com/BurntSushi/toml"
	"github.com/gorilla/sessions"
	"log"
)

var (
	confPth string
)

func init() {
	flag.StringVar(&confPth, "conf pth", "configs/configs.toml", "path to configs")
}

func main() {
	flag.Parse()
	config := &AppConfig{}
	_, err := toml.DecodeFile(confPth, config)
	if err != nil {
		log.Fatal(err)
	}
	app := Application{
		config:       config,
		server:       server.CreateServer(config.ServerPort),
		storage:      Db.NewStorage(config.DbConnectionStr),
		sessionStore: sessions.NewCookieStore([]byte(config.SessionKey)),
		Calendar:     calendar.NewCalendar(),
	}

	if err = app.storage.Open(); err != nil {
		log.Fatal(err)
	}
	defer app.storage.Close()
	app.Calendar.StartCal(app.storage)

	if err = app.StartServer(); err != nil {

		log.Fatal(err)
	}

	/**/

}
