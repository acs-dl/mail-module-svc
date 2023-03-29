package postgres

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/fatih/structs"
	"gitlab.com/distributed_lab/acs/mail-module/internal/data"
	"gitlab.com/distributed_lab/logan/v3/errors"

	sq "github.com/Masterminds/squirrel"
	"gitlab.com/distributed_lab/kit/pgdb"
)

const usersTableName = "users"

type UsersQ struct {
	db  *pgdb.DB
	sql sq.SelectBuilder
}

var selectedUsersTable = sq.Select("*").From(usersTableName)

var usersColumns = []string{
	usersTableName + ".id",
	usersTableName + ".email",
	usersTableName + ".mail_id",
	usersTableName + ".name",
	usersTableName + ".created_at",
}

func NewUsersQ(db *pgdb.DB) data.Users {
	return &UsersQ{
		db:  db.Clone(),
		sql: selectedUsersTable,
	}
}

func (q *UsersQ) New() data.Users {
	return NewUsersQ(q.db)
}

func (q *UsersQ) Upsert(user data.User) error {
	clauses := structs.Map(user)

	updateQuery := sq.Update(" ").
		Set("updated_at", time.Now())

	if user.Id != nil {
		updateQuery = updateQuery.Set("id", *user.Id)
	}

	updateStmt, args := updateQuery.MustSql()

	query := sq.Insert(usersTableName).SetMap(clauses).Suffix("ON CONFLICT (mail_id) DO "+updateStmt, args...)

	return q.db.Exec(query)
}

func (q *UsersQ) Delete(mailId string) error {
	var deleted []data.User

	query := sq.Delete(usersTableName).
		Where(sq.Eq{
			"mail_id": mailId,
		}).
		Suffix("RETURNING *")

	err := q.db.Select(&deleted, query)
	if err != nil {
		return err
	}
	if len(deleted) == 0 {
		return errors.Errorf("no rows with `%d` mail id", mailId)
	}

	return nil
}

func (q *UsersQ) Get() (*data.User, error) {
	var result data.User

	err := q.db.Get(&result, q.sql)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &result, err
}

func (q *UsersQ) Select() ([]data.User, error) {
	var result []data.User

	err := q.db.Select(&result, q.sql)

	return result, err
}

func (q *UsersQ) FilterById(id *int64) data.Users {
	stmt := sq.Eq{usersTableName + ".id": id}

	q.sql = q.sql.Where(stmt)

	return q
}

func (q *UsersQ) FilterByMailIds(mailIds ...string) data.Users {
	q.sql = q.sql.Where(sq.Eq{usersTableName + ".mail_id": mailIds})

	return q
}

func (q *UsersQ) SearchBy(search string) data.Users {
	search = strings.Replace(search, " ", "%", -1)
	search = fmt.Sprint("%", search, "%")

	q.sql = q.sql.Where(sq.ILike{"email": search})
	return q
}

func (q *UsersQ) Page(pageParams pgdb.OffsetPageParams) data.Users {
	q.sql = pageParams.ApplyTo(q.sql, "email")

	return q
}

func (q *UsersQ) Count() data.Users {
	q.sql = sq.Select("COUNT (*)").From(usersTableName)

	return q
}

func (q *UsersQ) GetTotalCount() (int64, error) {
	var count int64
	err := q.db.Get(&count, q.sql)

	return count, err
}

func (q *UsersQ) FilterByLowerTime(time time.Time) data.Users {
	q.sql = q.sql.Where(sq.Lt{usersTableName + ".updated_at": time})

	return q
}
