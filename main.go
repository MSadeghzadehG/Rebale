package main

import (
	"fmt"

	"net/http"

	"github.com/redis/go-redis/v9"
)

func Rserver() {
	c := redis.NewClient(&redis.Options{PoolSize: 50})
	s := &rServer{c}

	http.HandleFunc("/get", s.get)
	http.HandleFunc("/set", s.set)

	http.ListenAndServe(":8088", nil)
}

func main() {
	Rserver()
	c := &MyRebaleImpl{}
	c.Connect("127.0.0.1:6379", 50)
	fmt.Printf("redis connected\n")
	s := &RebaleServer{c}

	http.HandleFunc("/get", s.get)
	http.HandleFunc("/set", s.set)

	http.ListenAndServe(":8088", nil)
	fmt.Println("here")
}
