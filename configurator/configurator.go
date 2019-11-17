package configurator

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

func init() {
}

func LoadConfiguration(file string, v interface{}) (err error) {
	var configuration []byte
	var jsonFile *os.File
	if jsonFile, err = os.Open(file); err == nil {
		defer jsonFile.Close()
		if configuration, err = ioutil.ReadAll(jsonFile); err == nil {
			err = json.Unmarshal(configuration, &v)
		}
	}
	if err != nil {
		log.Println("unable to load configuration from", file, err)
	}
	return
}
