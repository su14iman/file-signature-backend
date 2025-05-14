package controllers

import (
	"os"
	"time"

	"file-signer/config"
	"file-signer/models"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// LoginAdmin godoc
// @Summary Login admin and get JWT
// @Tags auth
// @Accept json
// @Produce json
// @Param login body models.LoginInput true "Email and Password"
// @Success 200 {object} models.TokenResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Router /api/admin/login [post]
// @Security Bearer
func LoginAdmin(c *fiber.Ctx) error {
	var input models.LoginInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(models.ErrorResponse{Error: "Invalid input"})
	}

	var admin models.User
	result := config.DB.Where("email = ?", input.Email).First(&admin)
	if result.Error != nil {
		return c.Status(401).JSON(models.ErrorResponse{Error: "Invalid email or password"})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(input.Password)); err != nil {
		return c.Status(401).JSON(models.ErrorResponse{Error: "Invalid email or password"})
	}

	secret := os.Getenv("JWT_SECRET")
	claims := jwt.MapClaims{
		"admin_id": admin.Id,
		"role":     admin.Role,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return c.Status(500).JSON(models.ErrorResponse{Error: "Could not generate token"})
	}

	return c.JSON(models.TokenResponse{Token: signedToken})
}

// GetAdmin godoc
// @Summary Get admin by ID
// @Tags auth
// @Accept json
// @Produce json
// @Param id path int true "Admin ID"
// @Success 200 {object} models.UserResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /api/admin/{id} [get]
// @Security Bearer
func GetAdmin(c *fiber.Ctx) error {
	var admin models.User
	id := c.Params("id")
	result := config.DB.First(&admin, id)
	if result.Error != nil {
		return c.Status(404).JSON(models.ErrorResponse{Error: "Admin not found"})
	}

	return c.JSON(admin)
}

// GetAllAdmins godoc
// @Summary Get all admins
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {array} models.UserResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/admins [get]
// @Security Bearer
func GetAllAdmins(c *fiber.Ctx) error {
	var admins []models.User
	result := config.DB.Find(&admins)
	if result.Error != nil {
		return c.Status(500).JSON(models.ErrorResponse{Error: "Could not retrieve admins"})
	}

	return c.JSON(admins)
}

// DeleteAdmin godoc
// @Summary Delete admin by ID
// @Tags auth
// @Accept json
// @Produce json
// @Param id path int true "Admin ID"
// @Success 200 {object} models.SuccessResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/admin/{id} [delete]
// @Security Bearer
func DeleteAdmin(c *fiber.Ctx) error {
	id := c.Params("id")
	var admin models.User
	result := config.DB.First(&admin, id)
	if result.Error != nil {
		return c.Status(404).JSON(models.ErrorResponse{Error: "Admin not found"})
	}

	if err := config.DB.Delete(&admin).Error; err != nil {
		return c.Status(500).JSON(models.ErrorResponse{Error: "Could not delete admin"})
	}

	return c.JSON(models.SuccessResponse{Message: "Admin deleted successfully"})
}

// UpdateAdmin godoc
// @Summary Update admin by ID
// @Tags auth
// @Accept json
// @Produce json
// @Param id path int true "Admin ID"
// @Param admin body models.UpdateUserInput true "Admin data"
// @Success 200 {object} models.UserResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/admin/{id} [put]
// @Security Bearer
func UpdateAdmin(c *fiber.Ctx) error {
	id := c.Params("id")
	var admin models.User
	result := config.DB.First(&admin, id)
	if result.Error != nil {
		return c.Status(404).JSON(models.ErrorResponse{Error: "Admin not found"})
	}

	var input models.User
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(models.ErrorResponse{Error: "Invalid input"})
	}

	if err := config.DB.Model(&admin).Updates(input).Error; err != nil {
		return c.Status(500).JSON(models.ErrorResponse{Error: "Could not update admin"})
	}

	return c.JSON(admin)
}

// ChangePassword godoc
// @Summary Change admin password
// @Tags auth
// @Accept json
// @Produce json
// @Param id path int true "Admin ID"
// @Param password body models.ChangePasswordInput true "Old and New Password"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/admin/{id}/change-password [put]
// @Security Bearer
func ChangePassword(c *fiber.Ctx) error {
	id := c.Params("id")
	var admin models.User
	result := config.DB.First(&admin, id)
	if result.Error != nil {
		return c.Status(404).JSON(models.ErrorResponse{Error: "Admin not found"})
	}

	var input models.ChangePasswordInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(models.ErrorResponse{Error: "Invalid input"})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(input.OldPassword)); err != nil {
		return c.Status(401).JSON(models.ErrorResponse{Error: "Old password is incorrect"})
	}

	hashed, _ := bcrypt.GenerateFromPassword([]byte(input.NewPassword), 14)
	admin.Password = string(hashed)

	if err := config.DB.Save(&admin).Error; err != nil {
		return c.Status(500).JSON(models.ErrorResponse{Error: "Could not update password"})
	}

	return c.JSON(models.SuccessResponse{Message: "Password updated successfully"})
}

// CreateAdmin godoc
// @Summary Create a new admin
// @Tags auth
// @Accept json
// @Produce json
// @Param admin body models.User true "Admin data"
// @Success 201 {object} models.UserResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/admin [post]
// @Security Bearer
func CreateAdmin(c *fiber.Ctx) error {
	var input models.User
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(models.ErrorResponse{Error: "Invalid input"})
	}

	hashed, _ := bcrypt.GenerateFromPassword([]byte(input.Password), 14)
	input.Password = string(hashed)

	if err := config.DB.Create(&input).Error; err != nil {
		return c.Status(500).JSON(models.ErrorResponse{Error: "Could not create admin"})
	}

	return c.JSON(input)
}
