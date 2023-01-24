package main

import (
	"os"
	"regexp"
	"strconv"
	"strings"
)

const (
	inputFileName  = "exam42/input.txt"
	outputFileName = "exam42/output.txt"
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
	orders    [][][]byte
	queries   [][][]byte
	unpackedN int
	unpackedQ int
}

const (
	regexpStr = "(\\d+)[ ](\\d+)[ ](\\d+)|(\\d+)"
	rxiStart  = iota
	rxiEnd
	rxiCost
	rxiN
	rxiType = rxiCost
	rxiQ    = rxiN
)

func parseData(buf []byte) (result parsedDataType) {
	r := regexp.MustCompile(regexpStr)
	m := r.FindAllSubmatch(buf, -1)
	if len(m) < 2 {
		return
	}
	result.records = m
	result.unpackedN, _ = strconv.Atoi(string(m[0][rxiN]))
	result.orders = m[1 : result.unpackedN+1]
	if len(m) < result.unpackedN+3 {
		return
	}
	result.unpackedQ, _ = strconv.Atoi(string(m[result.unpackedN+1][rxiQ]))
	result.queries = m[result.unpackedN+2:]
	return
}

type ordersTableType struct {
	orders []orderType
}

type orderType struct {
	start int
	end   int
	cost  int
}

func newOrderTable(orders [][][]byte) (result ordersTableType) {
	result.orders = make([]orderType, len(orders))
	for i, order := range orders {
		result.orders[i].start, _ = strconv.Atoi(string(order[rxiStart]))
		result.orders[i].end, _ = strconv.Atoi(string(order[rxiEnd]))
		result.orders[i].cost, _ = strconv.Atoi(string(order[rxiCost]))
	}
	return
}

func (ot *ordersTableType) processQuery(query [][]byte) string {
	queryStart, _ := strconv.Atoi(string(query[rxiStart]))
	queryEnd, _ := strconv.Atoi(string(query[rxiEnd]))
	queryType, _ := strconv.Atoi(string(query[rxiType]))
	var resInt int
	switch queryType {
	case 1:
		sumCost := 0
		for _, ord := range ot.orders {
			if queryStart <= ord.start && ord.start <= queryEnd {
				sumCost += ord.cost
			}
		}
		resInt = sumCost
	case 2:
		sumDuration := 0
		for _, ord := range ot.orders {
			if queryStart <= ord.end && ord.end <= queryEnd {
				sumDuration += ord.end - ord.start
			}
		}
		resInt = sumDuration
	}
	return strconv.Itoa(resInt)
}

func main() {
	ioEnd := make(chan string)
	go readData(inputFileName, ioEnd)
	<-ioEnd

	parsedData := parseData(ioBuf.inBuf)
	ot := newOrderTable(parsedData.orders)
	res := make([]string, len(parsedData.queries))
	for i, query := range parsedData.queries {
		res[i] = ot.processQuery(query)
	}

	ioBuf.outBuf = []byte(strings.Join(res, " "))

	go writeData(outputFileName, ioEnd)
	<-ioEnd
}
