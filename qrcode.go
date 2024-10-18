package qrcode

import (
	"bytes"
	"errors"
	"fmt"
)

func getVersion(data []byte) *version {
	var numByteCharCountBits int

	for _, v := range versions {
		if v.version < 10 {
			numByteCharCountBits = 8
		} else {
			numByteCharCountBits = 16
		}
		if v.numDataBits() >= len(data)*8+8+numByteCharCountBits {
			return &v
		}
	}

	return nil
}

func encodeContent(data []byte, version *version) []byte {
	encoded := NewBitset()
	encoded.Write(0b0100, 4)

	if version.version < 10 {
		encoded.Write(uint(len(data)), 8)
	} else {
		encoded.Write(uint(len(data)), 16)
	}

	for _, b := range data {
		encoded.Write(uint(b), 8)
	}

	encoded.Write(0, 4)

	for encoded.Length < version.numDataBits() {
		encoded.Write(0b11101100, 8)
		if encoded.Length < version.numDataBits() {
			encoded.Write(0b00010001, 8)
		}
	}

	return encoded.Bytes
}

func encodeBlocks(data []byte, version *version) *Bitset {
	result := NewBitset()

	content := make([][]byte, 0)
	ecc := make([][]byte, 0)
	start := 0

	for _, g := range version.groups {
		numECCodeWords := g.numCodewords - g.numDataCodewords
		for j := 0; j < g.numBlocks; j++ {
			end := start + g.numDataCodewords
			content = append(content, data[start:end])
			ecc = append(ecc, getErrorCorrection(data[start:end], numECCodeWords))
			start = end
		}
	}

	// Interleave the blocks.
	working := true
	for i := 0; working; i += 1 {
		working = false
		for _, c := range content {
			if i < len(c) {
				result.Write(uint(c[i]), 8)
				working = true
			}
		}
	}

	working = true
	for i := 0; working; i += 1 {
		working = false
		for _, c := range ecc {
			if i < len(c) {
				result.Write(uint(c[i]), 8)
				working = true
			}
		}
	}

	result.Write(0, version.numRemainderBits)

	return result
}

func terminal(b *bitmap, padding int) string {
	var buf bytes.Buffer
	// if there is an odd number of rows, just shorten the final margin
	for y := -padding; y+1 < b.size+padding; y += 2 {
		for x := -padding; x < b.size+padding; x += 1 {
			if y < 0 || x < 0 || y >= b.size || x >= b.size {
				buf.WriteString("█")
			} else if y+1 == b.size {
				if b.get(x, y) {
					buf.WriteString("▄")
				} else {
					buf.WriteString("█")
				}
			} else if b.get(x, y) == b.get(x, y+1) {
				if b.get(x, y) {
					buf.WriteString(" ")
				} else {
					buf.WriteString("█")
				}
			} else {
				if b.get(x, y) {
					buf.WriteString("▄")
				} else {
					buf.WriteString("▀")
				}
			}
		}
		buf.WriteString("\n")
	}
	return buf.String()
}

func Print(content string) error {
	output, err := GetString(content)

	if err != nil {
		return err
	}

	fmt.Print(output)
	return nil
}

func GetString(content string) (string, error) {
	bytes := []byte(content)

	version := getVersion(bytes)
	if version == nil {
		return "", errors.New("content too long to encode")
	}

	encodedContent := encodeContent(bytes, version)
	encodedBlocks := encodeBlocks(encodedContent, version)
	bitmap := render(encodedBlocks, version)

	return terminal(bitmap, 2), nil
}
