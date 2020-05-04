package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

type square struct {
	Color string
	King  bool
	X     int64
	Y     int64
}

type request struct {
	Board [][]square
	Color string
}

func readJSON(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()

	var req request

	for key := range r.Form {
		json.Unmarshal([]byte(key), &req)
	}

	var montyMove = montecarlomove(req.Board, req.Color)

	printBoard(req.Board)

	//fmt.Println(montyMove)

	var jsonData, _ = json.Marshal(montyMove)

	fmt.Println(string(jsonData))
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(jsonData)

}

func main() {

	rand.Seed(time.Now().UTC().UnixNano())

	http.HandleFunc("/", readJSON)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
