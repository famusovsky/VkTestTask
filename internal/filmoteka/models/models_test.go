package models_test

import (
	"testing"
	"time"

	"github.com/famusovsky/VkTestTask/internal/filmoteka/models"
	"github.com/stretchr/testify/assert"
)

func TestActorInCheck(t *testing.T) {
	t.Run("valid actor", func(t *testing.T) {
		validActor := models.ActorIn{
			Name:        "John Doe",
			Gender:      "Male",
			DateOfBirth: time.Now().AddDate(-30, 0, 0), // 30 years ago
		}
		err := validActor.Check()
		assert.NoError(t, err)
	})

	t.Run("missing name", func(t *testing.T) {
		actorWithoutName := models.ActorIn{
			Gender:      "Female",
			DateOfBirth: time.Now().AddDate(-25, 0, 0), // 25 years ago
		}
		err := actorWithoutName.Check()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "name must not be null")
	})

	t.Run("missing genger", func(t *testing.T) {
		actorWithoutGender := models.ActorIn{
			Name:        "Jane Doe",
			DateOfBirth: time.Now().AddDate(-35, 0, 0), // 35 years ago
		}
		err := actorWithoutGender.Check()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "gender must not be null")
	})

	t.Run("missing date of birth", func(t *testing.T) {
		actorWithoutDOB := models.ActorIn{
			Name:   "James Smith",
			Gender: "Male",
		}
		err := actorWithoutDOB.Check()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "date of birth must not be null")
	})

	t.Run("all fields are missing", func(t *testing.T) {
		actorWithoutFields := models.ActorIn{}
		err := actorWithoutFields.Check()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "name must not be null")
		assert.Contains(t, err.Error(), "gender must not be null")
		assert.Contains(t, err.Error(), "date of birth must not be null")
	})
}

