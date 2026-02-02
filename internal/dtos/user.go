package dtos

type RegisterRequest struct {
	User struct {
		Username string `json:"username" binding:"required,min=3,max=50"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
	} `json:"user" binding:"required"`
}

type LoginRequest struct {
	User struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	} `json:"user" binding:"required"`
}

type UpdateUserRequest struct {
	User struct {
		Username *string `json:"username,omitempty"`
		Email    *string `json:"email,omitempty"`
		Password *string `json:"password,omitempty"`
		Bio      *string `json:"bio,omitempty"`
		Image    *string `json:"image,omitempty"`
	} `json:"user" binding:"required"`
}

type UserResponse struct {
	User UserData `json:"user"`
}

type UserData struct {
	Email    string `json:"email"`
	Token    string `json:"token"`
	Username string `json:"username"`
	Bio      string `json:"bio"`
	Image    string `json:"image"`
}
