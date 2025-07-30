package request

type Register struct {
	Name     string `json:"name" validator:"required,min=8,max=16"`
	Email    string `json:"email" validator:"required,email"`
	Password string `json:"password" validator:"required,min=8,max=24,password"`
}
type LogIn struct {
	Email    string `json:"email" validator:"required,email"`
	Password string `json:"password" validator:"required,min=8,max=24,password"`
}

type Create struct {
	Url string `json:"url" validator:"required,url"`
}
