package duplicatefinder

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
)

var once sync.Once

// Challenge 1
func Dupes() {
	var workers sync.WaitGroup

	const path = "./testdata"

	files, err := GetFiles(path)
	if err != nil {
		log.Fatalf("Error getting list of files from %s: %s", path, err)
	}

	codes := make(chan string)

	ctx, cancel := context.WithCancel(context.Background())

	for _, filename := range files {
		workers.Add(1)
		go Extract(ctx, filename, codes, &workers)
	}

	go Monitor(&workers, codes)

	done := make(chan string, 1)

	go ReportDuplicates(cancel, codes, done)

	duplicate := <-done
	if duplicate == "" {
		log.Printf("No Duplicates")
	} else {
		log.Printf("Found duplicate: %s", duplicate)
	}
}

// Extract codes from a specified file and stop either at end of file or if a duplicate was found
func Extract(ctx context.Context,
	filename string,
	codes chan<- string,
	wg *sync.WaitGroup) {
	defer wg.Done()

	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Error opening file %s: %s", filename, err)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)

	record, err := reader.Read() // read once to skip first header line
	var linecount = 1

	for {
		record, err = reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		linecount++

		select {
		case <-ctx.Done():
			fmt.Fprintf(os.Stderr, "cancelling on file %s after %d lines\n", filename, linecount)
			return
		case codes <- record[1]:
		}
	}

}

func Monitor(workers *sync.WaitGroup, codes chan string) {
	workers.Wait()
	close(codes)
}

func ReportDuplicates(cancel context.CancelFunc, codes <-chan string, done chan<- string) {
	registry := make(map[string]bool)

	for code := range codes {
		if _, ok := registry[code]; ok {
			cancel()
			done <- code
		}
		registry[code] = true
	}

	done <- ""
}
