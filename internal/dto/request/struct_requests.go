package request

type Register struct {
	Name     string `json:"name" validate:"required,min=8,max=16"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,password"`
}
type LogIn struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,password"`
}

type Add struct {
	Url string `json:"url" validate:"required,url,http"`
	TakeChart
}
type TakeChart struct {
	UserId int `json:"-"`
}
type Swap struct {
	Add
	NewUrl string `json:"newurl" validator:"required,url,http"`
}
type Delete struct {
	Add
}
