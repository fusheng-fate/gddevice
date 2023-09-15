package main

import (
	"errors"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"gddevice/device"
	"github.com/flopp/go-findfont"
	"golang.org/x/crypto/ssh"
	"log"
	"os"
	"strings"
)

const (
	initStatus     = "等待连接"
	successConnect = "连接成功"
	failConnect    = "连接失败"
)

var machine = device.Machine{IP: "42.192.200.28", Port: 22, Username: "lgb", Password: "lgb@1234"}

func main() {
	myApp := app.New()
	myApp.Settings().SetTheme(theme.DefaultTheme())
	myApp.SetIcon(theme.ComputerIcon())
	wd := myApp.NewWindow("配置工具GD")
	wd.CenterOnScreen()
	wd.Resize(fyne.NewSize(500, 500))
	content := makeContent()
	cv := wd.Canvas()
	cv.SetContent(content)
	//业务逻辑处理
	wd.ShowAndRun()
}

func makeContent() fyne.CanvasObject {
	descLabel := widget.NewLabel("移动运维配置工具: | 一键清除数据 | 一键清除指纹信息")
	label := widget.NewLabel("设备连接状态")
	var status = initStatus
	bindStatus := binding.BindString(&status)
	entryData := widget.NewLabelWithData(bindStatus)
	secBox := container.NewGridWithColumns(2, label, entryData)
	thirdBox := container.NewHScroll(container.NewGridWithColumns(3,
		&widget.Button{
			Alignment:  widget.ButtonAlignCenter,
			Importance: widget.MediumImportance,
			Text:       "连接设备",
			OnTapped:   func() { bindButtonConnectCaller(bindStatus) },
		},
		&widget.Button{
			Alignment:  widget.ButtonAlignCenter,
			Importance: widget.HighImportance,
			Text:       "清除指纹数据",
			OnTapped:   func() { bindButton2Test(bindStatus) },
		},
		&widget.Button{
			Alignment:  widget.ButtonAlignCenter,
			Importance: widget.DangerImportance,
			Text:       "清除设备数据",
			OnTapped:   func() { bindButton2Test(bindStatus) },
		}),
	)
	return container.NewVBox(widget.NewSeparator(), descLabel, widget.NewSeparator(), secBox, widget.NewSeparator(), thirdBox)
}

// 一键连接设备
func bindButtonConnectCaller(status binding.String) {
	st, _ := status.Get()
	if st == successConnect {
		log.Println("current device sshClient is connected")
		return
	}
	flag, msg := device.CheckConnectDevice(&machine)
	if !flag {
		log.Println(msg)
		_ = status.Set(failConnect)
	} else {
		_ = status.Set(successConnect)
	}
}

func bindButton2Test(status binding.String) {
	bashCmd := "cd / && ls -lth"
	_, err := execShell(status, bashCmd, true)
	if err != nil {
		return
	}
}

// 连接成功后执行shell
func execShell(status binding.String, bash string, isCmd bool) ([]byte, error) {
	st, _ := status.Get()
	var b []byte
	var err error
	if st == successConnect {
		flag, msg, sshClient := device.ConnectDevice(&machine)
		defer func(sshClient *ssh.Client) {
			_ = sshClient.Close()
		}(sshClient)
		if !flag {
			log.Println(msg)
			return nil, errors.New(msg)
		}
		if !isCmd {
			b, err = device.HandleBashShell(bash, sshClient)
		} else {
			b, err = device.HandleShell(bash, sshClient)
		}
	}
	return b, err
}

func init() {
	//设置中文字体:解决中文乱码问题
	fontPaths := findfont.List()
	for _, path := range fontPaths {
		if strings.Contains(path, "msyh.ttf") || strings.Contains(path, "simhei.ttf") || strings.Contains(path, "simsun.ttc") || strings.Contains(path, "simkai.ttf") {
			_ = os.Setenv("FYNE_FONT", path)
			break
		}
	}
}
