package model

import (
	"io/ioutil"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

//config yaml struct
type Service struct {
	APIVersion string `yaml:"apiVersion"`
	Spec       struct {
		Ports struct {
			Name string `yaml:"name"`
			Addr string `yaml:"bind_addr"`
			Crt  string `yaml:"crt"`
			Key  string `yaml:"key"`
		} `yaml:"ports"`
		DBpg struct {
			Name              string `yaml:"name"`
			Host              string `yaml:"host"`
			Port              uint16 `yaml:"port"`
			User              string `yaml:"user"`
			Password          string `yaml:"password"`
			Database          string `yaml:"database"`
			MaxConnLifetime   int    `yaml:"max_conn_lifetime"`
			MaxConnIdletime   int    `yaml:"max_conn_idletime"`
			MaxConns          int32  `yaml:"max_conns"`
			MinConns          int32  `yaml:"min_conns"`
			HealthCheckPeriod int    `yaml:"health_check_period"`
		} `yaml:"dbpg"`
		DBms struct {
			Name string `yaml:"name"`
			Url  string `yaml:"url"`
		} `yaml:"dbms"`
		Jwt struct {
			TokenDecode string `yaml:"token"`
			LifeTerm    int    `yaml:"term"`
		} `yaml:"jwt"`
		Client struct {
			UrlAzgazTest      string `yaml:"url_azgaz_test"`
			KeyAzgazTest      string `yaml:"key_azgaz_test"`
			UrlMailingService string `yaml:"url_mailing_service"`
		} `yaml:"client"`
		Queryies struct {
			Booking  string `yaml:"booking"`
			Logistic string `yaml:"logistic"`
		} `yaml:"queryies"`
	} `yaml:"spec"`
}

//New config
func NewConfig() (*Service, error) {

	var service *Service

	f, err := filepath.Abs("/root/config/server_http_rest_ar.yaml")
	if err != nil {
		return nil, err
	}

	y, err := ioutil.ReadFile(f)
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(y, &service); err != nil {
		return nil, err
	}

	return service, nil

}
