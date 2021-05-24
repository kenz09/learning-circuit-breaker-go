package main

import (
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	cb1 := &CircuitBreakerImpl{
		failureTreshold: 3,
		state:           "closed",
	}

	cb2 := &CircuitBreakerImpl{
		failureTreshold: 3,
		state:           "closed",
	}

	r.GET("/services/service1", func(c *gin.Context) {
		res, err := cb1.Execute(func() (interface{}, error) {
			resp, e := http.Get("localhost:8001/ping")
			if e != nil {
				return nil, e
			}

			return resp, nil
		})
		if err != nil {
			c.JSON(500, err.Error())
			return
		}

		resResp := res.(*http.Response)
		resBody, err := ioutil.ReadAll(resResp.Body)
		if err != nil {
			c.JSON(500, err.Error())
			return
		}

		c.JSON(200, string(resBody))
	})
	r.GET("/services/service2", func(c *gin.Context) {
		res, err := cb2.Execute(func() (interface{}, error) {
			resp, e := http.Get("http://localhost:8002/ping")
			if e != nil {
				return nil, e
			}

			return resp, nil
		})
		if err != nil {
			c.JSON(500, err.Error())
			return
		}

		resResp := res.(*http.Response)
		resBody, err := ioutil.ReadAll(resResp.Body)
		if err != nil {
			c.JSON(500, err.Error())
			return
		}

		c.JSON(200, string(resBody))
	})
	r.Run() //listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
