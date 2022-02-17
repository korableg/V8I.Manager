package user

//go:generate easyjson

type (
	//easyjson:json
	AddUserRequest struct {
		Name     string `json:"name" validate:"required"`
		Password string `json:"password" validate:"required"`
		Role     string `json:"role" validate:"required, oneof=admin reader"`
		Token    string `json:"token"`
	}

	//easyjson:json
	UpdateUserRequest struct {
		ID    int64  `json:"id" validate:"required"`
		Name  string `json:"name" validate:"required"`
		Role  string `json:"role" validate:"required, oneof=admin reader"`
		Token string `json:"token"`
	}

	//easyjson:json
	AddUserResponse struct {
		ID int64 `json:"id"`
	}

	//easyjson:json
	User struct {
		ID           int64  `json:"id"`
		Name         string `json:"name"`
		PasswordHash string `json:"-"`
		Token        string `json:"token"`
		Role         string `json:"role"`
	}
)