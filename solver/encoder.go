package solver

import (
	"fmt"
	"github.com/sonntuet1997/numberlink-sat/numberlink"
	"log"
	"time"
)

func GenerateCNFConstraints(board *numberlink.Board, algorithm string) CNFInterface {
	var cnf CNFInterface

	shouldUseParallel := false

	cnf = &CNF{
		Board:   board,
		Clauses: make([][]int, 0, board.Row*board.Column*(board.MaxValue+board.TotalMove)*3), // TODO: Check this
	}

	if shouldUseParallel {
		cnf = &CNFParallel{
			CNF: cnf.(*CNF),
		}
	}

	cnf.setInitialNbVar(board.NumCandidates)
	//cnf.initializeLits()

	if shouldUseParallel {
		cnf.(*CNFParallel).initWorkers()
	}
	exactly1 := cnfExactly1
	exactly2 := cnfExactly2
	atLeast1 := cnfAtLeast1
	//atMost1 := cnfAtMost1
	if algorithm == "product" {
		exactly1 = cnfExactly1Product
		atLeast1 = cnfAtLeast1
		//atMost1 = cnfAtMost1Product
	}
	fmt.Printf("%v%v%v\n", exactly1, exactly2, atLeast1)

	// log.Println("here", cnf.clauseLen())
	start := time.Now()
	//testExactly2(cnf,exactly2)
	initializeClauses(cnf) // Tested
	buildCNFAtLeast1Direction(cnf, atLeast1)
	buildCNFExact1ValuePerCell(cnf, exactly1)               // Tested
	buildCNFDirectionForNumberedCornerCell(cnf, exactly1)   // Tested
	buildCNFDirectionForNumberedBorderCell(cnf, exactly1)   // Tested
	buildCNFDirectionForNumberedInnerCell(cnf, exactly1)    // Tested
	buildCNFDirectionForUnNumberedCornerCell(cnf, atLeast1) // Tested
	//buildCNFDirectionForUnNumberedBorderCell(cnf, exactly2) // Tested
	buildCNFDirectionForUnNumberedBorderCell2(cnf, exactly1)
	buildCNFDirectionForUnNumberedInnerCell(cnf, exactly2) // Tested
	buildCNFConnectedValue(cnf, nil)
	fmt.Printf("sadsad%v\n", cnf.getClauses())
	elapsed := time.Since(start)
	log.Printf("Generating Clauses took %s", elapsed)

	if shouldUseParallel {
		cnf.(*CNFParallel).closeAndWait()
	}

	// if board.Size > 6 {
	// 	cnf.Simplify(SimplifyOptions{})
	// }

	return cnf
}

func initializeClauses(cnf CNFInterface) {
	b := cnf.getBoard()
	for r := 0; r < b.Row; r++ {
		for c := 0; c < b.Column; c++ {
			v := b.Lookup[r*b.Column+c]
			if v != 0 {
				cnf.addClause([]int{b.YLit(r, c, v)})
			}
		}
	}
}

func buildCNFAtLeast1Direction(cnf CNFInterface, builder CNFBuilder) {
	b := cnf.getBoard()
	for r := 0; r < b.Row; r++ {
		for c := 0; c < b.Column; c++ {
			lits := make([]int, b.TotalMove)
			for a, j := range numberlink.MoveType {
				lits[j-1] = b.XLit(r, c, a)
			}
			cnf.addFormula(filterZero(lits), builder)
		}
	}
}

