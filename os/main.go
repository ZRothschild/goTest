package main

import (
	"fmt"
	"log"
	"os/exec"
)

func main() {
	out, err := exec.Command("ipconfig", "/all").Output()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%v", string(out))
}
