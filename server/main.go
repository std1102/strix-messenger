package main

import (
	"fmt"
	"strix-server/persistence"
	"strix-server/router"
	"strix-server/system"
)

func main() {
	fmt.Println("  _  __     _________   __\n | | \\ \\   / /  __ \\ \\ / /\n | |  \\ \\_/ /| |  | \\ V / \n | |   \\   / | |  | |> <  \n | |____| |  | |__| / . \\ \n |______|_|  |_____/_/ \\_\\\n                          \n                          ")
	system.InitSystemConfig()
	system.InitLog()
	persistence.InitDb()
	persistence.InitBinary()
	router.Init()
}
