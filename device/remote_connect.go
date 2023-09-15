package device

import (
	"errors"
	"fmt"
	"golang.org/x/crypto/ssh"
	"log"
	"time"
)

// Machine 设备相关信息
type Machine struct {
	IP       string
	Port     int
	Username string
	Password string
}

// ConnectDevice 连接设备
// return flag: true/false 连接成功或者连接失败 msg: 提示信息
func ConnectDevice(m *Machine) (flag bool, msg string, sshClient *ssh.Client) {
	config := &ssh.ClientConfig{
		User:            m.Username,
		Auth:            []ssh.AuthMethod{ssh.Password(m.Password)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         10 * time.Second,
	}
	//地址:端口
	addr := fmt.Sprintf("%v:%v", m.IP, m.Port)
	sshClient, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		errMsg := fmt.Sprintf("unable create ssh conn, %v", err)
		return false, errMsg, nil
	}
	return true, "connected success!", sshClient
}

// CheckConnectDevice 绑定一键连接按钮
func CheckConnectDevice(m *Machine) (flag bool, msg string) {
	flag, msg, client := ConnectDevice(m)
	if flag {
		defer func(client *ssh.Client) {
			_ = client.Close()
		}(client)
	}
	return
}

// MonitorConnectStatus  复用一个ssh客户端，监控设备连接状态
func MonitorConnectStatus(sshClient *ssh.Client, flag *bool) {
	go func() {
		for range time.Tick(time.Second * 5) {
			_, _, err := sshClient.Conn.SendRequest("keepalive@openssh.com", true, nil)
			if err != nil {
				log.Printf("SSH Connect is closed-----%v", err)
				*flag = false
			} else {
				log.Printf("SSH Connect is still started...")
				*flag = true
			}
		}
	}()
}

// HandleShell 执行一个cmd命令
func HandleShell(cmd string, sshClient *ssh.Client) (bytes []byte, err error) {
	session, err := sshClient.NewSession()
	if err != nil {
		fmt.Println("unable create ssh conn", err)
		return nil, errors.New("session create fail")
	}
	defer func(session *ssh.Session) {
		_ = session.Close()
	}(session)
	flag, bytes := execCmd(cmd, session)
	if !flag {
		return nil, errors.New(fmt.Sprintf("exec cmd >>>[%v],result: [%v]", cmd, flag))
	}
	log.Printf("exec cmd >>>[%v],result: [%v]\n %v", cmd, flag, string(bytes))
	return
}

// HandleBashShell 执行一个shell脚本
// bashPath 脚本文件的路径
func HandleBashShell(bashPath string, sshClient *ssh.Client) (bytes []byte, err error) {
	session, err := sshClient.NewSession()
	if err != nil {
		fmt.Println("unable create ssh conn", err)
		return nil, errors.New("session create fail")
	}
	defer func(session *ssh.Session) {
		_ = session.Close()
	}(session)
	flag, bytes := execShellBash(bashPath, session)
	if !flag {
		return nil, errors.New(fmt.Sprintf("exec bash shell >>>[%v],result: [%v]", bashPath, flag))
	}
	log.Printf("exec shell bash >>>[%v],result: [%v]\n %v", bashPath, flag, string(bytes))
	return
}

// CloseSshClient 关闭ssh客户端连接
func CloseSshClient(client *ssh.Client) {
	_ = client.Close()
}

func execCmd(cmd string, session *ssh.Session) (flag bool, bytes []byte) {
	bytes, err := session.CombinedOutput(cmd)
	if err != nil {
		log.Println(err)
		return false, nil
	}
	return true, bytes
}

func execShellBash(bashPath string, session *ssh.Session) (flag bool, bytes []byte) {
	bytes, err := session.CombinedOutput(fmt.Sprintf("bash -c %v", bashPath))
	if err != nil {
		log.Println(err)
		return false, nil
	}
	return true, bytes
}
