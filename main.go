package main

import (
	"fmt"
	"net/http"
	// "github.com/redis/go-redis/v9"
)

// func Rserver() {
// 	c := redis.NewClient(&redis.Options{PoolSize: 1})
// 	s := &rServer{c}

// 	http.HandleFunc("/get", s.get)
// 	http.HandleFunc("/set", s.set)

// 	http.ListenAndServe(":8080", nil)
// }

func main() {
	c := &MyRebaleImpl{}
	c.Connect("127.0.0.1:6379")
	fmt.Printf("redis connected\n")
	s := &myServer{c}

	http.HandleFunc("/get", s.get)
	http.HandleFunc("/set", s.set)

	http.ListenAndServe(":8080", nil)
}
