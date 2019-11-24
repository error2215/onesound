package main

import (
	"onesound/server/api"
	_ "onesound/server/config"
)

func main() {
	api.Start()
}
