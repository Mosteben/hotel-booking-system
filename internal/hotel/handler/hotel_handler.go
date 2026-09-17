package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Mosteben/hotel-booking-system/internal/hotel/model"
	"github.com/Mosteben/hotel-booking-system/internal/hotel/service"
	"github.com/Mosteben/hotel-booking-system/pkg/response"
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

// CreateHotel godoc
// @Summary Create a hotel
// @Description Create a new hotel.
// @Tags Hotels
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param hotel body model.Hotel true "Hotel data"
// @Success 201 {object} map[string]interface{} "Hotel created successfully"
// @Failure 400 {object} map[string]interface{} "Invalid hotel data or validation failed"
// @Failure 401 {object} map[string]interface{} "User not authenticated"
// @Failure 403 {object} map[string]interface{} "Access denied"
// @Router /hotels [post]
func (h *HotelHandler) CreateHotel(c *gin.Context) {
	var hotel model.Hotel

	if err := c.ShouldBindJSON(&hotel); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "invalid hotel data",
			"error":   response.InvalidRequestMessage,
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

// GetAllHotels godoc
// @Summary Get all hotels
// @Description Retrieve all hotels.
// @Tags Hotels
// @Produce json
// @Success 200 {object} map[string]interface{} "Hotels retrieved successfully"
// @Failure 500 {object} map[string]interface{} "Failed to get hotels"
// @Router /hotels [get]
func (h *HotelHandler) GetAllHotels(c *gin.Context) {
	hotels, err := h.service.GetAllHotels()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "failed to get hotels",
			"error":   response.SanitizedError("GetAllHotels", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "hotels retrieved successfully",
		"data":    hotels,
	})
}

// SearchHotels godoc
// @Summary Search hotels
// @Description Search hotels using hotel and room filters.
// @Tags Hotels
// @Produce json
// @Param name query string false "Hotel name"
// @Param city query string false "City"
// @Param country query string false "Country"
// @Param stars query int false "Hotel star rating"
// @Param min_price query number false "Minimum room price"
// @Param max_price query number false "Maximum room price"
// @Param room_type query string false "Room type"
// @Param min_capacity query int false "Minimum room capacity"
// @Success 200 {object} map[string]interface{} "Hotel search completed successfully"
// @Failure 400 {object} map[string]interface{} "Invalid search parameters"
// @Router /hotels/search [get]
func (h *HotelHandler) SearchHotels(c *gin.Context) {
	var filters model.SearchRequest

	if err := c.ShouldBindQuery(&filters); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "invalid search parameters",
			"error":   response.InvalidRequestMessage,
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

// GetHotelByID godoc
// @Summary Get hotel by ID
// @Description Retrieve a hotel using its ID.
// @Tags Hotels
// @Produce json
// @Param id path int true "Hotel ID"
// @Success 200 {object} map[string]interface{} "Hotel retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Invalid hotel ID"
// @Failure 404 {object} map[string]interface{} "Hotel not found"
// @Failure 500 {object} map[string]interface{} "Failed to get hotel"
// @Router /hotels/{id} [get]
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
			"error":   response.SanitizedError("GetHotelByID", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "hotel retrieved successfully",
		"data":    hotel,
	})
}

// GetHotelDetails godoc
// @Summary Get hotel details
// @Description Retrieve a hotel with all of its rooms.
// @Tags Hotels
// @Produce json
// @Param id path int true "Hotel ID"
// @Success 200 {object} map[string]interface{} "Hotel details retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Invalid hotel ID"
// @Failure 404 {object} map[string]interface{} "Hotel not found"
// @Failure 500 {object} map[string]interface{} "Failed to get hotel details"
// @Router /hotels/{id}/details [get]
func (h *HotelHandler) GetHotelDetails(c *gin.Context) {
	idParam := c.Param("id")

	id, err := strconv.ParseUint(idParam, 10, 64)

	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "invalid hotel id",
		})
		return
	}

	hotel, err := h.service.GetHotelDetails(uint(id))

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
			"message": "failed to get hotel details",
			"error":   response.SanitizedError("GetHotelDetails", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "hotel details retrieved successfully",
		"data":    hotel,
	})
}

// UpdateHotel godoc
// @Summary Update a hotel
// @Description Update an existing hotel.
// @Tags Hotels
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Hotel ID"
// @Param hotel body model.Hotel true "Updated hotel data"
// @Success 200 {object} map[string]interface{} "Hotel updated successfully"
// @Failure 400 {object} map[string]interface{} "Invalid hotel data or validation failed"
// @Failure 401 {object} map[string]interface{} "User not authenticated"
// @Failure 403 {object} map[string]interface{} "Access denied"
// @Failure 404 {object} map[string]interface{} "Hotel not found"
// @Router /hotels/{id} [put]
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
			"error":   response.InvalidRequestMessage,
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

