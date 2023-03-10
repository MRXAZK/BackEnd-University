package middleware

import (
	"context"
	"strings"

	"github.com/MRXAZK/BackEnd-University/initializers"
	"github.com/MRXAZK/BackEnd-University/models"
	"github.com/MRXAZK/BackEnd-University/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func DeserializeUniversity(c *fiber.Ctx) error {
	var access_token string
	authorization := c.Get("Authorization")

	if strings.HasPrefix(authorization, "Bearer ") {
		access_token = strings.TrimPrefix(authorization, "Bearer ")
	} else if c.Cookies("access_token") != "" {
		access_token = c.Cookies("access_token")
	}

	if access_token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "fail", "message": "You are not logged in"})
	}

	config, _ := initializers.LoadConfig(".")

	tokenClaims, err := utils.ValidateToken(access_token, config.JWTTokenPublicKey)
	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	ctx := context.TODO()
	universityid, err := initializers.RedisClient.Get(ctx, tokenClaims.TokenUuid).Result()
	if err == redis.Nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": "fail", "message": "Token is invalid or session has expired"})
	}

	var university models.University
	err = initializers.DB.First(&university, "id = ?", universityid).Error

	if err == gorm.ErrRecordNotFound {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": "fail", "message": "The university belonging to this token no logger exists"})
	}

	c.Locals("university", models.FilterUniversityRecord(&university))
	c.Locals("access_token_uuid", tokenClaims.TokenUuid)

	return c.Next()
}
