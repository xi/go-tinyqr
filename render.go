package qrcode

type bitmap struct {
	module [][]bool
	isUsed [][]bool
	size   int
}

func newBitmap(size int) *bitmap {
	var b bitmap

	b.module = make([][]bool, size)
	b.isUsed = make([][]bool, size)

	for i := range b.module {
		b.module[i] = make([]bool, size)
		b.isUsed[i] = make([]bool, size)
	}

	b.size = size

	return &b
}

func (b *bitmap) norm(x int, y int) (int, int) {
	return (x + b.size) % b.size, (y + b.size) % b.size
}

func (b *bitmap) get(x int, y int) bool {
	x, y = b.norm(x, y)
	return b.module[y][x]
}

func (b *bitmap) empty(x int, y int) bool {
	x, y = b.norm(x, y)
	return !b.isUsed[y][x]
}

func (b *bitmap) set(x int, y int, v bool) {
	x, y = b.norm(x, y)
	b.module[y][x] = v
	b.isUsed[y][x] = true
}

func (b *bitmap) rect(x0 int, y0 int, width int, height int, value bool) {
	for y := y0; y < y0+height; y++ {
		for x := x0; x < x0+width; x++ {
			b.set(x, y, value)
		}
	}
}

func (b *bitmap) renderFinderPatterns() {
	// Top left Finder Pattern.
	b.rect(0, 0, 9, 9, false)
	b.rect(0, 0, 7, 7, true)
	b.rect(1, 1, 5, 5, false)
	b.rect(2, 2, 3, 3, true)

	// Top right Finder Pattern.
	b.rect(-8, 0, 8, 9, false)
	b.rect(-7, 0, 7, 7, true)
	b.rect(-6, 1, 5, 5, false)
	b.rect(-5, 2, 3, 3, true)

	// Bottom left Finder Pattern.
	b.rect(0, -8, 9, 8, false)
	b.rect(0, -7, 7, 7, true)
	b.rect(1, -6, 5, 5, false)
	b.rect(2, -5, 3, 3, true)
}

func (b *bitmap) renderAlignmentPatterns(version *version) {
	for _, x := range version.alignmentPatternCenter {
		for _, y := range version.alignmentPatternCenter {
			if !b.empty(x, y) {
				continue
			}

			b.rect(x-2, y-2, 5, 5, true)
			b.rect(x-1, y-1, 3, 3, false)
			b.set(x, y, true)
		}
	}
}

func (b *bitmap) renderTimingPatterns() {
	value := true

	for i := 7 + 1; i < b.size-7; i++ {
		b.set(i, 7-1, value)
		b.set(7-1, i, value)
		value = !value
	}
}

func (b *bitmap) renderFormatInfo() {
	b.set(8, 1, true)
	b.set(8, 4, true)
	b.set(8, -5, true)
	b.set(8, -3, true)
	b.set(8, -1, true)

	b.set(-2, 8, true)
	b.set(-5, 8, true)
	b.set(4, 8, true)
	b.set(2, 8, true)
	b.set(0, 8, true)

	b.set(8, -8, true)
}

func (b *bitmap) renderVersionInfo(v *version) {
	if v.version < 7 {
		return
	}

	for i := 0; i < 18; i++ {
		b.set(i/3, -11+i%3, (v.bitSequence>>i)&1 == 1)
		b.set(-11+i%3, i/3, (v.bitSequence>>i)&1 == 1)
	}
}

func (b *bitmap) renderData(data *Bitset) {
	xOffset := 1
	up := true

	x := b.size - 2
	y := b.size - 1

	for i := 0; i < data.Length; i++ {
		mask := (y+x+xOffset)%2 == 0

		// != is equivalent to XOR.
		b.set(x+xOffset, y, mask != data.At(i))

		if i == data.Length-1 {
			break
		}

		// Find next free bit in the symbol.
		for {
			if xOffset == 1 {
				xOffset = 0
			} else {
				xOffset = 1

				if up {
					if y > 0 {
						y--
					} else {
						up = false
						x -= 2
					}
				} else {
					if y < b.size-1 {
						y++
					} else {
						up = true
						x -= 2
					}
				}
			}

			// Skip over the vertical timing pattern entirely.
			if x == 5 {
				x--
			}

			if b.empty(x+xOffset, y) {
				break
			}
		}
	}
}

func render(data *Bitset, version *version) *bitmap {
	b := newBitmap(version.bitmapSize())
	b.renderFinderPatterns()
	b.renderAlignmentPatterns(version)
	b.renderTimingPatterns()
	b.renderFormatInfo()
	b.renderVersionInfo(version)
	b.renderData(data)
	return b
}
