package postgres

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/acs-dl/mail-module-svc/internal/data"
	"github.com/fatih/structs"
	"gitlab.com/distributed_lab/logan/v3/errors"

	sq "github.com/Masterminds/squirrel"
	"gitlab.com/distributed_lab/kit/pgdb"
)

const (
	usersTableName       = "users"
	usersIdColumn        = usersTableName + ".id"
	usersMailIdColumn    = usersTableName + ".mail_id"
	usersEmailColumn     = usersTableName + ".email"
	usersNameColumn      = usersTableName + ".name"
	usersCreatedAtColumn = usersTableName + ".created_at"
	usersUpdatedAtColumn = usersTableName + ".updated_at"
)

type UsersQ struct {
	db            *pgdb.DB
	selectBuilder sq.SelectBuilder
	deleteBuilder sq.DeleteBuilder
	updateBuilder sq.UpdateBuilder
}

var (
	usersColumns = []string{
		usersIdColumn,
		usersEmailColumn,
		usersMailIdColumn,
		usersNameColumn,
		usersCreatedAtColumn,
	}
	selectedUsersTable = sq.Select("*").From(usersTableName)
)

func NewUsersQ(db *pgdb.DB) data.Users {
	return &UsersQ{
		db:            db.Clone(),
		selectBuilder: selectedUsersTable,
		deleteBuilder: sq.Delete(usersTableName),
		updateBuilder: sq.Update(usersTableName),
	}
}

func (q UsersQ) New() data.Users {
	return NewUsersQ(q.db)
}

func (q UsersQ) Upsert(user data.User) error {
	clauses := structs.Map(user)

	updateQuery := sq.Update(" ").
		Set("name", user.Name).
		Set("email", user.Email).
		Set("updated_at", time.Now())

	if user.Id != nil {
		updateQuery = updateQuery.Set("id", *user.Id)
	}

	updateStmt, args := updateQuery.MustSql()

	query := sq.Insert(usersTableName).SetMap(clauses).Suffix("ON CONFLICT (mail_id) DO "+updateStmt, args...)

	return q.db.Exec(query)
}

func (q UsersQ) Delete() error {
	var deleted []data.User

	err := q.db.Select(&deleted, q.deleteBuilder.Suffix("RETURNING *"))
	if err != nil {
		return err
	}

	if len(deleted) == 0 {
		return errors.Errorf("no such data to delete")
	}

	return nil
}

func (q UsersQ) Get() (*data.User, error) {
	var result data.User

	err := q.db.Get(&result, q.selectBuilder)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &result, err
}

func (q UsersQ) Select() ([]data.User, error) {
	var result []data.User

	err := q.db.Select(&result, q.selectBuilder)

	return result, err
}

func (q UsersQ) FilterById(id *int64) data.Users {
	equalId := sq.Eq{usersIdColumn: id}

	q.selectBuilder = q.selectBuilder.Where(equalId)
	q.deleteBuilder = q.deleteBuilder.Where(equalId)
	q.updateBuilder = q.updateBuilder.Where(equalId)

	return q
}

func (q UsersQ) FilterByMailIds(mailIds ...string) data.Users {
	equalTelegramIds := sq.Eq{usersMailIdColumn: mailIds}

	q.selectBuilder = q.selectBuilder.Where(equalTelegramIds)
	q.deleteBuilder = q.deleteBuilder.Where(equalTelegramIds)
	q.updateBuilder = q.updateBuilder.Where(equalTelegramIds)

	return q
}

func (q UsersQ) FilterByEmail(email ...string) data.Users {
	equalEmails := sq.Eq{usersEmailColumn: email}

	q.selectBuilder = q.selectBuilder.Where(equalEmails)
	q.deleteBuilder = q.deleteBuilder.Where(equalEmails)
	q.updateBuilder = q.updateBuilder.Where(equalEmails)

	return q
}

func (q UsersQ) Page(pageParams pgdb.OffsetPageParams) data.Users {
	q.selectBuilder = pageParams.ApplyTo(q.selectBuilder, "username")

	return q
}

func (q UsersQ) Count() data.Users {
	q.selectBuilder = sq.Select("COUNT (*)").From(usersTableName)

	return q
}

func (q UsersQ) GetTotalCount() (int64, error) {
	var count int64
	err := q.db.Get(&count, q.selectBuilder)

	return count, err
}

func (q UsersQ) SearchBy(search string) data.Users {
	search = strings.Replace(search, " ", "%", -1)
	search = fmt.Sprint("%", search, "%")

	q.selectBuilder = q.selectBuilder.Where(sq.ILike{usersEmailColumn: search})

	return q
}

func (q UsersQ) FilterByLowerTime(time time.Time) data.Users {
	lowerTime := sq.Lt{usersUpdatedAtColumn: time}

	q.selectBuilder = q.selectBuilder.Where(lowerTime)
	q.deleteBuilder = q.deleteBuilder.Where(lowerTime)
	q.updateBuilder = q.updateBuilder.Where(lowerTime)

	return q
}
