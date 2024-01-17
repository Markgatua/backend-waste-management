package appsettings

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

func UnmarshalAppSettings(data []byte) (AppSettings, error) {
	var r AppSettings
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *AppSettings) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type AppSettings struct {
	Debug                    bool         `json:"debug"`
	DBMasterConnectionString string       `json:"db_master_connection_string"`
}


func GetAppSettings() (AppSettings, error) {
	jsonFile, err := os.Open(".app_settings.json")
	defer jsonFile.Close()
	if err == nil {
		byteValue, _ := ioutil.ReadAll(jsonFile)
		return UnmarshalAppSettings(byteValue)
	} else {
		return AppSettings{}, errors.New(fmt.Sprint("Error reading from json file ", err.Error()))
	}
}
