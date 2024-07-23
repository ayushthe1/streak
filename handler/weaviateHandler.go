package handler

import (
	"log"
	"net/http"

	"github.com/ayushthe1/streak/wv"
	"github.com/gofiber/fiber/v2"
)

func GetQueryDataFromWeaviate(c *fiber.Ctx) error {

	var data struct {
		Query string `json:"query" validate:"required"`
		From  string `json:"from" validate:"required"`
	}

	c.BodyParser(&data)

	validationErr := validate.Struct(data)
	if validationErr != nil {
		log.Println("error while validating")
		c.Status(http.StatusBadRequest)
		return c.JSON(fiber.Map{
			"error": validationErr,
		})
	}

	// get the data from weaviate
	result, err := wv.GetChatsRelatedToQuery(data.From, data.Query)

	if err != nil {
		c.Status(http.StatusBadRequest)
		return c.JSON(fiber.Map{
			"error": err,
		})
	}

	c.Status(http.StatusOK)
	return c.JSON(fiber.Map{
		"result": result,
	})
}
