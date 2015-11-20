# gofqueue
A Flip Queue implementation in Golang.

[GoDoc](https://godoc.org/github.com/swarvanusg/gofqueue) 


#### What is Flip Queue:
Flip queue is a lock less producer consumer queue that allows the producer and consumer concurrently use the same queue without any synchronization. It allows better performance in use cases where data in continious and collected rapidly

#### How it works:
Flip queue maintaines two sperate queue for producer and consumer. As per the requirment of consumtion of data it flips the queue.

**Flip :** A phenomena where queue changes its role. That is  producer queue becomes consumer queue and vice versa. 

#### gofqueue Interfaces:
gofqueue provides basic functionality of a queue. i.e Create, Insert and Get. 

###### Import Package
```
  import "github.com/swarvanusg/gofqueue"
```
###### Create FQueue
CreateFQueue() takes the max length of the queue. For making the dynamic length '0' should be specified 
```
  fqeueu := gofqueue.CreateFQueue(<length>)
```
###### Insert Into FQueue
Insert takes one param of type ```interface{}``` to accept any type of object
```
  err := fqueue.Insert(<value>)
```
###### Get From FQueue
Get returns the first most added data
```
  data, err := fqueue.Get()
```

#### Note:
It is only supported in single consumer and single producer model. In order to use same queue in multiple consumer and producer concurrently, explicite locking should be handled.

#### Current Status:
The flip queue build is success

The test cases are passing 


###### For any queries or concern add issues or write to : swarvanusg@gmail.com
