package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	inputFileName  = "exam12/input.txt"
	outputFileName = "exam12/output.txt"
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
	firstRecord []byte
	csvBlock    []byte
	unpackedN   int
}

const (
	regexpStr = "(\\d+).*|(\\w.*)"
	rxiN      = iota
	rxiCSV
)

func parseData(buf []byte) (result parsedDataType) {
	r := regexp.MustCompile(regexpStr)
	ixs := r.FindAllSubmatchIndex(buf, 2)
	if len(ixs) < 2 {
		return
	}
	result.firstRecord = buf[ixs[0][0]:ixs[0][1]]
	result.csvBlock = buf[ixs[1][0]:]
	result.unpackedN, _ = strconv.Atoi(string(buf[ixs[0][2*rxiN]:ixs[0][2*rxiN+1]]))
	return
}

type recordType struct {
	lastName   string
	firstName  string
	patronymic string
	birthDay   string
	birthMonth string
	birthYear  string
}

func (r *recordType) validate() bool {
	d, _ := strconv.Atoi(r.birthDay)
	m, _ := strconv.Atoi(r.birthMonth)
	y, _ := strconv.Atoi(r.birthYear)
	dateOfBirthStr := fmt.Sprintf("%02v.%02v.%04v", d, m, y)
	_, err := time.Parse("02.01.2006", dateOfBirthStr)
	if err != nil {
		return false
	}
	return true
}

func (r *recordType) getHash() string {
	var differentSymbols int = 0
	var sumDDMM int = 0
	var indexFirstChar int = 0

	chars := map[rune]struct{}{}
	for _, c := range r.lastName {
		chars[c] = struct{}{}
	}
	for _, c := range r.firstName {
		chars[c] = struct{}{}
	}
	for _, c := range r.patronymic {
		chars[c] = struct{}{}
	}
	differentSymbols = len(chars)

	for _, d := range r.birthDay {
		sumDDMM += int(d - '0')
	}
	for _, m := range r.birthMonth {
		sumDDMM += int(m - '0')
	}

	upperLastName := strings.ToUpper(r.lastName)
	indexFirstChar = int(upperLastName[0] - 'A' + 1)

	hashInt := differentSymbols + 64*sumDDMM + 256*indexFirstChar
	hash := fmt.Sprintf("%03X", hashInt)
	if len(hash) > 3 {
		hash = hash[len(hash)-3:]
	}
	return hash
}

func main() {
	ioEnd := make(chan string)
	go readData(inputFileName, ioEnd)
	<-ioEnd

	parsedData := parseData(ioBuf.inBuf)

	csvReader := csv.NewReader(bytes.NewReader(parsedData.csvBlock))
	csvRecords, err := csvReader.ReadAll()
	if err != nil {
		return
	}
	resultHashes := make([]string, 0, len(csvRecords))
	for _, csvRecord := range csvRecords {
		record := recordType{
			lastName:   csvRecord[0],
			firstName:  csvRecord[1],
			patronymic: csvRecord[2],
			birthDay:   csvRecord[3],
			birthMonth: csvRecord[4],
			birthYear:  csvRecord[5],
		}
		if record.validate() {
			resultHashes = append(resultHashes, record.getHash())
		}
	}
	ioBuf.outBuf = []byte(strings.Join(resultHashes, " "))

	go writeData(outputFileName, ioEnd)
	<-ioEnd
}
