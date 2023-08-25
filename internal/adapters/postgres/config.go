package postgres

type Config struct {
	Host     string
	Port     int
	UserName string
	Password string
	DBName   string
	PoolSize int
}
