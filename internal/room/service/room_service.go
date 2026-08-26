package service

import (
	"errors"
	"strings"

	"github.com/Mosteben/hotel-booking-system/internal/room/model"
	"github.com/Mosteben/hotel-booking-system/internal/room/repository"
)

type RoomService interface {
	CreateRoom(hotelID uint, room *model.Room) error
	GetRoomsByHotelID(hotelID uint) ([]model.Room, error)
	GetRoomByID(id uint) (*model.Room, error)
	UpdateRoom(id uint, room *model.Room) error
	DeleteRoom(id uint) error
}

type roomService struct {
	repo repository.RoomRepository
}

func NewRoomService(repo repository.RoomRepository) RoomService {
	return &roomService{
		repo: repo,
	}
}

func (s *roomService) CreateRoom(
	hotelID uint,
	room *model.Room,
) error {

	if hotelID == 0 {
		return errors.New("invalid hotel id")
	}

	if err := validateRoom(room); err != nil {
		return err
	}

	room.HotelID = hotelID
	room.RoomNumber = strings.TrimSpace(room.RoomNumber)
	room.Type = strings.TrimSpace(room.Type)
	room.Description = strings.TrimSpace(room.Description)
	room.Status = strings.TrimSpace(room.Status)

	if room.Status == "" {
		room.Status = "available"
	}

	exists, err := s.repo.ExistsByRoomNumber(
		hotelID,
		room.RoomNumber,
	)

	if err != nil {
		return err
	}

	if exists {
		return errors.New("room number already exists in this hotel")
	}

	return s.repo.Create(room)
}

func (s *roomService) GetRoomsByHotelID(
	hotelID uint,
) ([]model.Room, error) {

	if hotelID == 0 {
		return nil, errors.New("invalid hotel id")
	}

	return s.repo.GetAllByHotelID(hotelID)
}

func (s *roomService) GetRoomByID(
	id uint,
) (*model.Room, error) {

	if id == 0 {
		return nil, errors.New("invalid room id")
	}

	return s.repo.GetByID(id)
}

func (s *roomService) UpdateRoom(
	id uint,
	room *model.Room,
) error {

	if id == 0 {
		return errors.New("invalid room id")
	}

	if err := validateRoom(room); err != nil {
		return err
	}

	existingRoom, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	existingRoom.RoomNumber = strings.TrimSpace(room.RoomNumber)
	existingRoom.Type = strings.TrimSpace(room.Type)
	existingRoom.Description = strings.TrimSpace(room.Description)
	existingRoom.PricePerNight = room.PricePerNight
	existingRoom.Capacity = room.Capacity
	existingRoom.Status = strings.TrimSpace(room.Status)

	if existingRoom.Status == "" {
		existingRoom.Status = "available"
	}

	return s.repo.Update(existingRoom)
}

func (s *roomService) DeleteRoom(id uint) error {

	if id == 0 {
		return errors.New("invalid room id")
	}

	_, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	return s.repo.Delete(id)
}

func validateRoom(room *model.Room) error {

	if room == nil {
		return errors.New("room data is required")
	}

	if strings.TrimSpace(room.RoomNumber) == "" {
		return errors.New("room number is required")
	}

	if strings.TrimSpace(room.Type) == "" {
		return errors.New("room type is required")
	}

	if room.PricePerNight <= 0 {
		return errors.New("price per night must be greater than zero")
	}

	if room.Capacity <= 0 {
		return errors.New("room capacity must be greater than zero")
	}

	if room.Status != "" {
		switch strings.ToLower(strings.TrimSpace(room.Status)) {
		case "available", "occupied", "maintenance":
		default:
			return errors.New(
				"invalid room status",
			)
		}
	}

	return nil
}
