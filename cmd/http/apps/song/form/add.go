package form

import (
	"errors"
	"ffxvi-bard/domain/user"
	"ffxvi-bard/port/contract"
	"ffxvi-bard/port/dto"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
)

type NewSongFormView struct {
	EnsembleSize map[int]string
	Genres       []dto.DatabaseGenre
}

func NewAddNewSongFormView(song contract.SongInterface, genreRepository contract.GenreRepositoryInterface) (NewSongFormView, error) {
	form := NewSongFormView{}
	genres, err := genreRepository.FetchAll()
	if err != nil {
		return form, errors.New(fmt.Sprintf("failed to fetch available genres. Reason: %s", err))
	}
	form.Genres = genres
	form.EnsembleSize = song.GetDetailedEnsembleString()
	return form, nil
}

type SongFormSubmitted struct {
	Title        string
	Artist       string
	EnsembleSize int
	Genre        []int
	File         []byte
	User         user.User
}

func NewSongFormSubmitted(title string, artist string, ensembleSize string, genre []string, fileHeader *multipart.FileHeader, errorHandler contract.HttpErrorHandlerInterface, c *gin.Context) SongFormSubmitted {
	form := SongFormSubmitted{
		Title:  title,
		Artist: artist,
	}
	file, err := fileHeader.Open()
	if err != nil {
		errorHandler.RenderTemplate(err, http.StatusBadRequest, c)
		return form
	}
	defer file.Close()
	form.File, err = io.ReadAll(file)
	if err != nil {
		errorHandler.RenderTemplate(err, http.StatusBadRequest, c)
		return form
	}
	form.EnsembleSize, err = strconv.Atoi(ensembleSize)
	for _, genreStr := range genre {
		genreInt, err := strconv.Atoi(genreStr)
		if err != nil {
			errorHandler.RenderTemplate(err, http.StatusBadRequest, c)
			return form
		}
		form.Genre = append(form.Genre, genreInt)
	}
	storedUser, exists := c.Get("user")
	if exists && storedUser != nil {
		if userObj, ok := storedUser.(*user.User); ok { // Note the asterisk (*) indicating a pointer type
			form.User = *userObj // Dereference the pointer if you need the value type
		} else {
			errorHandler.RenderTemplate(errors.New("session user is not of the correct type"), http.StatusBadRequest, c)
			return form
		}
	}
	return form
}
