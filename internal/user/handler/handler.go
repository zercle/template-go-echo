package handler

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/zercle/template-go-echo/internal/config"
	"github.com/zercle/template-go-echo/internal/middleware"
	"github.com/zercle/template-go-echo/internal/user/domain"
	"github.com/zercle/template-go-echo/pkg"
)

// Handler handles user HTTP requests
type Handler struct {
	usecase domain.UserUsecase
}

// New creates a new user handler
func New(usecase domain.UserUsecase) *Handler {
	return &Handler{
		usecase: usecase,
	}
}

// RegisterRoutes registers user routes
func (h *Handler) RegisterRoutes(e *echo.Echo, jwtCfg *config.JWTConfig) {
	group := e.Group("/api/v1/users")

	// Public routes
	group.POST("/register", h.Register)
	group.POST("/login", h.Login)
	group.POST("/token/refresh", h.RefreshToken)

	// Protected routes
	group.GET("/:id", h.GetUser, middleware.JWTAuth(jwtCfg))
	group.GET("", h.ListUsers, middleware.JWTAuth(jwtCfg))
	group.PUT("/:id", h.UpdateProfile, middleware.JWTAuth(jwtCfg))
	group.POST("/:id/password", h.ChangePassword, middleware.JWTAuth(jwtCfg))
	group.DELETE("/:id", h.DeleteUser, middleware.JWTAuth(jwtCfg))
	group.POST("/logout", h.Logout, middleware.JWTAuth(jwtCfg))
	group.POST("/logout-all", h.LogoutAll, middleware.JWTAuth(jwtCfg))
}

// Register creates a new user account
// @Summary Register a new user
// @Description Create a new user account
// @Tags users
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Registration request"
// @Success 201 {object} pkg.JSendResponse{data=UserResponse}
// @Failure 400 {object} pkg.JSendResponse
// @Failure 409 {object} pkg.JSendResponse
// @Failure 500 {object} pkg.JSendResponse
// @Router /api/v1/users/register [post]
func (h *Handler) Register(c echo.Context) error {
	req := &RegisterRequest{}
	if err := c.Bind(req); err != nil {
		return pkg.Fail(c, http.StatusBadRequest, nil, "invalid request body")
	}

	user, err := h.usecase.RegisterUser(c.Request().Context(), req.Email, req.Name, req.Password)
	if err != nil {
		if domainErr, ok := err.(*pkg.DomainError); ok {
			return pkg.Error(c, http.StatusConflict, domainErr.Message, domainErr.Code)
		}
		return pkg.Error(c, http.StatusInternalServerError, err.Error(), pkg.ErrCodeInternalError)
	}

	return pkg.Success(c, http.StatusCreated, &UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	})
}

// Login authenticates a user
// @Summary Login a user
// @Description Authenticate a user and return access/refresh tokens
// @Tags users
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login request"
// @Success 200 {object} pkg.JSendResponse{data=LoginResponse}
// @Failure 400 {object} pkg.JSendResponse
// @Failure 401 {object} pkg.JSendResponse
// @Failure 500 {object} pkg.JSendResponse
// @Router /api/v1/users/login [post]
func (h *Handler) Login(c echo.Context) error {
	req := &LoginRequest{}
	if err := c.Bind(req); err != nil {
		return pkg.Fail(c, http.StatusBadRequest, nil, "invalid request body")
	}

	user, accessToken, refreshToken, err := h.usecase.LoginUser(
		c.Request().Context(),
		req.Email,
		req.Password,
		c.RealIP(),
		c.Request().UserAgent(),
	)
	if err != nil {
		if domainErr, ok := err.(*pkg.DomainError); ok {
			return pkg.Error(c, http.StatusUnauthorized, domainErr.Message, domainErr.Code)
		}
		return pkg.Error(c, http.StatusInternalServerError, err.Error(), pkg.ErrCodeInternalError)
	}

	return pkg.Success(c, http.StatusOK, &LoginResponse{
		User: &UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			Name:      user.Name,
			IsActive:  user.IsActive,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    3600, // 1 hour
	})
}

