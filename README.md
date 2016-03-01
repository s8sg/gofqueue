# gofqueue
A Flip Queue implementation in Golang.

[![Build Status](https://travis-ci.org/swarvanusg/gofqueue.svg?branch=master)](https://travis-ci.org/swarvanusg/gofqueue)
[![GoDoc](https://img.shields.io/badge/api-Godoc-blue.svg?style=flat-square)](https://godoc.org/github.com/swarvanusg/gofqueue)


#### What is Flip Queue:
Flip queue is a lock less producer consumer queue that allows the producer and consumer concurrently use the same queue without any synchronization. It allows better performance in use cases where data published continuously and collected rapidly

#### How it works:
Flip queue maintains two separate queue for producer and consumer. As per the requirement of consumption of data it flips the queue.

**Flip :** A phenomena where queue changes its role. That is the producer queue becomes consumer queue and vice versa. 

#### gofqueue Interfaces:
gofqueue provides basic functionality of a queue. i.e Create, Insert and Get. 

###### Import Package
```go
  import "github.com/swarvanusg/gofqueue"
```
###### Create FQueue
Createfqueue() takes the max length of the queue. For making the dynamic length ```0``` should be specified 
```go
  fqeueu := gofqueue.Createfqueue(<length>)
```
###### Insert Into FQueue
Insert takes one param of type ```interface{}``` to accept any type of object
```go
  err := fqueue.Insert(<value>)
```
###### Get From FQueue
Get returns the first-most added data
```gp
  data, err := fqueue.Get()
```
###### Get All Data From FQueue
Getall returns all the data as per the FIFO order. In concurrent queueing it tries to get as max added data, which sometimes make the call blocking if you have a nonstop producer
```gp
  data, err := fqueue.Getall()
```
###### Launch a puiblisher for the queue
A publisher could be very useful to dump data of the queue after a regular interval. It uses ```Publish``` interface type to call a ```publish()``` with all data after a specified interval. If ```0``` is specified the interval is set to default ```1 sec```. If queue is empty, call to ```publish()``` after the interval is ignored. Publisher doesn't carry out any locking, user should not make any explicit call to retrive data, or launch multiple publisher. 
```go
  fqueue.Startpublish(<interval>, publisher)
```
###### Stop a publisher for the queue
A publiser need to be stopped explicitly by calling Stoppublish(). After stopping publisher could be re launched again
```go
  fqueue.Stoppublish()
```

#### Note:
It is only supported in single consumer and single producer model. In order to use same queue in multiple consumer and producer concurrently, explicit locking should be handled.

#### Current Status:
The flip queue build is success

The test cases are passing 


###### For any queries or concern add issues or write to : swarvanusg@gmail.com
