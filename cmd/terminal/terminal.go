package terminal

import (
	"fmt"

	"github.com/williampring/media-transfer/config"
	"github.com/williampring/media-transfer/pkg/sftpclient"
)

// set up project with parameters
func Start(argsWithoutProg []string) {
	config, _ := config.GetConfig()
	fmt.Println(config)
	sftpclient.Start(argsWithoutProg[0], argsWithoutProg[1], config)
	fmt.Println("Hello")
}
