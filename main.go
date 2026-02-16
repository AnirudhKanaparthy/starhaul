package main

import (
	"fmt"
	"math"
	"os"

	"github.com/AnirudhKanaparthy/starhaul/sim"
)

type CostAction struct {
	cost    float64
	actions []sim.Action
}

var memo map[string]CostAction = make(map[string]CostAction)

func Search(simulation *sim.Simulation, visited map[string]bool) (float64, []sim.Action) {
	// This is a DFS. What I am doing here is to iterate through all the possible ways we can perform
	// actions and picking out the action which will take the least amount of cost

	currentState := simulation.State()
	if len(currentState.Tasks()) == 0 {
		return 0, []sim.Action{}
	}

	stateStr := currentState.Fingerprint()
	if v, exists := memo[stateStr]; exists {
		return v.cost, v.actions
	}

	if visited[stateStr] {
		return math.Inf(+1), nil
	}

	visited[stateStr] = true
	defer delete(visited, stateStr)

	smallestCost := math.Inf(+1)
	var smallestCostActions []sim.Action

	for _, action := range simulation.GenActions() {
		prevState := simulation.State().Copy()

		done, cost := action.Do(simulation)
		if !done {
			simulation.SetState(prevState)
			continue
		}

		searchCost, possibleActions := Search(simulation, visited)
		simulation.SetState(prevState)

		totalCost := cost + searchCost + 1.0 // 1 added to recommend smallest possible list
		if totalCost < smallestCost {
			smallestCost = totalCost

			// Build up actions
			smallestCostActions = append([]sim.Action{action}, possibleActions...)
		}
	}

	memo[stateStr] = CostAction{smallestCost, smallestCostActions}
	return smallestCost, smallestCostActions
}

func printUsage() {
	fmt.Printf("Usage: %v <filename>\n", os.Args[0])
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		return
	}

	filepath := os.Args[1]
	data, err := os.ReadFile(filepath)
	if err != nil {
		panic(err)
	}
	config, err := sim.DeserilizeJsonBytes(data)
	if err != nil {
		panic(err)
	}

	simulation, err := sim.MakeSimWithConfig(config)
	if err != nil {
		panic(err)
	}

	visited := make(map[string]bool, 0)

	estimatedCost, actions := Search(&simulation, visited)

	fmt.Println("Actions to take: ")
	actualCost := 0.0
	for step, action := range actions {
		_, cost := action.Do(&simulation)
		actualCost += cost + 1.0
		fmt.Printf("  %v. %v\n", step+1, action.Description())
	}

	fmt.Println("---")
	fmt.Printf("Lowest estimated cost: %v\n", estimatedCost)
	fmt.Printf("Actual cost          : %v\n", actualCost)
}
