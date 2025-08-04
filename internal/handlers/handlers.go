package handlers

import (
	"fmt"
	"scorer/internal/models"

	"github.com/gofiber/fiber/v2"
)

type servicesInterface interface {
	CalculateScore(req *models.Request) int
}

type Handler struct {
	service servicesInterface
}

func NewHandler(s servicesInterface) Handler {
	return Handler{service: s}
}

func (h Handler) HandleScore(c *fiber.Ctx) error {
	//parse json
	req := new(models.Request)
	if err := c.BodyParser(req); err != nil {
		return c.Status(400).SendString("couldn't parse body")
	}

	//check if request fully parsed, correct request
	err := req.Validate()
	if err != nil {
		return c.Status(400).SendString("Incorrect input, no parsing")
		//400 bad request
	}

	//compute the score
	score := h.service.CalculateScore(req)	

	//post the score
	

	return c.Status(200).SendString(fmt.Sprintf("%d sent to que", score))
}
