package service

import (
	"errors"
	"strings"

	"github.com/Mosteben/hotel-booking-system/internal/hotel/model"
	"github.com/Mosteben/hotel-booking-system/internal/hotel/repository"
	"gorm.io/gorm"
)

type HotelService interface {
	CreateHotel(hotel *model.Hotel) error
	GetAllHotels() ([]model.Hotel, error)
	GetHotelByID(id uint) (*model.Hotel, error)
	UpdateHotel(id uint, hotel *model.Hotel) error
	DeleteHotel(id uint) error
}

type hotelService struct {
	repo repository.HotelRepository
}

func NewHotelService(repo repository.HotelRepository) HotelService {
	return &hotelService{
		repo: repo,
	}
}

func (s *hotelService) CreateHotel(hotel *model.Hotel) error {
	if err := validateHotel(hotel); err != nil {
		return err
	}

	hotel.Name = strings.TrimSpace(hotel.Name)
	hotel.Address = strings.TrimSpace(hotel.Address)
	hotel.City = strings.TrimSpace(hotel.City)
	hotel.Country = strings.TrimSpace(hotel.Country)

	if hotel.Stars == 0 {
		hotel.Stars = 1
	}

	return s.repo.Create(hotel)
}

func (s *hotelService) GetAllHotels() ([]model.Hotel, error) {
	return s.repo.GetAll()
}

func (s *hotelService) GetHotelByID(id uint) (*model.Hotel, error) {
	if id == 0 {
		return nil, errors.New("invalid hotel id")
	}

	return s.repo.GetByID(id)
}

func (s *hotelService) UpdateHotel(id uint, hotel *model.Hotel) error {
	if id == 0 {
		return errors.New("invalid hotel id")
	}

	if err := validateHotel(hotel); err != nil {
		return err
	}

	existingHotel, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	existingHotel.Name = strings.TrimSpace(hotel.Name)
	existingHotel.Description = hotel.Description
	existingHotel.Address = strings.TrimSpace(hotel.Address)
	existingHotel.City = strings.TrimSpace(hotel.City)
	existingHotel.Country = strings.TrimSpace(hotel.Country)
	existingHotel.Phone = strings.TrimSpace(hotel.Phone)
	existingHotel.Email = strings.TrimSpace(hotel.Email)
	existingHotel.Stars = hotel.Stars

	return s.repo.Update(existingHotel)
}

func (s *hotelService) DeleteHotel(id uint) error {
	if id == 0 {
		return errors.New("invalid hotel id")
	}

	_, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	return s.repo.Delete(id)
}

func validateHotel(hotel *model.Hotel) error {
	if hotel == nil {
		return errors.New("hotel data is required")
	}

	if strings.TrimSpace(hotel.Name) == "" {
		return errors.New("hotel name is required")
	}

	if strings.TrimSpace(hotel.Address) == "" {
		return errors.New("hotel address is required")
	}

	if strings.TrimSpace(hotel.City) == "" {
		return errors.New("hotel city is required")
	}

	if strings.TrimSpace(hotel.Country) == "" {
		return errors.New("hotel country is required")
	}

	if hotel.Stars < 1 || hotel.Stars > 5 {
		return errors.New("hotel stars must be between 1 and 5")
	}

	if hotel.Email != "" && !strings.Contains(hotel.Email, "@") {
		return errors.New("invalid hotel email")
	}

	return nil
}

func IsNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}