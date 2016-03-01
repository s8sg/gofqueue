package gofqueue

import (
	"fmt"
	"testing"
	"time"
)

type Data struct {
	value int
}

func TestFqueueCreation(t *testing.T) {
	// test with length 100
	fqueue := Createfqueue(100)
	if fqueue == nil {
		t.Errorf("Failed to create fqueue")
	}
	fmt.Printf("FQueue created of length : %d\n", fqueue.Length)

	// test with length 0
	fqueue = Createfqueue(0)
	if fqueue == nil {
		t.Errorf("Failed to create fqueue")
	}
	fmt.Printf("FQueue created of length : %d\n", fqueue.Length)
}

func TestFqueueDataInsertion(t *testing.T) {
	var values [6]*Data = [...]*Data{&Data{1}, &Data{2}, &Data{3}, &Data{4}, &Data{5}, &Data{6}}

	// Create an fqueue with length 5
	fqueue := Createfqueue(5)
	if fqueue == nil {
		t.Errorf("Failed to create fqueue")
	}
	fmt.Printf("FQueue created of length : %d\n", fqueue.Length)

	// Insert 5 data from the data set to the queue
	for i := 0; i < 5; i++ {
		err := fqueue.Insert(values[i])
		if err != nil {
			t.Errorf("Failed to insert data : %v", err)
		}
	}
	// Try to insert the extra data
	err := fqueue.Insert(values[5])
	if err == nil {
		t.Errorf("Data more than length should not be inserted")
	}
}

func TestFqueueDataGet(t *testing.T) {
	var values [5]*Data = [...]*Data{&Data{1}, &Data{2}, &Data{3}, &Data{4}, &Data{5}}

	// Create an fqueue with length 5
	fqueue := Createfqueue(5)
	if fqueue == nil {
		t.Errorf("Failed to create fqueue")
	}
	fmt.Printf("FQueue created of length : %d\n", fqueue.Length)

	// Insert 5 data from the data set to the queue
	for i := 0; i < 5; i++ {
		err := fqueue.Insert(values[i])
		if err != nil {
			t.Errorf("Failed to insert data : %v", err)
		}
	}

	// Get 4 data from the queue
	for i := 0; i < 4; i++ {
		raw, err := fqueue.Get()
		if err != nil {
			t.Errorf("Failed to get %dth data : %v", i, err)
		}
		val := (raw.(*Data)).value
		if val != i+1 {
			t.Errorf("%dth value mismatched : %d", i, val)
		}
	}

	// Insert 3 more data to the queue
	for i := 0; i < 3; i++ {
		err := fqueue.Insert(values[i])
		if err != nil {
			t.Errorf("Failed to insert data : %v", err)
		}
	}

	// Get 4 more data from the queue
	for i := 0; i < 4; i++ {
		raw, err := fqueue.Get()
		if err != nil {
			t.Errorf("Failed to get %dth data : %v", i, err)
		}
		val := (raw.(*Data)).value
		fmt.Printf("%dth : %d\n", i, val)
	}

	// try to get One more value
	val, err := fqueue.Get()
	if err == nil {
		t.Errorf("Data more than inserted can't be retrived: %d", val)
	}
}

func TestFqueueGetAll(t *testing.T) {
	var values [5]*Data = [...]*Data{&Data{1}, &Data{2}, &Data{3}, &Data{4}, &Data{5}}

	// Create an fqueue with length 5
	fqueue := Createfqueue(5)
	if fqueue == nil {
		t.Errorf("Failed to create fqueue")
	}
	fmt.Printf("FQueue created of length : %d\n", fqueue.Length)

	// Insert 5 data from the data set to the queue
	for i := 0; i < 5; i++ {
		err := fqueue.Insert(values[i])
		if err != nil {
			t.Errorf("Failed to insert data : %v", err)
		}
	}

	// GetAll data from the queue
	alldata, err := fqueue.Getall()
	if err != nil {
		t.Errorf("All inserted should be retrived by getall: %v", err)
	}
	i := 0
	for _, raw := range alldata {
		val := (raw.(*Data)).value
		if val != i+1 {
			t.Errorf("%dth value mismatched : %d", i, val)
		}
		i++
	}
}

type MyPublisher struct {
	t *testing.T
}

var notifychan chan int

func (publisher *MyPublisher) Publish(data []interface{}) {
	i := 0
	for _, raw := range data {
		val := (raw.(*Data)).value
		if val != i+1 {
			publisher.t.Errorf("%dth value mismatched : %d", i, val)
		}
		i++
	}
	notifychan <- 0
}

func TestPublishData(t *testing.T) {
	var values [5]*Data = [...]*Data{&Data{1}, &Data{2}, &Data{3}, &Data{4}, &Data{5}}

	// Create an fqueue with length 50
	fqueue := Createfqueue(50)
	if fqueue == nil {
		t.Errorf("Failed to create fqueue")
	}
	fmt.Printf("FQueue created of length : %d\n", fqueue.Length)

	// Insert 5 data from the data set to the queue
	for i := 0; i < 5; i++ {
		err := fqueue.Insert(values[i])
		if err != nil {
			t.Errorf("Failed to insert data : %v", err)
		}
	}

	publisher := &MyPublisher{t}
	notifychan = make(chan int)
	// Start publisher
	fqueue.Startpublish(time.Millisecond*1000, publisher)

	<-notifychan

	// Insert 5 data from the data set to the queue
	for i := 0; i < 5; i++ {
		err := fqueue.Insert(values[i])
		if err != nil {
			t.Errorf("Failed to insert data : %v", err)
		}
	}
	<-notifychan
	fqueue.Stoppublish()
}
