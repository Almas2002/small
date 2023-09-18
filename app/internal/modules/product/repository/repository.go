package repository

import (
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"small/pkg/type/logger"
	"time"
)

func init() {
	viper.SetDefault("MIGRATIONS_DIR", "./app/migrations")
}

type repository struct {
	db      *pgxpool.Pool
	log     logger.Logger
	genSQL  squirrel.StatementBuilderType
	options Options
}

type Options struct {
	Timeout       time.Duration
	DefaultLimit  uint64
	DefaultOffset uint64
}

func New(db *pgxpool.Pool, o Options, logger logger.Logger) *repository {
	var r = &repository{
		db:     db,
		genSQL: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
		log:    logger,
	}
	r.SetOptions(o)
	return r
}

func (r *repository) SetOptions(options Options) {
	if options.DefaultLimit == 0 {
		options.DefaultLimit = 10

	}

	if options.Timeout == 0 {
		options.Timeout = time.Second * 30

	}

	if r.options != options {
		r.options = options
	}
}
