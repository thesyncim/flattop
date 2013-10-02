package main

import (
	"fmt"
	"os/exec"
	"strings"
	//"time"
)

type Network struct {
	PublicIp  string
	PrivateIp string
}

func (n *Network) AllocateNetwork(cid string) (err error) {
	n.setPrivateIp(cid)
	if err != nil {
		return err
	}
	n.setPublicIp(cid)
	if err != nil {
		return err
	}
	return nil
}

func (n *Network) resetPublicIp() (err error) {
	_, err = exec.Command("ip", "route", "flush", n.PublicIp).Output()
	if err != nil {
		return fmt.Errorf("failed to reset Public ip : %v", err)
	}
	return nil
}
func (n *Network) setPrivateIp(containerId string) error {

	//todo replace exit by fmt.Errorf
	//	time.Sleep(2 * time.Second)

	checkContainerExists := fmt.Sprintf("find /sys/fs/cgroup/devices -name %s* | wc -l", containerId)
	out, err := exec.Command("/bin/sh", "-c", checkContainerExists).Output()
	if err != nil {
		return fmt.Errorf("container doesnt exists %s", err)
	}

	if strings.TrimSpace(string(out)) != "1" {

		return fmt.Errorf(string(out), "container "+containerId+" doesnt exist")
	}

	getnspid := fmt.Sprintf("head -n 1 $(find /sys/fs/cgroup/devices -name %s* | head -n 1)/tasks", containerId)
	nspidbyte, err := exec.Command("/bin/sh", "-c", getnspid).Output()
	if err != nil {
		return fmt.Errorf("failed to get nspid %v", err)
	}

	nspid := strings.TrimSpace(string(nspidbyte))
	local_ifname := "veth1" + nspid
	guest_ifname := "vethg" + nspid

	if nspid == "" {
		return fmt.Errorf("no network nspid " + containerId + " ")
	}

	if _, err := exec.Command("/bin/sh", "-c", "mkdir -p /var/run/netns").Output(); err != nil {
		return err
	}

	if _, err := exec.Command("/bin/sh", "-c", fmt.Sprintf("rm -f /var/run/netns/%s", nspid)).Output(); err != nil {
		return err
	}

	if _, err := exec.Command("/bin/sh", "-c", fmt.Sprintf("ln -s /proc/%s/ns/net /var/run/netns/%s", nspid, nspid)).Output(); err != nil {
		return err
	}

	if _, err := exec.Command("/bin/sh", "-c", fmt.Sprintf("ip link add name %s type veth peer name %s", local_ifname, guest_ifname)).Output(); err != nil {
		return err
	}

	if _, err := exec.Command("/bin/sh", "-c", fmt.Sprintf("ip link set %s master docker0", local_ifname)).Output(); err != nil {
		return err
	}

	if _, err := exec.Command("/bin/sh", "-c", fmt.Sprintf("ip link set %s up", local_ifname)).Output(); err != nil {
		return err
	}

	if _, err := exec.Command("/bin/sh", "-c", fmt.Sprintf("ip link set %s netns %s", guest_ifname, nspid)).Output(); err != nil {
		return err
	}

	if _, err := exec.Command("/bin/sh", "-c", fmt.Sprintf("ip netns exec %s ip link set %s name eth1", nspid, guest_ifname)).Output(); err != nil {
		return err
	}

	if _, err := exec.Command("/bin/sh", "-c", fmt.Sprintf("ip netns exec %s ip addr add %s dev eth1", nspid, n.PrivateIp)).Output(); err != nil {
		return err
	}

	if _, err := exec.Command("/bin/sh", "-c", fmt.Sprintf("ip netns exec %s ip link set eth1 up", nspid)).Output(); err != nil {
		return err
	}

	return nil
}

func (n *Network) setPublicIp(containerId string) error {

	if n.PublicIp != "" {
		n.resetPublicIp()
		_, err := exec.Command("ip", "route", "add", "to", n.PublicIp, "via", n.PrivateIp).Output()
		if err != nil {
			return fmt.Errorf("failed to set public ip config %v", err)
		}

	}

	return nil
}
