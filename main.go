package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"
)

type square struct {
	Color string
	King  bool
	X     int64
	Y     int64
}

type request struct {
	Board      [][]square
	Color      string
	Iterations int
	Threads    int
}

func root(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello you are at root"))
}

func readJSON(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()

	var req request

	for key := range r.Form {
		json.Unmarshal([]byte(key), &req)
	}

	var montyMove = montecarlomove(req.Board, req.Color, req.Iterations, req.Threads)

	var jsonData, _ = json.Marshal(montyMove)

	fmt.Println(string(jsonData))
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(jsonData)

}

func main() {

	rand.Seed(time.Now().UTC().UnixNano())

	port := os.Getenv("PORT")

	if port == "" {
		fmt.Println("$PORT not set--setting to default 8080")
		port = "8080"
	}

	http.HandleFunc("/", readJSON)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		panic(err)
	}

}
