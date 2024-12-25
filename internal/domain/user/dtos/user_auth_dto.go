package dtos

type LoginDto struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type TokenDto struct {
	Token string `json:"token"`
	Type  string `json:"type"`
}

type RegisterDto struct {
	FirstName   string  `json:"first_name"`
	LastName    *string `json:"last_name"`
	PhoneNumber *string `json:"phone_number"`
	Email       string  `json:"email"`
	Password    string  `json:"password"`
}
