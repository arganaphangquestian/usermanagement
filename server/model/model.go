package model

// InputUser or Register model
type InputUser struct {
	Username string  `json:"username"`
	Name     string  `json:"name"`
	Email    string  `json:"email"`
	Password string  `json:"password"`
	Referral *string `json:"referral"`
}

// Login model
type Login struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// User model
type User struct {
	ID       string  `json:"id"`
	Username string  `json:"username"`
	Name     string  `json:"name"`
	Email    string  `json:"email"`
	Referral *string `json:"referral"`
	Role     string  `json:"role"`
}

// UserWithPassword struct
type UserWithPassword struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Username string  `json:"username"`
	Email    string  `json:"email"`
	Password string  `json:"password"`
	Referral *string `json:"referral"`
	Role     string  `json:"role"`
}

// UserToken model
type UserToken struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
