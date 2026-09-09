package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Mosteben/hotel-booking-system/internal/hotel/model"
	"github.com/Mosteben/hotel-booking-system/internal/hotel/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type HotelHandler struct {
	service service.HotelService
}

func NewHotelHandler(
	service service.HotelService,
) *HotelHandler {
	return &HotelHandler{
		service: service,
	}
}

// POST /hotels
func (h *HotelHandler) CreateHotel(c *gin.Context) {

	var hotel model.Hotel

	if err := c.ShouldBindJSON(&hotel); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "invalid hotel data",
			"error":   err.Error(),
		})
		return
	}

	if err := h.service.CreateHotel(&hotel); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "hotel created successfully",
		"data":    hotel,
	})
}

// GET /hotels
func (h *HotelHandler) GetAllHotels(c *gin.Context) {

	hotels, err := h.service.GetAllHotels()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "failed to get hotels",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "hotels retrieved successfully",
		"data":    hotels,
	})
}

// GET /hotels/search
func (h *HotelHandler) SearchHotels(c *gin.Context) {

	var filters model.SearchRequest

	if err := c.ShouldBindQuery(&filters); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "invalid search parameters",
			"error":   err.Error(),
		})
		return
	}

	hotels, err := h.service.SearchHotels(&filters)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "hotels search completed successfully",
		"data":    hotels,
	})
}

// GET /hotels/:id
func (h *HotelHandler) GetHotelByID(c *gin.Context) {

	idParam := c.Param("id")

	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "invalid hotel id",
		})
		return
	}

	hotel, err := h.service.GetHotelByID(uint(id))
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
			"message": "failed to get hotel",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "hotel retrieved successfully",
		"data":    hotel,
	})
}

// PUT /hotels/:id
func (h *HotelHandler) UpdateHotel(c *gin.Context) {

	idParam := c.Param("id")

	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "invalid hotel id",
		})
		return
	}

	var hotel model.Hotel

	if err := c.ShouldBindJSON(&hotel); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "invalid hotel data",
			"error":   err.Error(),
		})
		return
	}

	err = h.service.UpdateHotel(
		uint(id),
		&hotel,
	)

	if err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "hotel not found",
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
		"message": "hotel updated successfully",
	})
}

// DELETE /hotels/:id
func (h *HotelHandler) DeleteHotel(c *gin.Context) {

	idParam := c.Param("id")

	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "invalid hotel id",
		})
		return
	}

	err = h.service.DeleteHotel(uint(id))

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
			"message": "failed to delete hotel",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "hotel deleted successfully",
	})
}