package filmoteka

import (
	"log"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/famusovsky/VkTestTask/internal/filmoteka/postgres"
	"github.com/stretchr/testify/assert"
)

func TestCreateApp(t *testing.T) {
	t.Run("first case", func(t *testing.T) {
		app := CreateApp("addr", nil, nil, nil, false)
		assert.NotNil(t, app)
		assert.Equal(t, "addr", app.addr)
		assert.False(t, app.defAdmin)
		assert.Nil(t, app.dbHandler)
		assert.Nil(t, app.errorLog)
		assert.Nil(t, app.infoLog)
	})

	t.Run("second case", func(t *testing.T) {
		mockDB, _, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer mockDB.Close()
		var (
			addr         = "123"
			infoLog      = log.Logger{}
			errorLog     = log.Logger{}
			dbHandler, _ = postgres.GetHandler(mockDB, false)
			defAdmin     = true
		)
		app := CreateApp(addr, &infoLog, &errorLog, dbHandler, defAdmin)
		assert.NotNil(t, app)
		assert.Equal(t, addr, app.addr)
		assert.Equal(t, &infoLog, app.infoLog)
		assert.Equal(t, &errorLog, app.errorLog)
		assert.Equal(t, dbHandler, app.dbHandler)
		assert.True(t, app.defAdmin)
	})
}
