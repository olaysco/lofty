package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestGetUsers(t *testing.T) {
	testCases := []struct {
		url     string
		name    string
		prepare func() (*sql.DB, sqlmock.Sqlmock)
		assert  func(recoder *httptest.ResponseRecorder, mock sqlmock.SqlmockCommon)
	}{
		{
			name: "GetAllUsers",
			url:  "http://localhost/users",
			prepare: func() (*sql.DB, sqlmock.Sqlmock) {
				db, mock, err := sqlmock.New()
				if err != nil {
					t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
				}

				rows := sqlmock.NewRows([]string{"id", "first_name", "last_name", "email", "gender", "date_of_birth"}).
					AddRow(1, "doe", "ogene", "ogene@mail.com", "male", "2022-05-19 11:06:08+00").
					AddRow(2, "doe", "ogene", "ogene@mail.com", "male", "2022-05-19 11:06:08+00").
					AddRow(3, "doe", "ogene", "ogene@mail.com", "male", "2022-05-19 11:06:08+00").
					AddRow(4, "doe", "ogene", "ogene@mail.com", "male", "2022-05-19 11:06:08+00")

				mock.ExpectQuery("SELECT id, first_name, last_name, email, gender, date_of_birth FROM users LIMIT 10 OFFSET 0").WillReturnRows(rows)
				return db, mock
			},
			assert: func(recoder *httptest.ResponseRecorder, mock sqlmock.SqlmockCommon) {
				if err := mock.ExpectationsWereMet(); err != nil {
					t.Errorf("expectations are not met: %s", err)
				}
			},
		},
		{
			name: "GetTwoUsers",
			url:  "http://localhost/users?per_page=2",
			prepare: func() (*sql.DB, sqlmock.Sqlmock) {
				db, mock, err := sqlmock.New()
				if err != nil {
					t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
				}

				rows := sqlmock.NewRows([]string{"id", "first_name", "last_name", "email", "gender", "date_of_birth"}).
					AddRow(1, "doe", "ogene", "ogene1@mail.com", "male", "2022-05-19 11:06:08+00").
					AddRow(2, "doe", "ogene", "ogene2@mail.com", "male", "2022-05-19 11:06:08+00").
					AddRow(3, "doe", "ogene", "ogene3@mail.com", "male", "2022-05-19 11:06:08+00").
					AddRow(4, "doe", "ogene", "ogene4@mail.com", "male", "2022-05-19 11:06:08+00")

				mock.ExpectQuery("SELECT id, first_name, last_name, email, gender, date_of_birth FROM users LIMIT 2 OFFSET 0").WillReturnRows(rows)
				return db, mock
			},
			assert: func(recorder *httptest.ResponseRecorder, mock sqlmock.SqlmockCommon) {
				if err := mock.ExpectationsWereMet(); err != nil {
					t.Errorf("expectations are not met: %s", err)
				}

				var response PaginatedResponse
				err := json.Unmarshal(recorder.Body.Bytes(), &response)
				if err != nil {
					t.Errorf("Invalid response returned %s", bytes.NewBuffer(recorder.Body.Bytes()).String())
				}

				if response.Count != 2 {
					t.Errorf("Expected %v users but got %v instead", 2, response.Count)
				}
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			db, mock := tc.prepare()
			defer db.Close()
			api := &Api{DB: db}
			req, err := http.NewRequest("GET", tc.url, nil)
			if err != nil {
				t.Fatalf("an error '%s' was not expected while creating request", err)
			}
			w := httptest.NewRecorder()
			api.GetUsers(w, req)
			tc.assert(w, mock)
		})
	}
}

func TestGetUsersByEmail(t *testing.T) {
	testCases := []struct {
		url     string
		name    string
		prepare func() (*sql.DB, sqlmock.Sqlmock)
		assert  func(recoder *httptest.ResponseRecorder, mock sqlmock.SqlmockCommon)
	}{
		{
			name: "GetUserThatExist",
			url:  "http://localhost/users/ogene1@mail.com",
			prepare: func() (*sql.DB, sqlmock.Sqlmock) {
				db, mock, err := sqlmock.New()
				if err != nil {
					t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
				}

				rows := sqlmock.NewRows([]string{"id", "first_name", "last_name", "email", "gender", "date_of_birth"}).
					AddRow(1, "doe", "ogene", "ogene1@mail.com", "male", "2022-05-19 11:06:08+00").
					AddRow(2, "doe", "ogene", "ogene2@mail.com", "male", "2022-05-19 11:06:08+00").
					AddRow(3, "doe", "ogene", "ogene3@mail.com", "male", "2022-05-19 11:06:08+00").
					AddRow(4, "doe", "ogene", "ogene4@mail.com", "male", "2022-05-19 11:06:08+00")

				mock.ExpectQuery("SELECT id, first_name, last_name, email, gender, date_of_birth FROM users LIMIT 1 OFFSET 0").WillReturnRows(rows)
				return db, mock
			},
			assert: func(recorder *httptest.ResponseRecorder, mock sqlmock.SqlmockCommon) {
				if err := mock.ExpectationsWereMet(); err != nil {
					t.Errorf("expectations are not met: %s", err)
				}

				var response Response
				err := json.Unmarshal(recorder.Body.Bytes(), &response)
				if err != nil {
					t.Errorf("Invalid response returned %s", bytes.NewBuffer(recorder.Body.Bytes()).String())
				}

				user := response.Data.(map[string]interface{})
				if user["email"] != "ogene1@mail.com" {
					t.Errorf("Expected user %v but got %v instead", "ogene1@mail.com", user["email"])
				}
			},
		},
		{
			name: "GetUserThatDoesntExist",
			url:  "http://localhost/users/random",
			prepare: func() (*sql.DB, sqlmock.Sqlmock) {
				db, mock, err := sqlmock.New()
				if err != nil {
					t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
				}

				sqlmock.NewRows([]string{"id", "first_name", "last_name", "email", "gender", "date_of_birth"}).
					AddRow(1, "doe", "ogene", "ogene@mail.com", "male", "2022-05-19 11:06:08+00")

				mock.ExpectQuery("SELECT id, first_name, last_name, email, gender, date_of_birth FROM users LIMIT 1 OFFSET 0").WillReturnError(sql.ErrNoRows)
				return db, mock
			},
			assert: func(recorder *httptest.ResponseRecorder, mock sqlmock.SqlmockCommon) {
				if err := mock.ExpectationsWereMet(); err != nil {
					t.Errorf("expectations are not met: %s", err)
				}

				status := recorder.Result().StatusCode
				if status != http.StatusNotFound {
					t.Errorf("Expected status code %v but got %v instead", http.StatusNotFound, recorder.Result().StatusCode)
				}
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			db, mock := tc.prepare()
			defer db.Close()
			api := &Api{DB: db}
			req, err := http.NewRequest("GET", tc.url, nil)
			if err != nil {
				t.Fatalf("an error '%s' was not expected while creating request", err)
			}
			w := httptest.NewRecorder()
			api.GetUserByEmail(w, req)
			tc.assert(w, mock)
		})
	}
}
