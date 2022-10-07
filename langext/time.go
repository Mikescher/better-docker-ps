package langext

import (
	"math"
	"time"
)

func UnixFloatSeconds(v float64) time.Time {
	sec, dec := math.Modf(v)
	return time.Unix(int64(sec), int64(dec*(1e9)))
}
