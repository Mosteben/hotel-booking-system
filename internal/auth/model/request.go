package model

type UpdateProfileRequest struct {
	FirstName      string `json:"first_name" validate:"omitempty,min=2,max=50"`
	LastName       string `json:"last_name" validate:"omitempty,min=2,max=50"`
	Phone          string `json:"phone" validate:"omitempty,min=10,max=20"`
	Gender         string `json:"gender" validate:"omitempty,oneof=male female"`
	Nationality    string `json:"nationality" validate:"omitempty,max=100"`
	NationalID     string `json:"national_id" validate:"omitempty,max=50"`
	PassportNumber string `json:"passport_number" validate:"omitempty,max=50"`
	Address        string `json:"address" validate:"omitempty,max=255"`
	City           string `json:"city" validate:"omitempty,max=100"`
	State          string `json:"state" validate:"omitempty,max=100"`
	Country        string `json:"country" validate:"omitempty,max=100"`
	PostalCode     string `json:"postal_code" validate:"omitempty,max=20"`
}
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" validate:"required"`
}