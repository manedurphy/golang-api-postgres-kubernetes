package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/manedurphy/golang-start/api"
)

func main() {
	router := gin.Default()

	router.GET("/people", api.GetPeople)
	router.GET("/person/:id", api.GetPerson)
	router.POST("/person", api.CreatePerson)

	log.Fatal(router.Run(":8080"))

}
