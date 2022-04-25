package numberlink

import (
	"math"
)

type Board struct {
	Size       int
	Size2      int
	Candidates []bool // lit
	Lookup     []int  // idx

	NumCandidates     int
	rowCandidateCount []int
	colCandidateCount []int
	blkCandidateCount []int
	blkIdxMap         []int // idx -> blkIdx

	lit_cLit []int // lit -> compressed lit, 1-indexed
	cLit_lit []int // compressed lit -> lit, 1-indexed
}

func New(size int) *Board {
	size2 := size * size
	candidates := make([]bool, size2*size2*size2+1)
	rowCandidateCount := make([]int, size2*size2)
	colCandidateCount := make([]int, size2*size2)
	blkCandidateCount := make([]int, size2*size2)
	blkIdxMap := make([]int, size2*size2)

	board := &Board{
		Size:       size,
		Size2:      size2,
		Lookup:     make([]int, size2*size2),
		blkIdxMap:  blkIdxMap,
		Candidates: candidates,

		rowCandidateCount: rowCandidateCount,
		colCandidateCount: colCandidateCount,
		blkCandidateCount: blkCandidateCount,

		NumCandidates: len(candidates) - 1,
	}

	// determine block index of cells
	for r := 0; r < size2; r++ {
		for c := 0; c < size2; c++ {
			blkIdxMap[r*size2+c] = r/size*size + c/size
		}
	}

	// keep track numbers left in each row, column, block
	for i := 0; i < len(rowCandidateCount); i++ {
		rowCandidateCount[i] = size2
		colCandidateCount[i] = size2
		blkCandidateCount[i] = size2
	}

	//keep track all possible value of each cell
	for i := 1; i < len(candidates); i++ {
		candidates[i] = true
	}

	return board
}

func (b *Board) SetValue(row, col, val int) {
	blkIndex := b.blkIdxMap[b.Idx(row, col)]
	b.Lookup[b.Idx(row, col)] = val
	b.rowCandidateCount[row*b.Size2+val-1] = 1
	b.colCandidateCount[col*b.Size2+val-1] = 1
	b.blkCandidateCount[blkIndex*b.Size2+val-1] = 1
	blkRStart := b.Size * (row / b.Size)
	blkCStart := b.Size * (col / b.Size)

	for i := 0; i < b.Size2; i++ {
		if i+1 != val {
			b.SetValueFalse(row, col, i+1)
		}
		if i != row {
			b.SetValueFalse(i, col, val)
		}
		if i != col {
			b.SetValueFalse(row, i, val)
		}
	}

	for r := 0; r < b.Size; r++ {
		for c := 0; c < b.Size; c++ {
			blkR := blkRStart + r
			blkC := blkCStart + c
			if blkR != row && blkC != col {
				b.SetValueFalse(blkR, blkC, val)
			}
		}
	}
}

func (b *Board) GetUnresolvedCells() int {
	result := 0
	for r := 0; r < b.Size2; r++ {
		for c := 0; c < b.Size2; c++ {
			if b.Lookup[b.Idx(r, c)] == 0 {
				result++
			}
		}
	}
	return result
}

func (b *Board) SetValueFalse(row, col, val int) {
	blkIndex := b.blkIdxMap[b.Idx(row, col)]
	lit := b.Lit(row, col, val)
	prev := b.Candidates[lit]
	b.Candidates[lit] = false
	if prev {
		b.NumCandidates--
		b.rowCandidateCount[row*b.Size2+val-1] -= 1
		b.colCandidateCount[col*b.Size2+val-1] -= 1
		b.blkCandidateCount[blkIndex*b.Size2+val-1] -= 1
	}
}

// Lặp cho đến khi không tìm được thêm giá trị nào có thể điển
func (b *Board) BasicSolve() {
	restart := true
	for restart {
		restart = false
		restart = restart || b.NakedSingles()
		restart = restart || b.HiddenSingles()
	}
}

func (b *Board) NakedSingles() bool {
	restart := false
	for r := 0; r < b.Size2; r++ {
		for c := 0; c < b.Size2; c++ {
			if b.Lookup[b.Idx(r, c)] != 0 {
				continue
			}
			// naked singles
			last := 0
			for v := 1; v <= b.Size2; v++ {
				if !b.Candidates[b.Lit(r, c, v)] {
					continue
				}
				if last != 0 {
					last = 0
					break
				}
				last = v
			}
			if last != 0 {
				b.SetValue(r, c, last)
				restart = true
			}
		}
	}
	return restart
}

