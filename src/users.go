package main

import "net/http"

func (s *Server) HandleUsers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		return
	}
}
