package main

import (
	"fmt"
	"os/exec"
	"os"
	"runtime"
	"obliv/src/system"
	"obliv/src/front"
	"io"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
)

var Port string = ":2305"

func intro() {
	fmt.Println("obliv - control panel [BETA]")
	fmt.Println("https://github.com/radh1tya/obliv");
}

func runCmd(name string, arg ...string) {
    cmd := exec.Command(name, arg...)
    cmd.Stdout = os.Stdout
    cmd.Run()
}

func ClearTerminal() {
    switch runtime.GOOS {
    case "darwin":
        runCmd("clear")
    case "linux":
        runCmd("clear")
    case "windows":
        runCmd("cmd", "/c", "cls")
    default:
        runCmd("clear")
    }
}

func main() {
	ClearTerminal()
	intro()
	if err := system.CreateFile(); err != nil {
		fmt.Printf("Failed while creating a database file: %v\n", err)
		return
	}
	dtb := system.ConnectDatabase()
	defer dtb.Close()
	system.SetupDatabase(dtb)

	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = io.Discard
	r := gin.Default()
	store := cookie.NewStore([]byte("silampukau"))
	r.Use(sessions.Sessions("my-session", store))

	front.FrontSetup(r)
	go system.Shell()
	r.Run(Port)

}
