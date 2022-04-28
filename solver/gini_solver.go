package solver

import (
	"bufio"
	"fmt"
	"github.com/crillab/gophersat/explain"
	"github.com/crillab/gophersat/solver"
	"github.com/go-air/gini"
	"github.com/go-air/gini/z"
	"github.com/sonntuet1997/numberlink-sat/numberlink"
	"log"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"
)

func SolveWithGini(board *numberlink.Board, algorithm string) string {
	//log.Printf("Unresolved cells %d", board.GetUnesolvedCells())
	g := gini.New()
	cnf := GenerateCNFConstraints(board, algorithm)
	log.Printf("Var length %d", cnf.varLen())
	log.Printf("Clauses length %d", cnf.clauseLen())
	start := time.Now()
	giniAddConstraints(g, cnf.getClauses())
	giniSolve(g, board)
	elapsed := time.Since(start)
	log.Printf("Adding Clauses and Solving took %s", elapsed)
	return strconv.FormatInt(elapsed.Nanoseconds(), 10)
}

func GetInfo(board *numberlink.Board, algorithm string) {
	//board.BasicSolve()
	//log.Printf("Unresolved cells %d", board.GetUnresolvedCells())
	//g := gini.New()
	cnf := GenerateCNFConstraints(board, algorithm)
	log.Printf("Var length %d", cnf.varLen())
	log.Printf("Clauses length %d", cnf.clauseLen())
	//start := time.Now()
	//giniAddConstraints(g, cnf.getClauses())
	////giniSolve(g, board)
	//elapsed := time.Since(start)
	//log.Printf("Adding Clauses and Solving took %s", elapsed)

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

func giniSolve(g *gini.Gini, board *numberlink.Board) {
	// g.Write(os.Stdout)
	status := g.Solve()
	if status < 0 {
		ms := g.Why(nil)
		log.Println("ms", ms)
		panic("UNSAT")
	}
	model := make([]bool, board.NumCandidates)
	println(len(model))
	for i := 1; i <= len(model); i++ {
		model[i-1] = g.Value(z.Dimacs2Lit(i))
	}
	board.SolveWithModel(model)
}

func SolveWithCustomSolver(board *numberlink.Board, solver, algorithm string) {
	solverArgs := strings.Split(solver, " ")
	cmd := exec.Command(solverArgs[0], solverArgs[1:]...)
	stdin, _ := cmd.StdinPipe()
	stdout, _ := cmd.StdoutPipe()
	reader := bufio.NewScanner(stdout)
	writer := bufio.NewWriter(stdin)

	cmd.Start()
	defer cmd.Wait()

	//log.Printf("Unresolved cells %d", board.GetUnresolvedCells())
	cnf := GenerateCNFConstraints(board, algorithm)
	log.Printf("Var length %d", cnf.varLen())
	log.Printf("Clauses length %d", cnf.clauseLen())
	start := time.Now()
	cnf.Print(writer)
	writer.Flush()
	stdin.Close()
	elapsed := time.Since(start)
	log.Printf("Adding Clauses and Solving took %s", elapsed)

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
