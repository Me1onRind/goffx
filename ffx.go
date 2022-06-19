package goffx

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/binary"
	"fmt"
	"math"
	"math/big"
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
	if rounds <= 0 {
		rounds = defaultRounds
	}
	return &FFX{
		key:        key,
		rounds:     rounds,
		digestSize: sha1.Size, // only support sha persent
	}
}

func (f *FFX) encrypt(radix int, v []int) ([]int, error) {
	a, b := f.split(v)

	for i := 0; i < f.rounds; i++ {
		roundedB, err := f.round(radix, i, b, len(a))
		if err != nil {
			return nil, err
		}
		c := f.add(radix, a, roundedB)
		a, b = b, c
	}

	return append(a, b...), nil
}

func (f *FFX) decrypt(radix int, v []int) ([]int, error) {
	a, b := f.split(v)
	var c []int

	for i := f.rounds - 1; i > -1; i-- {
		b, c = a, b
		roundedB, err := f.round(radix, i, b, len(c))
		if err != nil {
			return nil, err
		}
		a = f.sub(radix, c, roundedB)
		//fmt.Println(b, c, a)
	}

	return append(a, b...), nil
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

func (f *FFX) sub(radix int, a, b []int) []int {
	length := len(a)
	if length > len(b) {
		length = len(b)
	}

	var res = make([]int, 0, length)
	for i := 0; i < length; i++ {
		tmpRes := (a[i] - b[i]) % radix
		if tmpRes < 0 {
			tmpRes = radix + tmpRes
		}
		res = append(res, tmpRes)
	}
	return res
}

func (f *FFX) round(radix, i int, v []int, seqLength int) ([]int, error) {
	list := make([]int, 0, len(v)+1)
	list = append(list, i)
	list = append(list, v...)
	packed, err := pack(list)
	if err != nil {
		return nil, fmt.Errorf("FFX pack failed:%w", err)
	}

	charsPerHash := int(float64(f.digestSize) * (math.Log(256) / math.Log(float64(radix))))
	index := 0
	res := make([]int, 0, seqLength)
	for {
		packed2, err := pack([]int{index})
		if err != nil {
			return nil, fmt.Errorf("FFX pack2 failed:%w", err)
		}
		msg := string(packed) + string(packed2)
		h := hmac.New(sha1.New, []byte(f.key))
		_, _ = h.Write([]byte(msg))
		d := bigInt16Str(fmt.Sprintf("%x", h.Sum(nil)))
		rd := big.NewInt(int64(radix))
		r := new(big.Int)
		for i := 0; i < charsPerHash; i++ {
			d, r = d.DivMod(d, rd, r)
			res = append(res, int(r.Int64()))
			if len(res) == seqLength {
				return res, nil
			}
		}
		packed = h.Sum(nil)
		index++
	}
}

func (f *FFX) split(v []int) ([]int, []int) {
	mid := len(v) / 2
	return v[:mid], v[mid:]
}

func pack(v []int) ([]byte, error) {
	var res []byte
	for _, value := range v {
		b, err := intToByte4(value)
		if err != nil {
			return nil, err
		}
		res = append(res, b...)
	}
	return res, nil
}

func intToByte4(n int) ([]byte, error) {
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

func bigInt16Str(in string) *big.Int {
	i := new(big.Int)
	i.SetString(in, 16)
	return i
}
