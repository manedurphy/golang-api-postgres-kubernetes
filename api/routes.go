package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/manedurphy/golang-start/db"
	"github.com/manedurphy/golang-start/db/models"
)

type errMsg struct {
	MSG string `json:"message"`
}

func GetPeople(c *gin.Context) {
	var people []models.Person
	db.DB.Find(&people)
	c.JSON(200, people)
}

func GetPerson(c *gin.Context) {
	var person models.Person
	id := c.Param("id")
	db.DB.First(&person, "id = ?", id)

	if person.ID == 0 {
		c.JSON(404, errMsg{MSG: "did not find person with that ID"})
		return
	}

	c.JSON(200, person)
}

func DeletePersion(c *gin.Context) {
	id := c.Param("id")
	db.DB.Delete(&models.Person{}, id)

	c.JSON(204, "")
}

func CreatePerson(c *gin.Context) {
	person, _ := ioutil.ReadAll(c.Request.Body)

	var p models.Person
	JsonErr := json.Unmarshal([]byte(person), &p)

	if JsonErr != nil {
		fmt.Println(JsonErr)
	}

	if !p.Degree {
		p.Degree = false
	}

	v := validator.New()
	validationErr := v.Struct(p)

	if validationErr != nil {

		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	db.DB.Create(&p)

	c.JSON(201, p)
}
