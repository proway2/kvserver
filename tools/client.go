package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/http/httptrace"
	"net/url"
	"sync"
	"time"
)

const (
	baseURL           = "http://localhost:8080/key"
	defaultNumWorkers = 3
	defaultTotalReqs  = 1000
)

func main() {
	numWorkers := flag.Int("workers", defaultNumWorkers, "Number of workers (concurrency) running in parallel")
	totalReqs := flag.Int("reqs", defaultTotalReqs, "Requests per worker")
	flag.Parse()

	var wg sync.WaitGroup
	wg.Add(*numWorkers)
	fmt.Printf("Running with concurrency=%v\n", *numWorkers)
	for i := 0; i < *numWorkers; i++ {
		go worker(i, *totalReqs, &wg)
	}

	wg.Wait()
	fmt.Println("All workers completed")
}

func worker(id, numReqs int, wg *sync.WaitGroup) {
	defer wg.Done()

	client := &http.Client{Timeout: 10 * time.Second}
	fmt.Printf("Starting worker %v with %v requests ...\n", id, numReqs)
	for i := 0; i < numReqs; i++ {
		key := fmt.Sprintf("key_%d_%d", id, rand.Intn(1000))
		value := fmt.Sprintf("%d", time.Now().Unix())

		// Store/update value (form data)
		_, err := storeValue(client, key, value)
		if err != nil {
			fmt.Printf("Worker %d: error storing key %s: %v\n", id, key, err)
			continue
		}

		// Get value
		if err := getValue(client, key, value); err != nil {
			fmt.Printf("Worker %d: error getting key %s: %v\n", id, key, err)
		}

		// Random delete (10% chance)
		if rand.Intn(10) == 0 {
			if err := deleteValue(client, key); err != nil {
				fmt.Printf("Worker %d: error deleting key %s: %v\n", id, key, err)
			}
		}

		if i%100 == 0 {
			fmt.Printf("Worker %d: completed %d requests\n", id, i)
		}
	}
}

func storeValue(client *http.Client, key, value string) (time.Duration, error) {
	api_url := baseURL + "/" + key
	formData := url.Values{}
	formData.Set("value", value)

	var t0 time.Time
	trace := &httptrace.ClientTrace{
		DNSStart: func(_ httptrace.DNSStartInfo) { t0 = time.Now() },
	}

	req, _ := http.NewRequest(http.MethodPost, api_url, bytes.NewBufferString(formData.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))
	resp, err := client.Do(req)
	if err != nil {
		return time.Duration(0), err
	}
	resp.Body.Close()
	diff := time.Since(t0)

	if resp.StatusCode != http.StatusOK {
		return time.Duration(0), fmt.Errorf("status %d", resp.StatusCode)
	}

	return diff, nil
}

func getValue(client *http.Client, key, refValue string) error {
	resp, err := client.Get(baseURL + "/" + key)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	value, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if string(value) != refValue {
		return fmt.Errorf("value incorrect, expected=%v, got=%v", refValue, value)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status %d", resp.StatusCode)
	}
	return nil
}

func deleteValue(client *http.Client, key string) error {
	req, err := http.NewRequest("POST", baseURL+"/"+key, nil)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status %d", resp.StatusCode)
	}
	return nil
}
