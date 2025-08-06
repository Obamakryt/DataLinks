package connect_pool

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
)

type PostgresUrl struct {
	Username string `yaml:"username"`
	Hostname string `yaml:"hostname"`
	Port     uint   `yaml:"port"`
	DBName   string `yaml:"dbname"`
	Pass     string `env:"PGPASSWORD" yaml:"-"`
}
type PostgresCfg struct {
	Data PostgresUrl `yaml:"Postgres"`
}

func NewParseStruct(yml string, env string) *DbParseStruct {
	return &DbParseStruct{PostgresUrl: PostgresCfg{}, YamlCfg: yml, EnvNamePass: env}
}

type DbParseStruct struct {
	PostgresUrl PostgresCfg
	YamlCfg     string
	EnvNamePass string
}

func (p *DbParseStruct) ParseConfig() error {
	err := cleanenv.ReadConfig(p.YamlCfg, &p.PostgresUrl)
	if err != nil {
		return fmt.Errorf("couldnt read yml cfg %w", err)
	}
	err = cleanenv.ReadEnv(&p.PostgresUrl)
	if err != nil {
		return fmt.Errorf("couldnt find postgres password in env %w", err)
	}
	return nil
}
