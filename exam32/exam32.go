package main

import (
	"os"
	"regexp"
	"strconv"
	"strings"
)

const (
	inputFileName  = "exam32/input.txt"
	outputFileName = "exam32/output.txt"
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
	records         [][]byte
	unpackedN       int
	unpackedQ       int
	unpackedChanges []int
}

const (
	regexpStr = "\\d+"
)

func parseData(buf []byte) (result parsedDataType) {
	r := regexp.MustCompile(regexpStr)
	m := r.FindAll(buf, -1)
	if len(m) < 2 {
		return
	}
	result.records = m
	result.unpackedN, _ = strconv.Atoi(string(m[0]))
	result.unpackedQ, _ = strconv.Atoi(string(m[1]))
	result.unpackedChanges = make([]int, len(m)-2)
	for i := 2; i < len(m); i++ {
		result.unpackedChanges[i-2], _ = strconv.Atoi(string(m[i]))
	}
	return
}

type treeType struct {
	root    int
	storage []treeNodeType
}

type treeNodeType struct {
	id         int
	parent     int
	leftChild  int
	rightChild int
}

func newTree(N int) treeType {
	res := treeType{
		root:    1,
		storage: make([]treeNodeType, N+1),
	}
	for i := 0; i < N; i++ {
		var lc, rc int = 0, 0
		if (i+1)*2 <= N {
			lc = (i + 1) * 2
		}
		if (i+1)*2+1 <= N {
			rc = (i+1)*2 + 1
		}
		res.storage[i+1] = treeNodeType{
			id:         i + 1,
			parent:     (i + 1) / 2,
			leftChild:  lc,
			rightChild: rc,
		}
	}
	return res
}

func (t *treeType) getLVRSlice(v int) []int {
	var res []int
	node := t.storage[v]
	if node.leftChild != 0 {
		res = t.getLVRSlice(node.leftChild)
	}
	res = append(res, node.id)
	if node.rightChild != 0 {
		res = append(res, t.getLVRSlice(node.rightChild)...)
	}
	return res
}

func (t *treeType) changeNode(v int) {
	p := t.storage[v].parent
	pp := t.storage[p].parent
	vl := t.storage[v].leftChild
	vr := t.storage[v].rightChild

	if p == 0 {
		return
	}

	if pp == 0 {
		t.root = v
	} else {
		if t.storage[pp].leftChild == p {
			t.storage[pp].leftChild = v
		} else {
			t.storage[pp].rightChild = v
		}
	}
	t.storage[v].parent = pp
	if t.storage[p].leftChild == v {
		t.storage[v].leftChild = p
		t.storage[p].parent = v
		if vl != 0 {
			t.storage[vl].parent = p
		}
		t.storage[p].leftChild = vl
	} else if t.storage[p].rightChild == v {
		t.storage[v].rightChild = p
		t.storage[p].parent = v
		if vr != 0 {
			t.storage[vr].parent = p
		}
		t.storage[p].rightChild = vr
	}
}

func main() {
	ioEnd := make(chan string)
	go readData(inputFileName, ioEnd)
	<-ioEnd

	parsedData := parseData(ioBuf.inBuf)

	tree := newTree(parsedData.unpackedN)
	for _, change := range parsedData.unpackedChanges {
		tree.changeNode(change)
	}

	LVRSlice := tree.getLVRSlice(tree.root)
	LVRStrSlice := make([]string, len(LVRSlice))
	for i, item := range LVRSlice {
		LVRStrSlice[i] = strconv.Itoa(item)
	}

	ioBuf.outBuf = []byte(strings.Join(LVRStrSlice, " "))

	go writeData(outputFileName, ioEnd)
	<-ioEnd
}
