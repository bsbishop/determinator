package atlas

import (
	"determinator/utils"
	"encoding/json"
	"errors"
	"math"
	"os"
)

type Region struct {
	Name        string  `json:"name"`
	Region      string  `json:"region"`
	Recommended bool    `json:"recommended"`
	Cost        float64 `json:"cost"`
}

type Regions struct {
	Regions []Region `json:"regions"`
}

type StorageCost struct {
	Min struct {
		Size units.Value `json:"size"`
		Cost units.Value `json:"cost"`
	} `json:"min"`
	Max struct {
		Size units.Value `json:"size"`
		Cost units.Value `json:"cost"`
	} `json:"max"`
}

type StorageIops struct {
	Storage units.Value `json:"storage"`
	Iops    IOPS        `json:"iops"`
}

type StorageIopsMinMax struct {
	Min StorageIops `json:"min"`
	Max StorageIops `json:"max"`
}

type Storage struct {
	Default units.Value       `json:"default"`
	Min     units.Value       `json:"min"`
	Max     units.Value       `json:"max"`
	Cost    StorageCost       `json:"cost"`
	Iops    StorageIopsMinMax `json:"iops"`
}

type IOPS struct {
	Value     float64 `json:"value"`
	Blocksize float64 `json:"blocksize"`
}

type Iops struct {
	Default IOPS        `json:"default"`
	Min     IOPS        `json:"min"`
	Max     IOPS        `json:"max"`
	Cost    units.Value `json:"cost"`
}

type ShardCost struct {
	Count float64     `json:"count"`
	Cost  units.Value `json:"cost"`
}

type Shards struct {
	Min ShardCost `json:"smin"`
	Max ShardCost `json:"smax"`
}

type RWIops struct {
	Read  IOPS `json:"read"`
	Write IOPS `json:"write"`
}

type Class struct {
	Name    string      `json:"name"`
	Cpus    float64     `json:"cpus"`
	Cost    units.Value `json:"cost"`
	Storage units.Value `json:"storage"`
	RWIops  RWIops      `json:"iops"`
}

type Tier struct {
	Tier        string      `json:"tier"`
	Skip        bool        `json:"skip"`
	Class       []Class     `json:"class"`
	Connections float64     `json:"connections"`
	Network     units.Value `json:"network"`
	Ram         units.Value `json:"ram"`
	Storage     Storage     `json:"storage"`
	Iops        IOPS        `json:"iops"`
	Shards      Shards      `json:"shards"`
}

type Tiers struct {
	Tiers []Tier `json:"tiers"`
}

type Atlas struct {
	Regions []Region `json:"regions"`
	Tiers   []Tier   `json:"tiers"`
}

func HowManyShards(tier Tier, shards float64) float64 {

	// If you have 0 to 1 shards needed then: replica set
	if shards >= 0 && shards <= 1 {
		shards = 0
	}

	// Too many shards! Can't be done (or contact Atlas team)
	if shards >= MaxShards(tier) {
		shards = -1
	}
	return math.Ceil(shards)
}

func MaxShards(tier Tier) float64 {
	var shards float64 = 0
	if (tier.Shards != Shards{}) {
		shards = tier.Shards.Max.Count
	}
	return shards
}

func FindClassIndex(tier Tier, class string) int {
	var result = -1
	for i := 0; i < len(tier.Class); i++ {
		if tier.Class[i].Name == class {
			result = i
			break
		}
	}
	return result
}

func Load(fn string) (Atlas, error) {
	bytes, err := os.ReadFile(fn)

	if err != nil {
		return Atlas{}, errors.New("Unable to load atlas.json")
	}

	var atlas Atlas
	err = json.Unmarshal(bytes, &atlas)

	if err != nil {
		return Atlas{}, errors.New("JSON decode error")
	}

	return atlas, nil
}
