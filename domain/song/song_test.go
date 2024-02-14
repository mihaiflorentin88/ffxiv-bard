package song

import (
	"ffxvi-bard/domain/user"
	"ffxvi-bard/mocks"
	"testing"
)

func getSongProcessorMock() {

}

func getFileSystemMock() *mocks.FilesystemMock {
	return &mocks.FilesystemMock{}
}

func getSongProcessorInterfaceMock() *mocks.MockSongProcessorInterface {
	return &mocks.MockSongProcessorInterface{}
}

func getSong() *song {
	return &song{
		storageID:     1,
		title:         "test",
		artist:        "test",
		ensembleSize:  Solo,
		fileCode:      "test",
		file:          []byte("test"),
		rating:        []Rating{},
		genre:         []Genre{},
		uploader:      user.User{},
		comments:      []contract.CommentInterface{},
		status:        Pending,
		statusMessage: "test",
		songProcessor: getSongProcessorInterfaceMock(),
		filesystem:    getFileSystemMock(),
	}
}

func TestNewSong(t *testing.T) {
	_, err := NewSong("", "artist", Solo, []Genre{}, []contract.CommentInterface{}, []byte{}, user.User{}, getSongProcessorInterfaceMock(), getFileSystemMock())
	if err == nil {
		t.Errorf("NewSong should return an error when title is empty")
	}
}

func TestSetTitle(t *testing.T) {
	title := "test"
	song := getSong()
	err := song.SetTitle("")
	if err == nil {
		t.Errorf("SetTitle should return an error when title is empty")
	}
	err = song.SetTitle(title)
	if err != nil {
		t.Errorf("Song should not return an error when title is not empty")
	}
	if song.GetTitle() != title {
		t.Errorf("Song should set title to %s", title)
	}
}

func TestSetArtist(t *testing.T) {
	artist := "test"
	song := getSong()
	err := song.SetArtist("")
	if err == nil {
		t.Errorf("SetArtist should return an error when artist is empty")
	}
	err = song.SetArtist(artist)
	if err != nil {
		t.Errorf("SetArtist should not return an error when artist is not empty")
	}
	if song.GetArtist() != artist {
		t.Errorf("SetArtist should set artist to %s", artist)
	}
}

func TestSetEnsembleSize(t *testing.T) {
	song := getSong()
	song.SetEnsembleSize(0)
	if song.ensembleSize != 0 {
		t.Errorf("SetEnsembleSize should set ensembleSize to 0")
	}
	song.SetEnsembleSize(Solo)
	if song.ensembleSize != Solo {
		t.Errorf("SetEnsembleSize should set ensembleSize to Solo")
	}
	song.SetEnsembleSize(Octet)
	if song.ensembleSize != Octet {
		t.Errorf("SetEnsembleSize should set ensembleSize to Octet")
	}
}

func TestSetGenre(t *testing.T) {
	song := getSong()
	genreRock := Genre{Name: "Rock"}
	genreClassical := Genre{Name: "Classical"}
	song.SetGenre([]Genre{})
	if len(song.genre) != 0 {
		t.Errorf("SetGenre should set genre to an empty slice")
	}
	song.SetGenre([]Genre{genreRock, genreClassical})
	if len(song.genre) != 2 {
		t.Errorf("SetGenre should set genre to a slice of length 2")
	}
	if song.genre[0] != genreRock {
		t.Errorf("SetGenre should set the first genre to genreRock")
	}
	if song.genre[1] != genreClassical {
		t.Errorf("SetGenre should set the second genre to genreClassical")
	}
}

func TestSetComments(t *testing.T) {
	song := getSong()
	comment1 := NewComment("title", "content", user.User{}, 0, 0)
	comment2 := NewComment("title", "content", user.User{}, 0, 0)
	song.SetComments([]contract.CommentInterface{})
	if len(song.comments) != 0 {
		t.Errorf("SetComments should set comments to an empty slice")
	}
	song.SetComments([]contract.CommentInterface{comment1, comment2})
	if len(song.comments) != 2 {
		t.Errorf("SetComments should set comments to a slice of length 2")
	}
	if song.comments[0] != comment1 {
		t.Errorf("SetComments should set the first comment to comment1")
	}
	if song.comments[1] != comment2 {
		t.Errorf("SetComments should set the second comment to comment2")
	}
}

func TestSetFile(t *testing.T) {
	song := getSong()
	song.SetFile([]byte{})
	if len(song.file) != 0 {
		t.Errorf("SetFile should set file to an empty slice")
	}
}

func TestGetStorageID(t *testing.T) {
	song := getSong()
	if song.GetStorageID() != 1 {
		t.Errorf("GetStorageID should return 1")
	}
}

func TestGetTitle(t *testing.T) {
	song := getSong()
	if song.GetTitle() != "test" {
		t.Errorf("GetTitle should return 'test'")
	}
}

func TestGetArtist(t *testing.T) {
	song := getSong()
	if song.GetArtist() != "test" {
		t.Errorf("GetArtist should return 'test'")
	}
}

func TestGetEnsembleSize(t *testing.T) {
	song := getSong()
	if song.GetEnsembleSize() != Solo {
		t.Errorf("GetEnsembleSize should return Solo")
	}
}

func TestGetGenre(t *testing.T) {
	song := getSong()
	if len(song.GetGenre()) != 0 {
		t.Errorf("GetGenre should return an empty slice")
	}
}

func TestGetUploader(t *testing.T) {
	song := getSong()
	user1 := user.User{}
	if song.GetUploader() != user1 {
		t.Errorf("GetUploader should return an empty user")
	}
}

