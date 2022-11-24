package api

import (
	"fmt"
	"moonspace/api/middleware"
	"moonspace/repository/types"
	"moonspace/service/payment"
)

type CORS struct {
	AllowOrigins  []string `yaml:"allow-origins"`
	AllowMethods  []string `yaml:"allow-methods"`
	AllowHeaders  []string `yaml:"allow-headers"`
	ExposeHeaders []string `yaml:"expose-headers"`
}

type OAuth struct {
	Realm     string                       `yaml:"realm"`
	ClientID  string                       `yaml:"clientId"`
	URL       string                       `yaml:"url"`
	AdminData *middleware.AdminCredentials `yaml:"adminCredentials"`
}

type Security struct {
	OAuth *OAuth `yaml:"oauth"`
	CORS  *CORS  `yaml:"cors"`
}

type Config struct {
	Host     string              `yaml:"host"`
	Port     int                 `yaml:"port"`
	Security *Security           `yaml:"security"`
	DB       *types.Config       `yaml:"persistence"`
	Payment  *payment.PaymentCfg `yaml:"payment"`
}

type ServerConfig struct {
	Server *Config `yaml:"server"`
}

func (cfg *ServerConfig) String() string {
	return fmt.Sprintf("\nHost: %s\nPort:%d\nSecurity: %+v\nDB: %+v\nPayment: %v",
		cfg.Server.Host, cfg.Server.Port, cfg.Server.Security, cfg.Server.DB, cfg.Server.Payment,
	)
}
