package dto

type UserRegister struct {
	Name          string  `json:"name" binding:"required;min=3;max=63"`
	Email         string  `json:"email" binding:"required;email;min=6;max=63"`
	Phone         string  `json:"phone" binding:"required;e164"`
	Password      string  `json:"password" binding:"required;max=127"`
	FavoriteSport *string `json:"favorite_sport"`
}

type UserLogin struct {
	Email    string `json:"email" binding:"required;email;min=3;max=63"`
	Password string `json:"password" binding:"required;max=63"`
}

type UserUpdate struct {
	Name          *string `json:"name" binding:"omitempty;min=3;max=63"`
	Email         *string `json:"email" binding:"omitempty;email;min=6;max=63"`
	Phone         *string `json:"phone" binding:"omitempty;e164"`
	OldPassword   *string `json:"old_password" binding:"omitempty;max=127"`
	NewPassword   *string `json:"new_password" binding:"omitempty;max=127"`
	FavoriteSport *string `json:"favorite_sport"`
}
type UserResponse struct {
	Name          string  `json:"name"`
	Email         string  `json:"email" binding:"required"`
	FavoriteSport *string `json:"favorite_sport"`
}
