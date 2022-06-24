package api

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	db "github.com/olaysco/evolve/db/model"
)

// GetUsers return multiple users, and can be filtered by column
func (a *Api) GetUsers(w http.ResponseWriter, r *http.Request) {
	email, emailExist := r.Form["email"]
	offset, _ := strconv.Atoi(r.FormValue("page"))
	limit, _ := strconv.Atoi(r.FormValue("per_page"))
	dateOfBirthTo, dateOfBirthToExist := r.Form["date_of_birth_to"]
	dateOfBirthFrom, dateOfBirthFromExist := r.Form["date_of_birth_from"]

	if limit <= 0 {
		limit = 10
	}
	if offset <= 0 {
		offset = 1
	}

	listUserArg := db.ListUserArg{
		Limit:           int32(limit),
		Offset:          int32(offset),
		Email:           db.NullString{Value: strings.Join(email, ","), Valid: emailExist},
		DateOfBirthTo:   db.NullString{Value: strings.Join(dateOfBirthTo, ""), Valid: dateOfBirthToExist},
		DateOfBirthFrom: db.NullString{Value: strings.Join(dateOfBirthFrom, ""), Valid: dateOfBirthFromExist},
	}

	users, err := db.ListUsers(a.DB, r.Context(), listUserArg)
	if err != nil {
		a.JSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	a.JSON(w, http.StatusOK, &PaginatedResponse{Data: users, Count: limit, Page: offset})
}

// GetUserByEmail returns single user by email
func (a *Api) GetUserByEmail(w http.ResponseWriter, r *http.Request) {
	email, emailExist := mux.Vars(r)["email"]

	listUserArg := db.ListUserArg{
		Limit:  1,
		Offset: 1,
		Email:  db.NullString{Value: email, Valid: emailExist},
	}

	users, _ := db.ListUsers(a.DB, r.Context(), listUserArg)

	if len(users) < 1 {
		a.JSON(w, http.StatusNotFound, ErrorResponse{Status: http.StatusNotFound, Message: "User not found"})
		return
	}

	a.JSON(w, http.StatusOK, Response{users[0]})
}
