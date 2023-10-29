package handler_test

import (
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/mager/bluedot/handler"
	"go.uber.org/zap"
)

func TestGetDataset(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"TestGetDataset"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mock GetDataset
			log := zap.NewExample().Sugar()
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()
			h := handler.NewGetDataset(log, db)

			mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, name, email, image, slug FROM "User" WHERE slug=$1`)).
				WithArgs("mager").
				WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email", "image", "slug"}).
					AddRow("1", "Mager", "testemail", "testimage", "mager"))

			mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, "userId", name, slug, source, description, created, updated FROM "Dataset" WHERE "userId" = '1' AND slug = 'iris'`)).
				WillReturnRows(sqlmock.NewRows([]string{"id", "userId", "name", "slug", "source", "description", "created", "updated"}).
					AddRow("1", "1", "iris", "iris", "testsource", "testdescription", "2020-01-01", "2020-01-01"))

			// TODO: Fix error:
			//     --- FAIL: TestGetDataset/TestGetDataset (0.00s)
			// panic: sql: Scan error on column index 6, name "created": unsupported Scan, storing driver.Value type string into type *time.Time [recovered]

			// Mock request
			req := httptest.NewRequest("GET", "/datasets/mager/iris", nil)

			// Mock the handler's ServeHTTP method
			rr := httptest.NewRecorder()
			h.ServeHTTP(rr, req)
		})
	}
}
