package main

import "net/http"

func (s *Server) HandleLinks(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		s.handleLinksGet(w, r)
		return
	}
}

func (s *Server) handleLinksGet(w http.ResponseWriter, r *http.Request) {

}
