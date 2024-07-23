package handler

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ayushthe1/streak/database"
	"github.com/ayushthe1/streak/models"
	util "github.com/ayushthe1/streak/utils"
	"github.com/ayushthe1/streak/wv"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var validate = validator.New()

func SignupHandler(c *fiber.Ctx) error {

	var data models.User
	var userData models.User
	if err := c.BodyParser(&data); err != nil {
		log.Println("Unable to parse bodyy")
		log.Println(err.Error())
		return err
	}

	validationErr := validate.Struct(data)
	if validationErr != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{"error": validationErr.Error()})

	}

	//Check if email already exist in database
	database.DB.Where("username=?", strings.TrimSpace(data.Username)).First(&userData)
	if userData.Id != 0 {
		c.Status(400)
		return c.JSON(fiber.Map{
			"message": "Email already exist",
		})

	}

	user := models.User{
		// FirstName: data.FirstName,
		Username: data.Username,
		// LastName:  data.LastName,
		// Email:     strings.TrimSpace(data.Email),
	}

	user.SetPassword(data.Password)
	result := database.DB.Create(&user)
	if result.Error != nil {
		log.Println(result.Error)
		return c.JSON(fiber.Map{
			"message": "Error creating user",
		})
	}

	token, _ := util.GenerateJwt(strconv.Itoa(int(user.Id)))

	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	err := wv.CreateNewTenant(user.Username)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"error creating tenant": err,
		})
	}

	return c.JSON(fiber.Map{
		"user":    user, // remove this later
		"message": "Account created successfully",
		"cookie":  cookie,
	})

}

func LoginHandler(c *fiber.Ctx) error {
	var data models.User

	if err := c.BodyParser(&data); err != nil {
		log.Println("Unable to parse body")
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{"error": err.Error()})
	}

	var user models.User
	database.DB.Where("username=?", data.Username).First(&user)
	if user.Id == 0 {
		c.Status(404)
		return c.JSON(fiber.Map{
			"message": " username doesn't exit, kindly create an account",
		})
	}
	if err := user.ComparePassword(data.Password); err != nil {
		c.Status(400)
		return c.JSON(fiber.Map{
			"message": "incorrect password",
		})
	}

	token, err := util.GenerateJwt(strconv.Itoa(int(user.Id)))
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return nil
	}

	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"user":    user, // remove this later
		"message": "Signed in Successfully",
		"cookie":  cookie,
	})

}

func LogoutHandler(c *fiber.Ctx) error {
	c.ClearCookie("jwt")
	log.Println("logged out")
	c.Status(200)
	return c.JSON(fiber.Map{
		"message": "Successully loggedout",
	})

}
