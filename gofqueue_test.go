package gofqueue

import (
	"fmt"
	"testing"
)

type data struct {
	value int
}

func TestFqueueCreation(t *testing.T) {
	// test with length 100
	fqueue := CreateFQueue(100)
	if fqueue == nil {
		t.Errorf("Failed to create fqueue")
	}
	fmt.Printf("FQueue created of length : %d\n", fqueue.Length)

	// test with length 0
	fqueue = CreateFQueue(0)
	if fqueue == nil {
		t.Errorf("Failed to create fqueue")
	}
	fmt.Printf("FQueue created of length : %d\n", fqueue.Length)
}

func TestFqueueDataInsertion(t *testing.T) {
	var values [6]*data = [...]*data{&data{1}, &data{2}, &data{3}, &data{4}, &data{5}, &data{6}}

	// Create an fqueue with length 5
	fqueue := CreateFQueue(5)
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
	var values [5]*data = [...]*data{&data{1}, &data{2}, &data{3}, &data{4}, &data{5}}

	// Create an fqueue with length 5
	fqueue := CreateFQueue(5)
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
		val := (raw.(*data)).value
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
		val := (raw.(*data)).value
		fmt.Printf("%dth : %d\n", i, val)
	}

	// try to get One more value
	val, err := fqueue.Get()
	if err == nil {
		t.Errorf("Data more than inserted can't be retrived: %d", val)
	}
}
