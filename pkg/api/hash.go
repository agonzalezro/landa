package api

import (
	"crypto/md5"
	"fmt"
	"io"
)

func hashCode(code string) string {
	h := md5.New()
	io.WriteString(h, code)
	return fmt.Sprintf("%x", h.Sum(nil))
}
