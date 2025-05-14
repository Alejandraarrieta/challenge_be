package tweet_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"challenge_be/internal/adapters/datasources/repositories/tweet"
	domain "challenge_be/internal/domain/tweet"
)

func TestGetTweetsByIDs(t *testing.T) {
	type fields struct {
		db   *sql.DB
		mock sqlmock.Sqlmock
	}

	tests := map[string]struct {
		inputIDs     []uint64
		prepare      func(f *fields)
		expectErr    error
		expectTweets []domain.Tweet
	}{
		"when no IDs are provided": {
			inputIDs:     []uint64{},
			expectTweets: []domain.Tweet{},
		},
		"when tweets are found for the given IDs": {
			inputIDs: []uint64{1, 2},
			prepare: func(f *fields) {
				rows := sqlmock.NewRows([]string{"id", "user_id", "content", "created_at"}).
					AddRow(1, 10, "Tweet 1", time.Now()).
					AddRow(2, 20, "Tweet 2", time.Now().Add(-time.Hour))
				f.mock.ExpectQuery("SELECT id, user_id, content, created_at FROM tweets WHERE user_id IN \\(\\$1,\\$2\\) ORDER BY created_at DESC").
					WithArgs(int64(1), int64(2)).
					WillReturnRows(rows)
			},
			expectTweets: []domain.Tweet{
				{ID: 1, UserID: 10, Content: "Tweet 1"},
				{ID: 2, UserID: 20, Content: "Tweet 2"},
			},
		},
		"when no tweets are found for the given IDs": {
			inputIDs: []uint64{3, 4},
			prepare: func(f *fields) {
				rows := sqlmock.NewRows([]string{"id", "user_id", "content", "created_at"})
				f.mock.ExpectQuery("SELECT id, user_id, content, created_at FROM tweets WHERE user_id IN \\(\\$1,\\$2\\) ORDER BY created_at DESC").
					WithArgs(int64(3), int64(4)).
					WillReturnRows(rows)
			},
			expectTweets: []domain.Tweet{},
		},
		"when query fails": {
			inputIDs: []uint64{5},
			prepare: func(f *fields) {
				f.mock.ExpectQuery("SELECT id, user_id, content, created_at FROM tweets WHERE user_id IN \\(\\$1\\) ORDER BY created_at DESC").
					WithArgs(int64(5)).
					WillReturnError(errors.New("query error"))
			},
			expectErr: errors.New("query error"),
		},
		"when scan fails": {
			inputIDs: []uint64{6},
			prepare: func(f *fields) {
				rows := sqlmock.NewRows([]string{"id", "user_id", "content", "created_at"}).
					AddRow("invalid", 30, "Invalid Tweet", time.Now())
				f.mock.ExpectQuery("SELECT id, user_id, content, created_at FROM tweets WHERE user_id IN \\(\\$1\\) ORDER BY created_at DESC").
					WithArgs(int64(6)).
					WillReturnRows(rows)
			},
			expectErr: errors.New("sql: Scan error on column index 0, name \"id\": converting driver.Value type string (\"invalid\") to a uint64: invalid syntax"),
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
			tweets, err := repo.GetTweetsByIDs(context.Background(), tc.inputIDs)

			if tc.expectErr != nil {
				assert.EqualError(t, err, tc.expectErr.Error())
				assert.Nil(t, tweets)
			} else {
				assert.NoError(t, err)
				assert.Len(t, tweets, len(tc.expectTweets))
				if len(tc.expectTweets) > 0 {
					for i, expected := range tc.expectTweets {
						assert.Equal(t, expected.ID, tweets[i].ID)
						assert.Equal(t, expected.UserID, tweets[i].UserID)
						assert.Equal(t, expected.Content, tweets[i].Content)
						// We don't compare CreatedAt precisely as it involves time.Now()
					}
				}
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
