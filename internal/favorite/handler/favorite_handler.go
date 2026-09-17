package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Mosteben/hotel-booking-system/internal/favorite/model"
	"github.com/Mosteben/hotel-booking-system/internal/favorite/service"
	"github.com/Mosteben/hotel-booking-system/pkg/response"
	"github.com/gin-gonic/gin"
)

type FavoriteHandler struct {
	service service.FavoriteService
}

func NewFavoriteHandler(
	service service.FavoriteService,
) *FavoriteHandler {
	return &FavoriteHandler{
		service: service,
	}
}

// =========================
// Add Favorite
// =========================

// AddFavorite godoc
// @Summary Add hotel to favorites
// @Description Add a hotel to the authenticated user's favorites.
// @Tags Favorites
// @Produce json
// @Param id path int true "Hotel ID"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Security BearerAuth
// @Router /hotels/{id}/favorite [post]
func (h *FavoriteHandler) AddFavorite(c *gin.Context) {

	hotelID, err := strconv.ParseUint(
		c.Param("id"),
		10,
		64,
	)

	if err != nil || hotelID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "invalid hotel id",
		})
		return
	}

	userIDValue, exists := c.Get("userID")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "user not authenticated",
		})
		return
	}

	userID, ok := userIDValue.(string)

	if !ok || userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "invalid user id",
		})
		return
	}

	favorite := &model.Favorite{
		UserID:  userID,
		HotelID: uint(hotelID),
	}

	if err := h.service.AddFavorite(favorite); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "hotel added to favorites successfully",
		"data":    favorite,
	})
}

// =========================
// Get My Favorites
// =========================

// GetMyFavorites godoc
// @Summary Get my favorite hotels
// @Description Get all hotels saved as favorites by the authenticated user.
// @Tags Favorites
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /favorites [get]
func (h *FavoriteHandler) GetMyFavorites(c *gin.Context) {

	userIDValue, exists := c.Get("userID")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "user not authenticated",
		})
		return
	}

	userID, ok := userIDValue.(string)

	if !ok || userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "invalid user id",
		})
		return
	}

	favorites, err := h.service.GetMyFavorites(userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "failed to get favorites",
			"error":   response.SanitizedError("GetMyFavorites", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "favorites retrieved successfully",
		"data":    favorites,
	})
}

// =========================
// Remove Favorite
// =========================

// RemoveFavorite godoc
// @Summary Remove hotel from favorites
// @Description Remove a hotel from the authenticated user's favorites.
// @Tags Favorites
// @Produce json
// @Param id path int true "Hotel ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Security BearerAuth
// @Router /hotels/{id}/favorite [delete]
func (h *FavoriteHandler) RemoveFavorite(c *gin.Context) {

	hotelID, err := strconv.ParseUint(
		c.Param("id"),
		10,
		64,
	)

	if err != nil || hotelID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "invalid hotel id",
		})
		return
	}

	userIDValue, exists := c.Get("userID")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "user not authenticated",
		})
		return
	}

	userID, ok := userIDValue.(string)

	if !ok || userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "invalid user id",
		})
		return
	}

	if err := h.service.RemoveFavorite(
		userID,
		uint(hotelID),
	); err != nil {

		if errors.Is(err, service.ErrFavoriteNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "hotel removed from favorites successfully",
	})
}
