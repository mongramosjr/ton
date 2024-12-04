package worker

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// PaymentManageWorkers dynamically manages the pool of workers based on load and worker status.
func PaymentManageWorkers(requests <-chan PaymentRequest, wg *sync.WaitGroup, workers *int32) {
	workerStatus := make(map[int]chan WorkerStatus)

	for {
		// listen for incoming request
		if len(requests) > 0 {
			idleWorkerFound := false
			for id, statusCh := range workerStatus {
				select {
				case status := <-statusCh:
					if status == IDLE {
						go func(id int, statusCh chan<- WorkerStatus) {
							PaymentWorker(id, requests, wg, statusCh)
						}(id, statusCh)
						idleWorkerFound = true
						continue
					} else {
						// Put status back if not IDLE
						statusCh <- status
					}
				default:
					// If the status channel is empty, keep checking
				}
			}

			if !idleWorkerFound && len(requests) > 0 {
				newWorkerID := int(atomic.AddInt32(workers, 1))
				wg.Add(1)
				statusCh := make(chan WorkerStatus)
				workerStatus[newWorkerID] = statusCh
				go PaymentWorker(newWorkerID, requests, wg, statusCh)
				fmt.Printf("Added new worker. Total workers: %d\n", newWorkerID+1)
			}
		}
		time.Sleep(500 * time.Millisecond) // Check every half second
	}
}
