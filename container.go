package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"net"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

type Container struct {
	Id               string
	Hostname         string
	DomainName       string
	WorkingDirectory string
	Entrypoint       string
	Cpu              string
	Image            string
	VolumesFrom      string
	User             string
	RunArgs          string
	Volumes          []string
	Dns              []string
	Environment      []string
	Port             []int
	Memory           int
	Privileged       bool
	Network          *Network
	Config           *Config
}

func ParseStartContext(c *cli.Context) (container *Container, err error) {

	if len(c.Args()) < 1 {
		return nil, fmt.Errorf("Error: you must provide a app name to run")
	}

	container = &Container{}
	container.Config = &Config{}
	container.Config.Filename = c.Args()[0]
	container.Config.Load(container)
	return container, nil

}

func ParseCreateContext(c *cli.Context) (container *Container, err error) {

	if len(c.Args()) < 2 {
		return nil, fmt.Errorf("Error: you must provide a app name and command to run")
	}

	if c.String("host") == "" {
		return nil, fmt.Errorf("Error: Hostname Required")
	}

	if c.String("image") == "" {
		return nil, fmt.Errorf("Error: Image Required")

	}
	//if we set dns servers we must verify if are valid ip addr
	if len(c.StringSlice("dns")) > 0 {
		for _, dns := range c.StringSlice("dns") {
			if net.ParseIP(dns) == nil {
				return nil, fmt.Errorf("Error: Dns not a valid Addr")
			}
		}
	}

	//if we set port we must verify if are valid ports
	if len(c.IntSlice("p")) > 0 {
		for _, port := range c.IntSlice("p") {
			if port < 1 && port > 65535 {
				return nil, fmt.Errorf("Error:  not a valid Port")
			}
		}
	}

	//private ip is required
	if c.String("privateip") == "" {
		return nil, fmt.Errorf("Error: Private Ip required")
	}

	//check if private ip is a valid ip
	if net.ParseIP(c.String("privateip")) == nil {
		return nil, fmt.Errorf("Error: Private Ip not a valid  Addr")
	}

	//if public ip is set must be a valid ip
	if net.ParseIP(c.String("publicip")) == nil && c.String("publicip") != "" {
		return nil, fmt.Errorf("Error: Public Ip not a valid Addr")
	}
	//todo check if folder exists
	if c.GlobalString("c") != "" {
		configDir = c.GlobalString("c")
	}

	appname := c.Args()[0]
	container = &Container{}
	container.Hostname = c.String("host")
	container.Dns = c.StringSlice("dns")
	container.Image = c.String("image")
	container.Cpu = c.String("c")
	container.Environment = c.StringSlice("e")
	container.User = c.String("u")
	container.WorkingDirectory = c.String("w")
	container.Port = c.IntSlice("p")
	container.Volumes = c.StringSlice("v")
	container.Network = &Network{}
	container.Network.PublicIp = c.String("publicip")
	container.Network.PrivateIp = c.String("privateip")
	container.Config = &Config{}
	container.Config.Filename = configDir + appname
	container.Config.Replace = c.Bool("r")
	container.RunArgs = strings.Join(c.Args()[1:], " ")
	return container, nil
}

func (c *Container) buildCmd() *exec.Cmd {
	//todo cpu|user|Environment....

	var cmd cmd

	cmd.add("docker", "run", "-d")

	if c.Memory != 0 {
		cmd.add("-m", string(c.Memory))

	}

	if c.Privileged {
		cmd.add("-privileged")

	}

	if len(c.Port) > 0 {
		for _, port := range c.Port {
			cmd.add("-p", strconv.Itoa(port))
		}
	}

	if len(c.Dns) > 0 {
		for _, dns := range c.Dns {
			cmd.add("-dns", dns)

		}
	}

	if len(c.Volumes) > 0 {
		for _, volumes := range c.Volumes {
			cmd.add("-v", volumes)

		}
	}

	cmd.add(c.Image)
	if c.Network.PublicIp != "" {
		cmd.add("/bin/sh -c", "'", "ip addr add dev eth0", c.Network.PublicIp+"/32", "&&", c.RunArgs, "'")
	} else {
		cmd.add(c.RunArgs)
	}

	return exec.Command("/bin/sh", "-c", cmd.String())

}

func (c *Container) Start() (err error) {

	output, err := c.buildCmd().Output()
	if err != nil {
		return fmt.Errorf("Error: error executing docker run comand %v ", err)
	}

	var re *regexp.Regexp
	var matchContanerId = "[0-9a-fA-F]{12}"

	re = regexp.MustCompile(matchContanerId)
	cid := re.FindString(string(output))

	if cid == "" {
		return fmt.Errorf("Error: no container Id provided by docker : %v", err)

	}

	err = c.Network.AllocateNetwork(cid)
	if err != nil {
		return err

	}
	c.Id = cid

	err = c.Config.Update(*c)
	if err != nil {
		return err

	}
	return nil

}

//maybe move to config
func (c *Container) Create() (err error) {

	err = c.Config.Save(*c)
	if err != nil {
		return err
	}
	return nil
}
