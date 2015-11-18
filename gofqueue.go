package gofqueue

import (
	"fmt"
	"sync"
)

const (
	defaultQLength = 100
)

// Structure to define a queue
type queue struct {
	queue  []interface{}
	length int
}

// Structure to define a FQueue
type FQueue struct {
	// The pointer to the producer queue
	producerq *queue
	// The pointer to the consumer queue
	consumerq *queue
	// The Flip sync lock
	flipLock *sync.Mutex
	// The Length of the queue
	Length int
}

// Method to create a Flip Queue
func CreateFQueue(length int) *FQueue {

	// Create the FQueue
	fqueue := &FQueue{}
	// Set length
	fqueue.Length = length
	fqueue.flipLock = &sync.Mutex{}

	// Create actual queue
	if length == 0 {
		length = defaultQLength
	}
	// Create the two queue of given length
	fqueue.producerq = &queue{}
	fqueue.producerq.queue = make([]interface{}, length)
	fqueue.producerq.length = 0
	fqueue.consumerq = &queue{}
	fqueue.consumerq.queue = make([]interface{}, length)
	fqueue.consumerq.length = 0

	return nil
}

// Method to insert data to a Flip Queue
func (fqueue *FQueue) Insert(data interface{}) error {
	// Take the lock
	fqueue.flipLock.Lock()
	defer fqueue.flipLock.Lock()

	producerq := fqueue.producerq

	// Check if the length is reached
	if fqueue.Length != 0 && producerq.length >= fqueue.Length {

		return fmt.Errorf("Reached fqueue max length")
	}
	// put the data in the location
	producerq.queue[producerq.length] = data

	// increase the length of the queue
	producerq.length++

	// Insert the data into
	return nil
}

// Method to Get data from a Flip Queue
func (fqueue *FQueue) Get() (interface{}, error) {

	consumerq := fqueue.consumerq

	// check if the length of the consumerq is more than 0
	if consumerq.length < 0 {
		fqueue.flip()
	}

	// check the length of the
	if consumerq.length < 0 {
		return nil, fmt.Errorf("Flip queue is empty")
	}

	data := consumerq.queue[consumerq.length-1]

	consumerq.length--

	return data, nil
}

// Internal method that flips a the queues functioanlity
func (fqueue *FQueue) flip() {

	var tempref *queue

	// Lock the flip queue
	fqueue.flipLock.Lock()
	defer fqueue.flipLock.Lock()

	// Swap the queue references
	tempref = fqueue.producerq
	fqueue.producerq = fqueue.consumerq
	fqueue.consumerq = tempref
}
