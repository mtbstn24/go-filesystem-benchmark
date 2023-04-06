package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

const (
	minfileSize = 1024 * 10         // 10KB
	maxFileSize = 1024 * 1024 * 100 // 100MB
)

var fileDir string

var (
	writeDurations              []map[string]interface{}
	readDurations               []map[string]interface{}
	finalDurations              []map[string]interface{}
	fibDurations                []map[string]interface{}
	writeDuration, readDuration float64
	filesizeInKB                int
	csvString                   string
	fibString                   string
	filePath                    string
	status                      bool
)

func writeProcess(fileSize int, filePath string) {

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

func readProcess(fileSize int, filePath string) {

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
	os.Remove(filePath)

	fmt.Println(readDurations)
	fmt.Printf("FileSize (KB): %d, AvgDuration (ms): %f\n", fileSize/1024, readDuration)
}

func fileProcess(filesize int) {
	filePath = filepath.Join(fileDir, fmt.Sprintf("file-%d", filesize))
	writeProcess(filesize, filePath)
	readProcess(filesize, filePath)
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
	header := []string{"FileSize (KB)", "Write Duration (ms)", "Read Duration (ms)", "Read and Write Duration (ms)"}
	rows := []string{strings.Join(header, ",")}

	for _, item := range finalDurations {
		row := []string{fmt.Sprintf("%d", item["Filesize"]), fmt.Sprintf("%f", item["WriteDuration"]), fmt.Sprintf("%f", item["ReadDuration"]), fmt.Sprintf("%f", item["ReadWriteDuration"])}
		rows = append(rows, strings.Join(row, ","))
	}
	csvString = strings.Join(rows, "\n")
	fmt.Println(csvString)
	csvfilePath := filepath.Join(fileDir, fmt.Sprintf("csvString.csv"))
	os.WriteFile(csvfilePath, []byte(csvString), os.ModePerm)
	status = true
}

func fibonacci(n int) int {
	if n <= 1 {
		return n
	}
	return fibonacci(n-1) + fibonacci(n-2)
}

func getFibString() {
	calDurationAvgs := make([]float64, 0)
	number := []int{20, 22, 24, 26, 28, 30, 32, 34, 36, 38, 40}
	fibonacciArray := make([]int, 0)
	fibDurations = make([]map[string]interface{}, 0)
	var result int
	var calDuration float64
	for i := 0; i < 11; i++ {
		calDurationSum := float64(0)
		for j := 0; j < 5; j++ {
			startTime := time.Now()
			result = fibonacci(number[i])
			calDuration = time.Since(startTime).Seconds() * 1000
			calDurationSum += calDuration
		}
		calDurationAvg := float64(calDurationSum) / 5.0
		calDurationAvgs = append(calDurationAvgs, calDurationAvg)
		fibonacciArray = append(fibonacciArray, result)
		fibDurations = append(fibDurations, map[string]interface{}{
			"FibonacciNumber": number[i],
			"FibonacciValue":  result,
			"CalDuration":     calDurationAvg,
		})
	}

	Fheader := []string{"Fibonacci Number", "Fibonacci Value", "Calculation Duration (ms)"}
	Frows := []string{strings.Join(Fheader, ",")}

	for _, item := range fibDurations {
		row := []string{fmt.Sprintf("%d", item["FibonacciNumber"]), fmt.Sprintf("%d", item["FibonacciValue"]), fmt.Sprintf("%f", item["CalDuration"])}
		Frows = append(Frows, strings.Join(row, ","))
	}
	fibString = strings.Join(Frows, "\n")
	fmt.Println(fibString)
}

func main() {

	err := godotenv.Load()
	if err != nil {
		println(".env file does not exist. Use the environment variables set by the deployment environment")
	}

	fileDir = os.Getenv("DIR")

	http.HandleFunc("/externalapi", func(w http.ResponseWriter, r *http.Request) {
		resp, err := http.Get("https://jsonplaceholder.typicode.com/users")
		if err != nil {
			fmt.Printf("Cannot fetch URL : %v", err)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/csv")
		var data []map[string]interface{}
		body, readErr := ioutil.ReadAll(resp.Body)
		if readErr != nil {
			fmt.Println(readErr)
			return
		}
		resp.Body.Close()
		err1 := json.Unmarshal(body, &data)
		if err1 != nil {
			fmt.Println(err1)
			return
		}
		jsonData, err2 := json.MarshalIndent(data, "", " ")
		if err2 != nil {
			fmt.Println(err2)
			return
		}
		fmt.Fprint(w, string(jsonData))
	})

	http.HandleFunc("/json", func(w http.ResponseWriter, r *http.Request) {
		sJson := append(sampleJson, map[string]interface{}{
			"time": time.Now(),
		})
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/csv")
		jsonData, err := json.MarshalIndent(sJson, "", " ")
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Fprint(w, string(jsonData))
	})

	http.HandleFunc("/file", func(w http.ResponseWriter, r *http.Request) {
		multipleFileProcess()
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/csv")
		w.Write([]byte(csvString))
	})

	http.HandleFunc("/response", func(w http.ResponseWriter, r *http.Request) {
		if status == true {
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "text/csv")
			w.Write([]byte(csvString))
		} else {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "Respond not found or Process not completed. \nMake a request to /file endpoint first. \nWait for some time and try again if you have already requested /file endpoint.\n")
		}
	})

	http.HandleFunc("/fibonacci", func(w http.ResponseWriter, r *http.Request) {
		getFibString()
		w.Header().Set("Content-Type", "application/csv")
		w.Write([]byte(fibString))

	})

	http.HandleFunc("/fibresponse", func(w http.ResponseWriter, r *http.Request) {
		if fibString != "" {
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "text/csv")
			w.Write([]byte(fibString))
		} else {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "Respond not found or Process not completed. \nMake a request to /fibonacci endpoint first. \nWait for some time and try again if you have already requested /file endpoint.\n")
		}
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Connection successful.")
		name, err := os.Hostname()
		if err != nil {
			fmt.Println("Error resolving hostname:", err)
			return
		} else {
			fmt.Fprintf(w, "Connection successful to the host: %s \nUse the /file endpoint to Benchmark the File oprations \nUse the /response endpoint to get the csv string of the response of Benchmarking the File oprations\nUse the /json endpoint to get static JSON content \nUse the /externalapi endpoint to get a sample json response from an external API \nUse the /fibonacci endpoint to get the 40 fibonacci numbers and Durations\nUse the /fibresponse endpoint to get the fibonacci Durations as csv\n\n", name)
		}

	})

	fmt.Println("App listening in port 8080.")
	http.ListenAndServe(":8080", nil)
}

