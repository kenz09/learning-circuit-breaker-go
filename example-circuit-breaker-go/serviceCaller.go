package main

import (
	"io/ioutil"
	"net/http"
)

type CircuitBreakerState interface {
	Entry() error
	Do() (*string, error)
	Exit() error
}

type CircuitBreakerClosed struct {
	httpClient      http.Client
	failureCount    int64
	failureTreshold int64
}

func (ths *CircuitBreakerClosed) Do() (*string, error) {
	resp, err := ths.httpClient.Get("localhost:8000")
	if err != nil {
		ths.failureCount++
		return nil, err
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		ths.failureCount++
		return nil, err
	}
	res := string(respBody)
	return &res, nil
}

type CircuitBreakerOpen struct {
}

type CircuitBreakerHalfOpen struct {
}

type ServiceCallerGateway struct {
}
