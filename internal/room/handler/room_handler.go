package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Mosteben/hotel-booking-system/internal/room/model"
	"github.com/Mosteben/hotel-booking-system/internal/room/service"
	"github.com/Mosteben/hotel-booking-system/pkg/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RoomHandler struct {
	service service.RoomService
}

func NewRoomHandler(roomService service.RoomService) *RoomHandler {
	return &RoomHandler{
		service: roomService,
	}
}

// CreateRoom godoc
// @Summary Create a room
// @Description Create a new room for a hotel.
// @Tags Rooms
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param hotel_id path int true "Hotel ID"
// @Param room body model.Room true "Room data"
// @Success 201 {object} map[string]interface{} "Room created successfully"
// @Failure 400 {object} map[string]interface{} "Invalid hotel ID, request body, or validation failed"
// @Failure 401 {object} map[string]interface{} "User not authenticated"
// @Failure 403 {object} map[string]interface{} "Access denied"
// @Failure 404 {object} map[string]interface{} "Hotel not found"
// @Router /hotels/{hotel_id}/rooms [post]
func (h *RoomHandler) CreateRoom(c *gin.Context) {
	hotelID, err := strconv.ParseUint(c.Param("hotel_id"), 10, 64)

	if err != nil || hotelID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "invalid hotel id",
		})
		return
	}

	var room model.Room

	if err := c.ShouldBindJSON(&room); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "invalid request body",
			"error":   response.InvalidRequestMessage,
		})
		return
	}

	if err := h.service.CreateRoom(uint(hotelID), &room); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "room created successfully",
		"data":    room,
	})
}

// GetRoomsByHotelID godoc
// @Summary Get rooms by hotel ID
// @Description Retrieve all rooms belonging to a specific hotel.
// @Tags Rooms
// @Produce json
// @Param hotel_id path int true "Hotel ID"
// @Success 200 {object} map[string]interface{} "Rooms retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Invalid hotel ID"
// @Failure 500 {object} map[string]interface{} "Failed to get rooms"
// @Router /hotels/{hotel_id}/rooms [get]
func (h *RoomHandler) GetRoomsByHotelID(c *gin.Context) {
	hotelID, err := strconv.ParseUint(c.Param("hotel_id"), 10, 64)

	if err != nil || hotelID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "invalid hotel id",
		})
		return
	}

	rooms, err := h.service.GetRoomsByHotelID(uint(hotelID))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "failed to get rooms",
			"error":   response.SanitizedError("GetRoomsByHotelID", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "rooms retrieved successfully",
		"data":    rooms,
	})
}

// GetRoomByID godoc
// @Summary Get room by ID
// @Description Retrieve a room using its ID.
// @Tags Rooms
// @Produce json
// @Param id path int true "Room ID"
// @Success 200 {object} map[string]interface{} "Room retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Invalid room ID"
// @Failure 404 {object} map[string]interface{} "Room not found"
// @Failure 500 {object} map[string]interface{} "Failed to get room"
// @Router /rooms/{id} [get]
func (h *RoomHandler) GetRoomByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)

	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "invalid room id",
		})
		return
	}

	room, err := h.service.GetRoomByID(uint(id))

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "room not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "failed to get room",
			"error":   response.SanitizedError("GetRoomByID", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "room retrieved successfully",
		"data":    room,
	})
}

// UpdateRoom godoc
// @Summary Update a room
// @Description Update an existing room.
// @Tags Rooms
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Room ID"
// @Param room body model.Room true "Updated room data"
// @Success 200 {object} map[string]interface{} "Room updated successfully"
// @Failure 400 {object} map[string]interface{} "Invalid room ID, request body, or validation failed"
// @Failure 401 {object} map[string]interface{} "User not authenticated"
// @Failure 403 {object} map[string]interface{} "Access denied"
// @Failure 404 {object} map[string]interface{} "Room not found"
// @Router /rooms/{id} [put]
func (h *RoomHandler) UpdateRoom(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)

	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "invalid room id",
		})
		return
	}

	var room model.Room

	if err := c.ShouldBindJSON(&room); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "invalid request body",
			"error":   response.InvalidRequestMessage,
		})
		return
	}

	if err := h.service.UpdateRoom(uint(id), &room); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "room not found",
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
		"message": "room updated successfully",
	})
}

