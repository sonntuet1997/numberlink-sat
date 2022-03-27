build: cmd sudoku sudokusolver
	go build -o ./bin/sudokusolver ./cmd/sudokusolver

install: build
	go install ./...

test: sudoku sudokusolver
	go test -cover ./...

bench: sudokusolver
	go test -run=XXX -benchmem -bench=. ./sudokusolver

benchprofile:
	go test -run=XXX -benchmem -cpu-profile=cpu.prof -mem-profile=mem.prof -bench=. ./sudokusolver

profile: build
	./bin/sudokusolver -cpu-profile=cpu.prof -mem-profile=mem.prof ${ARGS}

help:
	go run ./cmd/sudokusolver/main.go -help
	
run:
	go run ./cmd/sudokusolver/main.go -cpu-profile=./logs/cpu.prof -mem-profile=./logs/mem.prof ${ARGS}
	
run-product:
	go run ./cmd/sudokusolver/main.go -cpu-profile=./logs/cpu.prof -mem-profile=./logs/mem.prof -algorithm=product
cnf:
	go run ./cmd/sudokusolver/main.go -cpu-profile=./logs/cpu.prof -mem-profile=./logs/mem.prof -cnf=true ${ARGS}
cnf-product:
	go run ./cmd/sudokusolver/main.go -cpu-profile=./logs/cpu.prof -mem-profile=./logs/mem.prof -cnf=true -algorithm=product
