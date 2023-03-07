package main

import (
	"determinator/atlas"
	"determinator/resources"
	"determinator/utils"
	"fmt"
	"math"
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

	//fmt.Println(jsonResources)

	var jsonAtlas atlas.Atlas
	jsonAtlas, err = atlas.Load("atlas/atlas.json")

	if err != nil {
		fmt.Println(err)
	}

	//fmt.Println(jsonAtlas)

	for i := 0; i < len(jsonAtlas.Tiers); i++ {
		tier := jsonAtlas.Tiers[i]
		var res = jsonResources.Resources

		shardsCpus := math.Ceil(float64(res.Cpus) / float64(tier.Cpus))
		shardsRam := math.Ceil(float64(units.ToGB(res.Ram).Value) / float64(units.ToGB(tier.Ram).Value)) // GBs
		shardsIops := math.Ceil(float64(res.Iops.Reads.Value+res.Iops.Writes.Value) / float64(tier.Iops))
		shardsStorage := units.ToGB(res.Storage).Value // GBs

		if shardsCpus > float64(tier.Shards.Max) {
			println("Tier: ", tier.Tier, " - Won't fit (CPU)")
		} else {
			println("Tier: ", tier.Tier, " - ", shardsCpus, " shard(s). (CPU)")
		}

		if shardsRam > float64(tier.Shards.Max) {
			println("Tier: ", tier.Tier, " - Won't fit (RAM)")
		} else {
			println("Tier: ", tier.Tier, " - ", shardsCpus, " shard(s). (RAM)")
		}

		if shardsIops > float64(tier.Shards.Max) {
			println("Tier: ", tier.Tier, " - Won't fit (IOPS)")
		} else {
			println("Tier: ", tier.Tier, " - ", shardsCpus, " shard(s). (IOPS)")
		}

		if shardsStorage > float64(tier.Shards.Max) {
			println("Tier: ", tier.Tier, " - Won't fit (storage)")
		} else {
			println("Tier: ", tier.Tier, " - ", shardsCpus, " shard(s).")
		}
	}

}
