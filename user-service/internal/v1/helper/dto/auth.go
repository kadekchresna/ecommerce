package dto

type LoginUserRequest struct {
	LoginCredential string `json:"email_or_phone_number"`
	Password        string `json:"password"`
}

type LoginUserResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
