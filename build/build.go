package main

import (
	"log"
	"os"
	"todo/build/internal/bundle"
)

func main() {
	args := os.Args
	if len(args) < 2 {
		log.Fatalln("Missing version argument")
	}
	if err := bundle.CreateBundle(args[1]); err != nil {
		log.Fatalf("Failed to create bundle: %v\n", err)
	}
	log.Println("Bundle created")
}
