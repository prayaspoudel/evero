package database

import "context"

type NoSQL interface {
	Store(context.Context, string, interface{}) error
	Update(context.Context, string, interface{}, interface{}) error
	FindAll(context.Context, string, interface{}, interface{}) error
	FindOne(context.Context, string, interface{}, interface{}, interface{}) error
	StartSession() (Session, error)
}

type Session interface {
	WithTransaction(context.Context, func(context.Context) error) error
	EndSession(context.Context)
}

type SQL interface {
	ExecuteContext(context.Context, string, ...interface{}) error
	QueryContext(context.Context, string, ...interface{}) (Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) Row
	BeginTx(ctx context.Context) (Tx, error)
}

type Rows interface {
	Scan(dest ...interface{}) error
	Next() bool
	Err() error
	Close() error
}

type Row interface {
	Scan(dest ...interface{}) error
}

type Tx interface {
	ExecuteContext(context.Context, string, ...interface{}) error
	QueryContext(context.Context, string, ...interface{}) (Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) Row
	Commit() error
	Rollback() error
}
