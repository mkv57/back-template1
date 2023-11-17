package repo

import (
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"

	"github.com/ZergsLaw/back-template/cmd/user/internal/app"
)

type searchUsers struct {
	statuses       []string
	userName       string
	fullName       string
	email          string
	startCreatedAt time.Time
	endCreatedAt   time.Time
}

func newSearchUsers(params app.SearchParams) searchUsers {
	statuses := make([]string, 0, len(params.Statuses))
	for _, status := range params.Statuses {
		statuses = append(statuses, status.String())
	}

	return searchUsers{
		statuses:       statuses,
		userName:       params.Username,
		fullName:       params.FullName,
		email:          params.Email,
		startCreatedAt: params.StartCreatedAt,
		endCreatedAt:   params.EndCreatedAt,
	}
}

func getUsers(params app.SearchParams) (string, []interface{}, error) {
	sql := sq.Select("*").
		From("users").
		Where("id != ?", params.OwnerID).
		Limit(params.Limit).
		Offset(params.Offset)

	filters := newSearchUsers(params)

	sql = filters.getFilters(sql)

	query, args, err := sql.ToSql()
	if err != nil {
		return "", nil, fmt.Errorf("sql.ToSql: %w", err)
	}

	return query, args, nil
}

func getTotalUsers(params app.SearchParams) (string, []interface{}, error) {
	sql := sq.Select("count(*) over() as total").
		From("users").
		Where("id != ?", params.OwnerID)

	filters := newSearchUsers(params)

	sql = filters.getFilters(sql)

	query, args, err := sql.ToSql()
	if err != nil {
		return "", nil, fmt.Errorf("sql.ToSql: %w", err)
	}

	return query, args, nil
}

func (s *searchUsers) getFilters(sql sq.SelectBuilder) sq.SelectBuilder {
	if len(s.statuses) != 0 {
		sql = sql.Where(sq.Eq{"status": s.statuses})
	}

	if s.userName != "" || s.fullName != "" {
		sql = sql.Where(sq.Or{
			sq.Expr("full_name ilike ?", fmt.Sprintf("%%%s%%", s.fullName)),
			sq.Expr("name ilike ?", fmt.Sprintf("%%%s%%", s.userName)),
		})
	}

	if s.email != "" {
		sql = sql.Where("email ilike ?", fmt.Sprintf("%%%s%%", s.email))
	}

	if !s.startCreatedAt.IsZero() {
		sql = sql.Where("created_at >= ?", s.startCreatedAt)
	}

	if !s.endCreatedAt.IsZero() {
		sql = sql.Where("created_at <= ?", s.endCreatedAt)
	}

	return sql.PlaceholderFormat(sq.Dollar)
}
