package types

import "time"

type Metadata struct {
	Name        string
	Link        string
	Owner       string
	Size        int64
	Updated     time.Time
	ContentType string
	Created     time.Time
	Type        string
}
