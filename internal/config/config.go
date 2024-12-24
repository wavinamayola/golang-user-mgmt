package config

type Config struct {
	Server   Server
	Database Database
}

type Server struct {
	Port int
}

type Database struct {
	Driver       string
	Port         string
	Host         string
	User         string
	Password     string
	DBName       string
	MigrationDir string
}
