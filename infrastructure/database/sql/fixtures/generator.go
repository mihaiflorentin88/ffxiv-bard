package fixtures

import (
	"ffxvi-bard/port/dto"
	"fmt"
	"log"
	"math/rand"
	"os"
	"reflect"
	"strconv"
	"time"

	faker "github.com/go-faker/faker/v4"
	"gopkg.in/yaml.v3"
)

type FixtureEntry struct {
	Table  string                 `yaml:"table"`
	Pk     map[string]interface{} `yaml:"pk"`
	Fields map[string]interface{} `yaml:"fields"`
}

func GenerateFixtures(count int) {
	log.Println(fmt.Sprintf("Generating %v fixtures...", count))

	users := make([]dto.DatabaseUser, count)
	songs := make([]dto.DatabaseSong, count)
	//genres := make([]dto.DatabaseGenre, count)
	songGenres := make([]dto.DatabaseSongGenre, count)
	ratings := make([]dto.DatabaseRating, count)
	comments := make([]dto.DatabaseComment, count)

	generateFixtures(users, "user")
	generateFixtures(songs, "song")
	//generateFixtures(genres, "genre")
	generateFixtures(songGenres, "song_genre")
	generateFixtures(ratings, "rating")
	generateFixtures(comments, "comment")
}

func generateFixtures(slice interface{}, tableName string) {
	sliceVal := reflect.ValueOf(slice)

	for i := 0; i < sliceVal.Len(); i++ {
		log.Println(fmt.Sprintf("Generating fixture #%v...", i))
		elem := sliceVal.Index(i).Addr().Interface()

		if err := faker.FakeData(elem); err != nil {
			fmt.Printf("Error faking data: %v\n", err)
			continue
		}

		if user, ok := elem.(*dto.DatabaseUser); ok {
			user.ID = int64(i + 1)
			name := faker.Name()
			user.Username = faker.Name() + strconv.Itoa(i)
			user.Name = &name
			user.Email = faker.Email() + strconv.Itoa(i)
		}

		if song, ok := elem.(*dto.DatabaseSong); ok {
			song.ID = i + 1
			song.EnsembleSize = rand.Intn(7)
			song.Status = rand.Intn(4)
			song.Title = faker.Name()
			song.Artist = faker.Name()
			song.UploaderID = int64(i + 1)
		}

		if rating, ok := elem.(*dto.DatabaseRating); ok {
			rating.ID = i + 1
			rating.AuthorID = i + 1
			rating.SongID = i + 1
			rating.Rating = rand.Intn(9) + 1
		}

		if songGenre, ok := elem.(*dto.DatabaseSongGenre); ok {
			songGenre.SongID = i + 1
			songGenre.GenreID = rand.Intn(37) + 1
		}

		if comment, ok := elem.(*dto.DatabaseComment); ok {
			comment.SongID = i + 1
			comment.AuthorID = int64(i + 1)
		}
		log.Print("DONE!")
	}
	writeToFile(fmt.Sprintf("infrastructure/database/sql/fixtures/files/%s.yml", tableName), slice, tableName)
}

func writeToFile(filename string, data interface{}, tableName string) {
	var fixtureEntries []FixtureEntry

	reflectValue := reflect.ValueOf(data)
	for i := 0; i < reflectValue.Len(); i++ {
		fixtureEntries = append(fixtureEntries, structToFixtureEntry(reflectValue.Index(i).Interface(), tableName))
	}

	file, err := os.Create(filename)
	if err != nil {
		fmt.Printf("Error creating file %s: %v\n", filename, err)
		return
	}
	defer file.Close()

	encoder := yaml.NewEncoder(file)
	defer encoder.Close()

	if err := encoder.Encode(fixtureEntries); err != nil {
		fmt.Printf("Error encoding data to YAML for file %s: %v\n", filename, err)
	}
	log.Println(fmt.Sprintf("Finished writing `%s` to disk", filename))
}

func structToFixtureEntry(s interface{}, tableName string) FixtureEntry {
	val := reflect.ValueOf(s)
	typ := val.Type()

	entry := FixtureEntry{
		Table:  tableName,
		Pk:     make(map[string]interface{}),
		Fields: make(map[string]interface{}),
	}

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		dbTag := field.Tag.Get("db")

		if dbTag == "" {
			continue
		}
		fieldValue := val.Field(i).Interface()

		if field.Type == reflect.TypeOf(&time.Time{}) {
			if t, ok := fieldValue.(*time.Time); ok && t != nil {
				fieldValue = *t
			} else {
				fieldValue = time.Time{}
			}
		}

		if i == 0 {
			entry.Pk[dbTag] = fieldValue
		} else {
			entry.Fields[dbTag] = fieldValue
		}
	}
	return entry
}
