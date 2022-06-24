package db

type QueryBody string

type Query interface {
	ToSql()
}
