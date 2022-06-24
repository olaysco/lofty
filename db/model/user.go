package db

import (
	"context"
	"database/sql"
	"fmt"

	dbHelper "github.com/olaysco/evolve/db/helper"
)

type NullString struct {
	Value string
	Valid bool
}

type User struct {
	ID          int64  `json:"id"`
	Email       string `json:"email"`
	Gender      string `json:"gender"`
	LastName    string `json:"last_name"`
	FirstName   string `json:"first_name"`
	DateOfBirth string `json:"date_of_birth"`
}

type ListUserArg struct {
	Email           NullString `json:"email"`
	Limit           int32      `json:"limit"`
	Offset          int32      `json:"offset"`
	DateOfBirthTo   NullString `json:"dob_to"`
	DateOfBirthFrom NullString `json:"dob_from"`
}

// ListUsers from database
func ListUsers(db *sql.DB, ctx context.Context, args ListUserArg) ([]User, error) {
	query := new(dbHelper.SelectQuery)

	query.Select("users", "id, first_name, last_name, email, gender, date_of_birth")

	if args.Email.Valid {
		query.WhereEquals("email", args.Email.Value)
	}
	if args.DateOfBirthFrom.Valid {
		fmt.Println(args.DateOfBirthFrom.Valid)
		query.WhereGreaterOrEqual("date_of_birth", args.DateOfBirthFrom.Value)
	}
	if args.DateOfBirthTo.Valid {
		query.WhereLessOrEqual("date_of_birth", args.DateOfBirthTo.Value)
	}
	query.Limit(args.Limit)
	query.Offset((args.Offset - 1) * args.Limit)

	sql, filterValues := query.ToSql()

	rows, err := db.QueryContext(ctx, string(sql), filterValues...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	users := []User{}

	for rows.Next() {
		var u User
		if err := rows.Scan(
			&u.ID,
			&u.FirstName,
			&u.LastName,
			&u.Email,
			&u.DateOfBirth,
			&u.Gender,
		); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	if err := rows.Close(); err != nil {
		return nil, err
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
