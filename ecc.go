package qrcode

type Polynomial struct {
	term []gfElement
}

func newPolynomial(data []byte) Polynomial {
	result := Polynomial{term: make([]gfElement, len(data))}
	for j := 0; j < len(data); j += 1 {
		result.term[len(data) - j - 1] = gfElement(data[j])
	}
	return result
}

func newMonomial(term gfElement, degree int) Polynomial {
	if term == gfZero {
		return Polynomial{}
	}
	result := Polynomial{term: make([]gfElement, degree+1)}
	result.term[degree] = term
	return result
}

func (p Polynomial) data(numTerms int) []byte {
	result := make([]byte, numTerms)

	for j := len(p.term) - 1; j >= 0 && j < numTerms; j-- {
		result[numTerms - j - 1] = byte(p.term[j])
	}

	return result
}

func (p Polynomial) numTerms() int {
	return len(p.term)
}

func (p Polynomial) normalised() Polynomial {
	numTerms := p.numTerms()
	maxNonzeroTerm := numTerms - 1

	for i := numTerms - 1; i >= 0; i-- {
		if p.term[i] != 0 {
			break
		}

		maxNonzeroTerm = i - 1
	}

	if maxNonzeroTerm < 0 {
		return Polynomial{}
	} else if maxNonzeroTerm < numTerms-1 {
		p.term = p.term[0 : maxNonzeroTerm+1]
	}

	return p
}

func polyAdd(a, b Polynomial) Polynomial {
	numATerms := a.numTerms()
	numBTerms := b.numTerms()

	numTerms := numATerms
	if numBTerms > numTerms {
		numTerms = numBTerms
	}

	result := Polynomial{term: make([]gfElement, numTerms)}

	for i := 0; i < numTerms; i++ {
		switch {
		case numATerms > i && numBTerms > i:
			result.term[i] = a.term[i] ^ b.term[i]
		case numATerms > i:
			result.term[i] = a.term[i]
		default:
			result.term[i] = b.term[i]
		}
	}

	return result.normalised()
}

func polyMultiply(a, b Polynomial) Polynomial {
	numATerms := a.numTerms()
	numBTerms := b.numTerms()

	result := Polynomial{term: make([]gfElement, numATerms+numBTerms)}

	for i := 0; i < numATerms; i++ {
		for j := 0; j < numBTerms; j++ {
			if a.term[i] != 0 && b.term[j] != 0 {
				monomial := newMonomial(gfMultiply(a.term[i], b.term[j]), i + j)
				result = polyAdd(result, monomial)
			}
		}
	}

	return result.normalised()
}

func polyRemainder(numerator, denominator Polynomial) Polynomial {
	remainder := numerator

	for remainder.numTerms() >= denominator.numTerms() {
		degree := remainder.numTerms() - denominator.numTerms()
		coefficient := gfDivide(remainder.term[remainder.numTerms()-1],
			denominator.term[denominator.numTerms()-1])

		divisor := polyMultiply(denominator,
			newMonomial(coefficient, degree))

		remainder = polyAdd(remainder, divisor)
	}

	return remainder.normalised()
}

// ISO/IEC 18004 table 9 specifies the numECBytes required.
func getErrorCorrection(data []byte, numECBytes int) []byte {
	ecpoly := newPolynomial(data)
	ecpoly = polyMultiply(ecpoly, newMonomial(gfOne, numECBytes))

	generator := Polynomial{term: []gfElement{1}}
	for i := 0; i < numECBytes; i++ {
		nextPoly := Polynomial{term: []gfElement{gfExpTable[i], 1}}
		generator = polyMultiply(generator, nextPoly)
	}

	remainder := polyRemainder(ecpoly, generator)
	return remainder.data(numECBytes)
}