// DeleteRoom godoc
// @Summary Delete a room
// @Description Delete an existing room.
// @Tags Rooms
// @Security BearerAuth
// @Produce json
// @Param id path int true "Room ID"
// @Success 200 {object} map[string]interface{} "Room deleted successfully"
// @Failure 400 {object} map[string]interface{} "Invalid room ID"
// @Failure 401 {object} map[string]interface{} "User not authenticated"
// @Failure 403 {object} map[string]interface{} "Access denied"
// @Failure 404 {object} map[string]interface{} "Room not found"
// @Failure 500 {object} map[string]interface{} "Failed to delete room"
// @Router /rooms/{id} [delete]
func (h *RoomHandler) DeleteRoom(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)

	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "invalid room id",
		})
		return
	}

	if err := h.service.DeleteRoom(uint(id)); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "room not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "failed to delete room",
			"error":   response.SanitizedError("DeleteRoom", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "room deleted successfully",
	})
}

// UploadRoomImages godoc
// @Summary Upload room images
// @Description Upload one or more images for a room. Admin/Manager only.
// @Tags Rooms
// @Security BearerAuth
// @Accept multipart/form-data
// @Produce json
// @Param id path int true "Room ID"
// @Param images formData file true "Image files"
// @Success 201 {object} map[string]interface{} "Images uploaded successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 401 {object} map[string]interface{} "User not authenticated"
// @Failure 403 {object} map[string]interface{} "Access denied"
// @Failure 404 {object} map[string]interface{} "Room not found"
// @Router /rooms/{id}/images [post]
func (h *RoomHandler) UploadRoomImages(c *gin.Context) {
	idParam := c.Param("id")

	id, err := strconv.ParseUint(idParam, 10, 64)

	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "invalid room id",
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

	images, err := h.service.UploadRoomImages(uint(id), files)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "room not found",
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

// DeleteRoomImage godoc
// @Summary Delete a room image
// @Description Delete an image belonging to a room. Admin/Manager only.
// @Tags Rooms
// @Security BearerAuth
// @Produce json
// @Param id path int true "Room ID"
// @Param imageId path int true "Room Image ID"
// @Success 200 {object} map[string]interface{} "Image deleted successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 401 {object} map[string]interface{} "User not authenticated"
// @Failure 403 {object} map[string]interface{} "Access denied"
// @Failure 404 {object} map[string]interface{} "Image not found"
// @Router /rooms/{id}/images/{imageId} [delete]
func (h *RoomHandler) DeleteRoomImage(c *gin.Context) {
	roomID, err := strconv.ParseUint(c.Param("id"), 10, 64)

	if err != nil || roomID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "invalid room id",
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

	err = h.service.DeleteRoomImage(uint(roomID), uint(imageID))

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

// SetMainRoomImage godoc
// @Summary Set a room's main image
// @Description Mark an existing room image as the main/cover image. Admin/Manager only.
// @Tags Rooms
// @Security BearerAuth
// @Produce json
// @Param id path int true "Room ID"
// @Param imageId path int true "Room Image ID"
// @Success 200 {object} map[string]interface{} "Main image updated successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 401 {object} map[string]interface{} "User not authenticated"
// @Failure 403 {object} map[string]interface{} "Access denied"
// @Failure 404 {object} map[string]interface{} "Image not found"
// @Router /rooms/{id}/images/{imageId}/main [patch]
func (h *RoomHandler) SetMainRoomImage(c *gin.Context) {
	roomID, err := strconv.ParseUint(c.Param("id"), 10, 64)

	if err != nil || roomID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "invalid room id",
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

	err = h.service.SetMainRoomImage(uint(roomID), uint(imageID))

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
