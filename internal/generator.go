package generator

import (
	"errors"
	"math/rand"
	"strings"

	"github.com/ledex/passcard-generator/model"
)

const (
	UpperCase = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	LowerCase = "abcdefghijklmnopqrstuvwxyz"
	Numbers   = "0123456789"
	Special   = "!@#$%^&*()-_=+[]{}|;:,.<>?"
)

func GeneratePassCard(pci model.PassCardIdentifier, rows int, cols int) (*model.PassCard, error) {
	if pci.Version != 1 {
		return nil, errors.New("pass-card-identifier has invalid version")
	}

	var charset strings.Builder

	if getBit(pci.CharsetFlag, 0) {
		charset.WriteString(UpperCase)
	}
	if getBit(pci.CharsetFlag, 1) {
		charset.WriteString(LowerCase)
	}
	if getBit(pci.CharsetFlag, 2) {
		charset.WriteString(Numbers)
	}
	if getBit(pci.CharsetFlag, 3) {
		charset.WriteString(Special)
	}

	if charset.Len() == 0 {
		return nil, errors.New("no character set selected")
	}

	charsetString := charset.String()
	source := rand.NewSource(int64(pci.Seed.ID()))
	seededRand := rand.New(source)

	cardRows := make([]string, rows)
	for i := 0; i < rows; i++ {
		rowChars := make([]byte, cols)
		for j := range rowChars {
			rowChars[j] = charsetString[seededRand.Intn(len(charsetString))]
		}
		cardRows[i] = string(rowChars)
	}

	return &model.PassCard{
		rows,
		cols,
		&pci,
		cardRows,
	}, nil
}

func getBit(b byte, index uint) bool {
	return (b & (1 << index)) != 0
}
