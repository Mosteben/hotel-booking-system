package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Mosteben/hotel-booking-system/internal/room/model"
	"github.com/Mosteben/hotel-booking-system/internal/room/service"

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

// POST /hotels/:hotel_id/rooms
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
			"error":   err.Error(),
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

// GET /hotels/:hotel_id/rooms
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
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "rooms retrieved successfully",
		"data":    rooms,
	})
}

// GET /rooms/:id
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
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "room retrieved successfully",
		"data":    room,
	})
}

// PUT /rooms/:id
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
			"error":   err.Error(),
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

// DELETE /rooms/:id
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
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "room deleted successfully",
	})
}
