package main

import (
	"github.com/newm4n/mihp/internal"
	"io/ioutil"
	"os"
)

func SetupConfig(configFile string) (err error) {
	Mode := ""
	var Config *internal.MIHPConfig
	fInfo, err := os.Stat(configFile)
	if err != nil || fInfo.IsDir() {
		Mode = "NEW"
		Config = &internal.MIHPConfig{}
	} else {
		file, err := os.Open(configFile)
		if err != nil {
			return err
		}
		yamlBytes, err := ioutil.ReadAll(file)
		if err != nil {
			return err
		}
		cfg, err := internal.YAMLToMIHPConfig(yamlBytes)
		if err != nil {
			return err
		}
		Mode = "EDIT"
		Config = cfg
	}

}
