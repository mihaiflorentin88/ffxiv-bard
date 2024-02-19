package contract

import "ffxvi-bard/port/dto"

type MediaClientInterface interface {
	Search(track string, artist string) (dto.MediaResponse, error)
}
