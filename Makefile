build: cmd numberlink solver
	go build -o ./bin/numberlinksolver ./cmd/solver
	
build-server: cmd numberlink solver
	go build -o ./bin/server ./src

install: build
	go install ./...

test: numberlink solver
	go test -cover ./...

test-cnf: numberlink solver
	go test -cover -v ./solver/cnf_test.go


bench: solver
	go test -timeout=4h -run=XXX -benchmem -bench=. ./solver
	
benchprofile: solver
	go test -run=XXX -benchmem -cpu-profile=./cpu.prof -mem-profile=./mem.prof -bench=. ./solver

profile: build
	./bin/numberlinksolver -cpu-profile=cpu.prof -mem-profile=mem.prof ${ARGS}

help:
	go run ./cmd/solver/main.go -help
	
run:
	go run ./cmd/solver/main.go -cpu-profile=./logs/cpu.prof -mem-profile=./logs/mem.prof ${ARGS}
	
run-product:
	go run ./cmd/solver/main.go -cpu-profile=./logs/cpu.prof -mem-profile=./logs/mem.prof -algorithm=product
cnf:
	go run ./cmd/solver/main.go -cpu-profile=./logs/cpu.prof -mem-profile=./logs/mem.prof -cnf=true ${ARGS}
cnf-product:
	go run ./cmd/solver/main.go -cpu-profile=./logs/cpu.prof -mem-profile=./logs/mem.prof -cnf=true -algorithm=product


