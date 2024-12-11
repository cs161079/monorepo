package main

import "fmt"

func main() {
	appPtr, err := BuildInRuntime()
	if err != nil {
		fmt.Printf("Error Occured %v", err)
		return
	}

	appPtr.Boot()

}
