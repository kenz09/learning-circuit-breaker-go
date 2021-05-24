package main

import (
	"errors"
	"fmt"
	"time"
)

type CircuitBreaker interface {
	State() string
	Execute(func()) (interface{}, error)
	Reset()
	HalfOpen()
	IsClosed()
	Trip(err error)
	LastError()
}

type CircuitBreakerImpl struct {
	openStateTimeout time.Time
	lastError        error
	state            string
	failureCounter   int
	failureTreshold  int
}

func (this *CircuitBreakerImpl) Execute(f func() (interface{}, error)) (interface{}, error) {
	if !this.IsClosed() && this.State() != "half-open" {
		return nil, errors.New("Circuit breaker tripped!")
	}
	r, e := f()
	if e != nil {
		if this.State() == "half-open" {
			this.Trip(e)
		} else {
			this.failureCounter = this.failureCounter + 1
			if this.failureCounter > this.failureTreshold {
				this.Trip(e)
			}
		}
	} else {
		this.failureCounter = 0
		if this.State() == "half-open" {
			this.Reset()
		}
	}
	if e != nil {
		fmt.Println(e.Error())
	}
	return r, e
}

func (this *CircuitBreakerImpl) State() string {
	return this.state
}

func (this *CircuitBreakerImpl) Reset() {
	this.state = "closed"
	this.failureCounter = 0
	fmt.Println("Circuit Reset")
}

func (this *CircuitBreakerImpl) HalfOpen() {
	this.state = "half-open"
	fmt.Println("Circuit state changed to half-open")
}

func (this *CircuitBreakerImpl) Trip(err error) {
	fmt.Println("circuit tripped to open")
	fmt.Println(err.Error())
	this.lastError = err
	this.state = "open"
	this.failureCounter = 0
	this.openStateTimeout = time.Now().Add(30 * time.Second)

	go func() {
		time.Sleep(5 * time.Second)
		this.HalfOpen()
	}()
}

func (this *CircuitBreakerImpl) IsClosed() bool {
	return this.state == "closed"
}

func (this *CircuitBreakerImpl) LastError() error {
	return this.lastError
}
