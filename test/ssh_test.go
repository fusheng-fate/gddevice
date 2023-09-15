package test

import (
	"gddevice/device"
	"log"
	"testing"
	"time"
)

func TestSshConnect(t *testing.T) {
	machine := device.Machine{IP: "42.192.200.28", Port: 22, Username: "lgb", Password: "lgb@1234"}
	flag, msg, client := device.ConnectDevice(&machine)
	if !flag {
		log.Println(msg)
		return
	}
	var status = false
	device.MonitorConnectStatus(client, &status)
	defer device.CloseSshClient(client)
	for i := 0; i < 10; i++ {
		time.Sleep(time.Second * 3)
		log.Printf("monitor ssh connect status >>>[%t]", status)
	}
	_, err := device.HandleShell("cd / && ls -lth", client)
	if err != nil {
		log.Println(err)
	}
}
