package model

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

type PassCard struct {
	Rows int
	Cols int
	Pci  *PassCardIdentifier
	Data []string
}

type PassCardIdentifier struct {
	Version     int
	CharsetFlag byte
	Seed        uuid.UUID
}

func (pci *PassCardIdentifier) String() string {
	return fmt.Sprintf("v%d.%02x.%s", pci.Version, pci.CharsetFlag, pci.Seed.String())
}

func FromString(s string) (*PassCardIdentifier, error) {
	components := strings.Split(s, ".")
	if len(components) != 3 {
		return nil, errors.New("identifier contained more than three components")
	}

	// Version check
	if components[0][0] != 'v' {
		return nil, errors.New("the first component of the identifier has to start with 'v'")
	}

	version, err := strconv.Atoi(components[0][1:])
	if err != nil {
		return nil, errors.New("the first component of the identifier has to contin a valid integer after the 'v'")
	}

	charsetFlagArray, err := hex.DecodeString(components[1])
	if err != nil || len(charsetFlagArray) != 1 {
		return nil, errors.New("the second component of the identifier has to be a valid hex and not lager than '0xFF'")
	}

	seed, err := uuid.Parse(components[2])
	if err != nil {
		return nil, errors.New("the third component of the identifier has to be a valid UUID")
	}

	return &PassCardIdentifier{
		version,
		charsetFlagArray[0],
		seed,
	}, nil
}
