package etc

import (
	"github.com/joho/godotenv"
	"gopkg.in/yaml.v2"
	"io"
	"os"
	"sync"
	"urls/pkg/utils"
)

var (
	cnf  *Config
	once sync.Once
)

type Config struct {
	App      App      `yaml:"app"`
	Http     Http     `yaml:"http"`
	Rpc      Rpc      `yaml:"rpc"`
	Redis    Redis    `yaml:"redis"`
	Database Database `yaml:"database"`
	Hash     Hash     `yaml:"hash"`
}

type App struct {
	Mode string `yaml:"mode"`
	Host string `yaml:"host"`
}

type Http struct {
	Schema string `yaml:"schema"`
	Port   int    `yaml:"port"`
}

type Rpc struct {
	Port    int    `yaml:"port"`
	Network string `yaml:"network"`
}

type Redis struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Password string `yaml:"password"`
	DB       int    `yaml:"DB"`
}

type Database struct {
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Database string `yaml:"database"`
}

type Hash struct {
	Salt string `yaml:"salt"`
}

func GetConfig() *Config {
	once.Do(func() {
		initialise()
	})

	return cnf
}

func initialise() {
	err := godotenv.Load()
	if err != nil {
		GetLogger().Fatalf(".env read failed: %e\n", err)
	}

	f, err := os.Open("configs/app.yaml")
	if err != nil {
		panic("failed open config file")
	}

	cnfBytes, err := io.ReadAll(f)
	if err != nil {
		panic("failed read config file")
	}

	expanded := os.ExpandEnv(utils.B2S(cnfBytes))
	if err = yaml.Unmarshal([]byte(expanded), &cnf); err != nil {
		panic("failed parse app config")
	}
}
