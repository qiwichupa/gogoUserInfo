package main

import (
	"fmt"
	"net"
	"os/exec"
	"os/user"
	"runtime"
	"strings"
	"time"

	fqdn "github.com/Showmax/go-fqdn"
	"github.com/atotto/clipboard"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/layout"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func UNUSED(x ...interface{}) {}

func GetOutboundIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	UNUSED(err)
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	ip := localAddr.IP
	return ip.String()
}

func frmt(format string, args ...interface{}) string {
	args2 := make([]string, len(args))
	for i, v := range args {
		if i%2 == 0 {
			args2[i] = fmt.Sprintf("{%v}", v)
		} else {
			args2[i] = fmt.Sprint(v)
		}
	}
	r := strings.NewReplacer(args2...)
	out := r.Replace(format)
	return out
}

func getIp(content *widget.Label) {
	ip := GetOutboundIP()
	__contentstring__ = frmt(__contentstring__, "ip", ip)
	content.SetText(__contentstring__)
	return
}

func getUsername(content *widget.Label) {
	user, err := user.Current()
	UNUSED(err)
	username := user.Username
	__contentstring__ = frmt(__contentstring__, "username", username)
	content.SetText(__contentstring__)
	return
}

func getHostname(content *widget.Label) {
	hostname := fqdn.Get()
	__contentstring__ = frmt(__contentstring__, "hostname", hostname)
	content.SetText(__contentstring__)
	return
}
func addButtonEnable(button *widget.Button) {
	if __os__ == "linux" {
		fname, err := exec.LookPath("xclip")
		UNUSED(fname)
		if err != nil {
			button.SetText("Установите xclip чтобы копировать")
		} else {
			button.SetText("Секунду...")
			time.Sleep(1500 * time.Millisecond)
			button.Enable()
			button.SetText("Скопировать")
		}
	} else {
		button.SetText("Секунду...")
		time.Sleep(1500 * time.Millisecond)
		button.Enable()
		button.SetText("Скопировать")
	}
}

var __appname__ = "gogoUserInfo"
var __version__ = "1.1.1"
var __contentstring__ string = "Username: {username}\r\nHostname: {hostname}\r\nIP: {ip}"
var __os__ string = runtime.GOOS

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow(__appname__ + " v." + __version__ + " (running on " + __os__ + ")")

	content := widget.NewLabel("")

	add := widget.NewButton("Скопировать", func() {
		clipboard.WriteAll(__contentstring__)
		content.SetText(__contentstring__ + "\n\nСкопировано в буфер обмена")
	})
	add.Disable()

	exit := widget.NewButton("Закрыть", func() {
		myWindow.Close()
	})

	myWindow.SetContent(container.NewBorder(nil, container.New(layout.NewVBoxLayout(), content, add, exit), nil, nil))
	myWindow.Resize(fyne.NewSize(400, 0))
	myWindow.SetMaster()
	myWindow.CenterOnScreen()

	myWindow.Show()

	go getIp(content)
	go getUsername(content)
	go getHostname(content)
	go addButtonEnable(add)

	myApp.Run()

}
