package handler

import (
	"github.com/gofiber/fiber/v2"
	"midtrans-forwarder/internal/domain"
	"midtrans-forwarder/internal/usecase"
)

type MidtransHandler struct {
	midtransUseCase *usecase.MidtransUseCase
}

func NewMidtransHandler(muc *usecase.MidtransUseCase, app *fiber.App) *MidtransHandler {
	h := &MidtransHandler{
		midtransUseCase: muc,
	}

	h.SetupMidtransRoutes(app)
	return h
}

func (h *MidtransHandler) SetupMidtransRoutes(app *fiber.App) {
	app.Post("/midtrans/callback", h.HandleMidtransCallback)
}

func (h *MidtransHandler) HandleMidtransCallback(c *fiber.Ctx) error {
	var callback domain.MidtransCallback
	if err := c.BodyParser(&callback); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := h.midtransUseCase.ValidateCallback(callback); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid signature"})
	}

	if err := h.midtransUseCase.ForwardToReseller(c.Context(), callback); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to forward request"})
	}

	return c.SendStatus(fiber.StatusOK)
}
