package resources

import (
	"encoding/json"
	"errors"
	"os"
)

type Value struct {
	Value float64 `json:"value"`
	Unit  string  `json:"unit"`
}

type IOPS struct {
	Value     int `json:"value"`
	Blocksize int `json:"blocksize"`
}

type IOPSrw struct {
	Reads  IOPS `json:"reads"`
	Writes IOPS `json:"writes"`
}

type Resources struct {
	CPUs    int    `json:"cpus"`
	RAM     Value  `json:"ram"`
	Storage Value  `json:"storage"`
	IOPS    IOPSrw `json:"iops"`
}

type JSONResources struct {
	Organization string    `json:"organization"`
	Project      string    `json:"project"`
	Cluster      string    `json:"cluster"`
	Timestamp    string    `json:"timestamp"`
	Resources    Resources `json:"resources"`
}

func Load(fn string) (JSONResources, error) {
	bytes, err := os.ReadFile(fn)

	if err != nil {
		return JSONResources{}, errors.New("Unable to load file")
	}

	var jsonResources JSONResources
	err = json.Unmarshal(bytes, &jsonResources)

	if err != nil {
		return JSONResources{}, errors.New("JSON decode error")
	}

	return jsonResources, nil
}
