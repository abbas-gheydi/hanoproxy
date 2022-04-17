package main

import (
	"fmt"

	"github.com/Abbas-gheydi/hanoproxy/checker/updater"
	"github.com/Abbas-gheydi/hanoproxy/confs"
	"github.com/Abbas-gheydi/hanoproxy/dns/server"
)

func main() {
	banner()
	confs.Read()
	updater.FillRecordTables()
	go updater.HealthCheck()
	server.Start()

}

func banner() {
	fmt.Print(`
	_    _                      _____             __     __
	| |  | |   /\               |  __ \            \ \   / /
	| |__| |  /  \   _ __   ___ | |__) | __ _____  _\ \_/ / 
	|  __  | / /\ \ | '_ \ / _ \|  ___/ '__/ _ \ \/ /\   /  
	| |  | |/ ____ \| | | | (_) | |   | | | (_) >  <  | |   
	|_|  |_/_/    \_\_| |_|\___/|_|   |_|  \___/_/\_\ |_|   
															
															

	
`)
}
