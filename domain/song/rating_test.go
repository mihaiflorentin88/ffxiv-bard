package song

import (
	"ffxvi-bard/domain/user"
	"testing"
)

func getNewSongRating() *Rating {
	ranking, _ := NewSongRanking(getSong(), user.User{}, 5)
	return ranking
}

func TestNewSongRanking(t *testing.T) {
	_, err := NewSongRanking(getSong(), user.User{}, 11)
	if err == nil {
		t.Errorf("NewSongRanking should return an error when rating is greater than 10")
	}
}

func TestSetRanking(t *testing.T) {
	ranking := getNewSongRating()
	err := ranking.SetRanking(11)
	if err == nil {
		t.Errorf("SetRanking should return an error when rating is greater than 10")
	}
	err = ranking.SetRanking(5)
	if err != nil {
		t.Errorf("SetRanking should not return an error when rating is less than 10")
	}
	if ranking.GetRanking() != 5 {
		t.Errorf("SetRanking should set ranking to 5")
	}
}

func TestGetRanking(t *testing.T) {
	ranking := getNewSongRating()
	if ranking.GetRanking() != 5 {
		t.Errorf("GetRanking should return 5")
	}
}

func TestRankingsGetStorageID(t *testing.T) {
	ranking := getNewSongRating()
	if ranking.GetStorageID() != 0 {
		t.Errorf("GetStorageID should return 0")
	}
}

func TestSetStorageID(t *testing.T) {
	ranking := getNewSongRating()
	ranking.SetStorageID(1)
	if ranking.GetStorageID() != 1 {
		t.Errorf("SetStorageID should set storageID to 1")
	}
}

func TestGetSong(t *testing.T) {
	ranking := getNewSongRating()
	if song := ranking.GetSong(); song == nil {
		t.Errorf("GetSong should return a song")
	}
}

func TestGetAuthor(t *testing.T) {
	ranking := getNewSongRating()
	if ranking.GetAuthor() != (user.User{}) {
		t.Errorf("GetAuthor should return an author")
	}
}

func TestSetSong(t *testing.T) {
	ranking := getNewSongRating()
	ranking.SetSong(getSong())
	if ranking.GetSong() == nil {
		t.Errorf("SetSong should set song")
	}
}

func TestSetAuthor(t *testing.T) {
	ranking := getNewSongRating()
	rankingUser := user.User{Username: "test"}
	ranking.SetAuthor(rankingUser)
	result := ranking.GetAuthor()
	if result.Username != (rankingUser.Username) {
		t.Errorf("SetAuthor should set author")
	}
}

func TestLike(t *testing.T) {
	ranking := getNewSongRating()
	ranking.Like()
	if ranking.GetRanking() != 6 {
		t.Errorf("Like should increase ranking by 1")
	}
}

func TestDislike(t *testing.T) {
	ranking := getNewSongRating()
	ranking.Dislike()
	if ranking.GetRanking() != 4 {
		t.Errorf("Dislike should decrease ranking by 1")
	}
}