var sampleJson []map[string]interface{} = []map[string]interface{}{
	{
		"_id":        "641bdaaa0a5b4af1047d0dde",
		"index":      1,
		"guid":       "e097b36c-3a93-443e-9a3b-17da59dafcab",
		"isActive":   false,
		"balance":    "$2,865.23",
		"picture":    "http://placehold.it/32x32",
		"age":        29,
		"eyeColor":   "brown",
		"name":       "Marina Herrera",
		"gender":     "female",
		"company":    "BLEEKO",
		"email":      "marinaherrera@bleeko.com",
		"phone":      "+1 (913) 512-2676",
		"address":    "880 Amherst Street, Kenmar, California, 5070",
		"about":      "Proident sunt magna elit duis officia in esse labore tempor ipsum id ipsum. Sunt nisi nostrud anim veniam est nisi cupidatat ut minim esse laborum elit. Do cupidatat officia reprehenderit incididunt sit eiusmod excepteur dolor commodo esse nulla. Sit aute nisi veniam cillum aliqua.\r\n",
		"registered": "2015-11-12T01:02:54 -06:-30",
		"latitude":   65.668736,
		"longitude":  53.450258,
		"tags": []string{
			"cillum",
			"do",
			"cupidatat",
			"minim",
			"do",
			"sint",
			"ullamco",
		},
		"friends": []map[string]interface{}{
			{
				"id":   0,
				"name": "Newman Hamilton",
			},
			{
				"id":   1,
				"name": "Christi Bond",
			},
			{
				"id":   2,
				"name": "Nunez Saunders",
			},
		},
		"greeting":      "Hello, Marina Herrera! You have 2 unread messages.",
		"favoriteFruit": "banana",
	},
}
