package main

import (
	"fastgee/fastgee"
	"log"
	"net/http"
)

func main() {
	r := fastgee.New()
	r.GET("/", func(c *fastgee.Context) {
		c.HTML(http.StatusOK, "<h1>Hello FastGee</h1>")
	})
	r.GET("/hello", func(c *fastgee.Context) {
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
	})

	r.GET("/hello/:name", func(c *fastgee.Context) {
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
	})

	r.GET("/assets/*filepath", func(c *fastgee.Context) {
		c.JSON(http.StatusOK, fastgee.H{"filepath": c.Param("filepath")})
	})

	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
