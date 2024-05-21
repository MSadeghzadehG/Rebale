package main

import (
	"fmt"
	"io"
	"net/http"
)

type Server interface {
	get(w http.ResponseWriter, req *http.Request)
	set(w http.ResponseWriter, req *http.Request)
}

type myServer struct {
	c Rebale
}

func (s *myServer) get(w http.ResponseWriter, req *http.Request) {
	key := req.URL.Query().Get("key")
	r, err := s.c.Get(key)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(200)
	io.Copy(w, r)
}

func (s *myServer) set(w http.ResponseWriter, req *http.Request) {
	key := req.URL.Query().Get("key")
	err := s.c.Set(key, req.Body, int(req.ContentLength))
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(500)
		return
	}
	req.Body.Close()
	w.WriteHeader(201)
	w.Write([]byte(""))
}
