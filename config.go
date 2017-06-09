package main

import (
	"io/ioutil"
	"fmt"
	"github.com/RangelReale/osin"
	"github.com/ghodss/yaml"
	"net/http"
)

func parse(data []byte) *osin.ServerConfig {
	var config osin.ServerConfig

	err := yaml.Unmarshal(data, &config)
	if err != nil {
		panic(err)
	}
	fmt.Println("Loaded config: ", config)
	return &config

}

func GetLocalConfig(filename string) (*osin.ServerConfig, error) {
	source, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return parse(source), nil
}

func GetRemoteConfig(url string) (*osin.ServerConfig, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	bytes, err := ioutil.ReadAll(resp.Body)
	return parse(bytes), nil

}

func GetDefaultConfig() (*osin.ServerConfig, error) {
	return &osin.ServerConfig{
		AuthorizationExpiration:   250,
		AccessExpiration:          3600,
		TokenType:                 "Bearer",
		AllowedAuthorizeTypes:     osin.AllowedAuthorizeType{osin.CODE},
		AllowedAccessTypes:        osin.AllowedAccessType{osin.PASSWORD},
		ErrorStatusCode:           200,
		AllowClientSecretInParams: true,
		AllowGetAccessRequest:     false,
		RetainTokenAfterRefresh:   false,
	}, nil

}


