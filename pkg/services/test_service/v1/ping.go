package v1

import "github.com/gofiber/fiber/v2"

func (this *Test_service) Ping(c *fiber.Ctx) error {
	return c.Send([]byte("ping"))
}
