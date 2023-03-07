package atlas

import (
	"determinator/utils"
	"encoding/json"
	"errors"
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

type Range struct {
	Min units.Value `json:"min"`
	Max units.Value `json:"max"`
}

type Storage struct {
	Default units.Value `json:"default"`
	Range   Range       `json:"range"`
}

type Shards struct {
	Min  int         `json:"min"`
	Max  int         `json:"max"`
	Cost units.Value `json:"cost"`
}

type Tier struct {
	Tier        string      `json:"tier"`
	Cpus        int         `json:"cpus"`
	Iops        int         `json:"iops"`
	Connections int         `json:"connections"`
	Network     units.Value `json:"network"`
	Ram         units.Value `json:"ram"`
	Storage     Storage     `json:"storage"`
	Cost        units.Value `json:"cost"`
	Shards      Shards      `json:"shards"`
}

type Tiers struct {
	Tiers []Tier `json:"tiers"`
}

type Atlas struct {
	Regions []Region `json:"regions"`
	Tiers   []Tier   `json:"tiers"`
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
