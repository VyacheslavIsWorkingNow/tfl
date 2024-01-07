package oracle

type Oracle interface {
	BelongLanguage(word string) bool
	GetAlphabet() []byte
}
