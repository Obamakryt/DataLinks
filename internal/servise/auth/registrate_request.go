package auth

type RequestRegister struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
type StorageRegister struct {
	Name     string
	Email    string
	HashPass string
}

//hashpass := HashingPass(r.Password)
