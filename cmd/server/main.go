package main

import "github.com/gin-gonic/gin"

func main() {

}

func run() {
	router := gin.Default()

	err := router.Run()
	if err != nil {
		return
	}
}
