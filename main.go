package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const (
	fileDir     = "../tmp/"
	minfileSize = 1024 * 10         // 10KB
	maxFileSize = 1024 * 1024 * 100 // 100MB
)

var (
	writeDurations              []map[string]interface{}
	readDurations               []map[string]interface{}
	finalDurations              []map[string]interface{}
	writeDuration, readDuration float64
	filesizeInKB                int
	csvString                   string
	filePath                    string
)

func writeProcess(fileSize int) {
	filePath = filepath.Join(fileDir, fmt.Sprintf("file-%d", fileSize))

	writeDurations = make([]map[string]interface{}, 0)
	testData := make([]byte, fileSize)

	fmt.Println("Start file write")
	startTime := time.Now()
	os.WriteFile(filePath, testData, os.ModePerm)
	writeDuration = time.Since(startTime).Seconds() * 1000

	fmt.Printf("FileSize (KB): %d, AvgDuration (ms): %f\n", fileSize/1024, writeDuration)
}

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Connection successful.")
		name, err := os.Hostname()
		if err != nil {
			fmt.Println("Error resolving hostname:", err)
			return
		} else {
			fmt.Fprintf(w, "Connection successful to the host: %s \nUse the /file endpoint to Benchmark the File oprations", name)
		}

	})

	http.HandleFunc("/file", func(w http.ResponseWriter, r *http.Request) {
		writeProcess(minfileSize)
		fmt.Fprintf(w, "FileSize (KB): %d, AvgDuration (ms): %f\n", minfileSize/1024, writeDuration)
		writeProcess(maxFileSize)
		fmt.Fprintf(w, "FileSize (KB): %d, AvgDuration (ms): %f\n", maxFileSize/1024, writeDuration)
	})

	http.ListenAndServe(":8080", nil)
	fmt.Println("App listening in port 8080.")
}
