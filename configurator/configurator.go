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
	jsonFile, err = os.Open(file)
	defer jsonFile.Close()
	if err == nil {
		configuration, err = ioutil.ReadAll(jsonFile)
	}
	if err != nil {
		log.Println("unable to load configuration from", file, err)
	}
	json.Unmarshal([]byte(configuration), &v)
	return
}
