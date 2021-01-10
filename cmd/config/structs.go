package config

type AppConfigs struct {
	Api *ApiConfigs `json:"api"`
	Db  *DbConfigs  `json:"db"`
}

type ApiConfigs struct {
	Port string `json:"port" envconfig:"optional"`
}

type DbConfigs struct {
	Host     string `json:"host" envconfig:"optional"`
	Port     string `json:"port" envconfig:"optional"`
	Username string `json:"username" envconfig:"optional"`
	Password string `json:"password" envconfig:"optional"`
	Database string `json:"database" envconfig:"optional"`
}
