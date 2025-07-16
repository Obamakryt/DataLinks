package TakeConnect

type PostgresUrl struct {
	Username string `yaml:"username"`
	Hostname string `yaml:"hostname"`
	Port     uint   `yaml:"port"`
	DBName   string `yaml:"dbname"`
	Pass     string `env:"PSQL_PASS"`
}
type ParseStruct struct {
	PostgresUrl
	YamlCfg     string
	EnvNamePass string
}

func NewParseStruct(yml string, env string) *ParseStruct {
	return &ParseStruct{PostgresUrl{}, yml, env}
}

func NewPostgresUrl() *PostgresUrl {
	return &PostgresUrl{}
}
