package authentication

// LoginRequest object which is the request for Login function
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse object which is the response for Login function
type LoginResponse struct {
	AccessToken string `json:"access_token"`
}

// SignupRequest object which is the request for Signup function
type SignupRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Phone    string `json:"phone"`
}

// SignupReponse object which is the response for Signup function
type SignupReponse struct {
	Status bool
}

// Authentication object which is used to store in database
type Authentication struct {
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
	PK        string `json:"PK,omitempty"`
	ID        string `json:"id,omitempty"`
	Email     string `json:"email,omitempty"`
	Password  string `json:"password,omitempty"`
	Name      string `json:"name,omitempty"`
	Phone     string `json:"phone,omitempty"`
}
