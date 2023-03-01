package atlas

import (
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

type Value struct {
	Value float64 `json:"value"`
	Unit  string  `json:"unit"`
}

type Range struct {
	Min Value `json:"min"`
	Max Value `json:"max"`
}

type Storage struct {
	Default Value `json:"default"`
	Range   Range `json:"range"`
}

type Tier struct {
	Tier        string  `json:"tier"`
	Cpus        int     `json:"cpus"`
	Iops        int     `json:"iops"`
	Connections int     `json:"connections"`
	Network     Value   `json:"network"`
	Ram         Value   `json:"ram"`
	Storage     Storage `json:"storage"`
	Cost        Value   `json:"cost"`
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
