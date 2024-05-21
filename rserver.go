package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/redis/go-redis/v9"
)

type rServer struct {
	c *redis.Client
}

func (s *rServer) get(w http.ResponseWriter, req *http.Request) {
	key := req.URL.Query().Get("key")
	v, err := s.c.Get(req.Context(), key).Result()
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(200)
	w.Write([]byte(v))

}
func (s *rServer) set(w http.ResponseWriter, req *http.Request) {
	key := req.URL.Query().Get("key")
	b := new(bytes.Buffer)
	io.Copy(b, req.Body)
	req.Body.Close()
	err := s.c.Set(req.Context(), key, b.String(), 0).Err()
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(201)
	w.Write([]byte(""))
}
