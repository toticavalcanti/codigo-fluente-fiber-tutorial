package controllers

import (
	"fiber-project/database"
	"fiber-project/models"

	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}
	if data["password"] != data["confirm_password"] {
		c.Status(400)
		return c.JSON(fiber.Map{
			"message": "Passwords do not match!",
		})
	}

	password, _ := bcrypt.GenerateFromPassword([]byte(data["password"]), 14)

	user := models.User{
		FirstName: data["first_name"],
		LastName:  data["last_name"],
		Email:     data["email"],
		Password:  password,
	}

	database.DB.Create(&user)
	return c.JSON(user)
}

func Login(c *fiber.Ctx) error {
	//get the request parameter
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	var user models.User
	//get user by email
	database.DB.Where("email = ?", data["email"]).First(&user)

	//user not found
	if user.Id == 0 {
		c.Status(404)
		return c.JSON(fiber.Map{
			"message": "User not found!",
		})
	}

	//incorrect password
	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(data["password"])); err != nil {
		c.Status(400)
		return c.JSON(fiber.Map{
			"message": "Incorrect password!",
		})
	}

	claims := jwt.RegisteredClaims{
		ID:        strconv.Itoa(int(user.Id)),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := jwtToken.SignedString([]byte("secret"))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(24 * time.Hour),
		HTTPOnly: true,
	}
	c.Cookie(&cookie)
	return c.JSON(fiber.Map{
		"jwt": token,
	})

}
