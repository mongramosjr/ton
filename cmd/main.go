package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/tlb"

	"ton/internal/worker"
)

func main() {
	// Initial number of workers
	const initialWorkers = 8

	var (
		requests = make(chan worker.PaymentRequest, 100) // Buffered channel for requests
		wg       sync.WaitGroup
		workers  int32
	)

	// Spawn initial workers
	for i := 0; i < initialWorkers; i++ {
		wg.Add(1)
		statusCh := make(chan worker.WorkerStatus)
		go worker.PaymentWorker(i, requests, &wg, statusCh)
		atomic.AddInt32(&workers, 1)
	}

	// Manage the worker pool
	go worker.PaymentManageWorkers(requests, &wg, &workers)

	time.Sleep(5000 * time.Millisecond)

	// incoming requests
	go func() {
		for i := 0; i < 50; i++ {
			amount, err := tlb.FromTON(strconv.FormatFloat(0.1*rand.Float64(), 'f', -1, 32))
			if err != nil {
				continue
			}
			requests <- worker.PaymentRequest{
				ID:           i,
				Amount:       &amount,
				ReceiverAddr: address.MustParseAddr("0QC_Xqe8BRH4sr58r5z3AZuYSQRH-GhA0KnSYtUtWCYzKRTx")}
			time.Sleep(100 * time.Millisecond) // assume interval of payment request is 100 ms
		}
		close(requests) // Close channel
	}()

	// Wait for workers to finish
	wg.Wait()
	fmt.Println("All requests processed.")
}
