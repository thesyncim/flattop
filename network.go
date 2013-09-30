package main

import (
	"fmt"
	"log"
	"net"
	"os/exec"
	"strings"
	"time"
)

type Network struct {
	PublicIp  string
	PrivateIp string
}

func (n *Network) SetPrivateIp(containerId string) error {
	time.Sleep(2 * time.Second)

	checkContainerExists := fmt.Sprintf("find /sys/fs/cgroup/devices -name %s* | wc -l", containerId)
	getnspid := fmt.Sprintf("head -n 1 $(find /sys/fs/cgroup/devices -name %s* | head -n 1)/tasks", containerId)

	out, err := exec.Command("/bin/sh", "-c", checkContainerExists).Output()
	if err != nil {
		log.Fatal("container doesnt exists", err)
	}

	if strings.TrimSpace(string(out)) != "1" {

		log.Fatal(string(out), "container "+containerId+" doesnt exist")
	}

	nspidbyte, err := exec.Command("/bin/sh", "-c", getnspid).Output()
	if err != nil {
		log.Fatal("failed to get nspid", err)
	}

	nspid := strings.TrimSpace(string(nspidbyte))
	local_ifname := "veth1" + nspid
	guest_ifname := "vethg" + nspid

	if nspid == "" {
		log.Fatal("no network  " + containerId + " ")
	}

	if _, err := exec.Command("/bin/sh", "-c", "mkdir -p /var/run/netns").Output(); err != nil {
		log.Fatal("failed to create /var/run/netns", err)
	}

	if _, err := exec.Command("/bin/sh", "-c", fmt.Sprintf("rm -f /var/run/netns/%s", nspid)).Output(); err != nil {
		log.Fatal("failed to remove /var/run/netns"+string(nspid), err)
	}

	if _, err := exec.Command("/bin/sh", "-c", fmt.Sprintf("ln -s /proc/%s/ns/net /var/run/netns/%s", nspid, nspid)).Output(); err != nil {
		log.Fatal("a", err)
	}

	if _, err := exec.Command("/bin/sh", "-c", fmt.Sprintf("ip link add name %s type veth peer name %s", local_ifname, guest_ifname)).Output(); err != nil {
		log.Fatal("b", err)
	}

	if _, err := exec.Command("/bin/sh", "-c", fmt.Sprintf("ip link set %s master docker0", local_ifname)).Output(); err != nil {
		log.Fatal("c", err)
	}

	if _, err := exec.Command("/bin/sh", "-c", fmt.Sprintf("ip link set %s up", local_ifname)).Output(); err != nil {
		log.Fatal("d", err)
	}

	if _, err := exec.Command("/bin/sh", "-c", fmt.Sprintf("ip link set %s netns %s", guest_ifname, nspid)).Output(); err != nil {
		log.Fatal("e", err)
	}

	if _, err := exec.Command("/bin/sh", "-c", fmt.Sprintf("ip netns exec %s ip link set %s name eth1", nspid, guest_ifname)).Output(); err != nil {
		log.Fatal("f", err)
	}

	if _, err := exec.Command("/bin/sh", "-c", fmt.Sprintf("ip netns exec %s ip addr add %s dev eth1", nspid, n.PrivateIp)).Output(); err != nil {
		log.Fatal("g", err)
	}

	if _, err := exec.Command("/bin/sh", "-c", fmt.Sprintf("ip netns exec %s ip link set eth1 up", nspid)).Output(); err != nil {
		log.Fatal("h", err)
	}

	return nil
}

func (n *Network) setPublicIp(containerId string) error {
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
