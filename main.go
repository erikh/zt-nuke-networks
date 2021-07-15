package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/zerotier/go-ztcentral"
)

// given a token, nuke all networks under that token.
func main() {
	controllerToken := os.Getenv("ZEROTIER_CENTRAL_TOKEN")

	if len(os.Args) > 1 {
		controllerToken = os.Args[1]
	} else {
		content, err := ioutil.ReadFile("test-token.txt")
		if err == nil {
			controllerToken = strings.TrimSpace(string(content))
		}

	}

	if controllerToken == "" {
		panic("Please supply ZEROTIER_CENTRAL_TOKEN in the environment or test-token.txt on disk with the token assigned.")
	}

	fmt.Printf("This reaps all networks for the token: %q\n", controllerToken)
	c, err := ztcentral.NewClient(controllerToken)
	if err != nil {
		panic(err)
	}

	networks, err := c.GetNetworks(context.Background())
	if err != nil {
		panic(err)
	}

	for _, network := range networks {
		fmt.Println(*network.Config.Name)
	}

	fmt.Println("These networks will be deleted. Press enter to continue, ^C to cancel")
	buf := make([]byte, 1)
	os.Stdin.Read(buf)

	for _, network := range networks {
		if err := c.DeleteNetwork(context.Background(), *network.Id); err != nil {
			fmt.Fprintf(os.Stderr, "Error deleting network: %v; barrelling forward anyway", err)
		}
	}
}
