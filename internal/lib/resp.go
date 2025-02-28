package lib

import "github.com/gofiber/fiber/v3"

func TaskResponse(msg string, id int) fiber.Map {
	return fiber.Map{
		"msg":    msg,
		"taskID": id,
	}
}
