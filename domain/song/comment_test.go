package song

import (
	"ffxvi-bard/domain/user"
	"testing"
)

func getNewComment() *Comment {
	return NewComment("title", "content", user.User{Username: "test"}, 0, 0)
}

func TestNewComment(t *testing.T) {
	comment := getNewComment()
	if comment.Title != "title" {
		t.Errorf("NewComment should set title to 'title'")
	}
	if comment.Content != "content" {
		t.Errorf("NewComment should set content to 'content'")
	}
	if comment.Author.Username != "test" {
		t.Errorf("NewComment should set author to 'test'")
	}
	if comment.Likes != 0 {
		t.Errorf("NewComment should set likes to 0")
	}
	if comment.Dislikes != 0 {
		t.Errorf("NewComment should set dislikes to 0")
	}
}

func TestCommentLike(t *testing.T) {
	comment := getNewComment()
	comment.Like()
	if comment.Likes != 1 {
		t.Errorf("Like should increment likes by 1")
	}
}

func TestCommentDislike(t *testing.T) {
	comment := getNewComment()
	comment.Dislike()
	if comment.Dislikes != 1 {
		t.Errorf("Dislike should increment dislikes by 1")
	}
}

func TestCommentGetStorageID(t *testing.T) {
	comment := getNewComment()
	if comment.GetStorageID() != 0 {
		t.Errorf("GetStorageID should return 0")
	}
}

func TestCommentSetStorageID(t *testing.T) {
	comment := getNewComment()
	comment.SetStorageID(1)
	if comment.GetStorageID() != 1 {
		t.Errorf("SetStorageID should set storageID to 1")
	}
}

func TestCommentGetTitle(t *testing.T) {
	comment := getNewComment()
	if comment.Title != "title" {
		t.Errorf("GetTitle should return 'title'")
	}
}

func TestCommentGetContent(t *testing.T) {
	comment := getNewComment()
	if comment.Content != "content" {
		t.Errorf("GetContent should return 'content'")
	}
}

func TestCommentGetAuthor(t *testing.T) {
	comment := getNewComment()
	if comment.Author.Username != "test" {
		t.Errorf("GetAuthor should return 'test'")
	}
}

func TestCommentGetLikes(t *testing.T) {
	comment := getNewComment()
	if comment.Likes != 0 {
		t.Errorf("GetLikes should return 0")
	}
}

func TestCommentGetDislikes(t *testing.T) {
	comment := getNewComment()
	if comment.Dislikes != 0 {
		t.Errorf("GetDislikes should return 0")
	}
}

func TestCommentSetAuthor(t *testing.T) {
	comment := getNewComment()
	comment.Author = user.User{Username: "test2"}
	if comment.Author.Username != "test2" {
		t.Errorf("SetAuthor should set author to 'test2'")
	}
}

func TestCommentSetLikes(t *testing.T) {
	comment := getNewComment()
	comment.Likes = 1
	if comment.Likes != 1 {
		t.Errorf("SetLikes should set likes to 1")
	}
}

func TestCommentSetDislikes(t *testing.T) {
	comment := getNewComment()
	comment.Dislikes = 1
	if comment.Dislikes != 1 {
		t.Errorf("SetDislikes should set dislikes to 1")
	}
}

func TestCommentSetContent(t *testing.T) {
	comment := getNewComment()
	comment.Content = "newContent"
	if comment.Content != "newContent" {
		t.Errorf("SetContent should set content to 'newContent'")
	}
}

func TestCommentSetTitle(t *testing.T) {
	comment := getNewComment()
	comment.Title = "newTitle"
	if comment.Title != "newTitle" {
		t.Errorf("SetTitle should set title to 'newTitle'")
	}
}

func TestCommentGetID(t *testing.T) {
	comment := getNewComment()
	comment.SetStorageID(1)
	if comment.GetStorageID() != 1 {
		t.Errorf("GetID should return 1")
	}
}

func TestCommentSetID(t *testing.T) {
	comment := getNewComment()
	comment.SetStorageID(1)
	if comment.GetStorageID() != 1 {
		t.Errorf("SetID should set id to 1")
	}
}
