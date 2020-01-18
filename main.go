package main

import (
	"log"
	"sync"

	"github.com/IfCoffeeThenCode/Sparkfly_Challenges/duplicatefinder"
)

// Challenge 1
func dupes() {
	var workers sync.WaitGroup

	const path = "./testdata"

	files, err := duplicatefinder.GetFiles(path)
	if err != nil {
		log.Fatalf("Error getting list of files from %s: %s", path, err)
	}

	codes := make(chan string)

	for _, filename := range files {
		workers.Add(1)
		go duplicatefinder.Worker(filename, codes, &workers)
	}

	go duplicatefinder.Monitor(&workers, codes)

	done := make(chan string, 1)

	go duplicatefinder.Duplicates(codes, done)

	duplicate := <-done
	if duplicate == "" {
		log.Printf("No Duplicates")
	} else {
		log.Printf("Found duplicate: %s", duplicate)
	}
}

func main() {
	dupes()
}
