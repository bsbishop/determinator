package resources

import (
	"determinator/utils"
	"encoding/json"
	"errors"
	"os"
)

type IOPSrw struct {
	Reads  IOPS.IOPS `json:"reads"`
	Writes IOPS.IOPS `json:"writes"`
}

type Resources struct {
	Cpus    int         `json:"cpus"`
	Ram     units.Value `json:"ram"`
	Storage units.Value `json:"storage"`
	Iops    IOPSrw      `json:"iops"`
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
