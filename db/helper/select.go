package db

import (
	"database/sql"
	"fmt"
)

type SelectQuery struct {
	orderBy string
	column  string
	table   string
	wheres  string
	limit   sql.NullInt32
	offset  sql.NullInt32
	args    []interface{}
}

func (q *SelectQuery) Select(table string, column string) {
	q.table = table
	q.column = column
}

func (q *SelectQuery) Where(condition string, operator string, args ...interface{}) {
	q.args = append(q.args, args...)
	if q.wheres == "" {
		q.wheres = condition
		return
	}

	if operator == "" {
		operator = "AND"
	}

	q.wheres += fmt.Sprintf(" %s %s", operator, condition)
}

func (q *SelectQuery) WhereEquals(column string, value string) {
	sql := fmt.Sprintf("%s = $%d", column, len(q.args)+1)

	q.Where(sql, "", value)
}

func (q *SelectQuery) WhereGreaterOrEqual(column string, value string) {
	sql := fmt.Sprintf("%s >= $%d", column, len(q.args)+1)
	q.Where(sql, "", value)
}

func (q *SelectQuery) WhereLessOrEqual(column string, value string) {
	sql := fmt.Sprintf("%s <= $%d", column, len(q.args)+1)

	q.Where(sql, "", value)
}

func (q *SelectQuery) Limit(limit int32) {
	q.limit.Scan(limit)
}

func (q *SelectQuery) Offset(offset int32) {
	q.offset.Scan(offset)
}

func (q *SelectQuery) OrderBy(column string, order string) {
	if q.orderBy == "" {
		q.orderBy = fmt.Sprintf("%s %s", column, order)
		return
	}

	q.orderBy += fmt.Sprintf(", %s %s", column, order)
}

func (s *SelectQuery) ToSql() (QueryBody, []interface{}) {
	var query QueryBody = ""

	if s.table == "" {
		return query, nil
	}

	if s.column == "" {
		s.column = "*"
	}

	query = QueryBody(fmt.Sprintf("SELECT %s FROM %s", s.column, s.table))

	if s.wheres != "" {
		query += QueryBody(fmt.Sprintf(" WHERE %s", s.wheres))
	}

	if s.orderBy != "" {
		query += QueryBody(fmt.Sprintf(" ORDER BY %s", s.orderBy))
	}
	if limit, _ := s.limit.Value(); limit != nil {
		query += QueryBody(fmt.Sprintf(" LIMIT %d", limit))
	}

	if offset, _ := s.offset.Value(); offset != nil {
		query += QueryBody(fmt.Sprintf(" OFFSET %d", offset))
	}

	return query, s.args
}
