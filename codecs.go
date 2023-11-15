package goffx

import (
	"errors"
	"fmt"
	"strconv"
)

var (
	SequenceLengthNoMatch = errors.New("Sequence length is wrong")
	NonAlphabetCharacter  = errors.New("Non Alphabet Character")
)

type codec struct {
	radix int
	ffx   *FFX
}

func newCodec(ffx string, radix int) *codec {
	return &codec{
		radix: radix,
		ffx:   NewFFX(ffx, 0),
	}
}

func newSequence(ffx, alphabet string, length int) *sequence {
	s := &sequence{
		codec:    *newCodec(ffx, len(alphabet)),
		alphabet: alphabet,
		length:   length,
		packMap:  make(map[byte]int, len(alphabet)),
	}
	for k, v := range []byte(alphabet) {
		s.packMap[v] = k
	}
	return s
}

type sequence struct {
	codec
	alphabet string
	length   int
	packMap  map[byte]int
}

func (s *sequence) pack(v string) ([]int, error) {
	if len(v) != s.length {
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

func (s *sequence) unpack(v []int) string {
	res := make([]byte, 0, len(v))
	for _, i := range v {
		res = append(res, s.alphabet[i])
	}
	return string(res)
}

type StringCodecs struct {
	sequence
}

func String(ffx, alphabet string, length int) *StringCodecs {
	return &StringCodecs{
		sequence: *newSequence(ffx, alphabet, length),
	}
}

func (s *StringCodecs) Encrypt(v string) (string, error) {
	packed, err := s.pack(v)
	if err != nil {
		return "", err
	}
	encrypted, err := s.ffx.encrypt(s.radix, packed)
	if err != nil {
		return "", err
	}
	return s.unpack(encrypted), nil
}

func (s *StringCodecs) Decrypt(v string) (string, error) {
	packed, err := s.pack(v)
	if err != nil {
		return "", err
	}
	encrypted, err := s.ffx.decrypt(s.radix, packed)
	if err != nil {
		return "", err
	}
	return s.unpack(encrypted), nil
}

type IntegerCodecs struct {
	StringCodecs
}

func Integer(ffx string, length int) *IntegerCodecs {
	return &IntegerCodecs{
		StringCodecs: *String(ffx, "0123456789", length),
	}
}

func (i *IntegerCodecs) pack(v int) ([]int, error) {
	vStr := fmt.Sprintf(fmt.Sprintf("%%0%dd", i.length), v)
	return i.StringCodecs.pack(vStr)
}

func (i *IntegerCodecs) Encrypt(v int) (int, error) {
	packed, err := i.pack(v)
	if err != nil {
		return 0, err
	}
	encrypted, err := i.ffx.encrypt(i.radix, packed)
	if err != nil {
		return 0, err
	}
	str := i.unpack(encrypted)
	return strconv.Atoi(str)
}

func (i *IntegerCodecs) Decrypt(v int) (int, error) {
	packed, err := i.pack(v)
	if err != nil {
		return 0, err
	}
	encrypted, err := i.ffx.decrypt(i.radix, packed)
	if err != nil {
		return 0, err
	}
	str := i.unpack(encrypted)
	return strconv.Atoi(str)
}
