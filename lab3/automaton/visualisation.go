package automaton

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	dotFileName = "oracleName/machine.dot"
)

func (m *Machine) GetDotMachine(oracleName string) error {
	dot := m.toDOT()

	newDotFileName := strings.Replace(dotFileName, "oracleName", oracleName, 1)

	dir := filepath.Dir(newDotFileName)

	mErr := os.MkdirAll(dir, os.ModePerm)
	if mErr != nil {
		return fmt.Errorf("Ошибка при создании директорий %w\n", mErr)
	}

	dotFile, cErr := os.Create(newDotFileName)
	if cErr != nil {
		return fmt.Errorf("Ошибка создания DOT-файла: %w\n", cErr)
	}
	defer func() {
		_ = dotFile.Close()
	}()

	_, wErr := dotFile.WriteString(dot)
	if wErr != nil {
		return wErr
	}

	fmt.Println("файл автомата успешно записан")

	return nil
}

func (m *Machine) toDOT() string {

	dot := "digraph Automaton {\n"
	dot += "\trankdir=LR;\n"

	dot += fmt.Sprintf("\t%d -> %s;\n", -1, m.StartState)
	dot += fmt.Sprintf("\t%d [shape=\"point\"];\n", -1)

	for state := range m.FinalStates {
		dot += fmt.Sprintf("\t%s [shape=\"doublecircle\"];\n", state)
	}

	for fromState, transitions := range m.Transitions {
		for symbol, toState := range transitions {
			dot += fmt.Sprintf("\t%s -> %s [label=\"%s\"];\n", fromState, toState, string(symbol))
		}
	}

	dot += "}"
	return dot
}
