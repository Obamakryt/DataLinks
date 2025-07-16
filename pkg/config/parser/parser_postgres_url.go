package parser

import (
	db "DataLinks/internal/storages/postgreSQL/TakeConnect"
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
)

func ParseDBcfg(parseStruct *db.ParseStruct) (error, *db.PostgresUrl) {
	_, ok := os.LookupEnv("PSQL_PASS")
	if !ok {
		return fmt.Errorf("cant find PSQL_PASS"), &db.PostgresUrl{}
	}
	err := cleanenv.ReadConfig(parseStruct.YamlCfg, &parseStruct.PostgresUrl)
	if err != nil {
		return fmt.Errorf("couldnt read yml cfg %w", err), &db.PostgresUrl{}
	}
	err = cleanenv.ReadEnv(&parseStruct.PostgresUrl)
	if err != nil {
		return fmt.Errorf("couldnt find postgres password in env %w", err), &db.PostgresUrl{}
	}
	return nil, &parseStruct.PostgresUrl
}
