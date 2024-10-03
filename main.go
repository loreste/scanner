package main

import (
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"gopkg.in/yaml.v3"
)

// Structure for YAML data
type Config struct {
	URLs []string `yaml:"urls"`
}

func main() {
	config, err := readConfig("urls.yaml")
	if err != nil {
		fmt.Println("Error reading YAML:", err)
		return
	}

	var wg sync.WaitGroup
	for _, url := range config.URLs {
		wg.Add(1)
		go func(u string) {
			defer wg.Done()
			checkURL(u)
		}(url)
	}
	wg.Wait()
}

// Read the URLs from the YAML file using `os.ReadFile`
func readConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

// Function to check URL response and measure time
func checkURL(url string) {
	start := time.Now() // Start the timer
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Get(url)
	elapsed := time.Since(start) // Calculate elapsed time

	if err != nil {
		fmt.Printf("Failed to reach %s: %v (took %v)\n", url, err, elapsed)
		return
	}
	defer resp.Body.Close()

	fmt.Printf("Received %d from %s (took %v)\n", resp.StatusCode, url, elapsed)
}
