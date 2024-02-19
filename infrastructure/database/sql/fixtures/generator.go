package fixtures

import (
	"ffxvi-bard/port/dto"
	"fmt"
	"math/rand"
	"os"
	"reflect"
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

	users := make([]dto.DatabaseUser, count)
	songs := make([]dto.DatabaseSong, count)
	genres := make([]dto.DatabaseGenre, count)
	songGenres := make([]dto.DatabaseSongGenre, count)
	ratings := make([]dto.DatabaseRating, count)
	comments := make([]dto.DatabaseComment, count)

	generateFixtures(users, "user")
	generateFixtures(songs, "song")
	generateFixtures(genres, "genre")
	generateFixtures(songGenres, "song_genre")
	generateFixtures(ratings, "rating")
	generateFixtures(comments, "comment")
}

func generateFixtures(slice interface{}, tableName string) {
	sliceVal := reflect.ValueOf(slice)

	for i := 0; i < sliceVal.Len(); i++ {
		elem := sliceVal.Index(i).Addr().Interface()

		if err := faker.FakeData(elem); err != nil {
			fmt.Printf("Error faking data: %v\n", err)
			continue
		}

		if rating, ok := elem.(*dto.DatabaseRating); ok {
			rating.ID = i + 1
			rating.AuthorID = rand.Intn(sliceVal.Len()) + 1
			rating.SongID = rand.Intn(sliceVal.Len()) + 1
			rating.Rating = rand.Intn(9) + 1
		}
		if song, ok := elem.(*dto.DatabaseSong); ok {
			song.EnsembleSize = rand.Intn(7)
			song.Status = rand.Intn(4)
			song.Title = faker.Name()
			song.Artist = faker.Name()
		}

		if user, ok := elem.(*dto.DatabaseUser); ok {
			user.Username = faker.Name()
			user.Email = faker.Email()
		}
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
