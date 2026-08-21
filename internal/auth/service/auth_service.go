package service

import (
	"errors"
	"time"

	"github.com/google/uuid"

	authModel "github.com/Mosteben/hotel-booking-system/internal/auth/model"
	profileModel "github.com/Mosteben/hotel-booking-system/internal/profile/model"
	profileRepository "github.com/Mosteben/hotel-booking-system/internal/profile/repository"
	userModel "github.com/Mosteben/hotel-booking-system/internal/user/model"
	userRepository "github.com/Mosteben/hotel-booking-system/internal/user/repository"

	"github.com/Mosteben/hotel-booking-system/pkg/hash"
	"github.com/Mosteben/hotel-booking-system/pkg/jwt"
	validatorPkg "github.com/Mosteben/hotel-booking-system/pkg/validator"

	"gorm.io/gorm"
)

type AuthService struct {
	db *gorm.DB

	userRepository *userRepository.UserRepository

	profileRepository *profileRepository.ProfileRepository
}

func NewAuthService(
	db *gorm.DB,
	userRepo *userRepository.UserRepository,
	profileRepo *profileRepository.ProfileRepository,
) *AuthService {

	return &AuthService{
		db: db,

		userRepository: userRepo,

		profileRepository: profileRepo,
	}
}

// =========================
// Register
// =========================

