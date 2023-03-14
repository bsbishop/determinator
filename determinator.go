package main

import (
	"determinator/atlas"
	"determinator/resources"
	"determinator/utils"
	"fmt"
	"math"
)

type Shard struct {
	Tier              string
	ShardCountCpu     float64
	ShardCountRam     float64
	ShardCountStorage float64
	ShardStorageSize  units.Value
	ShardCountIops    float64
	ShardIops         atlas.IOPS
	Cost              units.Value
}

func calculateShardsCpu(res resources.Resources, tier atlas.Tier, class string) float64 {
	var shards float64 = -1
	i := atlas.FindClassIndex(tier, class)

	// How many shards?
	if i >= 0 {
		shards = res.Cpus / (tier.Class[i].Cpus)
	}

	return math.Ceil(shards)
}

func calculateShardsRam(res resources.Resources, tier atlas.Tier) float64 {
	// How many shards?
	var shards = units.ToGB(res.Ram).Value / units.ToGB(tier.Ram).Value
	return math.Ceil(shards)
}

func calculateStorageIops(tier atlas.Tier, storageSize units.Value) atlas.IOPS {
	// (givenStorage - minIopsStorage) * (maxIops - minIops) / (maxIopsStorage - minIopsStorage) + minIops
	blocksize := tier.Iops.Blocksize
	givenStorageGB := units.ToGB(storageSize).Value
	minIopsStorageGB := units.ToGB(tier.Storage.Iops.Min.Storage).Value
	maxIopsStorageGB := units.ToGB(tier.Storage.Iops.Max.Storage).Value
	maxIops := tier.Storage.Iops.Max.Iops.Value
	minIops := tier.Storage.Iops.Min.Iops.Value

	return atlas.IOPS{Value: (givenStorageGB-minIopsStorageGB)*(maxIops-minIops)/(maxIopsStorageGB-minIopsStorageGB) + minIops, Blocksize: blocksize}
}

func calculateIopsStorage(tier atlas.Tier, iops atlas.IOPS) units.Value {
	// Storage = (Iops - minIops) * (maxIopsStorage - minIopsStorage) / (maxIops - minIops) + minIopsStorage
	givenIops := iops.Value
	minIopsStorageGB := units.ToGB(tier.Storage.Iops.Min.Storage).Value
	maxIopsStorageGB := units.ToGB(tier.Storage.Iops.Max.Storage).Value
	maxIops := tier.Storage.Iops.Max.Iops.Value
	minIops := tier.Storage.Iops.Min.Iops.Value

	return units.Value{Value: (givenIops-minIops)*(maxIopsStorageGB-minIopsStorageGB)/(maxIops-minIops) + minIopsStorageGB, Unit: "GB"}
}

func calculateShardsStorage(res resources.Resources, tier atlas.Tier) (float64, units.Value) {
	// How many shards?
	var shards = math.Ceil(units.ToGB(res.Storage).Value / units.ToGB(tier.Storage.Max).Value)

	return shards, units.Value{Value: math.Ceil(units.ToGB(res.Storage).Value / shards), Unit: "GB"}
}

func calculateShardsIops(res resources.Resources, tier atlas.Tier) (float64, atlas.IOPS) {
	var totalIops float64 = res.Iops.Read.Value + res.Iops.Write.Value
	var shards float64 = math.Ceil(totalIops / tier.Storage.Iops.Max.Iops.Value)

	return shards, atlas.IOPS{Value: math.Ceil((totalIops / shards)), Blocksize: tier.Iops.Blocksize}
}

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

	var shards []Shard

	for i := 0; i < len(jsonAtlas.Tiers); i++ {
		tier := jsonAtlas.Tiers[i]
		if tier.Skip == true {
			continue
		}
		var res = jsonResources.Resources
		shard := Shard{Tier: tier.Tier}

		// RULES:
		// CPU
		//   - Only go Low-CPU when another resource demands more shards and you end up with
		//     > 2x the amount of CPUs needed. Otherwise, use lower tier
		// RAM
		//   - Keep it simple - just how many shards
		// Storage
		//   - buying more storage is cheaper than sharding / adding a shard.

		// CPUs (General)
		shard.ShardCountCpu = atlas.HowManyShards(tier, calculateShardsCpu(res, tier, "General"))
		shard.ShardCountRam = atlas.HowManyShards(tier, calculateShardsRam(res, tier))

		tmpShards, tmpStorageSize := calculateShardsStorage(res, tier)
		shard.ShardCountStorage = atlas.HowManyShards(tier, tmpShards)
		shard.ShardStorageSize = tmpStorageSize

		var totalIops float64 = res.Iops.Read.Value + res.Iops.Write.Value
		tmpShards, tmpIopsSize := calculateShardsIops(res, tier)
		shard.ShardCountIops = atlas.HowManyShards(tier, tmpShards)
		shard.ShardIops = tmpIopsSize

		if shard.ShardCountStorage > shard.ShardCountIops {
			// more storage shards needed than IOPS shards
			// make sure that there's enough IOPS across all shards
			// We know we have enough shards and storage for storage
			// Do we need to bump up the storage to satisfy IOPS at
			// this number of shards
			iopsPerShard := atlas.IOPS{Value: totalIops / shard.ShardCountStorage, Blocksize: tier.Iops.Blocksize}
			tmpIopsStorageSize := calculateIopsStorage(tier, iopsPerShard)
			if units.ToGB(tmpIopsStorageSize).Value > units.ToGB(shard.ShardStorageSize).Value {
				shard.ShardStorageSize = tmpIopsStorageSize
			}
			shard.ShardIops = iopsPerShard
		} else if shard.ShardCountStorage < shard.ShardCountIops {
			// more IOPS shards needed than storage shards
			storagePerShard := units.Value{Value: units.ToGB(res.Storage).Value / shard.ShardCountIops, Unit: "GB"}
			iopsStoragePerShard := calculateIopsStorage(tier, shard.ShardIops)
			if units.ToGB(iopsStoragePerShard).Value > units.ToGB(storagePerShard).Value {
				shard.ShardStorageSize = iopsStoragePerShard
			}
		} else {
			// storage and IOPS shards are equal
			// calculate storage
			iopsStoragePerShard := calculateIopsStorage(tier, shard.ShardIops)
			if units.ToGB(iopsStoragePerShard).Value > units.ToGB(shard.ShardStorageSize).Value {
				shard.ShardStorageSize = iopsStoragePerShard
			}
		}

		if shard.ShardCountCpu == -1 {
			println("Tier: ", tier.Tier, " - Won't fit (CPU)")
		} else {
			println("Tier: ", tier.Tier, " - ", shard.ShardCountCpu, " shard(s). (CPU)")
		}

		if shard.ShardCountRam == -1 {
			println("Tier: ", tier.Tier, " - Won't fit (RAM)")
		} else {
			println("Tier: ", tier.Tier, " - ", shard.ShardCountRam, " shard(s). (RAM)")
		}

		if shard.ShardCountIops == -1 {
			println("Tier: ", tier.Tier, " - Won't fit (IOPS)")
		} else {
			println("Tier: ", tier.Tier, " - ", shard.ShardCountIops, " shard(s). (IOPS)")
		}

		if shard.ShardCountStorage == -1 {
			println("Tier: ", tier.Tier, " - Won't fit (storage)")
		} else {
			println("Tier: ", tier.Tier, " - ", shard.ShardCountStorage, " shard(s).")
		}

		shards = append(shards, shard)
	}
	fmt.Println(shards)

}
