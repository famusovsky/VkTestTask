package filmoteka

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/famusovsky/VkTestTask/internal/filmoteka/postgres"
	"github.com/stretchr/testify/assert"
)

func TestSendJson(t *testing.T) {
	t.Run("send json success", func(t *testing.T) {
		logger := log.New(io.Discard, "test", log.LstdFlags)
		app := CreateApp(":8080", logger, logger, nil, false)
		w := httptest.NewRecorder()
		obj := map[string]interface{}{
			"key": "value",
		}

		app.sendJson(w, obj)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, obj, response)
	})

	t.Run("send json error", func(t *testing.T) {
		logger := log.New(io.Discard, "test", log.LstdFlags)
		app := CreateApp(":8080", logger, logger, nil, false)
		w := httptest.NewRecorder()
		obj := make(chan int)
		defer close(obj)

		app.sendJson(w, obj)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Equal(t, "json: unsupported type: chan int\n", w.Body.String())
	})
}

func TestAuthIsAdmin(t *testing.T) {
	logger := log.New(io.Discard, "test", log.LstdFlags)
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()
	dbHandler, _ := postgres.GetHandler(mockDB, false)
	app := CreateApp(":8080", logger, logger, dbHandler, true)

	t.Run("default admin", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.SetBasicAuth("admin", "admin")
		isAdmin, err := app.authIsAdmin(req)
		assert.NoError(t, err)
		assert.True(t, isAdmin)
	})

	t.Run("non-admin", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.SetBasicAuth("user", "password")
		mock.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
		isAdmin, err := app.authIsAdmin(req)
		assert.NoError(t, err)
		assert.False(t, isAdmin)
	})

	t.Run("admin", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.SetBasicAuth("user", "password")
		mock.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
		isAdmin, err := app.authIsAdmin(req)
		assert.NoError(t, err)
		assert.True(t, isAdmin)
	})

	t.Run("error parsing basic auth", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		isAdmin, err := app.authIsAdmin(req)
		assert.Error(t, err)
		assert.False(t, isAdmin)
	})
}

func TestHandleError(t *testing.T) {
	t.Run("first case", func(t *testing.T) {
		var buf bytes.Buffer
		logger := log.New(&buf, "", log.LstdFlags)
		w := httptest.NewRecorder()
		errTxt := "An error occurred"
		status := http.StatusInternalServerError

		handleError(logger, w, errTxt, status)
		errTxt += "\n"

		assert.Contains(t, buf.String(), errTxt)
		assert.Equal(t, status, w.Code)
		assert.Equal(t, errTxt, w.Body.String())
	})

	t.Run("second case", func(t *testing.T) {
		var buf bytes.Buffer
		logger := log.New(&buf, "", log.LstdFlags)
		w := httptest.NewRecorder()
		errTxt := "Something went wrong"
		status := http.StatusForbidden

		handleError(logger, w, errTxt, status)
		errTxt += "\n"

		assert.Contains(t, buf.String(), errTxt)
		assert.Equal(t, status, w.Code)
		assert.Equal(t, errTxt, w.Body.String())
	})
}
