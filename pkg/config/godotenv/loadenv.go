package godotenv

import (
	"fmt"
	"github.com/joho/godotenv"
)

func LoadEnv(path string) error {
	err := godotenv.Load(path)
	if err != nil {
		return fmt.Errorf("didnt load env")
	}

	return nil
}
