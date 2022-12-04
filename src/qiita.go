package main

import "net/http"

func (s *Server) HandleQiita(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		return
	}
}

func (s *Server) HandleQiitaGet(w http.ResponseWriter, r *http.Request) {

}
