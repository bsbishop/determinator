package main

import (
	"determinator/atlas"
	"fmt"
)

func main() {

	//	var cpus = 12
	//	var ram = 20
	//	var iops = 7500
	//	var storage = 3500

	var atlas, err = atlas.Load("atlas/atlas.json")

	if err != nil {
		fmt.Println(err)
	}

	for i := 0; i < len(atlas.Tiers); i++ {
		fmt.Println("Array:", i, atlas.Tiers[i])
	}

	// fmt.Println(atlas.Tiers[0])
}
