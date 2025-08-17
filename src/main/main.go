package main

import (
	"fmt"
	"obliv/src/system"
	"obliv/src/front"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
)

var Port string = ":2305"

func intro() {
	fmt.Println("obliv - control panel [BETA]")
	fmt.Println("https://github.com/radh1tya/obliv");
}

func main() {
	intro()
	if err := system.CreateFile(); err != nil {
		fmt.Printf("Failed while creating a database file: %v\n", err)
		return
	}
	dtb := system.ConnectDatabase()
	defer dtb.Close()
	system.SetupDatabase(dtb)

	r := gin.Default()

	store := cookie.NewStore([]byte("silampukau"))
	r.Use(sessions.Sessions("my-session", store))

	front.FrontSetup(r)
	r.Run(Port)
}
