package controllers

import (
	"github.com/MRXAZK/BackEnd-University/models"

	"github.com/gofiber/fiber/v2"
)

func GetMe(c *fiber.Ctx) error {
	university := c.Locals("university").(models.UniversityResponse)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "data": fiber.Map{"university": university}})
}