func (s *AuthService) Register(req authModel.RegisterRequest) error {

	// Validation
	if err := validatorPkg.ValidateStruct(req); err != nil {
		return err
	}

	// Confirm password
	if req.Password != req.ConfirmPassword {
		return errors.New("passwords do not match")
	}

	// Check email
	emailExists, err := s.userRepository.ExistsByEmail(req.Email)

	if err != nil {
		return err
	}

	if emailExists {
		return errors.New("email already exists")
	}

	// Check phone
	phoneExists, err := s.userRepository.ExistsByPhone(req.Phone)

	if err != nil {
		return err
	}

	if phoneExists {
		return errors.New("phone already exists")
	}

	// Hash password
	hashedPassword, err := hash.HashPassword(req.Password)

	if err != nil {
		return err
	}

	// Parse date of birth
	dateOfBirth, err := time.Parse(
		"2006-01-02",
		req.DateOfBirth,
	)

	if err != nil {
		return errors.New(
			"invalid date format (YYYY-MM-DD)",
		)
	}

	// =========================
	// Transaction
	// =========================

	tx := s.db.Begin()

	if tx.Error != nil {
		return tx.Error
	}

	// =========================
	// Create User
	// =========================

	user := userModel.User{

		FirstName: req.FirstName,

		LastName: req.LastName,

		Email: req.Email,

		Password: hashedPassword,

		Phone: req.Phone,

		// IMPORTANT:
		// Public registration always creates customers.
		// Admin/Manager must be created separately.
		Role: "customer",

		IsActive: true,

		IsVerified: false,
	}

	if err := s.userRepository.CreateTx(
		tx,
		&user,
	); err != nil {

		tx.Rollback()

		return err
	}

	// =========================
	// Create Profile
	// =========================

	profile := profileModel.Profile{

		UserID: user.ID,

		DateOfBirth: dateOfBirth,

		Gender: req.Gender,

		Nationality: req.Nationality,

		NationalID: req.NationalID,

		PassportNumber: req.PassportNumber,

		Address: req.Address,

		City: req.City,

		State: req.State,

		Country: req.Country,

		PostalCode: req.PostalCode,
	}

	if err := s.profileRepository.CreateTx(
		tx,
		&profile,
	); err != nil {

		tx.Rollback()

		return err
	}

	// =========================
	// Commit
	// =========================

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

// =========================
// Login
// =========================

func (s *AuthService) Login(
	req authModel.LoginRequest,
) (string, error) {

	// Validation
	if err := validatorPkg.ValidateStruct(req); err != nil {
		return "", err
	}

	// Find user
	user, err := s.userRepository.FindByEmail(req.Email)

	if err != nil || user == nil {
		return "", errors.New(
			"invalid email or password",
		)
	}

	// Check password
	if !hash.CheckPasswordHash(
		req.Password,
		user.Password,
	) {
		return "", errors.New(
			"invalid email or password",
		)
	}

	// Generate JWT
	// Includes:
	// user ID
	// user role
	token, err := jwt.GenerateToken(
		user.ID.String(),
		user.Role,
	)

	if err != nil {
		return "", err
	}

	return token, nil
}

// =========================
// Get Current User
// =========================

func (s *AuthService) GetCurrentUser(
	userID string,
) (*userModel.User, error) {

	user, err := s.userRepository.FindByID(userID)

	if err != nil || user == nil {
		return nil, errors.New(
			"user not found",
		)
	}

	return user, nil
}

// =========================
// Update Profile
// =========================

func (s *AuthService) UpdateProfile(
	userID string,
	req authModel.UpdateProfileRequest,
) error {

	// Validation
	if err := validatorPkg.ValidateStruct(req); err != nil {
		return err
	}

	// Find user
	user, err := s.userRepository.FindByID(userID)

	if err != nil || user == nil {
		return errors.New(
			"user not found",
		)
	}

	// =========================
	// Update User Data
	// =========================

	if req.FirstName != "" {
		user.FirstName = req.FirstName
	}

	if req.LastName != "" {
		user.LastName = req.LastName
	}

	if req.Phone != "" && req.Phone != user.Phone {

		exists, err := s.userRepository.ExistsByPhone(
			req.Phone,
		)

		if err != nil {
			return err
		}

		if exists {
			return errors.New(
				"phone already exists",
			)
		}

		user.Phone = req.Phone
	}

	// =========================
	// Update Profile Data
	// =========================

	profile := user.Profile

	if profile.UserID == uuid.Nil {

		profile = profileModel.Profile{
			UserID: user.ID,
		}
	}

	if req.Gender != "" {
		profile.Gender = req.Gender
	}

	if req.Nationality != "" {
		profile.Nationality = req.Nationality
	}

	if req.NationalID != "" {
		profile.NationalID = req.NationalID
	}

	if req.PassportNumber != "" {
		profile.PassportNumber = req.PassportNumber
	}

	if req.Address != "" {
		profile.Address = req.Address
	}

	if req.City != "" {
		profile.City = req.City
	}

	if req.State != "" {
		profile.State = req.State
	}

	if req.Country != "" {
		profile.Country = req.Country
	}

	if req.PostalCode != "" {
		profile.PostalCode = req.PostalCode
	}

	// =========================
	// Transaction
	// =========================

	tx := s.db.Begin()

	if tx.Error != nil {
		return tx.Error
	}

	// Update User
	if err := s.userRepository.UpdateTx(
		tx,
		user,
	); err != nil {

		tx.Rollback()

		return err
	}

	// Update Profile
	if profile.ID == 0 {

		if err := s.profileRepository.CreateTx(
			tx,
			&profile,
		); err != nil {

			tx.Rollback()

			return err
		}

	} else {

		if err := s.profileRepository.UpdateTx(
			tx,
			&profile,
		); err != nil {

			tx.Rollback()

			return err
		}
	}

	// Commit
	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

// =========================
// Change Password
// =========================

func (s *AuthService) ChangePassword(
	userID string,
	req authModel.ChangePasswordRequest,
) error {

	// Validation
	if err := validatorPkg.ValidateStruct(req); err != nil {
		return err
	}

	// Confirm password
	if req.NewPassword != req.ConfirmPassword {
		return errors.New(
			"passwords do not match",
		)
	}

	// Find user
	user, err := s.userRepository.FindByID(userID)

	if err != nil || user == nil {
		return errors.New(
			"user not found",
		)
	}

	// Check current password
	if !hash.CheckPasswordHash(
		req.CurrentPassword,
		user.Password,
	) {
		return errors.New(
			"current password is incorrect",
		)
	}

	// New password must be different
	if req.CurrentPassword == req.NewPassword {
		return errors.New(
			"new password must be different from current password",
		)
	}

	// Hash new password
	hashedPassword, err := hash.HashPassword(
		req.NewPassword,
	)

	if err != nil {
		return err
	}

	user.Password = hashedPassword

	// Update user
	if err := s.userRepository.Update(user); err != nil {
		return err
	}

	return nil
}
