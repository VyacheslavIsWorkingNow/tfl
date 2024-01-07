package tables

import (
	"context"
	"log"
	"time"

	"github.com/VyacheslavIsWorkingNow/tfl/lab3/automaton"
	"github.com/VyacheslavIsWorkingNow/tfl/lab3/mat"
	"github.com/VyacheslavIsWorkingNow/tfl/lab3/oracle"
)

type Angluin struct {
	suffix       map[string]struct{}
	prefix       map[string]struct{}
	extendPrefix map[string]struct{}
	table        map[string]bool
	extendTable  map[string]bool
	mat          mat.MAT
}

func New(o oracle.Oracle, maxLenWords int) *Angluin {
	a := Angluin{
		suffix:       make(map[string]struct{}),
		prefix:       make(map[string]struct{}),
		extendPrefix: make(map[string]struct{}),
		table:        make(map[string]bool),
		extendTable:  make(map[string]bool),
		mat:          *mat.New(o, maxLenWords),
	}

	a.table["_"] = a.mat.Membership("")

	a.prefix[""] = struct{}{}
	a.suffix[""] = struct{}{}

	a.addExtendPrefix("")

	return &a
}

func (a *Angluin) Run() *automaton.Machine {
	log.Println("start RUN")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			log.Println("Execution time exceeded")
			return nil
		default:
			a.PrintPrefix()
			a.PrintSuffix()
			a.PrintExtendPrefix()
			a.PrintTable()
			a.PrintExtendTable()
			var closed, consistent string
			closed = a.Closed()
			consistent = a.Consistent()
			if closed == "" && consistent == "" {
				m := automaton.New(a.mat.GetAlphabet())
				m.Translate(a.suffix, a.prefix, a.extendPrefix, a.table, a.extendTable)
				ok, newPrefix := a.mat.Equivalence(m)
				if ok {
					return m
				}
				a.AddPrefix(newPrefix)
			} else if closed != "" {
				log.Println("NOT CLOSED, ADD", closed)
				a.AddPrefix(closed)
			} else {
				log.Println("NOT CONSISTENT, ADD", consistent)
				a.AddSuffix(consistent)
			}
		}
	}
}

// Closed INFO: Closed: An observation table is called closed if for all t in S.A there exist an s` in S
// such that row(s’)=row(t).This states that every row(s.a) must be present in row(s).
// Если полна, то вернется пустая строка, если нет, то префикс
func (a *Angluin) Closed() string {

	suffixList := sortSet(a.suffix)
	tempTableRowMap := a.createTempTableRowMap(suffixList)
	extendPrefixList := sortSet(a.extendPrefix)
	for _, ep := range extendPrefixList {
		if _, ok := tempTableRowMap[a.getExtendTableRow(ep, suffixList)]; !ok {
			return ep
		}
	}

	return ""
}

func (a *Angluin) createTempTableRowMap(suffixList []string) map[string]struct{} {
	tempTableRowMap := make(map[string]struct{})

	for p := range a.prefix {
		tempTableRowMap[a.getTableRow(p, suffixList)] = struct{}{}
	}

	return tempTableRowMap
}

func (a *Angluin) getTableRow(prefix string, suffixList []string) string {
	row := ""
	for _, s := range suffixList {
		val := a.table[createTableKey(prefix, s)]
		if val {
			row += "1"
		} else {
			row += "0"
		}
	}
	return row
}

func (a *Angluin) getExtendTableRow(prefix string, suffixList []string) string {
	row := ""
	for _, s := range suffixList {
		val := a.extendTable[createTableKey(prefix, s)]
		if val {
			row += "1"
		} else {
			row += "0"
		}
	}
	return row
}

// Consistent INFO: Consistent: An observation table is said to be consistent if, whenever s1,s2 in S satisfy row(s1)=row(s2)
// then for every an in A must satisfy row(s1.a)=row(s2.a).
// Если консистентно, то вернется пустая строка, иначе - буква + суффикс
func (a *Angluin) Consistent() string {

	suffixList := sortSet(a.suffix)

	tablePrefixToRow := a.getDsForTablePrefixToRow(suffixList)
	extendTablePrefixToRow := a.getDsForExtendTablePrefixToRow(suffixList)

	equalTableRowToPrefix := getEqualRowForPrefix(tablePrefixToRow)

	return a.findConsistentForRowInTables(equalTableRowToPrefix, tablePrefixToRow, extendTablePrefixToRow, suffixList)
}

func (a *Angluin) getDsForTablePrefixToRow(suffixList []string) map[string]string {
	tableRowToPrefix := make(map[string]string)

	for p := range a.prefix {
		tableRowToPrefix[p] = a.getTableRow(p, suffixList)
	}

	return tableRowToPrefix
}

func (a *Angluin) getDsForExtendTablePrefixToRow(suffixList []string) map[string]string {
	extendTableRowToPrefix := make(map[string]string)

	for ep := range a.extendPrefix {
		extendTableRowToPrefix[ep] = a.getExtendTableRow(ep, suffixList)
	}

	return extendTableRowToPrefix
}

func getEqualRowForPrefix(prefixToRow map[string]string) map[string][]string {
	result := make(map[string][]string)
	for prefix, row := range prefixToRow {
		result[row] = append(result[row], prefix)
	}

	return result
}

func (a *Angluin) findConsistentForRowInTables(
	equalTableRowToPrefix map[string][]string,
	tablePrefixToRow, extendTablePrefixToRow map[string]string,
	suffixList []string) string {

	for _, prefixes := range equalTableRowToPrefix {
		for i, prefix1 := range prefixes {
			for _, prefix2 := range prefixes[i+1:] {
				for _, letter := range a.mat.GetAlphabet() {
					combined1 := getWordWithLetter(prefix1, letter)
					combined2 := getWordWithLetter(prefix2, letter)
					if ok, row1, row2 := a.consistentForPair(
						combined1, combined2, tablePrefixToRow, extendTablePrefixToRow,
					); !ok {
						return findSuffixAndLetterInRow(row1, row2, letter, suffixList)
					}
				}
			}
		}
	}

	return ""
}

func getWordWithLetter(p string, l byte) string {
	return p + string(l)
}

func (a *Angluin) consistentForPair(
	prefixLetter1, prefixLetter2 string, tablePrefixToRow, extendTablePrefixToRow map[string]string,
) (bool, string, string) {

	var row1, row2 string
	var ok1, ok2 bool

	row1, ok1 = tablePrefixToRow[prefixLetter1]
	if !ok1 {
		row1, ok1 = extendTablePrefixToRow[prefixLetter1]
		if !ok1 {
			log.Println("row1", prefixLetter1, "not find")
		}
	}

	row2, ok2 = tablePrefixToRow[prefixLetter2]
	if !ok2 {
		row2, ok2 = extendTablePrefixToRow[prefixLetter2]
		if !ok2 {
			log.Println("row2", prefixLetter2, "not find")
		}
	}

	return row1 == row2, row1, row2
}

func findSuffixAndLetterInRow(row1, row2 string, letter byte, suffixList []string) string {

	for i := 0; i < len(row1); i++ {
		if row1[i] != row2[i] {
			return string(letter) + suffixList[i]
		}
	}

	log.Println("This is can not happened. row1 == row2 but row1 != row2")

	return ""
}
