package main

import (
	"flag"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	client "valerius.me/yabt/pkg"
)

var (
	targetURL        string
	file             string
	connections      int
	numberOfRequests int
)

func init() {
	zerolog.TimeFieldFormat = time.RFC3339Nano
	flag.StringVar(&targetURL, "url", "http://localhost:8993/get", "URL to test for")
	flag.StringVar(&file, "file", "", "file to upload")
	flag.IntVar(&connections, "connections", 125, "number of connections to use")
	flag.IntVar(&numberOfRequests, "requests", 1000, "number of requests to make")
}

func main() {
	flag.Parse()
	c := client.NewClient(targetURL)
	guard := make(chan struct{}, connections)

	var wg sync.WaitGroup
	for i := 0; i < numberOfRequests; i++ {
		guard <- struct{}{}
		wg.Add(1)
		go func() {
			defer wg.Done()
			if file == "" {
				resp := c.Get()
				log.Info().Msg(fmt.Sprintf("Response: %d", resp.StatusCode()))
			} else {
				f, err := os.Open(file)
				if err != nil {
					panic(err)
				}
				resp := c.SendFile(f)
				log.Info().Msg(fmt.Sprintf("Response: %d", resp.StatusCode()))
			}
			<-guard
		}()
	}
	wg.Wait()
	log.Info().Msg("All requests completed")
}
