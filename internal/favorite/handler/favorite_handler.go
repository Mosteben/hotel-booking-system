package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Mosteben/hotel-booking-system/internal/favorite/model"
	"github.com/Mosteben/hotel-booking-system/internal/favorite/service"
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

// POST /hotels/:id/favorite
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

// GET /favorites
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
			"error":   err.Error(),
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

// DELETE /hotels/:id/favorite
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