func TestSetUploader(t *testing.T) {
	song := getSong()
	user1 := user.User{}
	song.SetUploader(user1)
	if song.uploader != user1 {
		t.Errorf("SetUploader should set uploader to user1")
	}
}

func TestAddComment(t *testing.T) {
	song := getSong()
	comment1 := NewComment("title", "content", user.User{}, 0, 0)
	comment1.SetStorageID(1)
	comment2 := NewComment("title", "content", user.User{}, 0, 0)
	comment2.SetStorageID(2)
	song.AddComment(comment1)
	song.AddComment(comment2)
	if len(song.comments) != 2 {
		t.Errorf("AddComment should add 2 comments to the slice")
	}
	if song.comments[0] != comment1 {
		t.Errorf("AddComment should add comment1 to the slice")
	}
	if song.comments[1] != comment2 {
		t.Errorf("AddComment should add comment2 to the slice")
	}
}

func TestRemoveComment(t *testing.T) {
	song := getSong()
	comment := NewComment("title", "content", user.User{}, 0, 0)
	song.AddComment(comment)
	song.RemoveComment(comment)
	if len(song.comments) != 0 {
		t.Errorf("RemoveComment should remove a comment from the slice")
	}
}

func TestGetComments(t *testing.T) {
	song := getSong()
	if len(song.GetComments()) != 0 {
		t.Errorf("GetComments should return an empty slice")
	}
}

func TestChangeStatus(t *testing.T) {
	song := getSong()
	song.ChangeStatus(Processing, "test")
	if song.status != Processing {
		t.Errorf("ChangeStatus should set status to Processing")
	}
	if song.statusMessage != "test" {
		t.Errorf("ChangeStatus should set statusMessage to 'test'")
	}
}

func TestGetStatus(t *testing.T) {
	song := getSong()
	if song.GetStatus() != Pending {
		t.Errorf("GetStatus should return Pending")
	}
}

func TestGetStatusMessage(t *testing.T) {
	song := getSong()
	if song.GetStatusMessage() != "test" {
		t.Errorf("GetStatusMessage should return 'test'")
	}
}

func TestProcessSong(t *testing.T) {
	song := getSong()
	err := song.ProcessSong()
	if err != nil {
		t.Errorf("ProcessSong should not return an error")
	}
}

func TestGetFileCode(t *testing.T) {
	song := getSong()
	if song.GetFileCode() != "test" {
		t.Errorf("GetFileCode should return 'test'")
	}
}

func TestGetFile(t *testing.T) {
	song := getSong()
	if string(song.GetFile()) != "test" {
		t.Errorf("GetFile should return 'test'")
	}
}

func TestGetAverageRating(t *testing.T) {
	song := getSong()
	if song.GetAverageRating() != 0 {
		t.Errorf("GetAverageRating should return 0")
	}
	song.AddRating(Rating{rating: 5})
	song.AddRating(Rating{rating: 7})
	song.AddRating(Rating{rating: 1})
	song.AddRating(Rating{rating: 8})
	song.AddRating(Rating{rating: 10})
	song.AddRating(Rating{rating: 10})
	result := song.GetAverageRating()
	if result != 6.83 {
		t.Errorf("GetAverageRating should return 6.83 not %v", result)
	}
}

func TestAddRating(t *testing.T) {
	song := getSong()
	song.AddRating(Rating{})
	if len(song.rating) != 1 {
		t.Errorf("AddRating should add a rating to the slice")
	}
}

func TestConstantAliases(t *testing.T) {
	if Pending != contract.Pending {
		t.Errorf("song.Pending (%d) is not equal to port.Pending (%d)", Pending, contract.Pending)
	}
	if Processing != contract.Processing {
		t.Errorf("song.Processing (%d) is not equal to port.Processing (%d)", Processing, contract.Processing)
	}
	if Processed != contract.Processed {
		t.Errorf("song.Processed (%d) is not equal to port.Processed (%d)", Processed, contract.Processed)
	}
	if Failed != contract.Failed {
		t.Errorf("song.Failed (%d) is not equal to port.Failed (%d)", Failed, contract.Failed)
	}
	if Deleted != contract.Deleted {
		t.Errorf("song.Deleted (%d) is not equal to port.Deleted (%d)", Deleted, contract.Deleted)
	}

	if Solo != contract.Solo {
		t.Errorf("song.Solo (%d) is not equal to port.Solo (%d)", Solo, contract.Solo)
	}
	if Duet != contract.Duet {
		t.Errorf("song.Duet (%d) is not equal to port.Duet (%d)", Duet, contract.Duet)
	}
	if Trio != contract.Trio {
		t.Errorf("song.Trio (%d) is not equal to port.Trio (%d)", Trio, contract.Trio)
	}
	if Quartet != contract.Quartet {
		t.Errorf("song.Quartet (%d) is not equal to port.Quartet (%d)", Quartet, contract.Quartet)
	}
	if Quintet != contract.Quintet {
		t.Errorf("song.Quintet (%d) is not equal to port.Quintet (%d)", Quintet, contract.Quintet)
	}
	if Sextet != contract.Sextet {
		t.Errorf("song.Sextet (%d) is not equal to port.Sextet (%d)", Sextet, contract.Sextet)
	}
	if Septet != contract.Septet {
		t.Errorf("song.Septet (%d) is not equal to port.Septet (%d)", Septet, contract.Septet)
	}
	if Octet != contract.Octet {
		t.Errorf("song.Octet (%d) is not equal to port.Octet (%d)", Octet, contract.Octet)
	}
}
