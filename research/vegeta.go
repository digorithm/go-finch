package main

import (
	"fmt"
	"sync"
	"time"

	vegeta "github.com/tsenart/vegeta/lib"
)

func AttackUserFlowDefault(wg *sync.WaitGroup) {
	rate := uint64(1) // per second
	duration := 1 * time.Second

	targeter := vegeta.NewStaticTargeter(vegeta.Target{
		Method: "POST",
		URL:    "http://localhost:8888/users",
		Body:   GenerateSignUp(),
	})
	attacker := vegeta.NewAttacker()

	var metrics vegeta.Metrics
	for res := range attacker.Attack(targeter, rate, duration) {
		metrics.Add(res)
	}
	metrics.Close()

	fmt.Printf("99th percentile: %s\n", metrics.Latencies.P99)
	wg.Done()
}
