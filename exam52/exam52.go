package main

import (
	"container/list"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const (
	inputFileName  = "exam52/input.txt"
	outputFileName = "exam52/output.txt"
)

type ioBufType struct {
	inBuf  []byte
	outBuf []byte
}

var ioBuf ioBufType

func readData(name string, ioEnd chan string) {
	var err error
	ioBuf.inBuf, err = os.ReadFile(name)
	if err != nil {
		panic(err)
	}
	ioEnd <- "ok"
}

func writeData(name string, ioEnd chan string) {
	var err error
	err = os.WriteFile(name, ioBuf.outBuf, 644)
	if err != nil {
		panic(err)
	}
	ioEnd <- "ok"
}

type parsedDataType struct {
	records   [][][]byte
	recipes   [][][]byte
	queries   [][][]byte
	unpackedN int
	unpackedQ int
}

const (
	regexpStr = "(\\d+)[ ]+(\\d+)[ ]*(\\d*)[ ]*([ \\d]*)|(\\d+)"
	rxiNComp  = iota
	rxiComp1
	rxiComp2
	rxiCompsEtc
	rxiN
	rxiQCA = iota - rxiN
	rxiQCB
	rxiQPotion
	rxiQ = rxiN
)

func parseData(buf []byte) (result parsedDataType) {
	r := regexp.MustCompile(regexpStr)
	m := r.FindAllSubmatch(buf, -1)
	if len(m) < 2 {
		return
	}
	result.records = m
	result.unpackedN, _ = strconv.Atoi(string(m[0][rxiN]))
	result.recipes = m[1 : result.unpackedN+1-2]
	if len(m) < result.unpackedN+1 {
		return
	}
	result.unpackedQ, _ = strconv.Atoi(string(m[result.unpackedN-1][rxiQ]))
	result.queries = m[result.unpackedN:]
	return
}

type queryType struct {
	reqPotion string
	ingA      string
	ingB      string
}

func processQuery(query queryType) bool {
	qA, _ := strconv.Atoi(query.ingA)
	qB, _ := strconv.Atoi(query.ingB)
	qRP, _ := strconv.Atoi(query.reqPotion)
	l := list.New()
	l.PushBack(recipesMap[qRP].components)
	for e := l.Front(); e != nil; {
		for _, c := range e.Value.([]int) {
			switch c {
			case qRP:
				return false
			case 1:
				qA--
			case 2:
				qB--
			default:
				l.PushBack(recipesMap[c].components)
			}
		}
		if qA < 0 || qB < 0 {
			return false
		}
		nextE := e.Next()
		l.Remove(e)
		e = nextE
	}
	return true
}

type recipesType struct {
	idPotion    int
	nComponents int
	components  []int
}

var recipesMap map[int]recipesType

func main() {
	ioEnd := make(chan string)
	go readData(inputFileName, ioEnd)
	<-ioEnd

	parsedData := parseData(ioBuf.inBuf)

	recipesMap = make(map[int]recipesType)
	for i, rawRecipe := range parsedData.recipes {
		nc, _ := strconv.Atoi(string(rawRecipe[rxiNComp]))
		cmps := make([]int, 0, nc)
		if nc >= 1 {
			cmp1, _ := strconv.Atoi(string(rawRecipe[rxiComp1]))
			cmps = append(cmps, cmp1)
			if nc >= 2 {
				cmp2, _ := strconv.Atoi(string(rawRecipe[rxiComp2]))
				cmps = append(cmps, cmp2)
				if nc >= 3 {
					rcmpe := strings.Fields(string(rawRecipe[rxiCompsEtc]))
					for _, s := range rcmpe {
						cmpe, _ := strconv.Atoi(s)
						cmps = append(cmps, cmpe)
					}
				}
			}
		}
		recipesMap[i+3] = recipesType{
			idPotion:    i + 3,
			nComponents: nc,
			components:  cmps,
		}
	}

	res := make([]string, 0, len(parsedData.queries))
	for _, query := range parsedData.queries {
		b := processQuery(queryType{
			reqPotion: string(query[rxiQPotion]),
			ingA:      string(query[rxiQCA]),
			ingB:      string(query[rxiQCB]),
		})
		if b {
			res = append(res, "1")
		} else {
			res = append(res, "0")
		}
	}

	ioBuf.outBuf = []byte(strings.Join(res, ""))

	go writeData(outputFileName, ioEnd)
	<-ioEnd
}
