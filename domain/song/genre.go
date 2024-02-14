package song

import "ffxvi-bard/domain/date"

type Genre struct {
	StorageID string
	Name      string
	Date      date.Date
}
