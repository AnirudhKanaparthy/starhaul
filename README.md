# Star Haul
It is an algorithm that calculates the most efficient way to complete a given set of tasks/contracts in Star Citizen. As of now it optimizes for distance and the amount of times you have to land on a station.

## Quick Start
```console
go build main.go
./main config.json
```

You can change the config as per your contracts. Just keep in mind the following:
1. The names in locations are basically IDs, you have to use them throughout the config to point to a specific location.
2. Make sure the capacity of the hauler is at least as big as the biggest item in any all tasks. If a task has an item with volume 8SCUs but your Hauler can only hold 4SCUs then there is no solution or actions for this scenario.

## Why?
I created this to maximize my productivity in Star Citizen hauling. I realised early on that I can optimize my routes by doing more than one contract at a time.

## Future work:
1. Proper CLI which takes in a yaml file as a description
2. Have an internal database of all the locations in Star Citizen and the distances between them.
3. Optimize for other parameters like, capital expenditures, danger of a trip etc