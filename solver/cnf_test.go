package solver_test

import (
	"bytes"
	"github.com/sonntuet1997/numberlink-sat/numberlink"
	"github.com/sonntuet1997/numberlink-sat/solver"
	"io/ioutil"
	"strings"
	"testing"
)

var aiEscargot1 = [...]string{
	"100007090030020008009600500005300900010080002600004000300000010041000007007000300",
	"162857493534129678789643521475312986913586742628794135356478219241935867897261354",
}

var hard11 = [...]string{
	"........8..3...4...9..2..6.....79.......612...6.5.2.7...8...5...1.....2.4.5.....3",
	"621943758783615492594728361142879635357461289869532174238197546916354827475286913",
}

var hard17clue1 = [...]string{
	"000000010400000000020000000000050407008000300001090000300400200050100000000806000",
	"693784512487512936125963874932651487568247391741398625319475268856129743274836159",
}

func cnfOneLiner(input, algorithm string) string {
	board := numberlink.NewFromString(input)
	// solver.Solve(board)
	solver.GetInfo(board, algorithm)
	var b bytes.Buffer
	board.PrintOneLine(&b)
	return strings.TrimSpace(b.String())
}

func TestCNFAiEscargot(b *testing.T) {
	cnfOneLiner(aiEscargot1[0], "normal")
}
func TestCNFAiEscargotWithProductEncoding(b *testing.T) {
	cnfOneLiner(aiEscargot1[0], "product")
}

func TestCNFHard9x9(b *testing.T) {
	cnfOneLiner(hard11[0], "normal")
}
func TestCNFHard9x9WithProductEncoding(b *testing.T) {
	cnfOneLiner(hard11[0], "product")
}

func TestCNF17clue9x9(b *testing.T) {
	cnfOneLiner(hard17clue1[0], "normal")
}
func TestCNF17clue9x9WithProductEncoding(b *testing.T) {
	cnfOneLiner(hard17clue1[0], "product")
}

func TestCNF25x25(b *testing.T) {
	bytes, _ := ioutil.ReadFile("../data/sudoku-25-1.txt")
	input := string(bytes)
	cnfOneLiner(input, "normal")
}
func TestCNF25x25WithProductEncoding(b *testing.T) {
	bytes, _ := ioutil.ReadFile("../data/sudoku-25-1.txt")
	input := string(bytes)
	cnfOneLiner(input, "product")
}

func TestCNF64x64(b *testing.T) {
	bytes, _ := ioutil.ReadFile("../data/sudoku-64-2.txt")
	input := string(bytes)
	cnfOneLiner(input, "normal")
}
func TestCNF64x64WithProductEncoding(b *testing.T) {
	bytes, _ := ioutil.ReadFile("../data/sudoku-64-2.txt")
	input := string(bytes)
	cnfOneLiner(input, "product")
}

//TLE (>11m)
func TestCNF64x64Hard(b *testing.T) {
	bytes, _ := ioutil.ReadFile("../data/sudoku-64-1.txt")
	input := string(bytes)
	cnfOneLiner(input, "normal")
}
func TestCNF64x64HardWithProductEncoding(b *testing.T) {
	bytes, _ := ioutil.ReadFile("../data/sudoku-64-1.txt")
	input := string(bytes)
	cnfOneLiner(input, "product")
}

func TestCNF81x81(b *testing.T) {
	bytes, _ := ioutil.ReadFile("../data/sudoku-81-1.txt")
	input := string(bytes)
	cnfOneLiner(input, "normal")
}
func TestCNF81x81WithProductEncoding(b *testing.T) {
	bytes, _ := ioutil.ReadFile("../data/sudoku-81-1.txt")
	input := string(bytes)
	cnfOneLiner(input, "product")
}

//
func TestCNF100x100(b *testing.T) {
	bytes, _ := ioutil.ReadFile("../data/sudoku-100-1.txt")
	input := string(bytes)
	cnfOneLiner(input, "normal")
}
func TestCNF100x100WithProductEncoding(b *testing.T) {
	bytes, _ := ioutil.ReadFile("../data/sudoku-100-1.txt")
	input := string(bytes)
	cnfOneLiner(input, "product")
}

//
func TestCNF144x144(b *testing.T) {
	bytes, _ := ioutil.ReadFile("../data/sudoku-144-1.txt")
	input := string(bytes)
	cnfOneLiner(input, "normal")
}
func TestCNF144x144WithProductEncoding(b *testing.T) {
	bytes, _ := ioutil.ReadFile("../data/sudoku-144-1.txt")
	input := string(bytes)
	cnfOneLiner(input, "product")
}

func TestCNF225x225(b *testing.T) {
	bytes, _ := ioutil.ReadFile("../data/sudoku-225-2.txt")
	input := string(bytes)
	cnfOneLiner(input, "normal")
}
func TestCNF225x225WithProductEncoding(b *testing.T) {
	bytes, _ := ioutil.ReadFile("../data/sudoku-225-2.txt")
	input := string(bytes)
	cnfOneLiner(input, "product")
}
