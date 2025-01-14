package main

import (
	"fmt"

	"github.com/cs161079/monorepo/webApplication/config"
)

func main() {
	appPtr, err := config.BuildInRuntime()
	if err != nil {
		fmt.Printf("Error Occured %v", err)
		return
	}

	appPtr.Boot()

}
