package godotenv

import (
	"fmt"
	"github.com/joho/godotenv"
)

func LoadEnv(path string) error {
	err := godotenv.Load(path)
	if err != nil {
		return fmt.Errorf("Didnt load env")
	}

	return nil
}
