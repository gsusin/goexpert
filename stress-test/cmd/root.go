/*
Copyright © 2024 Giancarlo Susin <giancarlosusin@gmail.com>

*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	url         string
	requests    int
	concurrency int
	stats       map[int]int
)

var rootCmd = &cobra.Command{
	Use:   "--url https://google.com --requests 100 --concurrency 10",
	Short: "Stress test a web server",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		statusCode := make(chan int)
		taken := make(chan int, requests)
		stats = make(map[int]int)
		client := http.Client{}
		startTime := time.Now()
		for i := 0; i < concurrency; i++ {
			go func(url string) {
				for {
					taken <- 1
					resp, err := client.Get(url)
					if err != nil {
						log.Println("Error in Get. Waiting 10 s before resuming.")
						time.Sleep(10 * time.Second)
						<-taken
						continue
					}
					statusCode <- (*resp).StatusCode
				}
			}(url)
		}
		for i := 0; i < requests; i++ {
			stats[<-statusCode]++
			log.Printf("Received %d responses with code 200\n", stats[200])
		}
		endTime := time.Now()
		duration := endTime.Sub(startTime)
		println("Duration: ", duration.String())
		fmt.Printf("Number of requests: %d\n", requests)
		for k, v := range stats {
			fmt.Printf("Number of responses with code %d: %d\n", k, v)
		}
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVar(&url, "url", "", "Server's URL")
	rootCmd.Flags().IntVar(&requests, "requests", 1, "Number of requests")
	rootCmd.Flags().IntVar(&concurrency, "concurrency", 1, "Number of concurrent calls")

}
