package service

import (
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/Mosteben/hotel-booking-system/internal/room/model"
	"github.com/Mosteben/hotel-booking-system/internal/room/repository"
	"github.com/Mosteben/hotel-booking-system/pkg/storage"
	"gorm.io/gorm"
)

const maxRoomImageSize = 5 * 1024 * 1024 // 5MB

var allowedRoomImageTypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/webp": true,
}

type RoomService interface {
	CreateRoom(hotelID uint, room *model.Room) error
	GetRoomsByHotelID(hotelID uint) ([]model.Room, error)
	GetRoomByID(id uint) (*model.Room, error)
	UpdateRoom(id uint, room *model.Room) error
	DeleteRoom(id uint) error

	UploadRoomImages(
		roomID uint,
		files []*multipart.FileHeader,
	) ([]model.RoomImage, error)
	DeleteRoomImage(roomID uint, imageID uint) error
	SetMainRoomImage(roomID uint, imageID uint) error
}

type roomService struct {
	repo    repository.RoomRepository
	storage storage.Storage
	db      *gorm.DB
}

func NewRoomService(
	repo repository.RoomRepository,
	storage storage.Storage,
	db *gorm.DB,
) RoomService {
	return &roomService{
		repo:    repo,
		storage: storage,
		db:      db,
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

	rooms, err := s.repo.GetAllByHotelID(hotelID)
	if err != nil {
		return nil, err
	}

	for i := range rooms {
		rooms[i].ResolveImageURL()
	}

	return rooms, nil
}

func (s *roomService) GetRoomByID(
	id uint,
) (*model.Room, error) {

	if id == 0 {
		return nil, errors.New("invalid room id")
	}

	room, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	room.ResolveImageURL()

	return room, nil
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

	room, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	// Capture the room's images before deleting anything, so they can
	// still be cleaned up from storage afterward. Deleting the image rows
	// and the room row happens atomically - a partial delete would either
	// orphan image rows pointing at a room that no longer exists, or
	// leave the room's own images dangling.
	images := room.Images

	if err := s.db.Transaction(func(tx *gorm.DB) error {
		txRepo := s.repo.WithTx(tx)

		for _, img := range images {
			if err := txRepo.DeleteImage(img.ID); err != nil {
				return err
			}
		}

		return txRepo.Delete(id)
	}); err != nil {
		return err
	}

	// The DB rows are already gone either way at this point, so a storage
	// provider hiccup here only risks (rarely) an orphaned remote file,
	// never an orphaned or dangling DB reference.
	for _, img := range images {
		if err := s.storage.Delete(img.PublicID); err != nil {
			log.Printf(
				"failed to delete stored room image %s: %v",
				img.PublicID,
				err,
			)
		}
	}

	return nil
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

// =========================
// Room Images
// =========================

func (s *roomService) UploadRoomImages(
	roomID uint,
	files []*multipart.FileHeader,
) ([]model.RoomImage, error) {

	if roomID == 0 {
		return nil, errors.New("invalid room id")
	}

	if len(files) == 0 {
		return nil, errors.New("at least one image file is required")
	}

	if _, err := s.repo.GetByID(roomID); err != nil {
		return nil, err
	}

	existingCount, err := s.repo.CountImagesByRoomID(roomID)
	if err != nil {
		return nil, err
	}

	uploaded := make([]model.RoomImage, 0, len(files))

	for _, fileHeader := range files {

		if fileHeader.Size > maxRoomImageSize {
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

		if !allowedRoomImageTypes[contentType] {
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

		image := &model.RoomImage{
			RoomID:   roomID,
			URL:      result.URL,
			PublicID: result.PublicID,
			// The very first image a room ever receives becomes its main
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

func (s *roomService) DeleteRoomImage(
	roomID uint,
	imageID uint,
) error {

	if roomID == 0 || imageID == 0 {
		return errors.New("invalid room or image id")
	}

	image, err := s.repo.GetImageByID(imageID)
	if err != nil {
		return err
	}

	if image.RoomID != roomID {
		return errors.New("image does not belong to this room")
	}

	wasMain := image.IsMain

	// Deleting the image and promoting a replacement main image (if any)
	// must be atomic, for the same reason as the hotel-image equivalent:
	// otherwise a crash or a concurrent SetMainRoomImage between the two
	// steps could leave the room with zero or two main images.
	if err := s.db.Transaction(func(tx *gorm.DB) error {
		txRepo := s.repo.WithTx(tx)

		if err := txRepo.DeleteImage(imageID); err != nil {
			return err
		}

		if !wasMain {
			return nil
		}

		remaining, err := txRepo.GetImagesByRoomID(roomID)
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
			"failed to delete stored room image %s: %v",
			image.PublicID,
			err,
		)
	}

	return nil
}

func (s *roomService) SetMainRoomImage(
	roomID uint,
	imageID uint,
) error {

	if roomID == 0 || imageID == 0 {
		return errors.New("invalid room or image id")
	}

	image, err := s.repo.GetImageByID(imageID)
	if err != nil {
		return err
	}

	if image.RoomID != roomID {
		return errors.New("image does not belong to this room")
	}

	if image.IsMain {
		return nil
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		txRepo := s.repo.WithTx(tx)

		if err := txRepo.UnsetMainImage(roomID); err != nil {
			return err
		}

		return txRepo.SetMainImage(imageID)
	})
}
