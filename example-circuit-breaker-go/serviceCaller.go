package main

import (
	"errors"
	"fmt"
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
	lastError                      error
	state                          string
	failureCounter                 int
	successCounter                 int
	failureTreshold                int
	successTreshold                int
	allowedHalfOpenRequestTreshold int
}

func (this *CircuitBreakerImpl) Execute(f func() (interface{}, error)) (interface{}, error) {
	if !this.IsClosed() {
		return nil, errors.New("Circuit breaker tripped!")
	}
	r, e := f()
	if e != nil {
		this.failureCounter++
		if this.failureCounter > this.failureTreshold {
			this.Trip(e)
		}
	}

	return r, e
}

func (this *CircuitBreakerImpl) State() string {
	return this.state
}

func (this *CircuitBreakerImpl) Reset() {
	this.state = "closed"
	this.failureCounter = 0
}

func (this *CircuitBreakerImpl) HalfOpen() {
	this.state = "half-open"
}

func (this *CircuitBreakerImpl) Trip(err error) {
	fmt.Println("circuit tripped")
	fmt.Println(err.Error())
	this.lastError = err
	this.state = "open"
	this.failureCounter = 0
	this.successCounter = 0
}

func (this *CircuitBreakerImpl) IsClosed() bool {
	return this.state == "closed"
}

func (this *CircuitBreakerImpl) LastError() error {
	return this.lastError
}
