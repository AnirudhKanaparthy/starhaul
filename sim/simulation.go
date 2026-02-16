package sim

import (
	"errors"
	"fmt"
	"maps"
	"sort"
	"strings"

	"github.com/AnirudhKanaparthy/starhaul/matrix"
)

/**
 * 1. The simulation is never in a wrong state
 * 2. You can perform actions on this state
 * 3. You can uniquely represent this state
 * 5. You can calculate the cost of any actions you perform on this state
 */

type ItemSet map[int]bool

type Constants struct {
	haulerCapacity            int
	distancesBetweenLocations matrix.SymmetricMatrix[int]
	items                     []int // Volume
}

type Task struct {
	fromLocation int
	toLocation   int
}

func MakeTask(fromLocation int, toLocation int) Task {
	return Task{
		fromLocation: fromLocation,
		toLocation:   toLocation,
	}
}

type State struct {
	haulerLocation   int
	haulerItems      ItemSet
	itemsInLocations []ItemSet

	tasks map[int]Task
}

func (s State) Tasks() map[int]Task {
	return s.tasks
}

func (s State) Copy() State {
	newState := State{}
	newState.haulerLocation = s.haulerLocation

	//
	newState.haulerItems = make(ItemSet)
	maps.Copy(newState.haulerItems, s.haulerItems)

	//
	newState.itemsInLocations = make([]ItemSet, 0, len(s.itemsInLocations))
	for range s.itemsInLocations {
		newState.itemsInLocations = append(newState.itemsInLocations, make(ItemSet))
	}

	for i := range s.itemsInLocations {
		newState.itemsInLocations[i] = make(ItemSet)
		maps.Copy(newState.itemsInLocations[i], s.itemsInLocations[i])
	}

	//
	newState.tasks = make(map[int]Task)
	maps.Copy(newState.tasks, s.tasks)

	return newState
}

func IsMapEqual[T comparable, E comparable](a map[T]E, b map[T]E) bool {
	if len(a) != len(b) {
		return false
	}

	for ak, av := range a {
		bv, ok := b[ak]
		if !ok {
			return false
		}
		if bv != av {
			return false
		}
	}

	for bk, bv := range b {
		av, ok := a[bk]
		if !ok {
			return false
		}
		if av != bv {
			return false
		}
	}

	return true
}

func (a *State) IsEqualTo(b *State) bool {
	if a.haulerLocation != b.haulerLocation {
		return false
	}

	if !IsMapEqual(a.haulerItems, b.haulerItems) {
		return false
	}
	if len(a.itemsInLocations) != len(b.itemsInLocations) {
		return false
	}
	for i := range a.itemsInLocations {
		if !IsMapEqual(a.itemsInLocations[i], b.itemsInLocations[i]) {
			return false
		}
	}
	if !IsMapEqual(a.tasks, b.tasks) {
		return false
	}
	return true
}

type Simulation struct {
	constants Constants
	state     State
}

func MakeSim(haulerCapacity int,
	distancesBetweenLocations matrix.SymmetricMatrix[int],
	items []int,
	haulerStartLocation int,
	tasks map[int]Task) (Simulation, error) {

	simulation := Simulation{}
	simulation.constants = Constants{
		haulerCapacity:            haulerCapacity,
		distancesBetweenLocations: distancesBetweenLocations,
		items:                     items,
	}
	simulation.state = State{
		haulerLocation:   haulerStartLocation,
		haulerItems:      make(ItemSet),
		itemsInLocations: make([]ItemSet, 0),
		tasks:            tasks,
	}

	locationItems := &(simulation.state.itemsInLocations)
	for i := 0; i < distancesBetweenLocations.Order(); i += 1 {
		*locationItems = append(*locationItems, make(ItemSet))
	}

	for item, task := range simulation.state.tasks {
		locItems, err := simulation.LocationItems(task.fromLocation)
		if err != nil {
			return Simulation{}, err
		}
		locItems[item] = true
	}

	return simulation, nil
}

func MakeSimWithConfig(config SimulationConfig) (Simulation, error) {
	return MakeSim(config.HaulerCapacity,
		config.DistancesBetweenLocations,
		config.Items,
		config.HaulerStartLocation,
		config.Tasks,
	)
}

func (sim *Simulation) State() State {
	return sim.state
}

func (sim *Simulation) SetState(state State) {
	sim.state = state
}

func (sim *Simulation) GenActions() []Action {
	actions := make([]Action, 0)

	for i := range sim.state.itemsInLocations {
		if i == sim.HaulerLocation() {
			continue
		}
		actions = append(actions, &ActionMove{toLocation: i})
	}

	haulerItemsVolume := sim.HaulerLoad()
	for itemIndex := range sim.CurrentLocationItems() {
		if haulerItemsVolume+sim.constants.items[itemIndex] > sim.HaulerCapacity() {
			continue
		}
		actions = append(actions, &ActionTake{itemIndex: itemIndex})
	}

	for itemIndex := range sim.HaulerItems() {
		actions = append(actions, &ActionPut{itemIndex: itemIndex})
	}

	return actions
}

func (sim *Simulation) HaulerLocation() int {
	return sim.state.haulerLocation
}

func (sim *Simulation) NumberOfLocations() int {
	return len(sim.state.itemsInLocations)
}

