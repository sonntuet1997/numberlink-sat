package sudokusolver

func cnfAtLeast1(c CNFInterface, lits []int) [][]int {
	for _, lit := range lits {
		if exists := c.lookup(lit); exists {
			return [][]int{}
		}
	}

	return [][]int{lits}
}

func cnfAtMost1(c CNFInterface, lits []int) [][]int {
	return _cnfAtMost1(c, lits, false)
}

func _cnfAtMost1(c CNFInterface, lits []int, pairwise bool) [][]int {
	filteredLits := make([]int, 0, len(lits))
	satisfyingLit := 0
	for _, lit := range lits {

		if !c.lookup(-lit) {
			filteredLits = append(filteredLits, lit)
		} else if c.lookup(lit) {
			satisfyingLit = lit
		}
	}

	if satisfyingLit != 0 {
		for _, lit := range lits {
			if satisfyingLit == lit || c.lookup(-lit) {
				continue
			}
			c.addLit(-lit)
		}
		return [][]int{}
	}

	if pairwise {
		return cnfAtMost1Pairwise(c, filteredLits)
	}

	// return cnfAtMost1Commander(c, filteredLits)
	return cnfAtMost1Bimander(c, filteredLits)
}

func cnfExactly1(c CNFInterface, lits []int) [][]int {
	return append(cnfAtMost1(c, lits), cnfAtLeast1(c, lits)...)
}

func cnfAtMost1Pairwise(c CNFInterface, lits []int) [][]int {
	n := len(lits)
	result := make([][]int, 0, n*n/2)
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			// each pair can't be true at the same time
			result = append(result, []int{-lits[i], -lits[j]})
		}
	}
	return result
}

// Will Klieber and Gihwon Kwon. 2007.
// Efficient CNF Encoding for Selecting 1 from N Objects.
func cnfAtMost1Commander(c CNFInterface, lits []int) [][]int {
	if len(lits) <= 3 {
		return _cnfAtMost1(c, lits, true)
	}

	factor := 3
	m := (len(lits) + factor - 1) / factor
	result := make([][]int, 0, len(lits)*m)

	groups := make([][]int, m)
	for i := 0; i < m; i++ {
		groups[i] = lits[factor*i : min(factor*(i+1), len(lits))]
		// 1. At most one variable in a group can be true
		result = append(result, _cnfAtMost1(c, groups[i], true)...)
	}

	commanders := c.requestLiterals(m)

	//  2. If the commander variable of a group is false,
	// then none of the variables in the group can be true
	for i, commander := range commanders {
		for _, lit := range groups[i] {
			// -commander -> -lit
			result = append(result, cnfAtLeast1(c, []int{commander, -lit})...)
		}
	}

	// 3. Exactly one of the commander variables is true
	result = append(result, cnfAtLeast1(c, commanders)...)
	result = append(result, cnfAtMost1Commander(c, commanders)...)

	return result
}

// Nguyen, Van-Hau, and Son T. Mai. 2015.
// A new method to encode the at-most-one constraint into SAT.
func cnfAtMost1Bimander(c CNFInterface, lits []int) [][]int {
	factor := 2
	n := len(lits)
	m := (n + factor - 1) / factor
	binLength := getBinLength(m)
	result := make([][]int, 0, n*m)

	groups := make([][]int, m)
	for i := 0; i < m; i++ {
		groups[i] = lits[factor*i : min(factor*(i+1), n)]
		result = append(result, _cnfAtMost1(c, groups[i], true)...)
	}

	auxVars := c.requestLiterals(binLength)

	for i := 0; i < m; i++ {
		for _, lit := range groups[i] {
			for j, aux := range auxVars {
				commanderLit := aux
				if (i & (1 << j)) == 0 {
					commanderLit = -commanderLit
				}
				// lit -> commander
				result = append(result, []int{-lit, commanderLit})
			}
		}
	}

	return result
}

func getBinLength(m int) int {
	len := 1
	for m > len {
		len <<= 1
	}
	return len
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func makeRange(min, max int) []int {
	a := make([]int, max-min+1)
	for i := range a {
		a[i] = min + i
	}
	return a
}
