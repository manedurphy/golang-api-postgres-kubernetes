package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/manedurphy/golang-start/db"
	"gorm.io/gorm"
)

type Person struct {
	gorm.Model
	Name   string `json:"name" validate:"required"`
	Age    int    `json:"age" validate:"required"`
	Degree bool   `json:"hasDegree"`
}

var people []Person

func GetPerson(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(500, gin.H{"message": "Invalid ID"})
		return
	}

	peopleLength := len(people)

	if peopleLength > 0 && id <= peopleLength && id > 0 {
		var index int = id - 1
		p := people[index]

		c.JSON(200, p)
		return
	}
	c.JSON(404, gin.H{"message": "No person with that ID"})
}

func GetPeople(c *gin.Context) {
	c.JSON(200, people)
}

func CreatePerson(c *gin.Context) {
	person, _ := ioutil.ReadAll(c.Request.Body)

	var p Person
	JsonErr := json.Unmarshal([]byte(person), &p)

	if JsonErr != nil {
		fmt.Println(JsonErr)
	}

	if p.Degree != true {
		p.Degree = false
	}

	v := validator.New()
	validationErr := v.Struct(p)

	if validationErr != nil {

		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	db.DB.Create(p)

	c.JSON(200, p)
}
