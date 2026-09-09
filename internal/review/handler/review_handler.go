package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Mosteben/hotel-booking-system/internal/review/model"
	"github.com/Mosteben/hotel-booking-system/internal/review/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ReviewHandler struct {
	service service.ReviewService
}

func NewReviewHandler(
	service service.ReviewService,
) *ReviewHandler {
	return &ReviewHandler{
		service: service,
	}
}

// POST /hotels/:id/reviews
func (h *ReviewHandler) CreateReview(c *gin.Context) {
	idParam := c.Param("id")

	hotelID, err := strconv.ParseUint(idParam, 10, 64)
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

	var request struct {
		Rating  int    `json:"rating"`
		Comment string `json:"comment"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "invalid review data",
			"error":   err.Error(),
		})
		return
	}

	review := &model.Review{
		HotelID: uint(hotelID),
		UserID:  userID,
		Rating:  request.Rating,
		Comment: request.Comment,
	}

	if err := h.service.CreateReview(review); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "review created successfully",
		"data":    review,
	})
}

// GET /hotels/:id/reviews
func (h *ReviewHandler) GetHotelReviews(c *gin.Context) {
	idParam := c.Param("id")

	hotelID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil || hotelID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "invalid hotel id",
		})
		return
	}

	reviews, err := h.service.GetHotelReviews(uint(hotelID))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "hotel not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "failed to get hotel reviews",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "hotel reviews retrieved successfully",
		"data":    reviews,
	})
}

// GET /hotels/:id/rating
func (h *ReviewHandler) GetHotelAverageRating(c *gin.Context) {
	idParam := c.Param("id")

	hotelID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil || hotelID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "invalid hotel id",
		})
		return
	}

	average, err := h.service.GetHotelAverageRating(uint(hotelID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "failed to get hotel rating",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "hotel rating retrieved successfully",
		"data": gin.H{
			"hotel_id":       hotelID,
			"average_rating": average,
		},
	})
}