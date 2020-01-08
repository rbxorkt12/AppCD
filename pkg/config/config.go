package config

import (
	"fmt"
	"github.com/spf13/viper"

)

//argo Rollout 기능은 좀더 고민해봐야 할것 같다.

type appoconfig struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Orders []order `json:"orders"`
}

type order struct {
	Destination string `json:"destination"`
	Charts []chart `json:"charts"`
}

type chart struct {
	Repository string `json:"repository"`
	Path 	string `json:"path"`
	Branch	string `json:"branch"`
}

func GetConfig(path string) (*appoconfig,error) {
	config := &appoconfig{}
	viper.SetConfigName("Appoconfig")
	viper.AddConfigPath(".")
	if path != "" {
		viper.AddConfigPath(path)
	}
	if err := viper.ReadInConfig(); err != nil {
		return nil,fmt.Errorf("Can't read config file: %s \n", err)
	}
	if err := viper.Unmarshal(config); err != nil {
		return nil,fmt.Errorf("config file format error: %s \n", err)
	}
	return config,nil
}