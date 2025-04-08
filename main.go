package main

import (
	"log"

	"github.com/rethinkdb/prometheus-exporter/cmd"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		log.Fatalf("error starting rethinkdb exporter: %v", err)
	}
}
