package main

import (
	"fmt"
	"net"
	"os/exec"
	"strings"
	"time"
)

type Network struct {
	PublicIp  string
	PrivateIp string
}

func (n *Network) resetPublicIp() {
	_, err := exec.Command("ip", "route", "flush", n.PublicIp).Output()
	if err != nil {
		exit("failed to flush public ip config", err)
	}

}
func (n *Network) SetPrivateIp(containerId string) error {
	time.Sleep(2 * time.Second)

	checkContainerExists := fmt.Sprintf("find /sys/fs/cgroup/devices -name %s* | wc -l", containerId)
	getnspid := fmt.Sprintf("head -n 1 $(find /sys/fs/cgroup/devices -name %s* | head -n 1)/tasks", containerId)

	out, err := exec.Command("/bin/sh", "-c", checkContainerExists).Output()
	if err != nil {
		exit("container doesnt exists", err)
	}

	if strings.TrimSpace(string(out)) != "1" {

		exit(string(out), "container "+containerId+" doesnt exist")
	}

	nspidbyte, err := exec.Command("/bin/sh", "-c", getnspid).Output()
	if err != nil {
		exit("failed to get nspid", err)
	}

	nspid := strings.TrimSpace(string(nspidbyte))
	local_ifname := "veth1" + nspid
	guest_ifname := "vethg" + nspid

	if nspid == "" {
		exit("no network  " + containerId + " ")
	}

	if _, err := exec.Command("/bin/sh", "-c", "mkdir -p /var/run/netns").Output(); err != nil {
		exit("failed to create /var/run/netns", err)
	}

	if _, err := exec.Command("/bin/sh", "-c", fmt.Sprintf("rm -f /var/run/netns/%s", nspid)).Output(); err != nil {
		exit("failed to remove /var/run/netns"+string(nspid), err)
	}

	if _, err := exec.Command("/bin/sh", "-c", fmt.Sprintf("ln -s /proc/%s/ns/net /var/run/netns/%s", nspid, nspid)).Output(); err != nil {
		exit("a", err)
	}

	if _, err := exec.Command("/bin/sh", "-c", fmt.Sprintf("ip link add name %s type veth peer name %s", local_ifname, guest_ifname)).Output(); err != nil {
		exit("b", err)
	}

	if _, err := exec.Command("/bin/sh", "-c", fmt.Sprintf("ip link set %s master docker0", local_ifname)).Output(); err != nil {
		exit("c", err)
	}

	if _, err := exec.Command("/bin/sh", "-c", fmt.Sprintf("ip link set %s up", local_ifname)).Output(); err != nil {
		exit("d", err)
	}

	if _, err := exec.Command("/bin/sh", "-c", fmt.Sprintf("ip link set %s netns %s", guest_ifname, nspid)).Output(); err != nil {
		exit("e", err)
	}

	if _, err := exec.Command("/bin/sh", "-c", fmt.Sprintf("ip netns exec %s ip link set %s name eth1", nspid, guest_ifname)).Output(); err != nil {
		exit("f", err)
	}

	if _, err := exec.Command("/bin/sh", "-c", fmt.Sprintf("ip netns exec %s ip addr add %s dev eth1", nspid, n.PrivateIp)).Output(); err != nil {
		exit("g", err)
	}

	if _, err := exec.Command("/bin/sh", "-c", fmt.Sprintf("ip netns exec %s ip link set eth1 up", nspid)).Output(); err != nil {
		exit("h", err)
	}

	return nil
}

func (n *Network) setPublicIp(containerId string) error {

	if n.PublicIp != "" {
		n.resetPublicIp()
		_, err := exec.Command("ip", "route", "add", "to", n.PublicIp, "via", n.PrivateIp).Output()
		if err != nil {
			exit("failed to set public ip config", err)
		}

	}

	return nil
}

func (n *Network) ValidateNetwork() error {

	//private ip is required
	if n.PrivateIp == "" {
		return fmt.Errorf("Error: Private Ip required")
	}

	//check if private ip is a valid ip
	if net.ParseIP(n.PrivateIp) == nil {
		return fmt.Errorf("Error: Private Ip not a valid  Addr")
	}

	//if public ip is set must be a valid ip
	if net.ParseIP(n.PublicIp) == nil && n.PublicIp != "" {
		return fmt.Errorf("Error: Public Ip not a valid Addr")
	}

	return nil
}
