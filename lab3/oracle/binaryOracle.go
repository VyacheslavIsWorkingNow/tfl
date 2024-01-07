package oracle

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
)

type BinaryOracle struct {
	path     string
	alphabet []byte
}

func NewBinaryOracle(alphabet []byte, path string) *BinaryOracle {
	return &BinaryOracle{
		path:     path,
		alphabet: alphabet,
	}
}

func (bo BinaryOracle) GetAlphabet() []byte {
	return bo.alphabet
}

func (bo BinaryOracle) BelongLanguage(word string) bool {
	return bo.matchedBinary(word)
}

func (bo BinaryOracle) matchedBinary(word string) bool {
	cmd := exec.Command(bo.path, word)
	output, err := cmd.Output()

	fmt.Printf("output: '%s'\n", strings.TrimSpace(string(output)))
	fmt.Println("word:", word)

	if err != nil {
		log.Fatalf("Binary oracul have problem %e", err)
	}

	return strings.TrimSpace(string(output)) == "yes"
}
