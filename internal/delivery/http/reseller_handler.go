package handler

import (
    "github.com/gofiber/fiber/v2"
    "strconv"
    "midtrans-forwarder/internal/domain"
    "midtrans-forwarder/internal/usecase"
)

type ResellerHandler struct {
    resellerUseCase *usecase.ResellerUseCase
}

func NewResellerHandler(resellerUseCase *usecase.ResellerUseCase, app *fiber.App) *ResellerHandler {
    h := &ResellerHandler{resellerUseCase: resellerUseCase}
    h.RegisterResellerRoutes(app)
    return h
}

func (h *ResellerHandler) RegisterResellerRoutes(app *fiber.App) {
    app.Post("/resellers", h.CreateReseller)
    app.Get("/resellers/:id", h.GetResellerByID)
    app.Get("/resellers", h.GetAllResellers)
    app.Put("/resellers/:id", h.UpdateReseller)
    app.Delete("/resellers/:id", h.DeleteReseller)
}

func (h *ResellerHandler) CreateReseller(c *fiber.Ctx) error {
    var reseller domain.Reseller
    if err := c.BodyParser(&reseller); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
    }
    err := h.resellerUseCase.CreateReseller(c.Context(), &reseller)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
    }
    return c.Status(fiber.StatusCreated).JSON(reseller)
}

func (h *ResellerHandler) GetResellerByID(c *fiber.Ctx) error {
    id, err := strconv.ParseInt(c.Params("id"), 10, 64)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid reseller ID"})
    }
    reseller, err := h.resellerUseCase.GetResellerByID(c.Context(), id)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
    }
    return c.JSON(reseller)
}

func (h *ResellerHandler) GetAllResellers(c *fiber.Ctx) error {
    resellers, err := h.resellerUseCase.GetAllResellers(c.Context())
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
    }
    return c.JSON(resellers)
}

func (h *ResellerHandler) UpdateReseller(c *fiber.Ctx) error {
    id, err := strconv.ParseInt(c.Params("id"), 10, 64)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid reseller ID"})
    }

    var reseller domain.Reseller
    if err := c.BodyParser(&reseller); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
    }
    reseller.ID = id

    err = h.resellerUseCase.UpdateReseller(c.Context(), &reseller)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
    }
    return c.JSON(reseller)
}

func (h *ResellerHandler) DeleteReseller(c *fiber.Ctx) error {
    id, err := strconv.ParseInt(c.Params("id"), 10, 64)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid reseller ID"})
    }
    err = h.resellerUseCase.DeleteReseller(c.Context(), id)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
    }
    return c.SendStatus(fiber.StatusNoContent)
}