package user

type (
	User struct {
		ID           int64  `json:"id"`
		Name         string `json:"name"`
		PasswordHash string `json:"-"`
		Token        string `json:"token"`
		Role         string `json:"role"`
	}
)
