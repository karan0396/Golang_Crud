package model

type User struct {
	ID         int    `json:"id,omitempty"`
	FName      string `json:"firstname,omitempty" validate:"required"`
	Lname      string `json:"lastname,omitempty" validate:"required"`
	Email      string `json:"email,omitempty",validate:"required,email"`
	Dob        string `json:"dob,omitempty"`
	Password   string `json:"password,omitempty",validate:"required,password"`
	Created_at string `json:"_"`
	Updated_at string `json:"_"`
	Archieved  bool   `json:"_"`
}

type Credential struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SoftDelete struct {
	ID    int    `json:"id,omitempty"`
	FName string `json:"firstname,omitempty" validate:"required"`
}
