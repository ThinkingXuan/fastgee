package main

import (
	"encoding/json"
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

	r.POST("/login", func(c *fastgee.Context) {
		c.JSON(http.StatusOK, fastgee.H{
			"username": c.PostForm("username"),
			"password": c.PostForm("password"),
		})
	})

	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}

}

func test(w http.ResponseWriter, req *http.Request) {
	obj := map[string]interface{}{
		"username": "fastgee",
		"password": "123456",
	}
	w.Header().Set("content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(obj); err != nil {
		http.Error(w, err.Error(), 500)
	}

}
