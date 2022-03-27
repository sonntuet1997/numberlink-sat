package sudokusolver

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/crillab/gophersat/explain"
	"github.com/crillab/gophersat/solver"
	"github.com/irifrance/gini"
	"github.com/irifrance/gini/z"
	"github.com/rkkautsar/sudoku-solver/sudoku"
)

func SolveWithGini(board *sudoku.Board, algorithm string) {
	//board.BasicSolve()
	g := gini.New()
	cnf := GenerateCNFConstraints(board, algorithm)
	start := time.Now()
	giniAddConstraints(g, cnf.getClauses())
	giniSolve(g, board)
	elapsed := time.Since(start)
	log.Printf("Adding Clauses and Solving took %s", elapsed)

}

func giniAddConstraints(g *gini.Gini, clauses [][]int) {
	for _, clause := range clauses {
		// log.Println("add clause", clause)
		for _, lit := range clause {
			g.Add(z.Dimacs2Lit(lit))
			// log.Println("add lit", lit)
		}
		g.Add(0)
	}
}

func giniSolve(g *gini.Gini, board *sudoku.Board) {
	// g.Write(os.Stdout)
	status := g.Solve()

	if status < 0 {
		ms := g.Why([]z.Lit{})
		log.Println(ms)
		panic("UNSAT")
	}
	model := make([]bool, board.NumCandidates)
	for i := 1; i <= len(model); i++ {
		model[i-1] = g.Value(z.Dimacs2Lit(i))
	}
	board.SolveWithModel(model)
}

func SolveWithCustomSolver(board *sudoku.Board, solver, algorithm string) {
	solverArgs := strings.Split(solver, " ")
	cmd := exec.Command(solverArgs[0], solverArgs[1:]...)
	stdin, _ := cmd.StdinPipe()
	stdout, _ := cmd.StdoutPipe()
	reader := bufio.NewScanner(stdout)
	writer := bufio.NewWriter(stdin)

	cmd.Start()
	defer cmd.Wait()
	board.BasicSolve()
	cnf := GenerateCNFConstraints(board, algorithm)
	cnf.Print(writer)
	writer.Flush()
	stdin.Close()

	model := make([]bool, board.NumCandidates)

	for reader.Scan() {
		line := reader.Text()

		if strings.HasPrefix(line, "s UNSATISFIABLE") {
			fmt.Println("UNSAT")
			return
		}

		if len(line) < 1 || !strings.HasPrefix(line, "v") {
			continue
		}

		values := strings.Split(line, " ")[1:]
		for _, val := range values {
			parsed, _ := strconv.Atoi(val)
			polarity := parsed > 0

			if parsed < 0 {
				parsed = -parsed
			}

			if parsed > 0 && parsed < len(model) {
				model[parsed-1] = polarity
			}
		}
	}

	board.SolveWithModel(model)
}

func ExplainUnsat(pb *solver.Problem) {
	fmt.Println("UNSAT")
	cnf := pb.CNF()

	unsatPb, err := explain.ParseCNF(strings.NewReader(cnf))
	if err != nil {
		panic(err)
	}

	mus, err := unsatPb.MUSDeletion()
	if err != nil {
		panic(err)
	}
	musCnf := mus.CNF()
	// Sort clauses so as to always have the same output
	lines := strings.Split(musCnf, "\n")
	sort.Sort(sort.StringSlice(lines[1:]))
	musCnf = strings.Join(lines, "\n")
	fmt.Println(musCnf)
}

func SolveManyGini(in io.Reader, out io.Writer, algorithm string) {
	scanner := bufio.NewScanner(in)
	writer := bufio.NewWriter(out)
	// base := GetBase9x9Clauses()
	// g := gini.New()
	// log.Println("new")
	// giniAddConstraints(g, base.getClauses())
	// log.Println("constraints")
	// board := base.Board

	// giniSolve(g, board)
	board := sudoku.New(3)

	// actLits := make([]z.Lit, 0, 81)

	for scanner.Scan() {
		// log.Println("start")
		input := scanner.Text()
		board.ReplaceWithSingleRowString(input, true)
		// board.BasicSolve()
		// board.NumCandidates = 729
		// // SolveWithGini(board)
		// // cnf := &CNF{Board: board, nbVar: base.nbVar}
		// actLits = actLits[:0]
		// for r := 0; r < 9; r++ {
		// 	for c := 0; c < 9; c++ {
		// 		for v := 1; v <= 9; v++ {
		// 			if board.Lookup[board.Idx(r, c)] == v {
		// 				g.Add(z.Dimacs2Lit(board.Lit(r, c, v)))
		// 				m := g.Activate()
		// 				// fmt.Println("activation", i, m.Dimacs())
		// 				actLits = append(actLits, m)
		// 			}
		// 		}
		// 	}
		// }
		// for i := 1; i <= 729; i++ {
		// 	if board.Lookup[(i-1)/9] == ((i-1)%9)+1 {
		// 		// log.Println("new:", i)
		// 		g.Add(z.Dimacs2Lit(i))
		// 		m := g.Activate()
		// 		// fmt.Println("activation", i, m.Dimacs())
		// 		actLits = append(actLits, m)
		// 	}
		// 	// else if !board.Candidates[i] {
		// 	// 	// log.Println("new:", -i)
		// 	// 	g.Add(z.Dimacs2Lit(-i))
		// 	// 	m := g.Activate()
		// 	// 	actLits = append(actLits, m)
		// 	// 	// fmt.Println("activation", -i, m.Dimacs())
		// 	// }
		// }
		// log.Println("assume")
		// log.Println(actLits)
		// g.Assume(actLits...)
		// giniSolve(g, board)
		// // log.Println("solve")
		// board.PrintOneLine(writer)
		// for _, m := range actLits {
		// 	g.Deactivate(m)
		// }
		// log.Println("end")

		board.ReplaceWithSingleRowString(input, false)
		SolveWithGini(board, algorithm)
		board.PrintOneLine(writer)
	}
	writer.Flush()
}
