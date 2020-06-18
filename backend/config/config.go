package config

import (
	"fmt"
	"sync"

	"github.com/kelseyhightower/envconfig"
)

type Vars struct {
	MySQLUser     string `envconfig:"MYSQL_USER"`
	MySQLPassword string `envconfig:"MYSQL_PASSWORD"`
	MySQLHost     string `envconfig:"MYSQL_HOST"`
	MySQLPort     string `envconfig:"MYSQL_PORT"`
	MySQLDatabase string `envconfig:"MYSQL_DATABASE"`
	HTTPPort      int    `envconfig:"PORT" default:"5001"`
}

func Process() (*Vars, error) {
	var vars Vars
	if err := envconfig.Process("", &vars); err != nil {
		return nil, err
	}

	return &vars, nil
}

func MustProcess() *Vars {
	vars, err := Process()
	if err != nil {
		panic(err)
	}
	return vars
}

var DefaultVars = &Vars{}
var once sync.Once

func MustProcessDefault() {
	once.Do(func() {
		DefaultVars = MustProcess()
	})
}

func (v *Vars) DBURL() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true&loc=UTC",
		v.MySQLUser, v.MySQLPassword, v.MySQLHost, v.MySQLPort, v.MySQLDatabase,
	)
}
