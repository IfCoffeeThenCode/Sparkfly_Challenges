package duplicatefinder

import (
	"context"
	"encoding/csv"
	"io"
	"log"
	"os"
	"sync"
)

// Challenge 1
func Dupes() {
	// Each file will be processed by a separate goroutine
	workers := &sync.WaitGroup{}

	// This could be an input to the function
	const path = "./testdata"

	filenames, err := GetFiles(path)
	if err != nil {
		log.Fatalf("Error getting list of files from %s: %s", path, err)
	}

	// Each of the worker routines will send their codes back to the `reportDuplicates` goroutine on this channel
	codes := make(chan string)

	// Block until `reportDuplicates` says we're done
	done := make(chan struct{})

	// `reportDuplicates` will use this context to cancel all worker routines in the event of a duplicate
	ctx, cancel := context.WithCancel(context.Background())

	// I'm assuming here that we're not handling an exceptionally large number of files
	for _, filename := range filenames {
		workers.Add(1)
		go extract(ctx, filename, codes, workers)
	}

	go monitor(workers, codes)

	go reportDuplicates(cancel, codes, done)

	<-done
}

// extract codes from a specified file and stop either at end of file or if a duplicate was found
func extract(ctx context.Context,
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

	for {
		record, err = reader.Read()
		if err == io.EOF {
			return
		}
		if err != nil {
			log.Fatal(err)
		}

		select {
		case <-ctx.Done():
			// `reportDuplicates` did in fact find a duplicate; quit
			return

		default:
			// Business as usual; send the code upstream
			codes <- record[1]
		}
	}
}

// monitor waits for the file readers to all finish, then closes the `codes` channel so
// that `reportDuplicates` knows to finish
func monitor(workers *sync.WaitGroup, codes chan string) {
	workers.Wait()
	close(codes)
}

// reportDuplicates keeps track of all codes seen by all file worker goroutines
func reportDuplicates(cancel context.CancelFunc, codes <-chan string, done chan<- struct{}) {
	registry := make(map[string]bool)

	for code := range codes {
		// Store each code unless we've already seen it; in which case tell
		// everyone to quit looking.
		if _, ok := registry[code]; ok {
			// NOTE: log.Fatalf() kills all goroutines when it exits; if this were just a quick
			// tool and not part of something greater, then that might be sufficient instead of
			// the complexity of cancel()'ing the goroutines manually.'
			log.Printf("Found duplicate: %s", code)
			cancel()
			done <- struct{}{}
			return
		}
		registry[code] = true
	}

	log.Printf("No Duplicates")
	done <- struct{}{}
}