func (b *Board) HiddenSingles() bool {
	restart := false
	for i := 0; i < b.Size2; i++ {
		for v := 1; v <= b.Size2; v++ {
			if b.rowCandidateCount[i*b.Size2+v-1] == 1 {
				for j := 0; j < b.Size2; j++ {
					if b.Candidates[b.Lit(i, j, v)] {
						if b.Lookup[b.Idx(i, j)] != v {
							b.SetValue(i, j, v)
							restart = true
						}
						break
					}
				}
			}
			if b.colCandidateCount[i*b.Size2+v-1] == 1 {
				for j := 0; j < b.Size2; j++ {
					if b.Candidates[b.Lit(j, i, v)] {
						if b.Lookup[b.Idx(j, i)] != v {
							b.SetValue(j, i, v)
							restart = true
						}
						break
					}
				}
			}
			if b.blkCandidateCount[i*b.Size2+v-1] == 1 {
				blkRStart := (i / b.Size) * b.Size
				blkCStart := (i % b.Size) * b.Size
				for r := 0; r < b.Size; r++ {
					for c := 0; c < b.Size; c++ {
						// block
						blkR := blkRStart + r
						blkC := blkCStart + c
						if b.Candidates[b.Lit(blkR, blkC, v)] {
							if b.Lookup[b.Idx(blkR, blkC)] != v {
								b.SetValue(blkR, blkC, v)
								restart = true
							}
							break
						}
					}
				}
			}
		}
	}

	return restart
}

// Xijk <=> b.Candidates = [0,0,0,0,   0,1,0,0,   0,0,0,0,     1,0,0,0,
//                          0,0,0,0,   0,0,0,0,   0,0,0,0,     0,0,0,0,
//                          0,0,0,0,   0,0,0,0,   0,0,0,0,     0,0,0,0,
//                          0,0,0,1,   0,0,0,0,   0,1,0,0,     0,0,0,0]
// Xijk <=> b.Candidates = [false false false true false  |  false true false false |  false false true true | true false false false
//          (thừa phần tử đầu)    true false true  false  |  true false true true   |   false false true true| false true true true
//                                 true true true false   |   true false true false |   true false true true | false false true true
//                                 false false false true |  true false true false  |   false true false false| false false true false]
// b.NumCandidates: Số lượng Xijk chưa được điền
// Thực hiện giảm số biến bằng cách lọc các phần tử bằng false. thực hiện ghi lại các vị trí bằng true trong candidate vào cLit_lit
// b.lit_cLit = [0 3 6 11 12 13 17 19 21 23 24 27 28 30 31 32 33 34 35 37 39 41 43 44 47 48 52 53 55 58 63] (bỏ phần tử đầu)
func (b *Board) InitCompressedLits() {
	b.lit_cLit = make([]int, len(b.Candidates)+1)
	b.cLit_lit = make([]int, b.NumCandidates+1)
	j := 1
	for i := 1; i < len(b.Candidates); i++ {
		if !b.Candidates[i] {
			continue
		}
		b.lit_cLit[i] = j
		b.cLit_lit[j] = i
		j++
	}
}

// func (b *Board) InitTriads() {

// }

// from model of compressed lits
// Xem phần tử nào == true thì thực hiện tìm ra vị trí thật của nó:
// lit := b.cLit_lit[i+1]
// Cập nhật kết quả vào mảng Lookup
func (b *Board) SolveWithModel(model []bool) {
	for i := 0; i < min(len(model), len(b.cLit_lit)-1); i++ {
		if !model[i] {
			continue
		}
		lit := b.cLit_lit[i+1]
		b.Lookup[(lit-1)/b.Size2] = 1 + (lit-1)%b.Size2
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Lit 1-indexed
// Cell candidate values
func (b *Board) Lit(row, col, val int) int {
	return 1 + b.Idx(row, col)*b.Size2 + (val - 1)
}

// CLit 1-indexed
func (b *Board) CLit(row, col, val int) int {
	return b.lit_cLit[b.Lit(row, col, val)]
}

// Idx 0-indexed
func (b *Board) Idx(row, col int) int {
	return row*b.Size2 + col
}

// func (b *Board) RowTriadLit(row, col, val int) int {
// 	return b.Size2*b.Size2*b.Size2 + b.TriadIdx(row, col, val)
// }

// func (b *Board) ColTriadLit(row, col, val int) int {
// 	return b.Size2*b.Size2*b.Size2 + b.Size2*b.Size*b.Size2 + b.TriadIdx(col, row, val)
// }

// // crossAxis = row for RowTriad
// //           = col for ColTriad
// func (b *Board) TriadIdx(mainAxis, crossAxis, value int) int {
// 	return mainAxis*b.Size2*b.Size + (crossAxis/b.Size)*b.Size + (value - 1)
// }

func getSize(size2 int) int {
	size := int(math.Sqrt(float64(size2)))
	if size2 != size*size {
		panic("Size is not a square")
	}
	return size
}
