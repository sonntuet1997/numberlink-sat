package numberlink

var MoveType = map[string]int{
	"Left":  1,
	"Up":    2,
	"Right": 3,
	"Down":  4,
}

type Board struct {
	Size  int
	Size2 int

	Row       int
	Column    int
	MaxValue  int
	TotalMove int

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

func New(row, column, maxValue int) *Board {
	size2 := 1
	size := 1
	candidates := make([]bool, row*column*(maxValue+len(MoveType))+1)

	rowCandidateCount := make([]int, size2*size2)
	colCandidateCount := make([]int, size2*size2)
	blkCandidateCount := make([]int, size2*size2)
	blkIdxMap := make([]int, size2*size2)

	board := &Board{
		Row:       row,
		Column:    column,
		MaxValue:  maxValue,
		TotalMove: len(MoveType),
		Lookup:    make([]int, row*column),

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
	b.Lookup[b.Idx(row, col)] = val
}

//func (b *Board) GetUnresolvedCells() int {
//	result := 0
//	for r := 0; r < b.Size2; r++ {
//		for c := 0; c < b.Size2; c++ {
//			if b.Lookup[b.Idx(r, c)] == 0 {
//				result++
//			}
//		}
//	}
//	return result
//}
// func (b *Board) InitTriads() {

// }

func (b *Board) SolveWithModel(model []bool) {
	//fmt.Printf("Model %v\n", model)
	for r := 0; r < b.Row; r++ {
		for c := 0; c < b.Column; c++ {
			if b.Lookup[r*b.Column+c] == 0 {
				for v := 1; v <= b.MaxValue; v++ {
					if model[b.YLit(r, c, v)-1] {
						b.SetValue(r, c, v)
						break
					}
				}
			}
		}
	}

	//// Print
	//for r := 0; r < b.Row; r++ {
	//	for c := 0; c < b.Column; c++ {
	//		for v := 1; v <= b.MaxValue; v++ {
	//			if model[b.YLit(r, c, v)-1] {
	//				print(v, " ")
	//			}
	//		}
	//		for k, _ := range MoveType {
	//			if model[b.XLit(r, c, k)-1] {
	//				print(k, "    ")
	//			}
	//		}
	//	}
	//	println()
	//}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// XLit 1-indexed
// Cell candidate values
func (b *Board) XLit(row, col int, director string) int {
	if MoveType[director] == 0 {
		panic("wrong direction!")
	}
	return 1 + b.Idx(row, col)*(b.MaxValue+b.TotalMove) + (MoveType[director] - 1)
}

// YLit 1-indexed
// Cell candidate values
func (b *Board) YLit(row, col, val int) int {
	return 1 + b.Idx(row, col)*(b.MaxValue+b.TotalMove) + b.TotalMove + (val - 1)
}

// Idx 0-indexed
func (b *Board) Idx(row, col int) int {
	return row*b.Column + col
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