func buildCNFExact1ValuePerCell(cnf CNFInterface, builder CNFBuilder) {
	b := cnf.getBoard()
	idx := 0
	for r := 0; r < b.Row; r++ {
		for c := 0; c < b.Column; c++ {
			idx++
			lits := make([]int, b.MaxValue)
			for v := 1; v <= b.MaxValue; v++ {
				lits[v-1] = b.YLit(r, c, v)
			}
			cnf.addFormula(filterZero(lits), builder)
		}
	}
}
func buildCNFDirectionForNumberedCornerCell(cnf CNFInterface, builder CNFBuilder) {
	b := cnf.getBoard()
	// Up Left
	if b.Lookup[0] != 0 {
		lits := make([]int, 2)
		lits[0] = b.XLit(0, 0, "Down")
		lits[1] = b.XLit(0, 0, "Right")
		cnf.addFormula(filterZero(lits), builder)
		cnf.addClause([]int{-b.XLit(0, 0, "Up")})
		cnf.addClause([]int{-b.XLit(0, 0, "Left")})
	}
	// Up right
	if b.Lookup[b.Column-1] != 0 {
		lits := make([]int, 2)
		lits[0] = b.XLit(0, b.Column-1, "Down")
		lits[1] = b.XLit(0, b.Column-1, "Left")
		cnf.addFormula(filterZero(lits), builder)
		cnf.addClause([]int{-b.XLit(0, b.Column-1, "Up")})
		cnf.addClause([]int{-b.XLit(0, b.Column-1, "Right")})
	}

	// Bottom left
	if b.Lookup[(b.Row-1)*b.Column] != 0 {
		lits := make([]int, 2)
		lits[0] = b.XLit(b.Row-1, 0, "Up")
		lits[1] = b.XLit(b.Row-1, 0, "Right")
		cnf.addFormula(filterZero(lits), builder)
		cnf.addClause([]int{-b.XLit(b.Row-1, 0, "Down")})
		cnf.addClause([]int{-b.XLit(b.Row-1, 0, "Left")})
	}

	// Bottom right
	if b.Lookup[b.Row*b.Column-1] != 0 {
		lits := make([]int, 2)
		lits[0] = b.XLit(b.Row-1, b.Column-1, "Up")
		lits[1] = b.XLit(b.Row-1, b.Column-1, "Left")
		cnf.addFormula(filterZero(lits), builder)
		cnf.addClause([]int{-b.XLit(b.Row-1, b.Column-1, "Down")})
		cnf.addClause([]int{-b.XLit(b.Row-1, b.Column-1, "Right")})
	}
}
func buildCNFDirectionForUnNumberedCornerCell(cnf CNFInterface, builder CNFBuilder) {
	b := cnf.getBoard()
	// Up Left
	if b.Lookup[0] == 0 {
		lits := make([]int, 1)
		lits[0] = b.XLit(0, 0, "Down")
		cnf.addFormula(filterZero(lits), builder)
		lits = make([]int, 1)
		lits[0] = b.XLit(0, 0, "Right")
		cnf.addFormula(filterZero(lits), builder)
		lits = make([]int, 1)
		lits[0] = -b.XLit(0, 0, "Up")
		cnf.addFormula(filterZero(lits), builder)
		lits = make([]int, 1)
		lits[0] = -b.XLit(0, 0, "Left")
		cnf.addFormula(filterZero(lits), builder)
	}
	// Up right
	if b.Lookup[b.Column-1] == 0 {
		lits := make([]int, 1)
		lits[0] = b.XLit(0, b.Column-1, "Down")
		cnf.addFormula(filterZero(lits), builder)
		lits = make([]int, 1)
		lits[0] = b.XLit(0, b.Column-1, "Left")
		cnf.addFormula(filterZero(lits), builder)
		lits = make([]int, 1)
		lits[0] = -b.XLit(0, b.Column-1, "Up")
		cnf.addFormula(filterZero(lits), builder)
		lits = make([]int, 1)
		lits[0] = -b.XLit(0, b.Column-1, "Right")
		cnf.addFormula(filterZero(lits), builder)

	}

	// Bottom left
	if b.Lookup[(b.Row-1)*b.Column] == 0 {
		lits := make([]int, 1)
		lits[0] = b.XLit(b.Row-1, 0, "Up")
		cnf.addFormula(filterZero(lits), builder)
		lits = make([]int, 1)
		lits[0] = b.XLit(b.Row-1, 0, "Right")
		cnf.addFormula(filterZero(lits), builder)
		lits = make([]int, 1)
		lits[0] = -b.XLit(b.Row-1, 0, "Left")
		cnf.addFormula(filterZero(lits), builder)
		lits = make([]int, 1)
		lits[0] = -b.XLit(b.Row-1, 0, "Down")
		cnf.addFormula(filterZero(lits), builder)
	}

	// Bottom right
	if b.Lookup[b.Row*b.Column-1] == 0 {
		lits := make([]int, 1)
		lits[0] = b.XLit(b.Row-1, b.Column-1, "Up")
		cnf.addFormula(filterZero(lits), builder)
		lits = make([]int, 1)
		lits[0] = b.XLit(b.Row-1, b.Column-1, "Left")
		cnf.addFormula(filterZero(lits), builder)
		lits = make([]int, 1)
		lits[0] = -b.XLit(b.Row-1, b.Column-1, "Right")
		cnf.addFormula(filterZero(lits), builder)
		lits = make([]int, 1)
		lits[0] = -b.XLit(b.Row-1, b.Column-1, "Down")
		cnf.addFormula(filterZero(lits), builder)
	}
}
func buildCNFDirectionForNumberedBorderCell(cnf CNFInterface, builder CNFBuilder) {
	b := cnf.getBoard()
	// Top
	for c := 1; c < b.Column-1; c++ {
		if b.Lookup[c] == 0 {
			continue
		}
		lits := make([]int, 3)
		lits[0] = b.XLit(0, c, "Right")
		lits[1] = b.XLit(0, c, "Left")
		lits[2] = b.XLit(0, c, "Down")
		cnf.addFormula(filterZero(lits), builder)
		lits = make([]int, 1)
		lits[0] = -b.XLit(0, c, "Up")
		cnf.addFormula(filterZero(lits), builder)
	}
	// Bottom
	for c := 1; c < b.Column-1; c++ {
		if b.Lookup[(b.Row-1)*b.Column+c] == 0 {
			continue
		}
		lits := make([]int, 3)
		lits[0] = b.XLit(b.Row-1, c, "Right")
		lits[1] = b.XLit(b.Row-1, c, "Left")
		lits[2] = b.XLit(b.Row-1, c, "Up")
		cnf.addFormula(filterZero(lits), builder)
		lits = make([]int, 1)
		lits[0] = -b.XLit(b.Row-1, c, "Down")
		cnf.addFormula(filterZero(lits), builder)
	}
	// Left
	for r := 1; r < b.Row-1; r++ {
		if b.Lookup[r*b.Column] == 0 {
			continue
		}
		lits := make([]int, 3)
		lits[0] = b.XLit(r, 0, "Up")
		lits[1] = b.XLit(r, 0, "Right")
		lits[2] = b.XLit(r, 0, "Down")
		cnf.addFormula(filterZero(lits), builder)
		lits = make([]int, 1)
		lits[0] = -b.XLit(r, 0, "Left")
		cnf.addFormula(filterZero(lits), builder)
	}
	// Right
	for r := 1; r < b.Row-1; r++ {
		if b.Lookup[(r+1)*b.Column-1] == 0 {
			continue
		}
		lits := make([]int, 3)
		lits[0] = b.XLit(r, b.Column-1, "Up")
		lits[1] = b.XLit(r, b.Column-1, "Left")
		lits[2] = b.XLit(r, b.Column-1, "Down")
		cnf.addFormula(filterZero(lits), builder)
		lits = make([]int, 1)
		lits[0] = -b.XLit(r, b.Column-1, "Right")
		cnf.addFormula(filterZero(lits), builder)
	}
}
func buildCNFDirectionForUnNumberedBorderCell(cnf CNFInterface, builder CNFBuilder) {
	b := cnf.getBoard()
	// Top
	for c := 1; c < b.Column-1; c++ {
		if b.Lookup[c] != 0 {
			continue
		}
		lits := make([]int, 3)
		lits[0] = b.XLit(0, c, "Right")
		lits[1] = b.XLit(0, c, "Left")
		lits[2] = b.XLit(0, c, "Down")
		cnf.addFormula(filterZero(lits), builder)
		cnf.addClause([]int{-b.XLit(0, c, "Up")})
	}
	// Bottom
	for c := 1; c < b.Column-1; c++ {
		if b.Lookup[(b.Row-1)*b.Column+c] != 0 {
			continue
		}
		lits := make([]int, 3)
		lits[0] = b.XLit(b.Row-1, c, "Right")
		lits[1] = b.XLit(b.Row-1, c, "Left")
		lits[2] = b.XLit(b.Row-1, c, "Up")
		cnf.addFormula(filterZero(lits), builder)
		cnf.addClause([]int{-b.XLit(b.Row-1, c, "Down")})
	}
	// Left
	for r := 1; r < b.Row-1; r++ {
		if b.Lookup[r*b.Column] != 0 {
			continue
		}
		lits := make([]int, 3)
		lits[0] = b.XLit(r, 0, "Up")
		lits[1] = b.XLit(r, 0, "Right")
		lits[2] = b.XLit(r, 0, "Down")
		cnf.addFormula(filterZero(lits), builder)
		cnf.addClause([]int{-b.XLit(r, 0, "Left")})
	}
	// Right
	for r := 1; r < b.Row-1; r++ {
		if b.Lookup[(r+1)*b.Column-1] != 0 {
			continue
		}
		lits := make([]int, 3)
		lits[0] = b.XLit(r, b.Column-1, "Up")
		lits[1] = b.XLit(r, b.Column-1, "Left")
		lits[2] = b.XLit(r, b.Column-1, "Down")
		cnf.addFormula(filterZero(lits), builder)
		cnf.addClause([]int{-b.XLit(r, b.Column-1, "Right")})
	}
}

