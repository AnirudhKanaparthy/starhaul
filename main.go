package main

import (
	"fmt"
	"math"

	"github.com/AnirudhKanaparthy/starhaul/matrix"
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

func main() {
	simulation, err := sim.MakeSimulation(
		4,
		matrix.MakeSymmetricMatrix([]int{0, 2, 2, 0, 3, 0}),
		[]int{2, 4, 1},
		0,
		map[int]sim.Task{
			0: sim.MakeTask(0, 1),
			1: sim.MakeTask(0, 2),
			2: sim.MakeTask(1, 2),
		},
	)
	if err != nil {
		panic(err)
	}
	_ = simulation

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
