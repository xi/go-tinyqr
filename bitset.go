package qrcode

type Bitset struct {
	Length int
	Bytes  []byte
}

func NewBitset() *Bitset {
	return &Bitset{Length: 0, Bytes: make([]byte, 0)}
}

func (b *Bitset) Write(value uint, length int) {
	for i := length - 1; i >= 0; i-- {
		if b.Length/8 == len(b.Bytes) {
			b.Bytes = append(b.Bytes, 0)
		}

		if value&(1<<i) != 0 {
			b.Bytes[b.Length/8] |= 0x80 >> (b.Length % 8)
		}

		b.Length++
	}
}

func (b *Bitset) At(index int) bool {
	return b.Bytes[index/8]&(0x80>>(index%8)) != 0
}