// GetUser retrieves a user by ID
// @Summary Get user by ID
// @Description Retrieve a user account by ID
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Success 200 {object} pkg.JSendResponse{data=UserResponse}
// @Failure 400 {object} pkg.JSendResponse
// @Failure 401 {object} pkg.JSendResponse
// @Failure 404 {object} pkg.JSendResponse
// @Router /api/v1/users/{id} [get]
func (h *Handler) GetUser(c echo.Context) error {
	userID := c.Param("id")
	if userID == "" {
		return pkg.Fail(c, http.StatusBadRequest, nil, "user id is required")
	}

	user, err := h.usecase.GetUser(c.Request().Context(), userID)
	if err != nil {
		if domainErr, ok := err.(*pkg.DomainError); ok {
			code := http.StatusNotFound
			if domainErr.Code == domain.ErrCodeUnauthorized {
				code = http.StatusUnauthorized
			}
			return pkg.Error(c, code, domainErr.Message, domainErr.Code)
		}
		return pkg.Error(c, http.StatusInternalServerError, err.Error(), pkg.ErrCodeInternalError)
	}

	return pkg.Success(c, http.StatusOK, &UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	})
}

// ListUsers retrieves a paginated list of users
// @Summary List users
// @Description Retrieve a paginated list of users
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Page limit (default: 10, max: 100)"
// @Param offset query int false "Page offset (default: 0)"
// @Success 200 {object} pkg.JSendResponse{data=UserListResponse}
// @Failure 401 {object} pkg.JSendResponse
// @Failure 500 {object} pkg.JSendResponse
// @Router /api/v1/users [get]
func (h *Handler) ListUsers(c echo.Context) error {
	limit := 10
	offset := 0

	if l := c.QueryParam("limit"); l != "" {
		if parsedLimit, err := strconv.Atoi(l); err == nil && parsedLimit > 0 {
			limit = parsedLimit
			if limit > 100 {
				limit = 100
			}
		}
	}

	if o := c.QueryParam("offset"); o != "" {
		if parsedOffset, err := strconv.Atoi(o); err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}

	users, total, err := h.usecase.ListUsers(c.Request().Context(), limit, offset)
	if err != nil {
		return pkg.Error(c, http.StatusInternalServerError, err.Error(), pkg.ErrCodeInternalError)
	}

	userResponses := make([]*UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = &UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			Name:      user.Name,
			IsActive:  user.IsActive,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		}
	}

	totalPages := (total + limit - 1) / limit
	if totalPages == 0 {
		totalPages = 1
	}

	return pkg.Success(c, http.StatusOK, &UserListResponse{
		Users:      userResponses,
		Total:      total,
		Limit:      limit,
		Offset:     offset,
		TotalPages: totalPages,
	})
}

// UpdateProfile updates user profile
// @Summary Update user profile
// @Description Update user name and email
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Param request body UpdateProfileRequest true "Update request"
// @Success 200 {object} pkg.JSendResponse{data=UserResponse}
// @Failure 400 {object} pkg.JSendResponse
// @Failure 401 {object} pkg.JSendResponse
// @Failure 404 {object} pkg.JSendResponse
// @Router /api/v1/users/{id} [put]
func (h *Handler) UpdateProfile(c echo.Context) error {
	userID := c.Param("id")
	if userID == "" {
		return pkg.Fail(c, http.StatusBadRequest, nil, "user id is required")
	}

	req := &UpdateProfileRequest{}
	if err := c.Bind(req); err != nil {
		return pkg.Fail(c, http.StatusBadRequest, nil, "invalid request body")
	}

	user, err := h.usecase.UpdateUserProfile(c.Request().Context(), userID, req.Name, req.Email)
	if err != nil {
		if domainErr, ok := err.(*pkg.DomainError); ok {
			code := http.StatusBadRequest
			switch domainErr.Code {
			case domain.ErrCodeUserNotFound:
				code = http.StatusNotFound
			case domain.ErrCodeUserExists:
				code = http.StatusConflict
			}
			return pkg.Error(c, code, domainErr.Message, domainErr.Code)
		}
		return pkg.Error(c, http.StatusInternalServerError, err.Error(), pkg.ErrCodeInternalError)
	}

	return pkg.Success(c, http.StatusOK, &UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	})
}

// ChangePassword changes user password
// @Summary Change password
// @Description Change user password with verification of old password
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Param request body ChangePasswordRequest true "Password change request"
// @Success 200 {object} pkg.JSendResponse
// @Failure 400 {object} pkg.JSendResponse
// @Failure 401 {object} pkg.JSendResponse
// @Failure 404 {object} pkg.JSendResponse
// @Router /api/v1/users/{id}/password [post]
func (h *Handler) ChangePassword(c echo.Context) error {
	userID := c.Param("id")
	if userID == "" {
		return pkg.Fail(c, http.StatusBadRequest, nil, "user id is required")
	}

	req := &ChangePasswordRequest{}
	if err := c.Bind(req); err != nil {
		return pkg.Fail(c, http.StatusBadRequest, nil, "invalid request body")
	}

	err := h.usecase.ChangePassword(c.Request().Context(), userID, req.OldPassword, req.NewPassword)
	if err != nil {
		if domainErr, ok := err.(*pkg.DomainError); ok {
			code := http.StatusBadRequest
			if domainErr.Code == domain.ErrCodeUserNotFound {
				code = http.StatusNotFound
			}
			return pkg.Error(c, code, domainErr.Message, domainErr.Code)
		}
		return pkg.Error(c, http.StatusInternalServerError, err.Error(), pkg.ErrCodeInternalError)
	}

	return pkg.SuccessWithMessage(c, http.StatusOK, nil, "password changed successfully")
}

