package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
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
	status                      bool
)

func writeProcess(fileSize int) {
	filePath = filepath.Join(fileDir, fmt.Sprintf("file-%d", fileSize))

	writeDurations = make([]map[string]interface{}, 0)
	testData := make([]byte, fileSize)
	var sum float64

	for i := 0; i < 10; i++ {
		startTime := time.Now()
		os.WriteFile(filePath, testData, os.ModePerm)
		writeDuration = time.Since(startTime).Seconds() * 1000

		writeDurations = append(writeDurations, map[string]interface{}{
			"size":          fileSize,
			"writeDuration": writeDuration,
		})

		sum = sum + writeDuration
	}

	writeDuration = sum / 10

	fmt.Println(writeDurations)
	fmt.Printf("FileSize (KB): %d, AvgDuration (ms): %f\n", fileSize/1024, writeDuration)
}

func readProcess(fileSize int) {
	filePath = filepath.Join(fileDir, fmt.Sprintf("file-%d", fileSize))

	readDurations = make([]map[string]interface{}, 0)
	var sum float64

	for i := 0; i < 10; i++ {
		startTime := time.Now()
		os.ReadFile(filePath)
		readDuration = time.Since(startTime).Seconds() * 1000

		readDurations = append(readDurations, map[string]interface{}{
			"size":         fileSize,
			"readDuration": readDuration,
		})

		sum = sum + readDuration
	}

	readDuration = sum / 10

	fmt.Println(readDurations)
	fmt.Printf("FileSize (KB): %d, AvgDuration (ms): %f\n", fileSize/1024, readDuration)
}

func fileProcess(filesize int) {
	writeProcess(filesize)
	readProcess(filesize)
	filesizeInKB = filesize / 1024
	finalDurations = append(finalDurations, map[string]interface{}{
		"Filesize":          filesizeInKB,
		"WriteDuration":     writeDuration,
		"ReadDuration":      readDuration,
		"ReadWriteDuration": writeDuration + readDuration,
	})
}

func multipleFileProcess() {
	status = false
	writeDurations = make([]map[string]interface{}, 0)
	readDurations = make([]map[string]interface{}, 0)
	for filesize := minfileSize; filesize <= maxFileSize; filesize = filesize + 1024*1024*2 {
		fileProcess(filesize)
	}
	fmt.Println(finalDurations)
	var csvData [][]string
	csvData = append(csvData, []string{
		"FileSize (KB)", "Write Duration (ms)", "Read Duration (ms)", "Read and Write Duration (ms)",
	})
	header := []string{"FileSize (KB)", "Write Duration (ms)", "Read Duration (ms)", "Read and Write Duration (ms)"}
	rows := []string{strings.Join(header, ",")}

	for _, item := range finalDurations {
		row := []string{fmt.Sprintf("%d", item["Filesize"]), fmt.Sprintf("%f", item["WriteDuration"]), fmt.Sprintf("%f", item["ReadDuration"])}
		rows = append(rows, strings.Join(row, ","))
	}
	csvString := strings.Join(rows, "\n")
	fmt.Println(csvString)
	csvfilePath := filepath.Join(fileDir, fmt.Sprintf("csvString.csv"))
	os.WriteFile(csvfilePath, []byte(csvString), os.ModePerm)
	status = true
}

func main() {

	http.HandleFunc("/file/", func(w http.ResponseWriter, r *http.Request) {
		multipleFileProcess()
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/csv")
		w.Write([]byte(csvString))
	})

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

	fmt.Println("App listening in port 8080.")
	http.ListenAndServe(":8080", nil)
}
