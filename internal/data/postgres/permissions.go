package postgres

import (
	"database/sql"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/fatih/structs"
	"gitlab.com/distributed_lab/acs/mail-module/internal/data"
	"gitlab.com/distributed_lab/acs/mail-module/internal/helpers"
	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

const permissionsTableName = "permissions"

type PermissionsQ struct {
	db  *pgdb.DB
	sql sq.SelectBuilder
}

var permissionsColumns = []string{
	permissionsTableName + ".request_id",
	permissionsTableName + ".mail_id",
	permissionsTableName + ".link",
	permissionsTableName + ".created_at",
	permissionsTableName + ".updated_at",
}

func NewPermissionsQ(db *pgdb.DB) data.Permissions {
	return &PermissionsQ{
		db:  db.Clone(),
		sql: sq.Select(permissionsColumns...).From(permissionsTableName),
	}
}

func (q *PermissionsQ) New() data.Permissions {
	return NewPermissionsQ(q.db)
}

func (q *PermissionsQ) Create(permission data.Permission) error {
	clauses := structs.Map(permission)

	query := sq.Insert(permissionsTableName).SetMap(clauses)

	return q.db.Exec(query)
}

func (q *PermissionsQ) Select() ([]data.Permission, error) {
	var result []data.Permission

	err := q.db.Select(&result, q.sql)

	return result, err
}

func (q *PermissionsQ) Upsert(permission data.Permission) error {
	updateStmt, args := sq.Update(" ").
		Set("updated_at", time.Now()).MustSql()

	query := sq.Insert(permissionsTableName).SetMap(structs.Map(permission)).
		Suffix("ON CONFLICT (mail_id, link) DO "+updateStmt, args...)

	return q.db.Exec(query)
}

func (q *PermissionsQ) Get() (*data.Permission, error) {
	var result data.Permission

	err := q.db.Get(&result, q.sql)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &result, err
}

func (q *PermissionsQ) Delete(mailId string, link string) error {
	var deleted []data.Permission

	query := sq.Delete(permissionsTableName).
		Where(sq.Eq{
			"mail_id": mailId,
			"link":    link,
		}).
		Suffix("RETURNING *")

	err := q.db.Select(&deleted, query)
	if err != nil {
		return err
	}
	if len(deleted) == 0 {
		return errors.Errorf("no rows with `%s` mail id", mailId)
	}

	return nil
}

func (q *PermissionsQ) FilterByMailIds(mailIds ...string) data.Permissions {
	stmt := sq.Eq{permissionsTableName + ".mail_id": mailIds}

	q.sql = q.sql.Where(stmt)

	return q
}

func (q *PermissionsQ) FilterByGreaterTime(time time.Time) data.Permissions {
	q.sql = q.sql.Where(sq.Gt{permissionsTableName + ".updated_at": time})

	return q
}

func (q *PermissionsQ) FilterByLowerTime(time time.Time) data.Permissions {
	q.sql = q.sql.Where(sq.Lt{permissionsTableName + ".updated_at": time})

	return q
}

func (q *PermissionsQ) Count() data.Permissions {
	q.sql = sq.Select("COUNT (*)").From(permissionsTableName)

	return q
}

func (q *PermissionsQ) GetTotalCount() (int64, error) {
	var count int64
	err := q.db.Get(&count, q.sql)

	return count, err
}

func (q *PermissionsQ) Page(pageParams pgdb.OffsetPageParams) data.Permissions {
	q.sql = pageParams.ApplyTo(q.sql, "link")

	return q
}

func (q *PermissionsQ) WithUsers() data.Permissions {
	q.sql = sq.Select().Columns(helpers.RemoveDuplicateColumn(append(permissionsColumns, usersColumns...))...).
		From(permissionsTableName).
		LeftJoin(fmt.Sprint(usersTableName, " ON ", usersTableName, ".mail_id = ", permissionsTableName, ".mail_id")).
		Where(sq.NotEq{permissionsTableName + ".request_id": nil}).
		GroupBy(helpers.RemoveDuplicateColumn(append(permissionsColumns, usersColumns...))...)

	return q
}

func (q *PermissionsQ) CountWithUsers() data.Permissions {
	q.sql = sq.Select("COUNT(*)").From(permissionsTableName).
		LeftJoin(fmt.Sprint(usersTableName, " ON ", usersTableName, ".mail_id = ", permissionsTableName, ".mail_id")).
		Where(sq.NotEq{permissionsTableName + ".request_id": nil})

	return q
}

func (q *PermissionsQ) FilterByUserIds(userIds ...int64) data.Permissions {
	stmt := sq.Eq{usersTableName + ".id": userIds}

	if len(userIds) == 0 {
		stmt = sq.Eq{usersTableName + ".id": nil}
	}

	q.sql = q.sql.Where(stmt)

	return q
}
