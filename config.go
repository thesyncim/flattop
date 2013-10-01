package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	//"os"
	//"os/user"
	//"path/filepath"
)

type Config struct {
	Filename string
	Replace  bool
}

func (c *Config) writeToFile(data interface{}) {

	b, err := json.Marshal(data)
	if err != nil {
		exit("unable to Marshal Config:", err)
	}

	prettyb, err := Pretty(b)
	if err != nil {
		exit("unable to Indent:", err)
	}

	err = ioutil.WriteFile(c.Filename, prettyb, 0644)
	if err != nil {
		exit("failed to write config to:"+c.Filename, err)
	}

}

func (c *Config) Save(data Container) (err error) {
	if fileExists(c.Filename) {
		if c.Replace {
			//reset replace
			c.Replace = false
			var confirm string
			fmt.Println("Do you want to continue [Y/n]?")
			fmt.Scan(&confirm)
			if (confirm == "y") || (confirm == "Y") {
				c.writeToFile(data)
				info("configFile " + c.Filename + " replaced")
			}
		} else {
			exit("already exists use -r to replace")
		}

	} else {
		c.writeToFile(data)
		info("new configFile " + c.Filename + " saved")

	}
	return nil

}

func (c *Config) Load(container *Container) {

	if fileExists(configDir + c.Filename) {

		b, err := ioutil.ReadFile(configDir + c.Filename)
		if err != nil {
			exit("failed to read " + c.Filename)
		}

		err = json.Unmarshal(b, container)
		if err != nil {
			exit("error to Unmarshal", err)
		}

	} else {
		exit("Error: missing configuration for " + c.Filename)
	}
}

//func (c *Config) LoadConfigFiles() error {
//
//	usr, err := user.Current()
//	if err != nil {
//		exit(err)
//	}
//
//	configFolder := filepath.Join(usr.HomeDir, ".dockercmd")
//	err = filepath.Walk(configFolder, visit)
//	if err != nil {
//		exit("failed to load config files")
//	}
//	return nil
//
//}

//func visit(path string, f os.FileInfo, err error) error {
//	fmt.Printf("Visited: %s\n", path)
//	return nil
//}
