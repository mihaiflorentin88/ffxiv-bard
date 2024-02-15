package main

import (
	"ffxvi-bard/cmd/cli"
	"ffxvi-bard/container"
)

//func main() {
//	config.NewConfig("test", "/tmp/bard/unprocessed", "/tmp/bard/processed", 1000000, filesystem.LocalFileSystem{})
//	fs := filesystem.LocalFileSystem{}
//	file, _ := fs.ReadFile(path.Join("/tmp/bard/test_songs", "test.mid"))
//	genre := []song.Genre{song.Genre{Name: "classical"}}
//	midiProcessor := songprocessor.NewMidiProcessor("/tmp/bard/unprocessed", "/tmp/bard/processed", fs)
//	comment := song.NewComment("test", "test", user.User{}, 0, 0)
//	comments := []song.Comment{*comment}
//	var commentInterfaces []_interface.CommentInterface = make([]_interface.CommentInterface, len(comments))
//	for i, comment := range comments {
//		commentInterfaces[i] = &comment
//	}
//	testsong, err := song.NewSong("test", "test", song.Solo, genre, commentInterfaces, file, user.User{}, midiProcessor, fs)
//	if err != nil {
//		panic(err)
//	}
//	err = testsong.ProcessSong()
//	if err != nil {
//		panic(err)
//	}
//}

func main() {
	cli.Execute()

	connection, err := container.GetDatabaseDriver()
	if err != nil {
		return
	}
	defer connection.Close()

	//signals := make(chan os.Signal, 1)
	//signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	//<-signals
}
