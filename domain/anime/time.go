package anime

import "time"

var now func() time.Time

func init() {
	now = time.Now().UTC
}

func SetNowFunc(newNow func() time.Time) {
	now = newNow
}
