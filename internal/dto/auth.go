package dto

type RegistrationInput struct {
	FullName     string `validate:"required,min=3,fullname"`
	Email        string `validate:"required,email"`
	Password     string `validate:"required,min=6,max=20"`
	CurrencyCode string `validate:"required,len=3"`
}

type AuthorizationInput struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required,min=6,max=20"`
}

type UpdatePasswordInput struct {
	Email       string `validate:"required,email" `
	OldPassword string `validate:"required,min=6,max=20"`
	NewPassword string `validate:"required,min=6,max=20"`
}