func (sim *Simulation) NumberOfTotalItems() int {
	return len(sim.constants.items)
}

func (sim *Simulation) MoveHaulerToLocation(locIndex int) bool {
	if locIndex < 0 || sim.NumberOfLocations() <= locIndex {
		// Invalid Location
		return false
	}
	sim.state.haulerLocation = locIndex
	return true
}

func (sim *Simulation) GetDistanceBetweenLocations(a int, b int) (int, error) {
	if a < 0 || sim.NumberOfLocations() <= a {
		return -1, errors.New("Invalid location index")
	}
	if b < 0 || sim.NumberOfLocations() <= b {
		return -1, errors.New("Invalid location index")
	}
	return sim.constants.distancesBetweenLocations.At(a, b), nil
}

func (sim *Simulation) VolumeOfItem(itemIndex int) int {
	return sim.constants.items[itemIndex]
}

func (sim *Simulation) HaulerCapacity() int {
	return sim.constants.haulerCapacity
}

func (sim *Simulation) CurrentLocationItems() ItemSet {
	return sim.state.itemsInLocations[sim.state.haulerLocation]
}

func (sim *Simulation) LocationItems(locationIndex int) (ItemSet, error) {
	if locationIndex < 0 || len(sim.state.itemsInLocations) <= locationIndex {
		return nil, errors.New("Invalid Location Index")
	}
	return sim.state.itemsInLocations[locationIndex], nil
}

func (sim *Simulation) HaulerItems() ItemSet {
	return sim.state.haulerItems
}

func (sim *Simulation) HaulerLoad() int {
	vol := 0
	for itemIndex := range sim.HaulerItems() {
		vol += sim.constants.items[itemIndex]
	}
	return vol
}

func (sim *Simulation) TakeItemIntoHauler(itemIndex int) bool {
	if itemIndex < 0 || sim.NumberOfTotalItems() <= itemIndex {
		// Invalid Location
		return false
	}
	if sim.HaulerLoad()+sim.VolumeOfItem(itemIndex) > sim.HaulerCapacity() {
		return false
	}

	locationItems := sim.CurrentLocationItems()

	// Delete item from location
	if _, ok := locationItems[itemIndex]; !ok {
		return false
	}

	delete(locationItems, itemIndex)

	// Add item to hauler
	haulerItems := sim.HaulerItems()
	haulerItems[itemIndex] = true
	return true
}

func (sim *Simulation) RemoveCompletedTasks() int {
	tasksToDelete := make([]int, 0)
	for itemIndex, task := range sim.state.tasks {
		locationItems, err := sim.LocationItems(task.toLocation)
		if err != nil {
			panic("Invalid state")
		}

		_, isPresent := locationItems[itemIndex]
		if isPresent {
			delete(locationItems, itemIndex)
			tasksToDelete = append(tasksToDelete, itemIndex)
		}
	}

	for _, itemIndex := range tasksToDelete {
		delete(sim.state.tasks, itemIndex)
	}

	return len(tasksToDelete)
}

func (sim *Simulation) PutItemIntoLocation(itemIndex int) bool {
	if itemIndex < 0 || sim.NumberOfTotalItems() <= itemIndex {
		// Invalid Location
		return false
	}
	haulerItems := sim.HaulerItems()

	// Delete item from location
	if _, ok := haulerItems[itemIndex]; !ok {
		return false
	}

	delete(haulerItems, itemIndex)

	// Add item to hauler
	locationItems := sim.CurrentLocationItems()
	locationItems[itemIndex] = true

	sim.RemoveCompletedTasks()
	return true
}

type Action interface {
	Do(*Simulation) (bool, float64)
	Description() string
}

type ActionMove struct {
	toLocation int
}

func (ma *ActionMove) Do(sim *Simulation) (bool, float64) {
	startLoc := sim.HaulerLocation()
	d, _ := sim.GetDistanceBetweenLocations(startLoc, ma.toLocation)
	return sim.MoveHaulerToLocation(ma.toLocation), float64(d) + 1.0
}

func (ma *ActionMove) Description() string {
	return fmt.Sprintf("move to location %v", ma.toLocation)
}

type ActionTake struct {
	itemIndex int
}

func (ta *ActionTake) Do(sim *Simulation) (bool, float64) {
	return sim.TakeItemIntoHauler(ta.itemIndex), 0.0
}

func (ta *ActionTake) Description() string {
	return fmt.Sprintf("take item %v", ta.itemIndex)
}

type ActionPut struct {
	itemIndex int
}

func (pa *ActionPut) Do(sim *Simulation) (bool, float64) {
	return sim.PutItemIntoLocation(pa.itemIndex), 0.0
}

func (pa *ActionPut) Description() string {
	return fmt.Sprintf("put item %v", pa.itemIndex)
}

func sortedKeys(m ItemSet) string {
	keys := make([]int, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	return fmt.Sprint(keys)
}

func (state State) Fingerprint() string {
	var sb strings.Builder
	for _, loc := range state.itemsInLocations {
		sb.WriteString(sortedKeys(loc))
		sb.WriteByte(';')
	}

	return fmt.Sprintf("%v|%v|%v",
		state.haulerLocation,
		sortedKeys(state.haulerItems),
		sb.String(),
	)
}
