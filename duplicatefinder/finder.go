package duplicatefinder

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"sync"
)

var once sync.Once

func Worker(filename string,
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
	for {
		record, err = reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		codes <- record[1]
	}
}

func Monitor(workers *sync.WaitGroup, codes chan string) {
	workers.Wait()
	close(codes)
}

func Duplicates(codes <-chan string, done chan<- string) {
	registry := make(map[string]bool)

	for code := range codes {
		if _, ok := registry[code]; ok {
			done <- code
		}
		registry[code] = true
	}

	done <- ""
}