func TestMovieInCheck(t *testing.T) {
	t.Run("valid movie", func(t *testing.T) {
		movie := models.MovieIn{
			Name:        "Interstellar",
			Description: "A team of explorers travel through a wormhole in space in an attempt to ensure humanity's survival.",
			ReleaseDate: time.Now().AddDate(-7, 0, 0), // 7 years ago
			Rating:      new(int),
		}
		*movie.Rating = 8
		err := movie.Check()
		assert.NoError(t, err)
	})

	t.Run("missing name", func(t *testing.T) {
		movie := models.MovieIn{
			Description: "A team of explorers travel through a wormhole in space in an attempt to ensure humanity's survival.",
			ReleaseDate: time.Now().AddDate(-7, 0, 0), // 7 years ago
			Rating:      new(int),
		}
		*movie.Rating = 8
		err := movie.Check()

		expectError := "name must not be null"
		assert.Error(t, err)
		assert.Contains(t, err.Error(), expectError)
	})

	t.Run("name is too long", func(t *testing.T) {
		movie := models.MovieIn{
			Description: "A team of explorers travel through a wormhole in space in an attempt to ensure humanity's survival.",
			ReleaseDate: time.Now().AddDate(-7, 0, 0), // 7 years ago
			Rating:      new(int),
		}
		for i := 0; i <= 150; i++ {
			movie.Name += "a"
		}
		*movie.Rating = 8
		err := movie.Check()

		expectError := "movie name must be less than 150 chars"
		assert.Error(t, err)
		assert.Contains(t, err.Error(), expectError)
	})

	t.Run("missing description", func(t *testing.T) {
		movie := models.MovieIn{
			Name:        "Interstellar",
			ReleaseDate: time.Now().AddDate(-7, 0, 0), // 7 years ago
			Rating:      new(int),
		}
		*movie.Rating = 8
		err := movie.Check()

		expectError := "description must not be null"
		assert.Error(t, err)
		assert.Contains(t, err.Error(), expectError)
	})

	t.Run("name is too long", func(t *testing.T) {
		movie := models.MovieIn{
			Name:        "Interstellar",
			ReleaseDate: time.Now().AddDate(-7, 0, 0), // 7 years ago
			Rating:      new(int),
		}
		for i := 0; i <= 1000; i++ {
			movie.Description += "a"
		}
		*movie.Rating = 8
		err := movie.Check()

		expectError := "movie description must be less than 1000 chars"
		assert.Error(t, err)
		assert.Contains(t, err.Error(), expectError)
	})

	t.Run("missing release date", func(t *testing.T) {
		movie := models.MovieIn{
			Name:        "Interstellar",
			Description: "A team of explorers travel through a wormhole in space in an attempt to ensure humanity's survival.",
			Rating:      new(int),
		}
		*movie.Rating = 8
		err := movie.Check()

		expectError := "date of release must not be null"
		assert.Error(t, err)
		assert.Contains(t, err.Error(), expectError)
	})

	t.Run("missing rating", func(t *testing.T) {
		movie := models.MovieIn{
			Name:        "Interstellar",
			Description: "A team of explorers travel through a wormhole in space in an attempt to ensure humanity's survival.",
			ReleaseDate: time.Now().AddDate(-7, 0, 0), // 7 years ago
		}
		err := movie.Check()

		expectError := "rating must not be null"
		assert.Error(t, err)
		assert.Contains(t, err.Error(), expectError)
	})

	t.Run("rating is too low", func(t *testing.T) {
		movie := models.MovieIn{
			Name:        "Interstellar",
			Description: "A team of explorers travel through a wormhole in space in an attempt to ensure humanity's survival.",
			ReleaseDate: time.Now().AddDate(-7, 0, 0), // 7 years ago
			Rating:      new(int),
		}
		*movie.Rating = -1
		err := movie.Check()

		expectError := "rating must in range 0 - 10"
		assert.Error(t, err)
		assert.Contains(t, err.Error(), expectError)
	})

	t.Run("rating is too high", func(t *testing.T) {
		movie := models.MovieIn{
			Name:        "Interstellar",
			Description: "A team of explorers travel through a wormhole in space in an attempt to ensure humanity's survival.",
			ReleaseDate: time.Now().AddDate(-7, 0, 0), // 7 years ago
			Rating:      new(int),
		}
		*movie.Rating = 11
		err := movie.Check()

		expectError := "rating must in range 0 - 10"
		assert.Error(t, err)
		assert.Contains(t, err.Error(), expectError)
	})

	t.Run("two errors", func(t *testing.T) {
		movie := models.MovieIn{
			Description: "A team of explorers travel through a wormhole in space in an attempt to ensure humanity's survival.",
			Rating:      new(int),
		}
		*movie.Rating = 11
		err := movie.Check()

		expectError := "name must not be null"
		assert.Error(t, err)
		assert.Contains(t, err.Error(), expectError)

		expectError = "date of release must not be null"
		assert.Contains(t, err.Error(), expectError)
	})

	t.Run("all fields are missing", func(t *testing.T) {
		movie := models.MovieIn{}
		err := movie.Check()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "name must not be null")
		assert.Contains(t, err.Error(), "description must not be null")
		assert.Contains(t, err.Error(), "date of release must not be null")
		assert.Contains(t, err.Error(), "rating must not be null")
	})
}

func TestUserCheck(t *testing.T) {
	t.Run("valid user", func(t *testing.T) {
		user := models.User{
			Nickname: "JohnDoe",
			Password: "password",
		}
		err := user.Check()
		assert.NoError(t, err)
	})

	t.Run("missing nickname", func(t *testing.T) {
		user := models.User{
			Password: "password",
		}
		err := user.Check()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "name must not be null")
	})

	t.Run("missing password", func(t *testing.T) {
		user := models.User{
			Nickname: "JohnDoe",
		}
		err := user.Check()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "password must not be null")
	})

	t.Run("all fields are missing", func(t *testing.T) {
		user := models.User{}
		err := user.Check()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "name must not be null")
		assert.Contains(t, err.Error(), "password must not be null")
	})
}
