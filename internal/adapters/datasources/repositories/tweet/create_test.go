package tweet_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"challenge_be/internal/adapters/datasources/repositories/tweet"
	domain "challenge_be/internal/domain/tweet"
)

func TestCreate(t *testing.T) {
	type fields struct {
		db   *sql.DB
		mock sqlmock.Sqlmock
	}

	tests := map[string]struct {
		inputTweet domain.Tweet
		prepare    func(f *fields)
		expectErr  error
		expectID   uint64 // Agregamos este campo para verificar el ID
	}{
		"when tweet is created successfully": {
			inputTweet: domain.Tweet{
				UserID:  1,
				Content: "Hello World",
			},
			prepare: func(f *fields) {
				f.mock.ExpectQuery("INSERT INTO tweets").
					WithArgs(1, "Hello World").
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1)) // Simulamos que retorna ID 1
			},
			expectID: 1, // Esperamos que el ID sea 1
		},
		"when prepare statement fails": {
			inputTweet: domain.Tweet{
				UserID:  2,
				Content: "Error test",
			},
			prepare: func(f *fields) {
				f.mock.ExpectQuery("INSERT INTO tweets").
					WillReturnError(errors.New("prepare error"))
			},
			expectErr: errors.New("prepare error"),
		},
		"when query fails": { // Cambiamos el nombre del caso de prueba
			inputTweet: domain.Tweet{
				UserID:  3,
				Content: "Query error",
			},
			prepare: func(f *fields) {
				f.mock.ExpectQuery("INSERT INTO tweets").
					WithArgs(3, "Query error").
					WillReturnError(errors.New("query error")) // Usamos un error genérico
			},
			expectErr: errors.New("query error"),
		},
		"when scan fails": {
			inputTweet: domain.Tweet{
				UserID:  4,
				Content: "Scan error",
			},
			prepare: func(f *fields) {
				rows := sqlmock.NewRows([]string{"id"}).AddRow("invalid") // Simulamos un valor no válido para el ID
				f.mock.ExpectQuery("INSERT INTO tweets").
					WithArgs(4, "Scan error").
					WillReturnRows(rows)
			},
			expectErr: errors.New("sql: Scan error on column index 0, name \"id\": converting driver.Value type string (\"invalid\") to a uint64: invalid syntax"), // El error esperado será un error de scan
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			f := fields{
				db:   db,
				mock: mock,
			}

			if tc.prepare != nil {
				tc.prepare(&f)
			}

			repo := tweet.NewRepository(f.db)
			id, err := repo.Create(context.Background(), tc.inputTweet)

			if tc.expectErr != nil {
				assert.EqualError(t, err, tc.expectErr.Error())
				assert.Zero(t, id) // Aseguramos que el ID es cero en caso de error
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectID, id) // Verificamos que el ID sea el esperado
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
