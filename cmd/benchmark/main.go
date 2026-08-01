package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/Go8089/govpn/internal/transport/udp"
)

func main() {
	server := flag.String("server", "127.0.0.1:51820", "VPN server address")
	count := flag.Int("count", 100, "number of encrypted ping requests")

	flag.Parse()

	client := udp.NewClient(*server)

	fmt.Println("========== GoVPN Benchmark ==========")
	fmt.Printf("Server           : %s\n", *server)
	fmt.Printf("Packets Requested: %d\n", *count)
	fmt.Println("-------------------------------------")

	start := time.Now()

	var (
		success    int
		failed     int
		totalDelay time.Duration
		minDelay   time.Duration
		maxDelay   time.Duration
	)

	for i := 0; i < *count; i++ {
		packetStart := time.Now()

		err := client.Send("PING")
		latency := time.Since(packetStart)

		if err != nil {
			failed++
			log.Printf("packet %d failed: %v\n", i+1, err)
			continue
		}

		success++

		totalDelay += latency

		if success == 1 || latency < minDelay {
			minDelay = latency
		}

		if latency > maxDelay {
			maxDelay = latency
		}
	}

	totalTime := time.Since(start)

	var avgDelay time.Duration
	if success > 0 {
		avgDelay = totalDelay / time.Duration(success)
	}

	rps := 0.0
	if totalTime > 0 {
		rps = float64(success) / totalTime.Seconds()
	}

	successRate := 0.0
	if *count > 0 {
		successRate = float64(success) * 100 / float64(*count)
	}

	fmt.Println()
	fmt.Println("============= Results =============")
	fmt.Printf("Packets Sent     : %d\n", *count)
	fmt.Printf("Packets Received : %d\n", success)
	fmt.Printf("Packets Failed   : %d\n", failed)
	fmt.Printf("Success Rate     : %.2f%%\n", successRate)
	fmt.Println()

	fmt.Printf("Average Latency  : %v\n", avgDelay)
	fmt.Printf("Minimum Latency  : %v\n", minDelay)
	fmt.Printf("Maximum Latency  : %v\n", maxDelay)
	fmt.Printf("Total Time       : %v\n", totalTime)
	fmt.Printf("Requests/Second  : %.2f\n", rps)

	fmt.Println("===================================")
}
