package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"net"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Container struct {
	Id               string
	Hostname         string
	DomainName       string
	WorkingDirectory string
	Entrypoint       string
	Cpu              string
	Image            string
	User             string
	RunArgs          string
	VolumesFrom      []string
	Volumes          []string
	Dns              []string
	Environment      []string
	Memory           int
	Port             []int
	Privileged       bool
	Network          *Network
	Config           *Config
}

func ParseReloadContext(c *cli.Context) (container *Container, err error) {
	if len(c.Args()) < 1 {
		return nil, fmt.Errorf("Error: no app name provided")
	}

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
	
	//private ip is required
	if c.String("privateip") == "" {
		return nil, fmt.Errorf("Error: Private Ip required")
	}
	
	//todo check if folder exists
	if c.GlobalString("c") != "" {
		configDir = c.GlobalString("c")
	}

	//check if private ip is a valid ip
	if net.ParseIP(c.String("privateip")) == nil {
		return nil, fmt.Errorf("Error: Private Ip not a valid  Addr")
	}

	//if public ip is set must be a valid ip
	if net.ParseIP(c.String("publicip")) == nil && c.String("publicip") != "" {
		return nil, fmt.Errorf("Error: Public Ip not a valid Addr")
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
	container.Entrypoint = c.String("entrypoint")
	container.VolumesFrom = c.StringSlice("-volumes-from")
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

	var cmd cmd

	cmd.add("docker", "run", "-d")

	if c.Memory != 0 {
		cmd.add("-m", string(c.Memory))
	}

	if c.WorkingDirectory != "" {
		cmd.add("-w", c.WorkingDirectory)
	}

	if c.Entrypoint != "" {
		cmd.add("-entrypoint", c.Entrypoint)
	}

	if c.Cpu != "" {
		cmd.add("-c", c.Cpu)
	}

	if c.User != "" {
		cmd.add("-u", c.User)
	}

	if c.Privileged {
		cmd.add("-privileged")
	}

	if len(c.Volumes) > 0 {
		for _, volumes := range c.Volumes {
			cmd.add("-v", volumes)
		}
	}

	if len(c.Environment) > 0 {
		for _, env := range c.Environment {
			cmd.add("-e", env)
		}
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

	if c.isStarted() {
		output, err := exec.Command("docker", "start", c.Id).Output()
		if err != nil {
			if strings.Contains(err.Error(), "exit status 1") {
				return fmt.Errorf("container %s already started", c.Id)
			}
			return fmt.Errorf("failed to start container: %v", err)
		}

		err = c.Network.AllocateNetwork(c.Id)
		if err != nil {
			return fmt.Errorf("failed to allocate network: %v", err)

		}

		if strings.Contains(string(output), c.Id) {
			info("started container " + c.Id)
		}
		return nil
	}

	output, err := c.buildCmd().Output()
	if err != nil {
		return fmt.Errorf("Error: error executing docker run comand %v ", err)
	}

	cid := c.matchContainerId(output)	

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
	info("container started with id " + c.Id)
	return nil
}

func (c *Container) Backup() {
	output,err

}

func (c *Container) () {
	
}

//maybe move to config
func (c *Container) Create() (err error) {

	err = c.Config.Save(*c)
	if err != nil {
		return err
	}
	return nil
}

func (c *Container) isStarted() bool {

	if c.Id != "" {
		return false
	}

	return true
}
func (c *Container) commitContainer() (err error){
	now:=time.Now()
	tag:= now.Year()+now.Month()+now.Day()
	output,err:=exec.Command("docker", "commit", c.Id, c.Image, tag).Output()
}

func (c *Container) matchContainerId(output []byte) (cid string) {
	regex := regexp.MustCompile("[0-9a-fA-F]{12}")
	cid = regex.FindString(string(output))
	return cid
}

