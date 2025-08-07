package response

// Login .
type Login struct {
	UserId      int    `json:"userId"`
	Status      int    `json:"status"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	LastLoginAt int64  `json:"lastLoginAt"`
	Token       string `json:"token"`
}

// ChangePassword .
type ChangePassword struct {
}

// ChangePasswordByPhone .
type ChangePasswordByPhone struct {
}

// Register .
type Register struct {
	UserId int    `json:"userId"`
	Token  string `json:"token"`
}

type Logout struct {
}
