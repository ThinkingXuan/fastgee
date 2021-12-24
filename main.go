package main

import (
	"fastgee/fastgee"
	"log"
	"net/http"
	"time"
)

func onlyForV2() fastgee.HandlerFunc {
	return func(c *fastgee.Context) {
		// Start timer
		t := time.Now()
		// if a server error occurred
		c.Fail(500, "Internal Server Error")
		// Calculate resolution time
		log.Printf("[%d] %s in %v for group v2", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}

func main() {
	r := fastgee.New()
	r.Use(fastgee.Logger()) // global midlleware
	r.GET("/", func(c *fastgee.Context) {
		c.HTML(http.StatusOK, "<h1>Hello Gee</h1>")
	})

	v2 := r.Group("/v2")
	v2.Use(onlyForV2()) // v2 group middleware
	{
		v2.GET("/hello/:name", func(c *fastgee.Context) {
			// expect /hello/geektutu
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
		})
	}
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
