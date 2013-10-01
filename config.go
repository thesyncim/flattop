package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"regexp"
)

type Config struct {
	Name string
	Network
	Container
}

func (c *Config) WriteToFile(filename string) {

	b, err := json.Marshal(*c)
	if err != nil {
		exit("unable to Marshal Config:", err)
	}

	prettyb, err := Pretty(b)
	if err != nil {
		exit("unable to Indent:", err)
	}

	err = ioutil.WriteFile(filename, prettyb, 0644)
	if err != nil {
		exit("failed to write config to:"+filename, err)
	}

}

func (c *Config) Save(filename string, repalce bool) {
	if fileExists(filename) {
		if repalce {
			var confirm string
			fmt.Println("are you sure ?:yes|no")
			fmt.Scan(&confirm)
			if confirm == "yes" {
				c.writeToFile(filename)
				info("configFile " + filename + " replaced")
			}
		} else {
			exit("already exists use -r to replace")
		}

	} else {
		c.writeToFile(filename)
		info("new configFile " + filename + " saved")

	}

}

func (c *Config) Run() (containerId string) {

	out, err := exec.Command("/bin/sh", "-c", c.buildDockerCmd(c.PublicIp)).Output()
	if err != nil {
		exit("failed to run ", c.buildDockerCmd(c.PublicIp), err)
	}

	var re *regexp.Regexp
	re = regexp.MustCompile("[0-9a-fA-F]{12}")
	containerId = re.FindString(string(out))
	if containerId == "" {
		exit("failed with error ", string(out))
	}
	c.Id = containerId

	return containerId
}

func (c *Config) LoadConfigFile(filename string) {

	if fileExists(filename) {

		b, err := ioutil.ReadFile(filename)
		if err != nil {
			exit("failed to read " + filename)
		}

		err = json.Unmarshal(b, c)
		if err != nil {
			exit("error to Unmarshal", err)
		}

	} else {
		exit("Error: missing configuration for " + filename)
	}
}

func (c *Config) LoadConfigFiles() error {

	usr, err := user.Current()
	if err != nil {
		exit(err)
	}

	configFolder := filepath.Join(usr.HomeDir, ".dockercmd")
	err = filepath.Walk(configFolder, visit)
	if err != nil {
		exit("failed to load config files")
	}
	return nil

}

func visit(path string, f os.FileInfo, err error) error {
	fmt.Printf("Visited: %s\n", path)
	return nil
}
