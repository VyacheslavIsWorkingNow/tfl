package mat

import (
	"log"
	"sort"

	"github.com/VyacheslavIsWorkingNow/tfl/lab3/automaton"
	"github.com/VyacheslavIsWorkingNow/tfl/lab3/oracle"
)

// Генератор получился отличный!

// MAT INFO: сгенерировать слова - трудоемкий процесс, при этом это всегда один набор слов (я решил так сделать),
// так что генерирую я слова один раз при инициализации - поэтому в начале может немного просесть программа
type MAT struct {
	oracle      oracle.Oracle
	maxLenWords int
	words       []string
}

func New(o oracle.Oracle, maxLenWords int) *MAT {
	log.Println("Generating words success")
	mat := &MAT{
		oracle:      o,
		maxLenWords: maxLenWords,
	}

	log.Println("Generating words...")
	mat.words = mat.GenerateWords()
	sort.Sort(automaton.ByLengthThenAlphabetical(mat.words))
	log.Println("Generating words success")
	mat.words = mat.GenerateWords()

	return mat
}

func (mat *MAT) GenerateWords() []string {

	words := make([]string, 0)

	for i := 1; i <= mat.maxLenWords; i++ {
		words = append(words, generateCombinations(mat.oracle.GetAlphabet(), i)...)
	}

	log.Printf("Generated %d words\n", len(words))

	return words
}

func generateCombinations(alphabet []byte, length int) []string {
	var resultByte [][]byte
	generateCombination(alphabet, length, []byte{}, &resultByte)

	resultString := make([]string, len(resultByte))

	for i := 0; i < len(resultByte); i++ {
		resultString[i] = string(resultByte[i])
	}

	return resultString
}

func generateCombination(alphabet []byte, remaining int, current []byte, result *[][]byte) {
	if remaining == 0 {
		*result = append(*result, append([]byte{}, current...))
		return
	}

	for _, char := range alphabet {
		generateCombination(alphabet, remaining-1, append(current, char), result)
	}
}

func (mat *MAT) Membership(word string) bool {
	return mat.oracle.BelongLanguage(word)
}

func (mat *MAT) Equivalence(m *automaton.Machine) (bool, string) {

	for _, word := range mat.words {
		if m.Membership(word) != mat.Membership(word) {
			log.Printf("word %s isn't equal\n", word)
			log.Printf("machine: %t | mat %t: \n", m.Membership(word), mat.Membership(word))
			return false, word
		}
	}

	return true, ""
}

func (mat *MAT) GetAlphabet() []byte {
	return mat.oracle.GetAlphabet()
}
