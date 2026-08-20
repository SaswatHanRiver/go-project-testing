package controllers

import (
	"net/http"

	"go-project-testing/models"
	"go-project-testing/services"

	"github.com/gin-gonic/gin"
)

// AuthController - equivalent to @RestController AuthController in Spring Boot
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
// @Success      201      {object}  models.AuthResponse
// @Failure      400      {object}  map[string]string
// @Failure      409      {object}  map[string]string
// @Router       /auth/register [post]
func (c *AuthController) Register(ctx *gin.Context) {
	var req models.RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := c.service.Register(&req)
	if err != nil {
		if err.Error() == "email already registered" {
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, response)
}

// Login godoc
// @Summary      Login
// @Description  Authenticates user and returns a JWT token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      models.LoginRequest  true  "Login credentials"
// @Success      200      {object}  models.AuthResponse
// @Failure      400      {object}  map[string]string
// @Failure      401      {object}  map[string]string
// @Router       /auth/login [post]
func (c *AuthController) Login(ctx *gin.Context) {
	var req models.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := c.service.Login(&req)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// Me godoc
// @Summary      Get current user
// @Description  Returns the currently authenticated user's info from the JWT token
// @Tags         auth
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]string
// @Router       /auth/me [get]
func (c *AuthController) Me(ctx *gin.Context) {
	// Read values set by AuthMiddleware - like SecurityContextHolder.getContext() in Spring Boot
	userID, _ := ctx.Get("userID")
	email, _ := ctx.Get("email")

	ctx.JSON(http.StatusOK, gin.H{
		"user_id": userID,
		"email":   email,
	})
}
