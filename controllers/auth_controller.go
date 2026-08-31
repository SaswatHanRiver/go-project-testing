package controllers

import (
	"net/http"

	"go-project-testing/models"
	"go-project-testing/services"
	"go-project-testing/utils"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	service *services.AuthService
}

func NewAuthController() *AuthController {
	return &AuthController{service: services.NewAuthService()}
}

// Register godoc
// @Summary      Register a new user
// @Description  Creates a new user account and returns a JWT token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      models.RegisterRequest  true  "Register credentials"
// @Success      201      {object}  utils.ApiResponse
// @Failure      400      {object}  utils.ApiResponse
// @Router       /auth/register [post]
func (c *AuthController) Register(ctx *gin.Context) {
	var req models.RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	response, err := c.service.Register(&req)
	if err != nil {
		if err.Error() == "email already registered" {
			utils.Error(ctx, http.StatusConflict, err.Error())
			return
		}
		utils.Error(ctx, http.StatusInternalServerError, "Registration failed")
		return
	}
	utils.Success(ctx, http.StatusCreated, "User registered successfully", response)
}

// Login godoc
// @Summary      Login
// @Description  Authenticates user and returns a JWT token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      models.LoginRequest  true  "Login credentials"
// @Success      200      {object}  utils.ApiResponse
// @Failure      401      {object}  utils.ApiResponse
// @Router       /auth/login [post]
func (c *AuthController) Login(ctx *gin.Context) {
	var req models.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	response, err := c.service.Login(&req)
	if err != nil {
		utils.Error(ctx, http.StatusUnauthorized, "Invalid email or password")
		return
	}
	utils.Success(ctx, http.StatusOK, "Login successful", response)
}

// Me godoc
// @Summary      Get current user
// @Description  Returns the currently authenticated user info
// @Tags         auth
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  utils.ApiResponse
// @Failure      401  {object}  utils.ApiResponse
// @Router       /api/auth/me [get]
func (c *AuthController) Me(ctx *gin.Context) {
	userID, _ := ctx.Get("userID")
	email, _ := ctx.Get("email")
	utils.Success(ctx, http.StatusOK, "User info fetched", gin.H{
		"user_id": userID,
		"email":   email,
	})
}
