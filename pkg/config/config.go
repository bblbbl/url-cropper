package config

import (
	"os"
	"strconv"
)

const defaultRedisDb = 0

var cnf *Config

type Config struct {
	App      App
	Http     Http
	Rpc      Rpc
	Redis    Redis
	Database Database
	Hash     Hash
}

type App struct {
	Mode string
	Host string
}

type Http struct {
	Schema string
	Port   string
}

type Rpc struct {
	Port    string
	Network string
}

type Redis struct {
	Host     string
	Port     string
	Password string
	DB       int
}

type Database struct {
	User     string
	Password string
	Host     string
	Port     string
	Database string
}

type Hash struct {
	Salt string
}

func InitConfig() {
	if cnf == nil {
		redisDb, err := strconv.Atoi(os.Getenv("REDIS_DB"))
		if err != nil {
			redisDb = defaultRedisDb
		}

		cnf = &Config{
			App: App{
				Mode: os.Getenv("GIN_MODE"),
				Host: os.Getenv("HOST_URL"),
			},
			Http: Http{
				Schema: os.Getenv("SCHEMA"),
				Port:   os.Getenv("HTTP_PORT"),
			},
			Rpc: Rpc{
				Port:    os.Getenv("RPC_PORT"),
				Network: os.Getenv("RPC_NETWORK"),
			},
			Redis: Redis{
				Host:     os.Getenv("REDIS_HOST"),
				Port:     os.Getenv("REDIS_PORT"),
				Password: os.Getenv("REDIS_PASSWORD"),
				DB:       redisDb,
			},
			Database: Database{
				User:     os.Getenv("DB_USER"),
				Password: os.Getenv("DB_PASSWORD"),
				Host:     os.Getenv("DB_HOST"),
				Port:     os.Getenv("DB_PORT"),
				Database: os.Getenv("DB_DATABASE"),
			},
			Hash: Hash{
				Salt: os.Getenv("HASH_SALT"),
			},
		}
	}
}

func GetConfig() *Config {
	return cnf
}
