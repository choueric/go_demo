package main

import (
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"sync"
)

const (
	audioUrl = "http://cdn5.lizhi.fm/audio/2015/04/08/19193551777309446_hd.mp3"
)

func parallelDownload(url, filename string, limit int) {
	var wg sync.WaitGroup
	// Get the content length from the header request
	res, _ := http.Head(url)
	maps := res.Header
	length, _ := strconv.Atoi(maps["Content-Length"][0])

	// Bytes for each Go-routine
	len_sub := length / limit
	// Get the remaining for the last request
	diff := length % limit
	// Make up a temporary array to hold the data to be written to the file
	fragments := make(map[int][]byte)
	for i := 0; i < limit; i++ {
		wg.Add(1)

		min := len_sub * i       // Min range
		max := len_sub * (i + 1) // Max range

		if i == limit-1 {
			max += diff // Add the remaining bytes in the last request
		}

		go func(min int, max int, i int) {
			client := &http.Client{}
			req, _ := http.NewRequest("GET", url, nil)
			// Add the data for the Range header of the form "bytes=0-100"
			range_header := "bytes=" + strconv.Itoa(min) + "-" + strconv.Itoa(max-1)
			req.Header.Add("Range", range_header)
			resp, _ := client.Do(req)
			defer resp.Body.Close()
			data, _ := ioutil.ReadAll(resp.Body)
			fragments[i] = data
			wg.Done()
		}(min, max, i)
	}
	wg.Wait()

	f, err := os.OpenFile(filename, os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	for i := 0; i < limit; i++ {
		f.Write(fragments[i])
	}
}

func main() {
	parallelDownload(audioUrl, "test.mp3", 10)
}
