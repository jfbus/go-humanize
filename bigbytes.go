package humanize

import (
	"fmt"
	"math/big"
	"strings"
	"unicode"
)

var (
	bigIECExp = big.NewInt(1024)

	// BigByte is one byte in bit.Ints
	BigByte = big.NewInt(1)
	// BigKiByte is 1,024 bytes in bit.Ints
	BigKiByte = (&big.Int{}).Mul(BigByte, bigIECExp)
	// BigMiByte is 1,024 k bytes in bit.Ints
	BigMiByte = (&big.Int{}).Mul(BigKiByte, bigIECExp)
	// BigGiByte is 1,024 m bytes in bit.Ints
	BigGiByte = (&big.Int{}).Mul(BigMiByte, bigIECExp)
	// BigTiByte is 1,024 g bytes in bit.Ints
	BigTiByte = (&big.Int{}).Mul(BigGiByte, bigIECExp)
	// BigPiByte is 1,024 t bytes in bit.Ints
	BigPiByte = (&big.Int{}).Mul(BigTiByte, bigIECExp)
	// BigEiByte is 1,024 p bytes in bit.Ints
	BigEiByte = (&big.Int{}).Mul(BigPiByte, bigIECExp)
	// BigZiByte is 1,024 e bytes in bit.Ints
	BigZiByte = (&big.Int{}).Mul(BigEiByte, bigIECExp)
	// BigYiByte is 1,024 z bytes in bit.Ints
	BigYiByte = (&big.Int{}).Mul(BigZiByte, bigIECExp)
)

var (
	bigSIExp = big.NewInt(1000)

	// BigSIByte is one SI byte in big.Ints
	BigSIByte = big.NewInt(1)
	// BigKByte is 1,000 SI bytes in big.Ints
	BigKByte = (&big.Int{}).Mul(BigSIByte, bigSIExp)
	// BigMByte is 1,000 SI k bytes in big.Ints
	BigMByte = (&big.Int{}).Mul(BigKByte, bigSIExp)
	// BigGByte is 1,000 SI m bytes in big.Ints
	BigGByte = (&big.Int{}).Mul(BigMByte, bigSIExp)
	// BigTByte is 1,000 SI g bytes in big.Ints
	BigTByte = (&big.Int{}).Mul(BigGByte, bigSIExp)
	// BigPByte is 1,000 SI t bytes in big.Ints
	BigPByte = (&big.Int{}).Mul(BigTByte, bigSIExp)
	// BigEByte is 1,000 SI p bytes in big.Ints
	BigEByte = (&big.Int{}).Mul(BigPByte, bigSIExp)
	// BigZByte is 1,000 SI e bytes in big.Ints
	BigZByte = (&big.Int{}).Mul(BigEByte, bigSIExp)
	// BigYByte is 1,000 SI z bytes in big.Ints
	BigYByte = (&big.Int{}).Mul(BigZByte, bigSIExp)
)

var bigBytesSizeTable = map[string]*big.Int{
	"b":   BigByte,
	"kib": BigKiByte,
	"kb":  BigKByte,
	"mib": BigMiByte,
	"mb":  BigMByte,
	"gib": BigGiByte,
	"gb":  BigGByte,
	"tib": BigTiByte,
	"tb":  BigTByte,
	"pib": BigPiByte,
	"pb":  BigPByte,
	"eib": BigEiByte,
	"eb":  BigEByte,
	"zib": BigZiByte,
	"zb":  BigZByte,
	"yib": BigYiByte,
	"yb":  BigYByte,
	// Without suffix
	"":   BigByte,
	"ki": BigKiByte,
	"k":  BigKByte,
	"mi": BigMiByte,
	"m":  BigMByte,
	"gi": BigGiByte,
	"g":  BigGByte,
	"ti": BigTiByte,
	"t":  BigTByte,
	"pi": BigPiByte,
	"p":  BigPByte,
	"ei": BigEiByte,
	"e":  BigEByte,
	"z":  BigZByte,
	"zi": BigZiByte,
	"y":  BigYByte,
	"yi": BigYiByte,
}

var ten = big.NewInt(10)

func humanateBigBytes(s, base *big.Int, sizes []string) string {
	if s.Cmp(ten) < 0 {
		return fmt.Sprintf("%dB", s)
	}
	c := (&big.Int{}).Set(s)
	val, mag := oomm(c, base, len(sizes)-1)
	suffix := sizes[mag]
	f := "%.0f%s"
	if val < 10 {
		f = "%.1f%s"
	}

	return fmt.Sprintf(f, val, suffix)

}

// BigBytes produces a human readable representation of an SI size.
//
// BigBytes(82854982) -> 83MB
func (h *BaseHumanizer) BigBytes(s *big.Int) string {
	sizes := []string{"B", "KB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB"}
	return humanateBigBytes(s, bigSIExp, sizes)
}

// BigBytes produces a human readable representation of an SI size.
//
// BigBytes(82854982) -> 83MB
func BigBytes(s *big.Int) string {
	return Default.BigBytes(s)
}

// BigIBytes produces a human readable representation of an IEC size.
//
// BigIBytes(82854982) -> 79MiB
func (h *BaseHumanizer) BigIBytes(s *big.Int) string {
	sizes := []string{"B", "KiB", "MiB", "GiB", "TiB", "PiB", "EiB", "ZiB", "YiB"}
	return humanateBigBytes(s, bigIECExp, sizes)
}

// BigIBytes produces a human readable representation of an IEC size.
//
// BigIBytes(82854982) -> 79MiB
func BigIBytes(s *big.Int) string {
	return Default.BigIBytes(s)
}

// ParseBigBytes parses a string representation of bytes into the number
// of bytes it represents.
//
// ParseBigBytes("42MB") -> 42000000, nil
// ParseBigBytes("42mib") -> 44040192, nil
func (h *BaseHumanizer) ParseBigBytes(s string) (*big.Int, error) {
	lastDigit := 0
	for _, r := range s {
		if !(unicode.IsDigit(r) || r == '.') {
			break
		}
		lastDigit++
	}

	val := &big.Rat{}
	_, err := fmt.Sscanf(s[:lastDigit], "%f", val)
	if err != nil {
		return nil, err
	}

	extra := strings.ToLower(strings.TrimSpace(s[lastDigit:]))
	if m, ok := bigBytesSizeTable[extra]; ok {
		mv := (&big.Rat{}).SetInt(m)
		val.Mul(val, mv)
		rv := &big.Int{}
		rv.Div(val.Num(), val.Denom())
		return rv, nil
	}

	return nil, fmt.Errorf("unhandled size name: %v", extra)
}

// ParseBigBytes parses a string representation of bytes into the number
// of bytes it represents.
//
// ParseBigBytes("42MB") -> 42000000, nil
// ParseBigBytes("42mib") -> 44040192, nil
func ParseBigBytes(s string) (*big.Int, error) {
	return Default.ParseBigBytes(s)
}
