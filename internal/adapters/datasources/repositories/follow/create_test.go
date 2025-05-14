package follow_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"challenge_be/internal/adapters/datasources/repositories/follow"
	domain "challenge_be/internal/domain/follow"
)

func TestCreate(t *testing.T) {
	type fields struct {
		db   *sql.DB
		mock sqlmock.Sqlmock
	}

	tests := map[string]struct {
		inputFollow domain.Follow
		prepare     func(f *fields)
		expectErr   error
	}{
		"when follow is created successfully": {
			inputFollow: domain.Follow{
				FollowerID: 1,
				FolloweeID: 2,
			},
			prepare: func(f *fields) {
				f.mock.ExpectQuery("INSERT INTO follows").
					WithArgs(1, 2).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
			},
			expectErr: nil,
		},
		"when query fails": {
			inputFollow: domain.Follow{
				FollowerID: 5,
				FolloweeID: 6,
			},
			prepare: func(f *fields) {
				f.mock.ExpectQuery("INSERT INTO follows").
					WithArgs(5, 6).
					WillReturnError(errors.New("query error"))
			},
			expectErr: errors.New("query error"),
		},
		"when scan fails": {
			inputFollow: domain.Follow{
				FollowerID: 7,
				FolloweeID: 8,
			},
			prepare: func(f *fields) {
				f.mock.ExpectQuery("INSERT INTO follows").
					WithArgs(7, 8).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("invalid")) // Simulate a non-integer ID
			},
			expectErr: errors.New("sql: Scan error"), // A more generic scan error.  The actual error will vary by DB.
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

			repo := follow.NewRepository(f.db)
			err = repo.Create(context.Background(), tc.inputFollow)

			if tc.expectErr != nil {
				assert.Error(t, err)
				if tc.expectErr.Error() != "sql: Scan error" { // পরা evitar comparar el string "sql: Scan error"
					assert.EqualError(t, err, tc.expectErr.Error())
				}
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
