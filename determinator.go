package main

import (
	"determinator/resources"
	"fmt"
)

func main() {

	//	var cpus = 12
	//	var ram = 20
	//	var iops = 7500
	//	var storage = 3500

	var jsonResources, err = resources.Load("resources.json")

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(jsonResources)
}
