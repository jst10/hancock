package config

import (
	"encoding/json"
	"github.com/vrischmann/envconfig"
	"io/ioutil"
	custom_errors "made.by.jst10/outfit7/hancock/cmd/custom_errors"
	"os"
)

func LoadConfig(appConfigs *AppConfigs)  *custom_errors.CustomError {
	file, err := os.Open("config.json")
	if err != nil {
		return custom_errors.GetErrorLoadingConfigs(err,"Opening config file")
	}
	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return custom_errors.GetErrorLoadingConfigs(err,"Failed to read appsettings")
	}
	err = json.Unmarshal(data, appConfigs)
	if err != nil {
		return custom_errors.GetErrorLoadingConfigs(err, "Failed to unmarshal appsettings")
	}
	err = envconfig.Init(appConfigs)
	if err != nil {
		return custom_errors.GetErrorLoadingConfigs(err, "Failed to update with env vars")
	}
	return nil
}
