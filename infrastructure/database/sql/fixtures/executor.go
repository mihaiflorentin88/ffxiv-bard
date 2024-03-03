package fixtures

import (
	database "ffxvi-bard/infrastructure/database/sql"
	"ffxvi-bard/port/contract"
	"github.com/RichardKnop/go-fixtures"
	"log"
)

type Fixture struct {
	driver contract.DatabaseDriverInterface
}

func NewFixtures(driver contract.DatabaseDriverInterface) Fixture {
	return Fixture{driver: driver}
}

func (f Fixture) GetFixtureFiles(directory string) []string {
	return []string{
		"infrastructure/database/sql/fixtures/files/user.yml",
		//"infrastructure/database/sql/fixtures/files/genre.yml",
		"infrastructure/database/sql/fixtures/files/song.yml",
		"infrastructure/database/sql/fixtures/files/song_genre.yml",
		"infrastructure/database/sql/fixtures/files/rating.yml",
		"infrastructure/database/sql/fixtures/files/comment.yml",
	}
}

func (f Fixture) Execute() {
	log.Println("Executing fixtures...")
	db := database.Instance
	fixtureFiles := f.GetFixtureFiles("infrastructure/database/sql/fixtures/files/")

	err := fixtures.LoadFiles(fixtureFiles, db, "sqlite")
	if err != nil {
		log.Fatal(err)
	}
	log.Print("DONE!")
}
