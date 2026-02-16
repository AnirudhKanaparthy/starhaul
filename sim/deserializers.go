package sim

import (
	"encoding/json"

	"github.com/AnirudhKanaparthy/starhaul/matrix"
)

type SimulationConfig struct {
	HaulerCapacity            int
	HaulerStartLocation       int
	DistancesBetweenLocations matrix.SymmetricMatrix[int]
	Items                     []int
	Tasks                     map[int]Task
}

type JsonItem struct {
	Volume int
}

type JsonTask struct {
	From  string
	To    string
	Items []JsonItem
}

type JsonConfig struct {
	HaulerCapacity      int
	Locations           []string
	HaulerStartLocation string
	Distances           map[string]map[string]int
	Tasks               []JsonTask
}

type JsonDeserializer struct{}

func DeserilizeJsonString(text string) (SimulationConfig, error) {
	return DeserilizeJsonBytes([]byte(text))
}

func DeserilizeJsonBytes(data []byte) (SimulationConfig, error) {
	config := JsonConfig{}
	err := json.Unmarshal(data, &config)
	if err != nil {
		return SimulationConfig{}, err
	}

	// Location Indices
	locIndex := make(map[string]int)
	for i, loc := range config.Locations {
		locIndex[loc] = i
	}

	// Distances
	distances := matrix.MakeSymmWithOrder[int](len(config.Locations))
	for loc, dmap := range config.Distances {
		from := locIndex[loc]

		for l, d := range dmap {
			to := locIndex[l]

			distances.Set(from, to, d)
		}
	}

	// Iems and Tasks
	items := make([]int, 0)
	tasks := make(map[int]Task)
	for _, task := range config.Tasks {
		for _, item := range task.Items {
			items = append(items, item.Volume)

			tasks[len(items)-1] = Task{
				fromLocation: locIndex[task.From],
				toLocation:   locIndex[task.To],
			}
		}
	}

	return SimulationConfig{
		HaulerCapacity:            config.HaulerCapacity,
		HaulerStartLocation:       locIndex[config.HaulerStartLocation],
		DistancesBetweenLocations: distances,
		Items:                     items,
		Tasks:                     tasks,
	}, nil
}
