package configurator

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

func init() {
}

func LoadConfiguration(file string, v interface{}) {
	var configuration []byte
	var err error
	var jsonFile *os.File
	jsonFile, err = os.Open(file)
	defer jsonFile.Close()
	if err == nil {
		configuration, err = ioutil.ReadAll(jsonFile)
	}
	if err != nil {
		fmt.Println("Unable to load file from", file, err)
	}
	json.Unmarshal([]byte(configuration), &v)
}
