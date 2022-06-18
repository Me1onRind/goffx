package goffx

import (
	"errors"
	"fmt"
)

var (
	SequenceLengthNoMatch = errors.New("Sequence length is wrong")
	NonAlphabetCharacter  = errors.New("Non Alphabet Character")
)

type Codec struct {
	ffx   string
	radix int
	f     *FFX
}

type Sequence struct {
	Codec
	alphabet string
	length   int
	packMap  map[byte]int
}

func (s *Sequence) pack(v string) ([]int, error) {
	if len([]rune(v)) != s.length {
		return nil, fmt.Errorf("%w, Sequence length must be %d", SequenceLengthNoMatch, s.length)
	}

	res := make([]int, 0, len(v))
	for _, c := range []byte(v) {
		if i, ok := s.packMap[c]; ok {
			res = append(res, i)
		} else {
			return nil, fmt.Errorf("%w: %c", NonAlphabetCharacter, c)
		}
	}
	return res, nil
}

func (s *Sequence) unpack(v []int) string {
	res := make([]byte, 0, len(v))
	for _, i := range v {
		res = append(res, s.alphabet[i])
	}
	return string(res)
}

type String struct {
	Sequence
}

type Integer struct {
	String
}

func (i *Integer) pack(v int) ([]int, error) {
	vStr := fmt.Sprintf(fmt.Sprintf("%%0%dd", i.length), i.length)
	return i.String.pack(vStr)
}

func (i *Integer) Encrypt(v int) int {
	return i.String.unpack()
	return 0
}
