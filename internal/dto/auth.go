package dto

type RegistrationInput struct {
	FullName     string `validate:"required,min=3,fullname"`
	Email        string `validate:"required,email"`
	Password     string `validate:"required,min=6"`
	CurrencyCode string `validate:"required,len=3"`
}

type AuthorizationInput struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required,min=6"`
}
