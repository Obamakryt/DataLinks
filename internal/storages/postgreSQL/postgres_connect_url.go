package postgreSQL

type PostgresUrl struct {
	Username string `env:"username"`
	Hostname string `env:"hostname"`
	Port     uint   `env:"port"`
	DBName   string `env:"dbname"`
}

func NewPostgresUrl() *PostgresUrl {
	return &PostgresUrl{}
}
