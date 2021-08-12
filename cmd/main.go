package main

import (
	"os"

	"github.com/kousuketk/GogRPC2/client"
)

func main() {
	os.Exit(client.NewReversi().Run())
}
