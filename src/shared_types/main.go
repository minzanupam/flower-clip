package shared_types

import "time"

type SVG struct {
	ID        int
	Name      string
	File      []byte
	CreatedAt time.Time
}