// DeleteUser deletes a user account
// @Summary Delete user
// @Description Delete a user account permanently
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Success 204
// @Failure 400 {object} pkg.JSendResponse
// @Failure 401 {object} pkg.JSendResponse
// @Failure 404 {object} pkg.JSendResponse
// @Router /api/v1/users/{id} [delete]
func (h *Handler) DeleteUser(c echo.Context) error {
	userID := c.Param("id")
	if userID == "" {
		return pkg.Fail(c, http.StatusBadRequest, nil, "user id is required")
	}

	err := h.usecase.DeleteUser(c.Request().Context(), userID)
	if err != nil {
		if domainErr, ok := err.(*pkg.DomainError); ok {
			return pkg.Error(c, http.StatusNotFound, domainErr.Message, domainErr.Code)
		}
		return pkg.Error(c, http.StatusInternalServerError, err.Error(), pkg.ErrCodeInternalError)
	}

	return c.NoContent(http.StatusNoContent)
}

// Logout invalidates current session
// @Summary Logout
// @Description Invalidate current user session
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} pkg.JSendResponse
// @Failure 401 {object} pkg.JSendResponse
// @Router /api/v1/users/logout [post]
func (h *Handler) Logout(c echo.Context) error {
	sessionID := c.QueryParam("session_id")
	if sessionID == "" {
		slog.Warn("logout called without session_id")
		return pkg.SuccessWithMessage(c, http.StatusOK, nil, "logged out successfully")
	}

	err := h.usecase.LogoutUser(c.Request().Context(), sessionID)
	if err != nil {
		slog.Warn("logout failed", slog.String("error", err.Error()))
	}

	return pkg.SuccessWithMessage(c, http.StatusOK, nil, "logged out successfully")
}

// LogoutAll invalidates all sessions for user
// @Summary Logout all sessions
// @Description Invalidate all user sessions
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} pkg.JSendResponse
// @Failure 401 {object} pkg.JSendResponse
// @Router /api/v1/users/logout-all [post]
func (h *Handler) LogoutAll(c echo.Context) error {
	userID := middleware.GetUserID(c)
	if userID == "" {
		return pkg.Error(c, http.StatusUnauthorized, "user not authenticated", pkg.ErrCodeUnauthorized)
	}

	err := h.usecase.LogoutAllSessions(c.Request().Context(), userID)
	if err != nil {
		return pkg.Error(c, http.StatusInternalServerError, err.Error(), pkg.ErrCodeInternalError)
	}

	return pkg.SuccessWithMessage(c, http.StatusOK, nil, "all sessions logged out successfully")
}

// RefreshToken generates a new access token
// @Summary Refresh token
// @Description Generate a new access token using refresh token
// @Tags users
// @Accept json
// @Produce json
// @Param request body RefreshTokenRequest true "Refresh token request"
// @Success 200 {object} pkg.JSendResponse{data=TokenResponse}
// @Failure 400 {object} pkg.JSendResponse
// @Failure 401 {object} pkg.JSendResponse
// @Router /api/v1/users/token/refresh [post]
func (h *Handler) RefreshToken(c echo.Context) error {
	req := &RefreshTokenRequest{}
	if err := c.Bind(req); err != nil {
		return pkg.Fail(c, http.StatusBadRequest, nil, "invalid request body")
	}

	accessToken, err := h.usecase.RefreshToken(c.Request().Context(), req.RefreshToken)
	if err != nil {
		if domainErr, ok := err.(*pkg.DomainError); ok {
			return pkg.Error(c, http.StatusUnauthorized, domainErr.Message, domainErr.Code)
		}
		return pkg.Error(c, http.StatusInternalServerError, err.Error(), pkg.ErrCodeInternalError)
	}

	return pkg.Success(c, http.StatusOK, &TokenResponse{
		AccessToken: accessToken,
		ExpiresIn:   3600,
	})
}
