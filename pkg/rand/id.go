package rand

import (
	"crypto/rand"
	"fmt"
)

func GenerateTransactionID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return fmt.Sprintf("%x", b)
}
