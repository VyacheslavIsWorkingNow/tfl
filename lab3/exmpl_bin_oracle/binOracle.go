package main

import (
	"fmt"
	"os"
)

// INFO: собираем файл командой go build <имя файла>

func main() {
	if len(os.Args) < 2 {
		return
	}

	word := os.Args[1]
	result := processWord(word)
	fmt.Println(result)
}

func processWord(word string) string {

	if len(word) < 2 || word[len(word)-2] != 'b' {
		return "no"
	}
	return "yes"
}