// DeleteHotel godoc
// @Summary Delete a hotel
// @Description Delete an existing hotel.
// @Tags Hotels
// @Security BearerAuth
// @Produce json
// @Param id path int true "Hotel ID"
// @Success 200 {object} map[string]interface{} "Hotel deleted successfully"
// @Failure 400 {object} map[string]interface{} "Invalid hotel ID"
// @Failure 401 {object} map[string]interface{} "User not authenticated"
// @Failure 403 {object} map[string]interface{} "Access denied"
// @Failure 404 {object} map[string]interface{} "Hotel not found"
// @Failure 500 {object} map[string]interface{} "Failed to delete hotel"
// @Router /hotels/{id} [delete]
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
			"error":   response.SanitizedError("DeleteHotel", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "hotel deleted successfully",
	})
}

// UploadHotelImages godoc
// @Summary Upload hotel images
// @Description Upload one or more images for a hotel. Admin/Manager only.
// @Tags Hotels
// @Security BearerAuth
// @Accept multipart/form-data
// @Produce json
// @Param id path int true "Hotel ID"
// @Param images formData file true "Image files"
// @Success 201 {object} map[string]interface{} "Images uploaded successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 401 {object} map[string]interface{} "User not authenticated"
// @Failure 403 {object} map[string]interface{} "Access denied"
// @Failure 404 {object} map[string]interface{} "Hotel not found"
// @Router /hotels/{id}/images [post]
func (h *HotelHandler) UploadHotelImages(c *gin.Context) {
	idParam := c.Param("id")

	id, err := strconv.ParseUint(idParam, 10, 64)

	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "invalid hotel id",
		})
		return
	}

	form, err := c.MultipartForm()

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "invalid multipart form",
			"error":   response.InvalidRequestMessage,
		})
		return
	}

	files := form.File["images"]

	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "at least one image file is required in the \"images\" field",
		})
		return
	}

	images, err := h.service.UploadHotelImages(uint(id), files)

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
			"data":    images,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "images uploaded successfully",
		"data":    images,
	})
}

// DeleteHotelImage godoc
// @Summary Delete a hotel image
// @Description Delete an image belonging to a hotel. Admin/Manager only.
// @Tags Hotels
// @Security BearerAuth
// @Produce json
// @Param id path int true "Hotel ID"
// @Param imageId path int true "Hotel Image ID"
// @Success 200 {object} map[string]interface{} "Image deleted successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 401 {object} map[string]interface{} "User not authenticated"
// @Failure 403 {object} map[string]interface{} "Access denied"
// @Failure 404 {object} map[string]interface{} "Image not found"
// @Router /hotels/{id}/images/{imageId} [delete]
func (h *HotelHandler) DeleteHotelImage(c *gin.Context) {
	hotelID, err := strconv.ParseUint(c.Param("id"), 10, 64)

	if err != nil || hotelID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "invalid hotel id",
		})
		return
	}

	imageID, err := strconv.ParseUint(c.Param("imageId"), 10, 64)

	if err != nil || imageID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "invalid image id",
		})
		return
	}

	err = h.service.DeleteHotelImage(uint(hotelID), uint(imageID))

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "image not found",
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
		"message": "image deleted successfully",
	})
}

// SetMainHotelImage godoc
// @Summary Set a hotel's main image
// @Description Mark an existing hotel image as the main/cover image. Admin/Manager only.
// @Tags Hotels
// @Security BearerAuth
// @Produce json
// @Param id path int true "Hotel ID"
// @Param imageId path int true "Hotel Image ID"
// @Success 200 {object} map[string]interface{} "Main image updated successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 401 {object} map[string]interface{} "User not authenticated"
// @Failure 403 {object} map[string]interface{} "Access denied"
// @Failure 404 {object} map[string]interface{} "Image not found"
// @Router /hotels/{id}/images/{imageId}/main [patch]
func (h *HotelHandler) SetMainHotelImage(c *gin.Context) {
	hotelID, err := strconv.ParseUint(c.Param("id"), 10, 64)

	if err != nil || hotelID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "invalid hotel id",
		})
		return
	}

	imageID, err := strconv.ParseUint(c.Param("imageId"), 10, 64)

	if err != nil || imageID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "invalid image id",
		})
		return
	}

	err = h.service.SetMainHotelImage(uint(hotelID), uint(imageID))

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "image not found",
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
		"message": "main image updated successfully",
	})
}
