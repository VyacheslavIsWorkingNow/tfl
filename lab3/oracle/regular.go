package oracle

import (
	"log"
	"regexp"
)

type RegularOracle struct {
	regex    string
	alphabet []byte
}

func NewRegularOracle(regex string, alphabet []byte) *RegularOracle {
	return &RegularOracle{
		regex:    regex,
		alphabet: alphabet,
	}
}

func (ro RegularOracle) BelongLanguage(word string) bool {
	return matched(word, ro.regex)
}

func matched(word, regex string) bool {
	match, err := regexp.MatchString("^"+regex+"$", word)
	if err != nil {
		log.Fatalf("can't match word %s in regex %s in oracle", word, regex)
	}
	return match
}

func (ro RegularOracle) GetAlphabet() []byte {
	return ro.alphabet
}
