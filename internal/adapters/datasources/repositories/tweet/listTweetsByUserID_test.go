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

func TestListTweetsByUserID(t *testing.T) {
	type fields struct {
		db   *sql.DB
		mock sqlmock.Sqlmock
	}

	const (
		defaultLimit  = 20
		defaultOffset = 0
		testUserID    = 123
	)

	tests := map[string]struct {
		userID       uint64
		prepare      func(f *fields)
		expectErr    error
		expectTweets []domain.Tweet
	}{
		"when tweets are found for the user": {
			userID: testUserID,
			prepare: func(f *fields) {
				rows := sqlmock.NewRows([]string{"id", "content", "user_id", "created_at"}).
					AddRow(1, "Tweet one", testUserID, time.Now()).
					AddRow(2, "Tweet two", testUserID, time.Now().Add(-time.Hour))
				f.mock.ExpectQuery("SELECT id, content, user_id, created_at FROM tweets WHERE user_id = \\$1 ORDER BY created_at DESC LIMIT \\$2 OFFSET \\$3").
					WithArgs(testUserID, defaultLimit, defaultOffset).
					WillReturnRows(rows)
			},
			expectTweets: []domain.Tweet{
				{ID: 1, Content: "Tweet one", UserID: testUserID},
				{ID: 2, Content: "Tweet two", UserID: testUserID},
			},
		},
		"when no tweets are found for the user": {
			userID: testUserID,
			prepare: func(f *fields) {
				rows := sqlmock.NewRows([]string{"id", "content", "user_id", "created_at"})
				f.mock.ExpectQuery("SELECT id, content, user_id, created_at FROM tweets WHERE user_id = \\$1 ORDER BY created_at DESC LIMIT \\$2 OFFSET \\$3").
					WithArgs(testUserID, defaultLimit, defaultOffset).
					WillReturnRows(rows)
			},
			expectTweets: []domain.Tweet{},
		},
		"when query fails": {
			userID: testUserID,
			prepare: func(f *fields) {
				f.mock.ExpectQuery("SELECT id, content, user_id, created_at FROM tweets WHERE user_id = \\$1 ORDER BY created_at DESC LIMIT \\$2 OFFSET \\$3").
					WithArgs(testUserID, defaultLimit, defaultOffset).
					WillReturnError(errors.New("query error"))
			},
			expectErr: errors.New("query error"),
		},
		"when scan fails": {
			userID: testUserID,
			prepare: func(f *fields) {
				rows := sqlmock.NewRows([]string{"id", "content", "user_id", "created_at"}).
					AddRow("invalid", "Invalid Tweet", testUserID, time.Now())
				f.mock.ExpectQuery("SELECT id, content, user_id, created_at FROM tweets WHERE user_id = \\$1 ORDER BY created_at DESC LIMIT \\$2 OFFSET \\$3").
					WithArgs(testUserID, defaultLimit, defaultOffset).
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
			tweets, err := repo.ListTweetsByUserID(context.Background(), tc.userID)

			if tc.expectErr != nil {
				assert.EqualError(t, err, tc.expectErr.Error())
				assert.Nil(t, tweets)
			} else {
				assert.NoError(t, err)
				assert.Len(t, tweets, len(tc.expectTweets))
				if len(tc.expectTweets) > 0 {
					for i, expected := range tc.expectTweets {
						assert.Equal(t, expected.ID, tweets[i].ID)
						assert.Equal(t, expected.Content, tweets[i].Content)
						assert.Equal(t, expected.UserID, tweets[i].UserID)
						// We don't compare CreatedAt precisely
					}
				}
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
