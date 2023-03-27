package main

import (
	"encoding/json"
	"fmt"
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
	writeDuration, readDuration float64
	filesizeInKB                int
	csvString                   string
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

func main() {

	err := godotenv.Load()
	if err != nil {
		println(".env file does not exist. Use the environment variables set by the deployment environment")
	}

	fileDir = os.Getenv("DIR")

	http.HandleFunc("/json", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/csv")
		jsonData, err := json.MarshalIndent(sampleJson, "", " ")
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

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Connection successful.")
		name, err := os.Hostname()
		if err != nil {
			fmt.Println("Error resolving hostname:", err)
			return
		} else {
			fmt.Fprintf(w, "Connection successful to the host: %s \nUse the /file endpoint to Benchmark the File oprations \nUse the /response endpoint to get the csv string of the response of Benchmarking the File oprations\nUse the /json endpoint to get static JSON content \nUse the /externalapi endpoint to get a sample json response from an external API\n\n", name)
		}

	})

	fmt.Println("App listening in port 8080.")
	http.ListenAndServe(":8080", nil)
}

var sampleJson []map[string]interface{} = []map[string]interface{}{
	{
		"_id":        "641bdaaaf8e98d12763f22d4",
		"index":      0,
		"guid":       "b06a15d1-526a-473f-9209-e0ca22b7f90a",
		"isActive":   false,
		"balance":    "$2,086.82",
		"picture":    "http://placehold.it/32x32",
		"age":        40,
		"eyeColor":   "green",
		"name":       "Arline Dudley",
		"gender":     "female",
		"company":    "LUNCHPOD",
		"email":      "arlinedudley@lunchpod.com",
		"phone":      "+1 (987) 549-3761",
		"address":    "111 Moore Street, Walton, Texas, 8265",
		"about":      "Ad nulla eiusmod voluptate laborum in consectetur mollit pariatur officia aliqua adipisicing est ad. Quis anim aliqua exercitation sit cillum irure nisi quis labore. Dolor nulla enim elit qui amet ipsum exercitation. Eiusmod adipisicing culpa dolore duis est voluptate nulla. Qui sint irure qui irure excepteur laborum pariatur ipsum ipsum Lorem ut irure anim.\r\n",
		"registered": "2022-06-09T10:09:42 -06:-30",
		"latitude":   -78.980394,
		"longitude":  33.205681,
		"tags": []string{
			"deserunt",
			"officia",
			"id",
			"pariatur",
			"sunt",
			"excepteur",
			"consectetur",
		},
		"friends": []map[string]interface{}{
			{
				"id":   0,
				"name": "Julianne Wright",
			},
			{
				"id":   1,
				"name": "David Kane",
			},
			{
				"id":   2,
				"name": "Todd Holman",
			},
		},
		"greeting":      "Hello, Arline Dudley! You have 7 unread messages.",
		"favoriteFruit": "banana",
	},
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
