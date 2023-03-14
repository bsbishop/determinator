package resources

import (
	"determinator/utils"
	"encoding/json"
	"errors"
	"os"
)

type IOPS struct {
	Value     float64 `json:"value"`
	Blocksize float64 `json:"blocksize"`
}

type RWIops struct {
	Read  IOPS `json:"read"`
	Write IOPS `json:"write"`
}

type Resources struct {
	Cpus    float64     `json:"cpus"`
	Ram     units.Value `json:"ram"`
	Storage units.Value `json:"storage"`
	Iops    RWIops      `json:"iops"`
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
		return JSONResources{}, errors.New("unable to load file")
	}

	var jsonResources JSONResources
	err = json.Unmarshal(bytes, &jsonResources)

	if err != nil {
		return JSONResources{}, errors.New("JSON decode error")
	}

	return jsonResources, nil
}
