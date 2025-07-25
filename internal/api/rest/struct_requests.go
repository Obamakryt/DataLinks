package rest

type RequestRegister struct {
	Name     string `json:"name" validator:"required,min=8,max=16"`
	Email    string `json:"email" validator:"required,email"`
	Password string `json:"password" validator:"required,min=8,max=24,password"`
}
type RequestLogIn struct {
	Email    string `json:"email" validator:"required,email"`
	Password string `json:"password" validator:"required,min=8,max=24,password"`
}
