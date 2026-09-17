package service

import (
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/Mosteben/hotel-booking-system/internal/hotel/model"
	"github.com/Mosteben/hotel-booking-system/internal/hotel/repository"
	"github.com/Mosteben/hotel-booking-system/pkg/storage"
	"gorm.io/gorm"
)

const maxHotelImageSize = 5 * 1024 * 1024 // 5MB

var allowedHotelImageTypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/webp": true,
}

type HotelService interface {
	CreateHotel(hotel *model.Hotel) error
	GetAllHotels() ([]model.Hotel, error)
	GetHotelByID(id uint) (*model.Hotel, error)
	GetHotelDetails(id uint) (*model.Hotel, error)
	UpdateHotel(id uint, hotel *model.Hotel) error
	DeleteHotel(id uint) error
	SearchHotels(filters *model.SearchRequest) ([]model.Hotel, error)

	UploadHotelImages(
		hotelID uint,
		files []*multipart.FileHeader,
	) ([]model.HotelImage, error)
	DeleteHotelImage(hotelID uint, imageID uint) error
	SetMainHotelImage(hotelID uint, imageID uint) error
}

type hotelService struct {
	repo    repository.HotelRepository
	storage storage.Storage
	db      *gorm.DB
}

func NewHotelService(
	repo repository.HotelRepository,
	storage storage.Storage,
	db *gorm.DB,
) HotelService {
	return &hotelService{
		repo:    repo,
		storage: storage,
		db:      db,
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
	hotels, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	for i := range hotels {
		hotels[i].ResolveImageURL()
	}

	return hotels, nil
}

func (s *hotelService) GetHotelByID(id uint) (*model.Hotel, error) {
	if id == 0 {
		return nil, errors.New("invalid hotel id")
	}

	hotel, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	hotel.ResolveImageURL()

	return hotel, nil
}

func (s *hotelService) GetHotelDetails(id uint) (*model.Hotel, error) {
	if id == 0 {
		return nil, errors.New("invalid hotel id")
	}

	hotel, err := s.repo.GetDetails(id)
	if err != nil {
		return nil, err
	}

	hotel.ResolveImageURL()

	// Rooms are preloaded with their own Images here (Preload("Rooms.Images")
	// in the repository), but each Room's computed ImageURL is never set by
	// GORM itself - same as the hotel's own ImageURL above, it only exists
	// after ResolveImageURL() runs, and the room-specific service methods
	// aren't in the call path for a hotel's nested room list.
	for i := range hotel.Rooms {
		hotel.Rooms[i].ResolveImageURL()
	}

	return hotel, nil
}

func (s *hotelService) UpdateHotel(
	id uint,
	hotel *model.Hotel,
) error {

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

func (s *hotelService) SearchHotels(
	filters *model.SearchRequest,
) ([]model.Hotel, error) {

	if filters == nil {
		return nil, errors.New("search filters are required")
	}

	if filters.Stars < 0 || filters.Stars > 5 {
		return nil, errors.New(
			"hotel stars must be between 1 and 5",
		)
	}

	if filters.MinPrice < 0 {
		return nil, errors.New(
			"minimum price cannot be negative",
		)
	}

	if filters.MaxPrice < 0 {
		return nil, errors.New(
			"maximum price cannot be negative",
		)
	}

	if filters.MinPrice > 0 &&
		filters.MaxPrice > 0 &&
		filters.MinPrice > filters.MaxPrice {

		return nil, errors.New(
			"minimum price cannot be greater than maximum price",
		)
	}

	if filters.MinCapacity < 0 {
		return nil, errors.New(
			"minimum capacity cannot be negative",
		)
	}

	filters.Name = strings.TrimSpace(filters.Name)
	filters.City = strings.TrimSpace(filters.City)
	filters.Country = strings.TrimSpace(filters.Country)
	filters.RoomType = strings.TrimSpace(filters.RoomType)

	hotels, err := s.repo.Search(filters)
	if err != nil {
		return nil, err
	}

	for i := range hotels {
		hotels[i].ResolveImageURL()
	}

	return hotels, nil
}

// =========================
// Hotel Images
// =========================

func (s *hotelService) UploadHotelImages(
	hotelID uint,
	files []*multipart.FileHeader,
) ([]model.HotelImage, error) {

	if hotelID == 0 {
		return nil, errors.New("invalid hotel id")
	}

	if len(files) == 0 {
		return nil, errors.New("at least one image file is required")
	}

	if _, err := s.repo.GetByID(hotelID); err != nil {
		return nil, err
	}

	existingCount, err := s.repo.CountImagesByHotelID(hotelID)
	if err != nil {
		return nil, err
	}

	uploaded := make([]model.HotelImage, 0, len(files))

	for _, fileHeader := range files {

		if fileHeader.Size > maxHotelImageSize {
			return uploaded, fmt.Errorf(
				"%s exceeds the 5MB size limit",
				fileHeader.Filename,
			)
		}

		file, err := fileHeader.Open()
		if err != nil {
			return uploaded, err
		}

		// Sniff the real content type from the file's own bytes - the
		// client-supplied filename and Content-Type header are never
		// trusted for this check, since either can be spoofed.
		sniffBuf := make([]byte, 512)
		n, err := file.Read(sniffBuf)

		if err != nil && err != io.EOF {
			file.Close()
			return uploaded, err
		}

		contentType := http.DetectContentType(sniffBuf[:n])

		if !allowedHotelImageTypes[contentType] {
			file.Close()
			return uploaded, fmt.Errorf(
				"%s is not a supported image type (only JPEG, PNG, and WebP are allowed)",
				fileHeader.Filename,
			)
		}

		if _, err := file.Seek(0, io.SeekStart); err != nil {
			file.Close()
			return uploaded, err
		}

		result, err := s.storage.Upload(file, fileHeader.Filename)
		file.Close()

		if err != nil {
			return uploaded, err
		}

		image := &model.HotelImage{
			HotelID:  hotelID,
			URL:      result.URL,
			PublicID: result.PublicID,
			// The very first image a hotel ever receives becomes its main
			// image automatically; later uploads default to gallery-only
			// until an admin explicitly changes the main image.
			IsMain: existingCount == 0 && len(uploaded) == 0,
		}

		if err := s.repo.CreateImage(image); err != nil {
			return uploaded, err
		}

		uploaded = append(uploaded, *image)
	}

	return uploaded, nil
}

func (s *hotelService) DeleteHotelImage(
	hotelID uint,
	imageID uint,
) error {

	if hotelID == 0 || imageID == 0 {
		return errors.New("invalid hotel or image id")
	}

	image, err := s.repo.GetImageByID(imageID)
	if err != nil {
		return err
	}

	if image.HotelID != hotelID {
		return errors.New("image does not belong to this hotel")
	}

	wasMain := image.IsMain

	// Deleting the image and promoting a replacement main image (if any)
	// must be atomic - otherwise a crash or concurrent request between the
	// two steps could leave the hotel with zero main images, or a
	// concurrent SetMainHotelImage could race the promotion below.
	if err := s.db.Transaction(func(tx *gorm.DB) error {
		txRepo := s.repo.WithTx(tx)

		if err := txRepo.DeleteImage(imageID); err != nil {
			return err
		}

		if !wasMain {
			return nil
		}

		remaining, err := txRepo.GetImagesByHotelID(hotelID)
		if err != nil {
			return err
		}

		if len(remaining) > 0 {
			return txRepo.SetMainImage(remaining[0].ID)
		}

		return nil
	}); err != nil {
		return err
	}

	// The DB row is already gone either way at this point, so a storage
	// provider hiccup here only risks (rarely) an orphaned remote file,
	// never an orphaned or dangling DB reference.
	if err := s.storage.Delete(image.PublicID); err != nil {
		log.Printf(
			"failed to delete stored hotel image %s: %v",
			image.PublicID,
			err,
		)
	}

	return nil
}

func (s *hotelService) SetMainHotelImage(
	hotelID uint,
	imageID uint,
) error {

	if hotelID == 0 || imageID == 0 {
		return errors.New("invalid hotel or image id")
	}

	image, err := s.repo.GetImageByID(imageID)
	if err != nil {
		return err
	}

	if image.HotelID != hotelID {
		return errors.New("image does not belong to this hotel")
	}

	if image.IsMain {
		return nil
	}

	// Unsetting the previous main image and setting the new one must be
	// atomic - otherwise a concurrent request (another SetMainHotelImage,
	// or a DeleteHotelImage promoting a replacement) could observe or
	// leave the hotel with zero or two main images.
	return s.db.Transaction(func(tx *gorm.DB) error {
		txRepo := s.repo.WithTx(tx)

		if err := txRepo.UnsetMainImage(hotelID); err != nil {
			return err
		}

		return txRepo.SetMainImage(imageID)
	})
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
		return errors.New(
			"hotel stars must be between 1 and 5",
		)
	}

	if hotel.Email != "" &&
		!strings.Contains(hotel.Email, "@") {

		return errors.New("invalid hotel email")
	}

	return nil
}

func IsNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}
