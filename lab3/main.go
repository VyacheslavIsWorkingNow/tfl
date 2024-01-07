package main

import (
	"fmt"
	"github.com/VyacheslavIsWorkingNow/tfl/lab3/oracle"
	"github.com/VyacheslavIsWorkingNow/tfl/lab3/tables"
	"log"
)

func main() {

	fmt.Println("START")

	fmt.Println("Введите максимальную длину слов для тестирования MAT:")
	fmt.Println("Например, длина 3 означает, что MAT ограничится всеми перестановками длины 1, 2, 3")
	var lenMaxMATWord int
	_, _ = fmt.Scanln(&lenMaxMATWord)

	fmt.Println("Введите путь до бинарного файла, который отвечает за оракул:")
	var oraclePath string
	_, _ = fmt.Scanln(&oraclePath)

	fmt.Println("Введите длину алфавита и алфавит, который использует оракул," +
		" через enter - буква алфавита - byte:")

	var alphabetCount int
	_, _ = fmt.Scanln(&alphabetCount)

	alphabet := make([]byte, alphabetCount)
	for i := 0; i < alphabetCount; i++ {
		var letter string
		fmt.Println("letter:", letter)
		_, _ = fmt.Scanln(&letter)
		alphabet[i] = []byte(letter)[0]
	}

	binOracle := oracle.NewBinaryOracle(alphabet, oraclePath)

	fmt.Println("Получен оракул:")
	fmt.Println("lenMaxMATWord:", lenMaxMATWord, "\noraclePath:", oraclePath, "\nalphabet:", alphabet)

	angluin := tables.New(binOracle, lenMaxMATWord)
	auto := angluin.Run()

	err := auto.GetDotMachine("mainOracle")
	if err != nil {
		log.Fatalf("failed get dot machine %e\n", err)
	}

	fmt.Println("SUCCESS")

}

// INFO: для визуализации можно установить dot и командой dot -Tpng machine.dot -o machine.png в папке с
// dot-файлом получить изображение
