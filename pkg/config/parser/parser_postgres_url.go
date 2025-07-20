package parser

import (
	db "DataLinks/internal/storages/postgreSQL/TakeConnect"
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
)

func ParseDBcfg(parseStruct *db.ParseStruct) (error, *db.PostgresUrl) {
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