func buildCNFDirectionForNumberedInnerCell(cnf CNFInterface, builder CNFBuilder) {
	b := cnf.getBoard()
	for r := 1; r < b.Row-1; r++ {
		for c := 1; c < b.Column-1; c++ {
			if b.Lookup[b.Column*r+c] == 0 {
				continue
			}
			lits := make([]int, b.TotalMove)
			for a, j := range numberlink.MoveType {
				lits[j-1] = b.XLit(r, c, a)
			}
			cnf.addFormula(filterZero(lits), builder)
		}
	}
}
func buildCNFDirectionForUnNumberedInnerCell(cnf CNFInterface, builder CNFBuilder) {
	b := cnf.getBoard()
	for r := 1; r < b.Row-1; r++ {
		for c := 1; c < b.Column-1; c++ {
			if b.Lookup[b.Column*r+c] != 0 {
				continue
			}
			println(r, " ", c)
			lits := make([]int, b.TotalMove)
			for a, j := range numberlink.MoveType {
				lits[j-1] = b.XLit(r, c, a)
			}
			fmt.Printf("%v\n", lits)
			cnf.addFormula(filterZero(lits), builder)
		}
	}
}
func testExactly2(cnf CNFInterface, builder CNFBuilder) {
	cnf.addClause([]int{-1})
	cnf.addClause([]int{-2})
	cnf.addClause([]int{-3})
	cnf.addClause([]int{-4})
	cnf.addClause([]int{-1, -2, -3, -4})
	cnf.addClause([]int{1, 2, 3, 4})
	cnf.addFormula([]int{1, 2, 3, 4}, builder)
}
func buildCNFDirectionForUnNumberedBorderCell2(cnf CNFInterface, builder CNFBuilder) {
	b := cnf.getBoard()
	// Top
	for c := 1; c < b.Column-1; c++ {
		if b.Lookup[c] != 0 {
			continue
		}
		lits := make([]int, 3)
		lits[0] = -b.XLit(0, c, "Right")
		lits[1] = -b.XLit(0, c, "Left")
		lits[2] = -b.XLit(0, c, "Down")
		cnf.addFormula(filterZero(lits), builder)
		cnf.addClause([]int{-b.XLit(0, c, "Up")})
	}
	// Bottom
	for c := 1; c < b.Column-1; c++ {
		if b.Lookup[(b.Row-1)*b.Column+c] != 0 {
			continue
		}
		lits := make([]int, 3)
		lits[0] = -b.XLit(b.Row-1, c, "Right")
		lits[1] = -b.XLit(b.Row-1, c, "Left")
		lits[2] = -b.XLit(b.Row-1, c, "Up")
		cnf.addFormula(filterZero(lits), builder)
		cnf.addClause([]int{-b.XLit(b.Row-1, c, "Down")})
	}
	// Left
	for r := 1; r < b.Row-1; r++ {
		if b.Lookup[r*b.Column] != 0 {
			continue
		}
		lits := make([]int, 3)
		lits[0] = -b.XLit(r, 0, "Up")
		lits[1] = -b.XLit(r, 0, "Right")
		lits[2] = -b.XLit(r, 0, "Down")
		cnf.addFormula(filterZero(lits), builder)
		cnf.addClause([]int{-b.XLit(r, 0, "Left")})
	}
	// Right
	for r := 1; r < b.Row-1; r++ {
		if b.Lookup[(r+1)*b.Column-1] != 0 {
			continue
		}
		lits := make([]int, 3)
		lits[0] = -b.XLit(r, b.Column-1, "Up")
		lits[1] = -b.XLit(r, b.Column-1, "Left")
		lits[2] = -b.XLit(r, b.Column-1, "Down")
		cnf.addFormula(filterZero(lits), builder)
		cnf.addClause([]int{-b.XLit(r, b.Column-1, "Right")})
	}
}
func buildCNFConnectedValue(cnf CNFInterface, builder CNFBuilder) {
	b := cnf.getBoard()
	for r := 0; r < b.Row; r++ {
		for c := 0; c < b.Column; c++ {
			for v := 1; v <= b.MaxValue; v++ {
				if c-1 >= 0 {
					clause := make([]int, 3)
					clause[0] = -b.YLit(r, c, v)
					clause[1] = -b.XLit(r, c, "Left")
					clause[2] = b.YLit(r, c-1, v)
					cnf.addClause(clause)
					clause = make([]int, 3)
					clause[0] = -b.YLit(r, c-1, v)
					clause[1] = -b.XLit(r, c, "Left")
					clause[2] = b.YLit(r, c, v)
					cnf.addClause(clause)
					clause = make([]int, 2)
					clause[0] = -b.XLit(r, c-1, "Right")
					clause[1] = b.XLit(r, c, "Left")
					cnf.addClause(clause)
				}
				if c+1 < b.Column {
					clause := make([]int, 3)
					clause[0] = -b.YLit(r, c, v)
					clause[1] = -b.XLit(r, c, "Right")
					clause[2] = b.YLit(r, c+1, v)
					cnf.addClause(clause)
					clause = make([]int, 3)
					clause[0] = -b.YLit(r, c+1, v)
					clause[1] = -b.XLit(r, c, "Right")
					clause[2] = b.YLit(r, c, v)
					cnf.addClause(clause)
					clause = make([]int, 2)
					clause[0] = -b.XLit(r, c+1, "Left")
					clause[1] = b.XLit(r, c, "Right")
					cnf.addClause(clause)
				}
				if r+1 < b.Row {
					clause := make([]int, 3)
					clause[0] = -b.YLit(r, c, v)
					clause[1] = -b.XLit(r, c, "Down")
					clause[2] = b.YLit(r+1, c, v)
					cnf.addClause(clause)
					clause = make([]int, 3)
					clause[0] = -b.YLit(r+1, c, v)
					clause[1] = -b.XLit(r, c, "Down")
					clause[2] = b.YLit(r, c, v)
					cnf.addClause(clause)
					clause = make([]int, 2)
					clause[0] = -b.XLit(r+1, c, "Up")
					clause[1] = b.XLit(r, c, "Down")
					cnf.addClause(clause)
				}
				if r-1 >= 0 {
					clause := make([]int, 3)
					clause[0] = -b.YLit(r, c, v)
					clause[1] = -b.XLit(r, c, "Up")
					clause[2] = b.YLit(r-1, c, v)
					cnf.addClause(clause)
					clause = make([]int, 3)
					clause[0] = -b.YLit(r-1, c, v)
					clause[1] = -b.XLit(r, c, "Up")
					clause[2] = b.YLit(r, c, v)
					cnf.addClause(clause)
					clause = make([]int, 2)
					clause[0] = -b.XLit(r-1, c, "Down")
					clause[1] = b.XLit(r, c, "Up")
					cnf.addClause(clause)
				}
			}
		}
	}
}

func filterZero(slice []int) []int {
	//newSlice := slice[:0]
	//for _, x := range slice {
	//	if x != 0 {
	//		newSlice = append(newSlice, x)
	//	}
	//}
	return slice
}
