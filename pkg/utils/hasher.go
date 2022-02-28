package utils

import (
	"crypto/sha256"
	"encoding/hex"
)

// HashStrings concatenates data strings and cal
func HashStrings(data ...string) (hashedEncoded [sha256.Size * 2]byte, err error) {

	var hasher = sha256.New()
	var hashbuf [1024]byte

	for i := range data {
		j := 0
		for j < len(data[i]) {
			copied := copy(hashbuf[:], data[i][j:])
			j += copied

			_, err = hasher.Write(hashbuf[:copied])
			if err != nil {
				return
			}
		}
	}

	var hashedBytes [sha256.Size]byte
	hasher.Sum(hashedBytes[:0])

	hex.Encode(hashedEncoded[:], hashedBytes[:])

	return
}
