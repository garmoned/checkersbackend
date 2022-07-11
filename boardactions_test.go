package main

import (
	"encoding/json"
	"log"
	"os"
	"testing"
)

func Test_MonteCarlo(t *testing.T) {
	dat, err := os.ReadFile("test_board.json")
	if err != nil {
		log.Fatal(err)
	}
	var req request
	json.Unmarshal(dat, &req)
	montecarlomove(req.Board, req.Color, req.Iterations, req.Threads)
}
