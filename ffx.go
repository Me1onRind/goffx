package goffx

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/binary"
	"fmt"
	"math/big"
	"strconv"
)

var (
	defaultRounds = 10
)

type FFX struct {
	key        string
	rounds     int
	digestmod  int
	digestSize int
}

func NewFFX(key string, rounds int) *FFX {
	if rounds < 0 {
		rounds = defaultRounds
	}
	return &FFX{
		key:    key,
		rounds: rounds,
	}
}

func (f *FFX) encrypt(radix int, v []int) (string, error) {
	a, b := f.split(v)

	for i := 0; i < f.rounds; i++ {
		roundedB, err := f.round(i, b)
		if err != nil {
			return "", err
		}
		c := f.add(radix, a, roundedB)
		a, b = b, c
	}

	tmpRes := append(a, b...)
	buf := make([]byte, 0, len(tmpRes)*5)
	for _, v := range tmpRes {
		buf = strconv.AppendInt(buf, int64(v), 10)
		buf = append(buf)
	}
	return string(buf), nil
}

func (f *FFX) add(radix int, a, b []int) []int {
	length := len(a)
	if length > len(b) {
		length = len(b)
	}

	var res = make([]int, 0, length)
	for i := 0; i < length; i++ {
		res = append(res, (a[i]+b[i])%radix)
	}
	return res
}

func (f *FFX) round(i int, v []int) ([]int, error) {
	list := make([]int, 0, len(v)+1)
	list = append(list, i)
	list = append(list, v...)
	packed, err := pack(v)
	if err != nil {
		return nil, fmt.Errorf("FFX pack failed:%w", err)
	}
	packed2, err := pack(v)
	if err != nil {
		return nil, fmt.Errorf("FFX pack2 failed:%w", err)
	}
	msg := string(packed) + string(packed2)

	h := hmac.New(sha1.New, []byte(f.key))
	_, _ = h.Write([]byte(msg))

	d := bigInt16Str(fmt.Sprintf("%x", h.Sum(nil)))
	s := d[len(d)-6:]
	var res []int
	for i := 0; i < 6; i++ {
		v := s[len(s)-1-i]
		n, err := strconv.Atoi(string(v))
		if err != nil {
			return nil, err
		}
		res = append(res, n)
	}
	return res, nil
}

func (f *FFX) split(v []int) ([]int, []int) {
	mid := len(v) / 2
	return v[:mid], v[mid:]
}

func pack(v []int) ([]byte, error) {
	var res []byte
	for _, value := range v {
		b, err := IntToBytes4(value)
		if err != nil {
			return nil, err
		}
		res = append(res, b...)
	}
	return res, nil
}

func IntToBytes4(n int) ([]byte, error) {
	m := int32(n)
	bytesBuffer := bytes.NewBuffer([]byte{})
	err := binary.Write(bytesBuffer, binary.BigEndian, m)
	if err != nil {
		return nil, err
	}

	data := bytesBuffer.Bytes()
	k := 4
	x := len(data)
	nb := make([]byte, k)
	for i := 0; i < k; i++ {
		nb[i] = data[x-i-1]
	}
	return nb, nil
}

func bigInt16Str(in string) string {
	i := new(big.Int)
	i.SetString(in, 16)
	return i.String()
}
