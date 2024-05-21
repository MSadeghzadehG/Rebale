package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"net/http"

	"github.com/redis/go-redis/v9"
)

func (s *myServer) get(w http.ResponseWriter, req *http.Request) {
	key := req.URL.Query().Get("key")
	r, err := s.c.Get(key)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(500)
		return
	}
	//fmt.Println(length)
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

type Server interface {
	get(w http.ResponseWriter, req *http.Request)
	set(w http.ResponseWriter, req *http.Request)
}

type myServer struct {
	c Rebale
}

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
	err := s.c.Set(req.Context(), key, string(b.Bytes()), 0).Err()
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(201)
	w.Write([]byte(""))
}

func main() {
	fmt.Println("hi")

	client := flag.Bool("client", false, "my-client")
	flag.Parse()

	var s Server
	if *client {
		c := &MyRebaleImpl{}
		c.Connect("127.0.0.1:6379")
		s = &myServer{c}
	} else {
		c := redis.NewClient(&redis.Options{PoolSize: 1})
		s = &rServer{c}
	}
	fmt.Println("go")

	http.HandleFunc("/get", s.get)
	http.HandleFunc("/set", s.set)

	http.ListenAndServe(":8080", nil)
}
