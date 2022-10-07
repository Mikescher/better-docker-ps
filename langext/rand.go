package langext

import (
	"crypto/rand"
	"io"
	"math"
	"math/big"
)

func RandBytes(size int) []byte {
	b := make([]byte, size)
	_, err := io.ReadFull(rand.Reader, b)
	if err != nil {
		panic(err)
	}
	return b
}

func RandBase62(rlen int) string {
	ecs := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

	bi52 := big.NewInt(int64(len(ecs)))

	randMax := big.NewInt(math.MaxInt64)

	r := ""

	for i := 0; i < rlen; i++ {
		v, err := rand.Int(rand.Reader, randMax)
		if err != nil {
			panic(err)
		}

		r += string(ecs[v.Mod(v, bi52).Int64()])
	}

	return r
}
