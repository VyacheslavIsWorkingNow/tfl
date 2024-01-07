package automaton

import (
	"log"
	"sort"
)

type State string
type StateTransitions map[byte]State

type Machine struct {
	StartState  State
	Transitions map[State]StateTransitions
	FinalStates map[State]struct{}
	Alphabet    []byte
}

func New(alphabet []byte) *Machine {
	return &Machine{
		StartState:  "",
		Transitions: make(map[State]StateTransitions),
		FinalStates: make(map[State]struct{}),
		Alphabet:    alphabet,
	}
}

func (m *Machine) Translate(suffix, prefix, extendPrefix map[string]struct{}, table, extendTable map[string]bool) {
	log.Println("START TRANSLATE")

	states, prefixToRow, rowToPrefix := m.getStates(suffix, prefix, extendPrefix, table, extendTable)

	m.setStartState(prefixToRow)

	m.setTransitions(states, prefixToRow, rowToPrefix)

	m.setFinalStates(prefixToRow, table)

	log.Println(m.StartState)
	log.Println(m.Transitions)
	log.Println(m.FinalStates)

	log.Println("END TRANSLATE")
}

func (m *Machine) setStartState(prefixToRow map[string]State) {
	m.StartState = prefixToRow[""]
}

func (m *Machine) getStates(
	suffix, prefix, extendPrefix map[string]struct{}, table, extendTable map[string]bool,
) (map[State]struct{}, map[string]State, map[State]string) {

	states := make(map[State]struct{})
	suffixSort := getSortArrayFromSet(suffix)
	prefixSort := getSortArrayFromSet(prefix)

	prefixToRow := make(map[string]State)
	rowToPrefix := make(map[State]string)

	for _, p := range prefixSort {
		newRow := State(getTableRow(p, suffixSort, table))
		if _, ok := rowToPrefix[newRow]; !ok {
			rowToPrefix[newRow] = p
			prefixToRow[p] = newRow
			states[newRow] = struct{}{}
		}
	}

	for ep := range extendPrefix {
		if _, ok := prefixToRow[ep]; !ok {
			prefixToRow[ep] = State(getTableRow(ep, suffixSort, extendTable))
		}
	}

	return states, prefixToRow, rowToPrefix
}

func getSortArrayFromSet(set map[string]struct{}) []string {
	sortArr := make([]string, 0, len(set))

	for s := range set {
		sortArr = append(sortArr, s)
	}

	sort.Sort(ByLengthThenAlphabetical(sortArr))

	return sortArr
}

type ByLengthThenAlphabetical []string

func (s ByLengthThenAlphabetical) Len() int {
	return len(s)
}

func (s ByLengthThenAlphabetical) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s ByLengthThenAlphabetical) Less(i, j int) bool {
	if len(s[i]) == len(s[j]) {
		return s[i] < s[j]
	}
	return len(s[i]) < len(s[j])
}

func getTableRow(prefix string, suffixList []string, table map[string]bool) string {
	row := ""
	for _, s := range suffixList {
		val := table[createTableKey(prefix, s)]
		if val {
			row += "1"
		} else {
			row += "0"
		}
	}
	return row
}

func createTableKey(prefix, suffix string) string {
	return prefix + "_" + suffix
}

func (m *Machine) setTransitions(
	states map[State]struct{}, prefixToRow map[string]State, rowToPrefix map[State]string) {
	for state := range states {
		statePrefix := rowToPrefix[state]
		for _, letter := range m.Alphabet {
			nextPrefixTransition := getWord(statePrefix, letter)
			nextState := prefixToRow[nextPrefixTransition]
			m.addTransition(state, nextState, letter)
		}
	}
}

func getWord(prefix string, letter byte) string {
	return prefix + string(letter)
}

func (m *Machine) setFinalStates(prefixToRow map[string]State, table map[string]bool) {
	for prefix, state := range prefixToRow {
		if val := table[createTableKey(prefix, "")]; val {
			m.addFinalState(state)
		}
	}
}

func (m *Machine) addFinalState(s State) {
	m.FinalStates[s] = struct{}{}
}

func (m *Machine) addTransition(from, to State, letter byte) {
	if to == "" {
		return
	}
	if _, ok := m.Transitions[from]; ok {
		m.Transitions[from][letter] = to
	} else {
		m.Transitions[from] = make(StateTransitions)
		m.Transitions[from][letter] = to
	}
}

func (m *Machine) Membership(word string) bool {
	currentState := m.StartState

	for _, symbol := range word {
		nextState, exists := m.Transitions[currentState][byte(symbol)]
		if !exists {
			return false
		}
		currentState = nextState
	}

	_, isFinal := m.FinalStates[currentState]
	return isFinal
}
