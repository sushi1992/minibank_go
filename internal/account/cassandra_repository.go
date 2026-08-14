package account

import (
	"time"

	gocql "github.com/apache/cassandra-gocql-driver/v2"
)

type CassandraRepository struct {
	session *gocql.Session
}

func NewCassandraRepository(host string) (*CassandraRepository, error) {
	cluster := gocql.NewCluster(host)
	cluster.Keyspace = "minibank"

	session, err := cluster.CreateSession()
	if err != nil {
		return nil, err
	}

	return &CassandraRepository{
		session: session,
	}, nil
}

func (repo *CassandraRepository) SaveAccountAndOutboxEvent(account *Account, event OutboxEvent) error {
	batch := repo.session.Batch(gocql.LoggedBatch)
	accountQuery := `
		INSERT INTO accounts (
			id,
			owner,
			balance_pence,
			currency
		)
		VALUES (?, ?, ?, ?)
	`

	batch.Query(accountQuery,
		account.ID,
		account.Owner,
		account.BalancePence,
		string(account.Currency),
	)

	outboxQuery := `
		INSERT INTO outbox_pending (
			bucket,
			created_at,
			event_id,
			account_id,
			event_type,
			payload
		)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	batch.Query(outboxQuery,
		"default",
		time.Now(),
		event.EventID,
		event.AccountID,
		event.Type,
		event.Payload)

	err := batch.Exec()
	if err != nil {
		return err
	}

	return nil
}

func (repo *CassandraRepository) Save(account *Account) error {
	query := `
		INSERT INTO accounts (
			id,
			owner,
			balance_pence,
			currency
		)
		VALUES (?, ?, ?, ?)
	`

	err := repo.session.Query(
		query,
		account.ID,
		account.Owner,
		account.BalancePence,
		string(account.Currency),
	).Exec()

	if err != nil {
		return err
	}

	return nil
}

func (repo *CassandraRepository) Get(id string) (*Account, error) {
	query := `
		SELECT id, owner, balance_pence, currency FROM accounts WHERE id = ?
	`

	var account Account
	var currency string

	err := repo.session.Query(query, id).Scan(
		&account.ID,
		&account.Owner,
		&account.BalancePence,
		&currency)

	if err != nil {
		return nil, err
	}

	account.Currency = Currency(currency)
	return &account, nil
}
