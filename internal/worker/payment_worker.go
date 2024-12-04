package worker

import (
	"fmt"
	"sync"
	"time"

	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/tlb"
	// "github.com/xssnick/tonutils-go/ton"
	// "github.com/xssnick/tonutils-go/ton/wallet"
)

// PaymentRequest
type PaymentRequest struct {
	ID           int
	Amount       *tlb.Coins
	ReceiverAddr *address.Address
}

// WorkerStatus represents the state of a worker
type WorkerStatus int

const (
	IDLE WorkerStatus = iota
	RUNNING
	DONE
)

// PaymentWorker simulates processing a request with status tracking
func PaymentWorker(id int, requests <-chan PaymentRequest, wg *sync.WaitGroup, statusCh chan<- WorkerStatus) {
	defer wg.Done()
	status := IDLE
	statusCh <- status // Initial status

	for {
		select {
		case req, ok := <-requests:
			if !ok {
				// Channel closed, worker is done
				status = DONE
				statusCh <- status
				return
			}
			// Start processing
			status = RUNNING
			statusCh <- status
			fmt.Printf("PaymentWorker is processing %d an amount of %d from %s\n", id, req.Amount.Decimals(), req.ReceiverAddr)
			time.Sleep(10 * time.Second) // simulate work time of 10 seconds
			status = IDLE
			statusCh <- status
		default:
			// PaymentWorker is idle, waiting for requests
			if status != IDLE {
				status = IDLE
				statusCh <- status
			}
			time.Sleep(100 * time.Millisecond) // sleep to not hog CPU while idle
		}
	}
}
