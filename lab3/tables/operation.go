package tables

import (
	"fmt"
	"sort"
)

func (a *Angluin) AddSuffix(suffix string) {
	a.suffix[suffix] = struct{}{}

	a.updateTableBySuffix(suffix)
	a.updateExtendTableBySuffix(suffix)
}

func (a *Angluin) AddPrefix(prefix string) {
	for i := 1; i < len(prefix)+1; i++ {
		a.prefix[prefix[:i]] = struct{}{}
		a.updateTableByPrefix(prefix[:i])

		a.addExtendPrefix(prefix[:i])
	}

}

func (a *Angluin) addExtendPrefix(prefix string) {
	for _, letter := range a.mat.GetAlphabet() {
		newPrefix := prefix + string(letter)

		a.extendPrefix[newPrefix] = struct{}{}
		a.deleteAllPrefixesPrefixFromExtendPrefix(newPrefix)

		a.updateExtendTableByAddPrefix(newPrefix)
	}
}

func (a *Angluin) deleteAllPrefixesPrefixFromExtendPrefix(extendPrefix string) {
	for i := 1; i < len(extendPrefix); i++ {
		delete(a.extendPrefix, extendPrefix[:i])
		a.updateExtendTableByDeletePrefix(extendPrefix[:i])
	}
}

func (a *Angluin) PrintPrefix() {
	fmt.Printf("PREFIX\n")
	for p := range a.prefix {
		fmt.Printf("'%s' ", p)
	}
	fmt.Printf("\n")
}

func (a *Angluin) PrintExtendPrefix() {
	fmt.Printf("EXTEND PREFIX\n")
	for ep := range a.extendPrefix {
		fmt.Printf("'%s' ", ep)
	}
	fmt.Printf("\n")
}

func (a *Angluin) PrintSuffix() {
	fmt.Printf("SUFFIX\n")
	for s := range a.suffix {
		fmt.Printf("'%s' ", s)
	}
	fmt.Printf("\n")
}

// INFO: Как я храню таблицу? У меня есть два unordered set: префиксы и суффиксы, соответственно я сделаю map с ключами
// prefix_suffix, в которой будет храниться входит ли это слово в язык.
// Чтобы получить строку - надо пройтись по всем suffix от 1 до N и получим: prefix_suffix1...prefix_suffixN
// Столбец получать не нужно
// Так как множество не упорядочено, то и таблица не упорядочена, но ключи уникальны
// Во время вывода таблицы я буду сортировать оба множества для удобства отладки

func createTableKey(prefix, suffix string) string {
	return prefix + "_" + suffix
}

func createWord(prefix, suffix string) string {
	return prefix + suffix
}

func (a *Angluin) updateTableByPrefix(prefix string) {
	for s := range a.suffix {
		a.table[createTableKey(prefix, s)] = a.mat.Membership(createWord(prefix, s))
	}
}

func (a *Angluin) updateTableBySuffix(suffix string) {
	for p := range a.prefix {
		a.table[createTableKey(p, suffix)] = a.mat.Membership(createWord(p, suffix))
	}
}

func (a *Angluin) updateExtendTableByAddPrefix(prefix string) {
	for s := range a.suffix {
		a.extendTable[createTableKey(prefix, s)] = a.mat.Membership(createWord(prefix, s))
	}
}

func (a *Angluin) updateExtendTableByDeletePrefix(prefix string) {
	for s := range a.suffix {
		delete(a.extendTable, createTableKey(prefix, s))
	}
}

func (a *Angluin) updateExtendTableBySuffix(suffix string) {
	for p := range a.extendPrefix {
		a.extendTable[createTableKey(p, suffix)] = a.mat.Membership(createWord(p, suffix))
	}
}

func sortSet(set map[string]struct{}) []string {
	list := make([]string, 0, len(set))
	for s := range set {
		list = append(list, s)
	}
	sort.Strings(list)
	return list
}

func sortSets(prefixSet, suffixSet map[string]struct{}) ([]string, []string) {

	suffixList := sortSet(suffixSet)
	prefixList := sortSet(prefixSet)

	return prefixList, suffixList
}

func printFromLists(prefixList, suffixList []string, table map[string]bool) {
	fmt.Printf("%-10s", "lambda")

	for _, s := range suffixList {
		fmt.Printf("%-10s", s)
	}
	fmt.Printf("\n")

	for _, p := range prefixList {
		fmt.Printf("%-10s", p)
		for _, s := range suffixList {
			key := createTableKey(p, s)
			val := table[key]
			if val {
				fmt.Printf("%-10v", 1)
			} else {
				fmt.Printf("%-10v", 0)
			}
		}
		fmt.Printf("\n")
	}
}

func (a *Angluin) PrintTable() {
	fmt.Printf("TABLE\n")

	prefixList, suffixList := sortSets(a.prefix, a.suffix)

	printFromLists(prefixList, suffixList, a.table)
}

func (a *Angluin) PrintExtendTable() {
	fmt.Printf("EXTEND TABLE\n")

	prefixList, suffixList := sortSets(a.extendPrefix, a.suffix)

	printFromLists(prefixList, suffixList, a.extendTable)
}
