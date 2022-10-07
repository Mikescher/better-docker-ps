package langext

import (
	"fmt"
	"time"
)

func FormatNaturalDurationEnglish(iv time.Duration) string {
	if sec := int64(iv.Seconds()); sec < 180 {
		if sec == 1 {
			return "1 second ago"
		} else {
			return fmt.Sprintf("%d seconds ago", sec)
		}
	}

	if min := int64(iv.Minutes()); min < 180 {
		return fmt.Sprintf("%d minutes ago", min)
	}

	if hours := int64(iv.Hours()); hours < 72 {
		return fmt.Sprintf("%d hours ago", hours)
	}

	if days := int64(iv.Hours() / 24.0); days < 21 {
		return fmt.Sprintf("%d days ago", days)
	}

	if weeks := int64(iv.Hours() / 24.0 / 7.0); weeks < 12 {
		return fmt.Sprintf("%d weeks ago", weeks)
	}

	if months := int64(iv.Hours() / 24.0 / 7.0 / 30); months < 36 {
		return fmt.Sprintf("%d months ago", months)
	}

	years := int64(iv.Hours() / 24.0 / 7.0 / 365)
	return fmt.Sprintf("%d years ago", years)
}
