package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"
)

type Config struct {
	Name string
	Network
	Container
}

func (c *Config) writeToFile(filename string) {
	b, err := json.Marshal(*c)
	if err != nil {
		log.Fatalln("unable to Marshal Config:", err)
	}
	prettyb, err := Pretty(b)
	if err != nil {
		log.Fatalln("unable to Indent:", err)
	}

	err = ioutil.WriteFile(filename, prettyb, 0644)
	if err != nil {
		log.Fatalln("failed to write config to:"+filename, err)
	}

}

func (c *Config) Save(filename string, repalce bool) {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		c.writeToFile(filename)
		info("new configFile " + filename + " saved")
	} else {
		if repalce {
			var confirm string
			fmt.Println("are you sure ?:yes|no")
			fmt.Scan(&confirm)
			if confirm == "yes" {
				c.writeToFile(filename)
				info("configFile " + filename + " replaced")
			}
		} else {
			log.Fatal("already exists use -r to replace")
		}
	}
}

func (c *Config) LoadConfigFiles() error {

	usr, err := user.Current()
	if err != nil {
		log.Fatalln(err)
	}

	configFolder := filepath.Join(usr.HomeDir, ".dockercmd")
	err = filepath.Walk(configFolder, visit)
	if err != nil {
		log.Fatalln("failed to load config files")
	}
	return nil

}

func visit(path string, f os.FileInfo, err error) error {
	fmt.Printf("Visited: %s\n", path)
	return nil
}
