package main

import (
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

const (
	inputFileName  = "exam22/input.txt"
	outputFileName = "exam22/output.txt"
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
	unpackedN int
}

const (
	regexpStr = "(\\d+)\\s+(\\d+)\\s+(\\d+)\\s+(\\d+)\\s+([ABCS])|(\\d+)"
	rxiDay    = iota
	rxiHour
	rxiMinute
	rxiID
	rxiStatus
	rxiN
)

func parseData(buf []byte) (result parsedDataType) {
	r := regexp.MustCompile(regexpStr)
	m := r.FindAllSubmatch(buf, -1)
	if len(m) < 2 {
		return
	}
	result.records = m
	result.unpackedN, _ = strconv.Atoi(string(m[0][rxiN]))
	return
}

type logRecordType struct {
	parsedRecord [][]byte
	timeDHM      int
}

func getFlyTime(logRecords []logRecordType) string {
	for i, logRecord := range logRecords {
		d, _ := strconv.Atoi(string(logRecord.parsedRecord[rxiDay]))
		h, _ := strconv.Atoi(string(logRecord.parsedRecord[rxiHour]))
		m, _ := strconv.Atoi(string(logRecord.parsedRecord[rxiMinute]))
		logRecords[i].timeDHM = d*24*60 + h*60 + m
	}
	sort.SliceStable(logRecords, func(i, j int) bool {
		return logRecords[i].timeDHM < logRecords[j].timeDHM
	})
	var timeA, flyTime int
	for _, record := range logRecords {
		switch string(record.parsedRecord[rxiStatus]) {
		case "A":
			timeA = record.timeDHM
		case "S":
			flyTime += record.timeDHM - timeA
		case "C":
			flyTime += record.timeDHM - timeA
		}
	}

	return strconv.Itoa(flyTime)
}

func main() {
	ioEnd := make(chan string)
	go readData(inputFileName, ioEnd)
	<-ioEnd

	rocketsLogMap := make(map[int][]logRecordType)
	parsedData := parseData(ioBuf.inBuf)
	for i, parsedRecord := range parsedData.records {
		if i == 0 {
			continue
		}
		id, _ := strconv.Atoi(string(parsedRecord[rxiID]))
		rocketsLogMap[id] = append(rocketsLogMap[id], logRecordType{
			parsedRecord: parsedRecord,
			timeDHM:      0,
		})
	}

	rocketsIDSlice := make([]int, 0, len(rocketsLogMap))
	flyTimes := make([]string, len(rocketsLogMap))
	for i, _ := range rocketsLogMap {
		rocketsIDSlice = append(rocketsIDSlice, i)
	}
	sort.Ints(rocketsIDSlice)
	for i, id := range rocketsIDSlice {
		flyTimes[i] = getFlyTime(rocketsLogMap[id])
	}

	ioBuf.outBuf = []byte(strings.Join(flyTimes, " "))

	go writeData(outputFileName, ioEnd)
	<-ioEnd
}
