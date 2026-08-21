package model

type RegisterRequest struct {

	FirstName string `json:"first_name" validate:"required,min=2,max=100"`

	LastName string `json:"last_name" validate:"required,min=2,max=100"`

	Email string `json:"email" validate:"required,email"`

	Password string `json:"password" validate:"required,min=8"`

	ConfirmPassword string `json:"confirm_password" validate:"required"`

	Phone string `json:"phone" validate:"required"`

	DateOfBirth string `json:"date_of_birth" validate:"required"`

	Gender string `json:"gender" validate:"required,oneof=male female"`

	Nationality string `json:"nationality" validate:"required"`

	NationalID string `json:"national_id"`

	PassportNumber string `json:"passport_number"`

	Address string `json:"address" validate:"required"`

	City string `json:"city" validate:"required"`

	State string `json:"state"`

	Country string `json:"country" validate:"required"`

	PostalCode string `json:"postal_code"`
}