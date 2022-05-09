package solver_test

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/sonntuet1997/numberlink-sat/numberlink"
	"github.com/sonntuet1997/numberlink-sat/solver"
)

func TestMain(m *testing.M) {
	//log.SetOutput(ioutil.Discard)
	os.Exit(m.Run())
}

//
//func BenchmarkSolve5x5Gini(b *testing.B) {
//	for i := 0; i < b.N; i++ {
//		solveOneLiner(solver.Data5x5, "normal")
//	}
//}

func BenchmarkSolve5x5Cadical(b *testing.B) {
	for i := 0; i < b.N; i++ {
		customSolveOneLiner(solver.Data5x5, "/home/sv1/codes/cadical/build/cadical -q", "normal")
	}
}

//func BenchmarkSolve8x8Gini(b *testing.B) {
//	for i := 0; i < b.N; i++ {
//		solveOneLiner(solver.Data8x8, "normal")
//	}
//}
//
func BenchmarkSolve8x8Cadical(b *testing.B) {
	for i := 0; i < b.N; i++ {
		customSolveOneLiner(solver.Data8x8, "/home/sv1/codes/cadical/build/cadical -q", "normal")
	}
}

//
//func BenchmarkSolve10x10Gini(b *testing.B) {
//	for i := 0; i < b.N; i++ {
//		solveOneLiner(solver.Data10x10, "normal")
//	}
//}
//
func BenchmarkSolve10x10Cadical(b *testing.B) {
	for i := 0; i < b.N; i++ {
		customSolveOneLiner(solver.Data10x10, "/home/sv1/codes/cadical/build/cadical -q", "normal")
	}
}

//func BenchmarkSolve12x12Gini(b *testing.B) {
//	for i := 0; i < b.N; i++ {
//		solveOneLiner(solver.Data12x12, "normal")
//	}
//}
//
func BenchmarkSolve12x12Cadical(b *testing.B) {
	for i := 0; i < b.N; i++ {
		customSolveOneLiner(solver.Data12x12, "/home/sv1/codes/cadical/build/cadical -q", "normal")
	}
}

//func BenchmarkSolve15x15Gini(b *testing.B) {
//	for i := 0; i < b.N; i++ {
//		solveOneLiner(solver.Data15x15, "normal")
//	}
//}
//
func BenchmarkSolve15x15Cadical(b *testing.B) {
	for i := 0; i < b.N; i++ {
		customSolveOneLiner(solver.Data15x15, "/home/sv1/codes/cadical/build/cadical -q", "normal")
	}
}

//func BenchmarkSolve20x20Gini(b *testing.B) {
//	for i := 0; i < b.N; i++ {
//		solveOneLiner(solver.Data20x20, "normal")
//	}
//}

func BenchmarkSolve20x20Cadical(b *testing.B) {
	for i := 0; i < b.N; i++ {
		customSolveOneLiner(solver.Data20x20, "/home/sv1/codes/cadical/build/cadical -q", "normal")
	}
}

func solveOneLiner(input, algorithm string) string {
	board := numberlink.NewFromString(input)
	// solver.Solve(board)
	solver.SolveWithGini(board, algorithm)
	return ""
	//
	//var b bytes.Buffer
	//board.PrintOneLine(&b)
	//return strings.TrimSpace(b.String())
}

func customSolveOneLiner(input, satSolver, algorithm string) string {
	board := numberlink.NewFromString(input)
	solver.SolveWithCustomSolver(board, satSolver, algorithm)
	var b bytes.Buffer
	board.PrintOneLine(&b)
	return strings.TrimSpace(b.String())
}

//func solveMany(inputFile string) {
//	solveManyWithGini(inputFile)
//}
//
//func solveManyWithGophersat(inputFile string) {
//	file, _ := os.Open(inputFile)
//	solver.SolveManyGophersat(file, ioutil.Discard)
//}
//
//func solveManyWithGini(inputFile string) {
//	file, _ := os.Open(inputFile)
//	solver.SolveManyGini(file, ioutil.Discard)
//}
