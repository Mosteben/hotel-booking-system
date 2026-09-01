package repository

import (
	"github.com/Mosteben/hotel-booking-system/internal/room/model"
	"gorm.io/gorm"
)

type RoomRepository interface {
	Create(room *model.Room) error
	GetAllByHotelID(hotelID uint) ([]model.Room, error)
	GetByID(id uint) (*model.Room, error)
	Update(room *model.Room) error
	Delete(id uint) error
	ExistsByRoomNumber(hotelID uint, roomNumber string) (bool, error)
}

type roomRepository struct {
	db *gorm.DB
}

func NewRoomRepository(db *gorm.DB) RoomRepository {
	return &roomRepository{
		db: db,
	}
}

func (r *roomRepository) Create(room *model.Room) error {
	return r.db.Create(room).Error
}

func (r *roomRepository) GetAllByHotelID(hotelID uint) ([]model.Room, error) {
	var rooms []model.Room

	err := r.db.
		Where("hotel_id = ?", hotelID).
		Order("id DESC").
		Find(&rooms).Error

	return rooms, err
}

func (r *roomRepository) GetByID(id uint) (*model.Room, error) {
	var room model.Room

	err := r.db.First(&room, id).Error
	if err != nil {
		return nil, err
	}

	return &room, nil
}

func (r *roomRepository) Update(room *model.Room) error {
	return r.db.Save(room).Error
}

func (r *roomRepository) Delete(id uint) error {
	return r.db.Delete(&model.Room{}, id).Error
}

func (r *roomRepository) ExistsByRoomNumber(
	hotelID uint,
	roomNumber string,
) (bool, error) {

	var count int64

	err := r.db.
		Model(&model.Room{}).
		Where("hotel_id = ? AND room_number = ?", hotelID, roomNumber).
		Count(&count).Error

	return count > 0, err
}
