package numberlink

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"regexp"
	"strconv"
	"strings"
)

var SPACE_REGEX = regexp.MustCompile(`  +`)

/*
Parse newline and space separated sudoku problem
0 0 3 ...
9 0 0 ...
0 0 1 ...
...
*/
func NewFromString(input string) *Board {
	input = strings.Trim(input, " \n\t")
	input = SPACE_REGEX.ReplaceAllString(input, " ")
	input = strings.ReplaceAll(input, ".", "0")

	stringArr := strings.Split(input, "\n")
	info := stringArr[0]
	infoArr := strings.Split(info, " ")
	column, _ := strconv.Atoi(infoArr[0])
	row, _ := strconv.Atoi(infoArr[1])

	cells := make([][]int, row)
	r := bufio.NewReader(strings.NewReader(strings.Join(stringArr[1:], "\n")))

	for i := 0; i < row; i++ {
		cells[i] = make([]int, column)
		for j := 0; j < column; j++ {
			fmt.Fscan(r, &cells[i][j])
		}
	}

	return NewFromArray(cells)
}

func (b *Board) ReplaceWithSingleRowString(input string, skipCandidateElimination bool) {
	size2 := 9
	b.NumCandidates = len(b.Candidates) - 1

	for i := 0; i < len(b.Lookup); i++ {
		b.Lookup[i] = 0
	}

	if !skipCandidateElimination {
		for i := 1; i < len(b.Candidates); i++ {
			b.Candidates[i] = true
		}
		for i := 0; i < len(b.rowCandidateCount); i++ {
			b.rowCandidateCount[i] = size2
			b.colCandidateCount[i] = size2
			b.blkCandidateCount[i] = size2
		}
	}

	for i, c := range input {
		if c != '0' && c != '.' {
			if skipCandidateElimination {
				b.Lookup[b.Idx(i/size2, i%size2)] = int(c - '0')
			} else {
				b.SetValue(i/size2, i%size2, int(c-'0'))
			}
		}
	}
}

func NewFromArray(cells [][]int) *Board {
	row := len(cells)
	column := len(cells[0])
	maxValue := 1
	for _, row := range cells {
		for _, val := range row {
			if val > maxValue {
				maxValue = val
			}
		}
	}
	board := New(row, column, maxValue)
	for r, row := range cells {
		for c, val := range row {
			if val < 1 {
				continue
			}
			board.SetValue(r, c, val)
		}
	}
	return board
}

func (s *Board) Print(w io.Writer) {
	charLen := int(math.Floor(math.Log10(float64(s.Row * s.Column * (s.MaxValue + s.TotalMove)))))
	formatter := fmt.Sprintf("%%s%%%dd%%s", charLen)
	for r := 0; r < s.Row; r++ {
		fmt.Fprintf(w, formatter, "", s.Lookup[s.Idx(r, 0)], "")
		for c := 1; c < s.Column-1; c++ {
			fmt.Fprintf(w, formatter, " ", s.Lookup[s.Idx(r, c)], "")
		}
		fmt.Fprintf(w, formatter, " ", s.Lookup[s.Idx(r, s.Column-1)], "\n")
	}
}

func (s *Board) PrintOneLine(w io.Writer) {
	for i := 0; i < s.Size2*s.Size2; i++ {
		fmt.Fprint(w, s.Lookup[i])
	}
	fmt.Fprintln(w)
}
