package main

import (
	"bufio"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"runtime/pprof"

	"github.com/sonntuet1997/numberlink-sat/numberlink"
	"github.com/sonntuet1997/numberlink-sat/solver"
)

var (
	isCNFMode bool
	//isSolveMode  bool
	//isManyMode bool
	cpuProfile string
	memProfile string
	algorithm  string
)

func init() {
	flag.BoolVar(&isCNFMode, "cnf", false, "Generate CNF")
	//flag.BoolVar(&isSolveMode, "solve", true, "Solve with SAT solver")
	//flag.BoolVar(&isManyMode, "many", false, "Solve many one-line 9x9 sudoku w/ gophersat")
	flag.StringVar(&algorithm, "algorithm", "normal", "Normal or Product algorithm")
	flag.StringVar(&cpuProfile, "cpu-profile", "", "Write CPU profile to a file")
	flag.StringVar(&memProfile, "mem-profile", "", "Write memory profile to a file")
	flag.Parse()
}

func main() {
	if cpuProfile != "" {
		f, err := os.Create(cpuProfile)
		if err != nil {
			log.Fatal(err)
		}
		err = pprof.StartCPUProfile(f)
		if err != nil {
			return
		}
		defer pprof.StopCPUProfile()
	}

	mode := "solve"
	if isCNFMode {
		mode = "cnf"
	}
	bytes, _ := ioutil.ReadAll(os.Stdin)
	input := string(bytes)
	solve(mode, algorithm, input)
	if memProfile != "" {
		f, err := os.Create(memProfile)
		if err != nil {
			log.Fatal("could not create memory profile: ", err)
		}
		defer func(f *os.File) {
			err := f.Close()
			if err != nil {

			}
		}(f) // error handling omitted for example
		runtime.GC() // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Fatal("could not write memory profile: ", err)
		}
	}
}

func solve(mode, algorithm, input string) {
	board := numberlink.NewFromString(input)

	if mode == "cnf" {
		cnf := solver.GenerateCNFConstraints(board, algorithm)
		writer := bufio.NewWriter(os.Stdout)
		cnf.Print(writer)
		err := writer.Flush()
		if err != nil {
			return
		}
		return
	}

	if mode == "solve" {
		println("cadical")
		solver.SolveWithCustomSolver(board, "cadical -q", algorithm)
	}

	board.Print(os.Stdout)
}
