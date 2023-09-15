package sftp

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"log"
	"time"
)

// Machine 设备相关信息
type Machine struct {
	IP       string
	Port     string
	Username string
	Password string
}

// ConnectDevice 测试连接设备
// return flag: true/false 连接成功或者连接失败 msg: 提示信息
func ConnectDevice(m *Machine) (flag bool, msg string) {
	config := &ssh.ClientConfig{
		User:            m.Username,
		Auth:            []ssh.AuthMethod{ssh.Password(m.Password)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         10 * time.Second,
	}
	//地址:端口
	addr := fmt.Sprintf("%v:%v", m.IP, m.Port)
	_, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		errMsg := fmt.Sprintf("unable create ssh conn, %v", err)
		return false, errMsg
	}
	return true, "connected success!"
}

func RemoteConnect(m *Machine) error {
	//初始化连接信息
	config := &ssh.ClientConfig{
		User:            m.Username,
		Auth:            []ssh.AuthMethod{ssh.Password(m.Password)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         10 * time.Second,
	}
	addr := fmt.Sprintf("%v:%v", m.IP, m.Port)
	sshClient, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		fmt.Println("unable create ssh conn", err)
		return err
	}
	defer sshClient.Close()
	session, err := sshClient.NewSession()
	if err != nil {
		fmt.Println("unable create ssh conn", err)
		return err
	}
	defer session.Close()
	flag := removeFile("/home/lgb/aa.txt", session)
	log.Printf("移除一个文件 >>>> %v\n", flag)
	return err
}

func removeFile(filePath string, session *ssh.Session) (flag bool) {
	err := session.Run(fmt.Sprintf("rm -rf %v", filePath))
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}
