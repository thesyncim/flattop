package main

import (
	"fmt"
	"net"

	"strconv"
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
	Command          string
	Volumes          []string
	Dns              []string
	Environment      []string
	Port             []int
	Memory           int
	Privileged       bool
}

func (c *Container) ValidateContainer() error {

	if c.Hostname == "" {
		return fmt.Errorf("Error: Hostname Required")
	}

	if c.Image == "" {
		return fmt.Errorf("Error: Image Required")

	}
	//if we set dns servers we must verify if are valid ip addr
	if len(c.Dns) > 0 {
		for _, dns := range c.Dns {
			if net.ParseIP(dns) == nil {
				return fmt.Errorf("Error: Dns not a valid Addr")
			}
		}
	}

	//if we set port we must verify if are valid ports
	if len(c.Port) > 0 {
		for _, port := range c.Port {
			if port < 1 && port > 65535 {
				return fmt.Errorf("Error:  not a valid Port")
			}
		}
	}

	return nil
}

func (c *Container) buildDockerCmd(ip string) string {

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
	if ip != "" {
		cmd.add("/bin/sh -c", "'", "ip addr add dev eth0", ip+"/32", "&&", c.Command, "'")
	} else {
		cmd.add(c.Command)
	}

	return cmd.String()

}
