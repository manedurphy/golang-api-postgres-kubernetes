package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/manedurphy/golang-start/api"
	"github.com/manedurphy/golang-start/db"
)

func main() {

	db.Connect()

	router := gin.Default()

	router.GET("/people", api.GetPeople)
	router.GET("/person/:id", api.GetPerson)
	router.POST("/person", api.CreatePerson)
	router.DELETE("/person/:id", api.DeletePersion)

	log.Fatal(router.Run(":8080"))

}
