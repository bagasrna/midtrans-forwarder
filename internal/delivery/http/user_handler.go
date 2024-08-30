package handler

import (
    "github.com/gofiber/fiber/v2"
    "strconv"
    "midtrans-forwarder/internal/domain"
    "midtrans-forwarder/internal/usecase"
)

type UserHandler struct {
    userUseCase *usecase.UserUseCase
}

func NewUserHandler(userUseCase *usecase.UserUseCase, app *fiber.App) *UserHandler {
    h := &UserHandler{userUseCase: userUseCase}
    h.RegisterUserRoutes(app)
    return h
}

func (h *UserHandler) RegisterUserRoutes(app *fiber.App) {
    app.Post("/users", h.CreateUser)
    app.Get("/users/:id", h.GetUserByID)
    app.Get("/users", h.GetAllUsers)
    app.Put("/users/:id", h.UpdateUser)
    app.Delete("/users/:id", h.DeleteUser)
}

func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
    var user domain.User
    if err := c.BodyParser(&user); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
    }
    err := h.userUseCase.CreateUser(c.Context(), &user)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
    }
    return c.Status(fiber.StatusCreated).JSON(user)
}

func (h *UserHandler) GetUserByID(c *fiber.Ctx) error {
    id, err := strconv.ParseInt(c.Params("id"), 10, 64)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
    }
    user, err := h.userUseCase.GetUserByID(c.Context(), id)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
    }
    return c.JSON(user)
}

func (h *UserHandler) GetAllUsers(c *fiber.Ctx) error {
    users, err := h.userUseCase.GetAllUsers(c.Context())
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
    }
    return c.JSON(users)
}

func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
    id, err := strconv.ParseInt(c.Params("id"), 10, 64)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
    }

    var user domain.User
    if err := c.BodyParser(&user); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
    }
    user.ID = id

    err = h.userUseCase.UpdateUser(c.Context(), &user)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
    }
    return c.JSON(user)
}

func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
    id, err := strconv.ParseInt(c.Params("id"), 10, 64)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
    }
    err = h.userUseCase.DeleteUser(c.Context(), id)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
    }
    return c.SendStatus(fiber.StatusNoContent)
}