package helpers

import "time"

// ParseTimestamp 文字列で渡されたタイムスタンプをtime.Timeに変換する。
func ParseTimestamp(ts string) (time.Time, error) {
	return time.Parse(time.RFC3339, ts)
}
