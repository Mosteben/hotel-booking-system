package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Mosteben/hotel-booking-system/internal/review/model"
	"github.com/Mosteben/hotel-booking-system/internal/review/service"
	"github.com/Mosteben/hotel-booking-system/pkg/response"
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

// =========================
// Create Review
// =========================

// CreateReview godoc
// @Summary Create a hotel review
// @Description Create a review and rating for a hotel. The authenticated user can submit one review per hotel.
// @Tags Reviews
// @Accept json
// @Produce json
// @Param id path int true "Hotel ID"
// @Param request body object{rating=int,comment=string} true "Review data"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Security BearerAuth
// @Router /hotels/{id}/reviews [post]
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
			"error":   response.InvalidRequestMessage,
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

// =========================
// Get Hotel Reviews
// =========================

// GetHotelReviews godoc
// @Summary Get hotel reviews
// @Description Get all reviews submitted for a specific hotel.
// @Tags Reviews
// @Produce json
// @Param id path int true "Hotel ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /hotels/{id}/reviews [get]
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
			"error":   response.SanitizedError("GetHotelReviews", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "hotel reviews retrieved successfully",
		"data":    reviews,
	})
}

// =========================
// Get All Reviews (Admin/Manager)
// =========================

// GetAllReviews godoc
// @Summary Get all reviews
// @Description Retrieve every review across all hotels. Admin and manager only.
// @Tags Reviews
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /reviews [get]
func (h *ReviewHandler) GetAllReviews(c *gin.Context) {
	reviews, err := h.service.GetAllReviews()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "failed to get reviews",
			"error":   response.SanitizedError("GetAllReviews", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "reviews retrieved successfully",
		"data":    reviews,
	})
}

// =========================
// Delete Review (Admin/Manager)
// =========================

// DeleteReview godoc
// @Summary Delete a review
// @Description Delete a review for moderation purposes. Admin and manager only.
// @Tags Reviews
// @Security BearerAuth
// @Produce json
// @Param id path int true "Review ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /reviews/{id} [delete]
func (h *ReviewHandler) DeleteReview(c *gin.Context) {
	idParam := c.Param("id")

	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "invalid review id",
		})
		return
	}

	if err := h.service.DeleteReview(uint(id)); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "review not found",
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
		"message": "review deleted successfully",
	})
}

// =========================
// Get Hotel Average Rating
// =========================

// GetHotelAverageRating godoc
// @Summary Get hotel average rating
// @Description Get the average rating of a specific hotel based on submitted reviews.
// @Tags Reviews
// @Produce json
// @Param id path int true "Hotel ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /hotels/{id}/rating [get]
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
			"error":   response.SanitizedError("GetHotelAverageRating", err),
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
