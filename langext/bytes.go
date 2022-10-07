package langext

import (
	"fmt"
	"github.com/joomcode/errorx"
)

func BytesXOR(a []byte, b []byte) ([]byte, error) {
	if len(a) != len(b) {
		return nil, errorx.InternalError.New("length mismatch")
	}

	r := make([]byte, len(a))

	for i := 0; i < len(a); i++ {
		r[i] = a[i] ^ b[i]
	}

	return r, nil
}

func FormatBytes(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB", float64(b)/float64(div), "KMGTPE"[exp])
}
