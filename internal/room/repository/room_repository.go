package repository

import (
	"github.com/Mosteben/hotel-booking-system/internal/room/model"
	"gorm.io/gorm"
)

type RoomRepository interface {
	Create(room *model.Room) error
	GetAllByHotelID(hotelID uint) ([]model.Room, error)
	GetByID(id uint) (*model.Room, error)
	GetByIDs(ids []uint) ([]model.Room, error)
	Update(room *model.Room) error
	Delete(id uint) error
	ExistsByRoomNumber(hotelID uint, roomNumber string) (bool, error)

	CreateImage(image *model.RoomImage) error
	GetImageByID(id uint) (*model.RoomImage, error)
	GetImagesByRoomID(roomID uint) ([]model.RoomImage, error)
	CountImagesByRoomID(roomID uint) (int64, error)
	DeleteImage(id uint) error
	UnsetMainImage(roomID uint) error
	SetMainImage(id uint) error

	// Create a repository that uses the provided transaction.
	WithTx(tx *gorm.DB) RoomRepository
}

type roomRepository struct {
	db *gorm.DB
}

func NewRoomRepository(db *gorm.DB) RoomRepository {
	return &roomRepository{
		db: db,
	}
}

func (r *roomRepository) WithTx(
	tx *gorm.DB,
) RoomRepository {
	return &roomRepository{
		db: tx,
	}
}

func (r *roomRepository) Create(room *model.Room) error {
	return r.db.Create(room).Error
}

func (r *roomRepository) GetAllByHotelID(hotelID uint) ([]model.Room, error) {
	var rooms []model.Room

	err := r.db.
		Preload("Images").
		Where("hotel_id = ?", hotelID).
		Order("id DESC").
		Find(&rooms).Error

	return rooms, err
}

func (r *roomRepository) GetByID(id uint) (*model.Room, error) {
	var room model.Room

	err := r.db.
		Preload("Images").
		First(&room, id).Error

	if err != nil {
		return nil, err
	}

	return &room, nil
}

func (r *roomRepository) GetByIDs(ids []uint) ([]model.Room, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	var rooms []model.Room

	err := r.db.
		Where("id IN ?", ids).
		Find(&rooms).Error

	return rooms, err
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

// =========================
// Room Images
// =========================

func (r *roomRepository) CreateImage(image *model.RoomImage) error {
	return r.db.Create(image).Error
}

func (r *roomRepository) GetImageByID(id uint) (*model.RoomImage, error) {
	var image model.RoomImage

	err := r.db.First(&image, id).Error
	if err != nil {
		return nil, err
	}

	return &image, nil
}

func (r *roomRepository) GetImagesByRoomID(
	roomID uint,
) ([]model.RoomImage, error) {

	var images []model.RoomImage

	err := r.db.
		Where("room_id = ?", roomID).
		Order("id ASC").
		Find(&images).Error

	return images, err
}

func (r *roomRepository) CountImagesByRoomID(
	roomID uint,
) (int64, error) {

	var count int64

	err := r.db.
		Model(&model.RoomImage{}).
		Where("room_id = ?", roomID).
		Count(&count).Error

	return count, err
}

func (r *roomRepository) DeleteImage(id uint) error {
	return r.db.Delete(&model.RoomImage{}, id).Error
}

func (r *roomRepository) UnsetMainImage(roomID uint) error {
	return r.db.
		Model(&model.RoomImage{}).
		Where("room_id = ? AND is_main = ?", roomID, true).
		Update("is_main", false).Error
}

func (r *roomRepository) SetMainImage(id uint) error {
	return r.db.
		Model(&model.RoomImage{}).
		Where("id = ?", id).
		Update("is_main", true).Error
}
