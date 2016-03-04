// Copyright © 2014 Swarvanu Sengupta <swarvanusg.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package gofqueue

import (
	"fmt"
	"sync"
	"time"
)

const (
	defaultQLength = 100
)

// Structure to define a queue
type queue struct {
	queue  chan interface{}
	length int
}

// Publisher defines the queue publish function
type Publisher interface {
	Publish([]interface{})
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
	// The quit channel
	publisherquitchan chan int
}

// Method to create a Flip Queue
func Createfqueue(length int) *FQueue {

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
	fqueue.producerq.queue = make(chan interface{}, length)
	fqueue.producerq.length = 0
	fqueue.consumerq = &queue{}
	fqueue.consumerq.queue = make(chan interface{}, length)
	fqueue.consumerq.length = 0

	return fqueue
}

// Method to insert data into a Flip Queue
func (fqueue *FQueue) Insert(data interface{}) error {
	// Take the lock
	fqueue.flipLock.Lock()
	defer fqueue.flipLock.Unlock()

	producerq := fqueue.producerq

	// Check if the length is reached
	if fqueue.Length != 0 && producerq.length >= fqueue.Length {
		return fmt.Errorf("Reached fqueue max length")
	}

	// put the data in the location
	producerq.queue <- data

	// increase the length of the queue
	producerq.length++

	return nil
}

// Method to Get data from a Flip Queue
func (fqueue *FQueue) Get() (interface{}, error) {

	consumerq := fqueue.consumerq

	// check if the length of the consumerq is more than 0
	if consumerq.length <= 0 {
		fqueue.flip()
	}

	consumerq = fqueue.consumerq
	// check the length of the
	if consumerq.length <= 0 {
		return nil, fmt.Errorf("Flip queue is empty")
	}

	data := <-consumerq.queue

	consumerq.length--

	return data, nil
}

/* Method to Get all data from a Flip Queue
   In a concurrent queueing it tries to get you max possible data
   which sometimes makes the call blocking if you have a nonstop producer */
func (fqueue *FQueue) Getall() ([]interface{}, error) {

	var datalist []interface{}

	consumerq := fqueue.consumerq

	// check if the length of the consumerq is more than 0
	if consumerq.length <= 0 {
		fqueue.flip()
	}

	consumerq = fqueue.consumerq
	// check the length of the
	if consumerq.length <= 0 {
		return nil, fmt.Errorf("Flip queue is empty")
	}

	fmt.Println(consumerq.length)
LOOP1:
	for {
	LOOP2:
		for {
			select {
			case data := <-consumerq.queue:
				datalist = append(datalist, data)
				consumerq.length--
			default:
				break LOOP2
			}
		}
		// Invoke a flip
		fqueue.flip()
		consumerq = fqueue.consumerq
		fmt.Println(consumerq.length)
		// cheeck if flipped queue is empty
		if consumerq.length <= 0 {
			break LOOP1
		}
	}

	return datalist, nil
}

/* Method to create publisher to publish data after a certain interval
   If interval specified by 0 it will take default interval (1 SEC) */
func (fqueue *FQueue) Startpublish(interval time.Duration, publisher Publisher) {
	if interval == 0 {
		interval = time.Second
	}
	fqueue.publisherquitchan = make(chan int)
	go func() {
		tickchan := time.Tick(interval)
	LOOP:
		for {
			select {
			case <-tickchan:
				data, err := fqueue.Getall()
				if err == nil {
					// Call the publish on publisher
					publisher.Publish(data)
				}
			case <-fqueue.publisherquitchan:
				break LOOP
			}
		}
		fqueue.publisherquitchan <- 0
	}()
}

// Method to stop the publisher method
func (fqueue *FQueue) Stoppublish() {
	// Tell the publisher to stop
	fqueue.publisherquitchan <- 0
	// Wait till publisher is really stopped
	<-fqueue.publisherquitchan
}

// Internal method that flips a the queues functioanlity
func (fqueue *FQueue) flip() {

	var tempref *queue

	// Lock the flip queue
	fqueue.flipLock.Lock()
	defer fqueue.flipLock.Unlock()

	// Swap the queue references
	tempref = fqueue.producerq
	fqueue.producerq = fqueue.consumerq
	fqueue.consumerq = tempref
}
