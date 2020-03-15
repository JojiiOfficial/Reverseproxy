package units

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/JojiiOfficial/gaw"
)

// Datasize represents a unit of data size (in bits, bit)
type Datasize float32

// ...
const (
	// base 10 (SI prefixes)
	Bit Datasize = 1e0

	Byte     = Bit * 8
	Kilobyte = Byte * 1e3
	Megabyte = Byte * 1e6
	Gigabyte = Byte * 1e9
	Terabyte = Byte * 1e12
	Petabyte = Byte * 1e15
	Exabyte  = Byte * 1e18
)

// Bits returns the datasize in bit
func (b Datasize) Bits() float64 {
	return float64(b)
}

// Bytes returns the datasize in B
func (b Datasize) Bytes() float64 {
	return float64(b / Byte)
}

// Kilobytes returns the datasize in kB
func (b Datasize) Kilobytes() float64 {
	return float64(b / Kilobyte)
}

// Megabytes returns the datasize in MB
func (b Datasize) Megabytes() float64 {
	return float64(b / Megabyte)
}

// Gigabytes returns the datasize in GB
func (b Datasize) Gigabytes() float64 {
	return float64(b / Gigabyte)
}

// Terabytes returns the datasize in TB
func (b Datasize) Terabytes() float64 {
	return float64(b / Terabyte)
}

// Petabytes returns the datasize in PB
func (b Datasize) Petabytes() float64 {
	return float64(b / Petabyte)
}

// Exabytes returns the datasize in EB
func (b Datasize) Exabytes() float64 {
	return float64(b / Exabyte)
}

func (b Datasize) String() string {
	if b > Exabyte {
		return fmt.Sprintf("%dEB", int(b.Exabytes()))
	}

	if b > Petabyte {
		return fmt.Sprintf("%dPB", int(b.Petabytes()))
	}

	if b > Terabyte {
		return fmt.Sprintf("%dTB", int(b.Terabytes()))
	}

	if b > Gigabyte {
		return fmt.Sprintf("%dGB", int(b.Gigabytes()))
	}

	if b > Megabyte {
		return fmt.Sprintf("%dMB", int(b.Megabytes()))
	}

	if b > Kilobyte {
		return fmt.Sprintf("%dKB", int(b.Kilobytes()))
	}

	return fmt.Sprintf("%dB", int(b.Bytes()))
}

// ParseDatasize parses datasize
func ParseDatasize(str string) (float32, error) {
	numBuff := ""
	runes := []rune(str)
	i := 0

	for _, r := range runes {
		s := string(r)

		if isInt(s) {
			numBuff += s
		} else {
			break
		}

		i++
	}

	if len(numBuff) == 0 {
		return 0, errors.New("No number provided")
	}

	unit := strings.ToLower(string(runes[i:]))
	num, _ := strconv.Atoi(numBuff)

	if !gaw.IsInStringArray(unit, []string{"b", "kb", "mb", "gb", "tb", "pb", "eb"}) {
		return 0, errors.New("Invaild unit")
	}

	switch unit {
	case "b":
		return float32(Datasize(num) * Byte), nil
	case "kb":
		return float32(Datasize(num) * Kilobyte), nil
	case "mb":
		return float32(Datasize(num) * Megabyte), nil
	case "gb":
		return float32(Datasize(num) * Gigabyte), nil
	case "tb":
		return float32(Datasize(num) * Terabyte), nil
	case "pb":
		return float32(Datasize(num) * Petabyte), nil
	case "eb":
		return float32(Datasize(num) * Exabyte), nil
	}

	return 0, nil
}

func isInt(text string) bool {
	_, err := strconv.Atoi(text)
	return err == nil
}

//UnmarshalText unmashal data
func (b *Datasize) UnmarshalText(text []byte) error {
	duration, err := ParseDatasize(string(text))
	if err == nil {
		*b = Datasize(duration)
	}
	return err
}

// MarshalText implements encoding.TextMarshaler
func (b Datasize) MarshalText() ([]byte, error) {
	return []byte(b.String()), nil
}
