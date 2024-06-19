package server

import (
	"github.com/gorilla/mux"
)

type ApiServer struct {
	Port   string
	Router *mux.Router
}

func CreateServer(port string) *ApiServer {

	return &ApiServer{
		Port:   port,
		Router: mux.NewRouter(),
	}

}
