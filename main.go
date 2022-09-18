package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"sync"
	"time"
)

func separateErrorMessage(err error) string {
	splittedError := strings.Split(err.Error(), "")
	if pos := strings.Index(err.Error(), ":"); pos != -1 {
		return strings.Join(splittedError[pos+2:], "")
	}

	return ""
}

func handleFileErrorMessage(errorMessage string) {
	switch errorMessage {
	case "no such file or directory":
		log.Fatal("This path does not reference nothing on the memory.")
	default:
		fmt.Println("Error not identified. Please try again.")
	}
}

var buffer = make(chan struct{}, 20)

func handleDirOpening(dir string) []os.DirEntry {
	buffer <- struct{}{}
	defer func() {
		<-buffer
	}()

	items, err := os.ReadDir(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "dir error: %v", err)
	}

	return items
}

func LoopThroughDir(dir string, wg *sync.WaitGroup, filesChannel chan int64, directoriesNumber *int64) {
	time.Sleep(time.Millisecond * 100)

	defer wg.Done()
	for _, item := range handleDirOpening(dir) {
		if item.IsDir() {
			*directoriesNumber++

			wg.Add(1)
			subdir := path.Join(dir, item.Name())
			go LoopThroughDir(subdir, wg, filesChannel, directoriesNumber)
		} else {
			info, err := item.Info()
			if err != nil {
				panic(err)
			}

			filesChannel <- info.Size()
		}
	}
}

func main() {
	fileName := strings.Trim(os.Args[1], "")
	if len(fileName) > 1 && strings.Split(fileName, "")[0] == "/" {
		log.Fatalf("%s is an invalid name. Please insert another.", fileName)
	}

	info, err := os.Stat(fileName)
	if err != nil {
		handleFileErrorMessage(separateErrorMessage(err))
	} else {
		if info.IsDir() {
			fileSizes := make(chan int64)
			var wg sync.WaitGroup
			var filesNumber, dirNumber, occupiedMemory int64

			wg.Add(1)
			go LoopThroughDir(fileName, &wg, fileSizes, &dirNumber)
			go func() {
				wg.Wait()
				close(fileSizes)
			}()

			for fs := range fileSizes {
				filesNumber++
				occupiedMemory += fs
			}

			fmt.Printf("The directory %s has %d directories, %d files, and uses %fMB of memory\n",
				fileName, dirNumber, filesNumber, float64(occupiedMemory)/(1e6))
		} else {
			fmt.Printf("%s uses %.2fKB of memory\n", info.Name(), float64(info.Size())/(1e3))
		}
	}
}
