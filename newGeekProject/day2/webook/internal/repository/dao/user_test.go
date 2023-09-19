package dao

import (
	"context"
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	gormMysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
)

func TestUserDao_InsertUser(t *testing.T) {
	testCase := []struct {
		name string
		//这里不是go mock
		mock    func(t *testing.T) *sql.DB
		userDB  UserDB
		wantErr error
	}{
		{
			name: "插入成功",
			mock: func(t *testing.T) *sql.DB {
				mockDB, mock, err := sqlmock.New()
				res := sqlmock.NewResult(3, 1)
				mock.ExpectExec("INSERT INTO `user_dbs` .*").
					WillReturnResult(res)
				require.NoError(t, err)
				return mockDB
			},
			userDB: UserDB{
				Email: sql.NullString{
					Valid:  true,
					String: "1233@qqq.com",
				},
			},
		},
		{
			name: "邮箱冲突",
			mock: func(t *testing.T) *sql.DB {
				mockDB, mock, err := sqlmock.New()
				// 这边预期的是正则表达式
				// 这个写法的意思就是，只要是 INSERT 到 users 的语句
				mock.ExpectExec("INSERT INTO `user_dbs` .*").
					WillReturnError(&mysql.MySQLError{
						Number: 1062,
					})
				require.NoError(t, err)
				return mockDB
			},
			userDB:  UserDB{},
			wantErr: ErrUserDuplicate,
		},
		{
			name: "数据库错误",
			mock: func(t *testing.T) *sql.DB {
				mockDB, mock, err := sqlmock.New()
				// 这边预期的是正则表达式
				// 这个写法的意思就是，只要是 INSERT 到 users 的语句
				mock.ExpectExec("INSERT INTO `user_dbs` .*").
					WillReturnError(errors.New("数据库错误"))
				require.NoError(t, err)
				return mockDB
			},
			userDB:  UserDB{},
			wantErr: errors.New("数据库错误"),
		},
		// TODO: Add test cases.
	}
	for _, tt := range testCase {
		t.Run(tt.name, func(t *testing.T) {
			db, err := gorm.Open(gormMysql.New(gormMysql.Config{
				Conn:                      tt.mock(t),
				SkipInitializeWithVersion: true,
			}), &gorm.Config{
				DisableAutomaticPing:   true,
				SkipDefaultTransaction: true,
			})
			require.NoError(t, err)
			ud := NewUserDao(db)
			err = ud.InsertUser(context.Background(), tt.userDB)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
