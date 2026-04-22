package databases

import (
	_ "embed"
	"errors"
	"log"

	"github.com/goccy/go-yaml"
)

type MysqlConfigLoader struct{
	Environment environment
}

type DatabaseConfig struct{
	Host string `yaml:"host"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
	Port int `yaml:"port"`
}

type DbEnvironments struct{
	Development DatabaseConfig `yaml:"development"`
	Test DatabaseConfig `yaml:"test"`
	Production DatabaseConfig `yaml:"production"`
}

type environment int

const (
	Production  environment = iota + 1
	Test
	Development
)

//go:embed database.yml
var databaseConfig string
func(ml *MysqlConfigLoader) LoadMysqlConfiguration() (*DatabaseConfig,error) {
	configobj := &DbEnvironments{}
	if err := yaml.Unmarshal([]byte(databaseConfig), configobj); err != nil{
		log.Print(err.Error())
		return nil, err
	}
	switch(ml.Environment){
	case Production:
		return &configobj.Production, nil
	case Test:
		return &configobj.Test, nil
	case Development:
		return &configobj.Development, nil
	}
	return nil, errors.New("not found a suitable configuration")
}