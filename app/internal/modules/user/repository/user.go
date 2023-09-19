package repository

import (
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"
	"github.com/opentracing/opentracing-go"
	"small/internal/models"
	"small/internal/modules/user/repository/dao"
	"small/pkg/tools/transaction"
	"small/pkg/tracing"
)

func (r *repository) CreateUser(ctx context.Context, user *models.User) (uint32, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "repository.CreateUser")
	defer span.Finish()

	r.log.Info("repository.CreateUser")

	ctx, cancel := context.WithTimeout(ctx, r.options.Timeout)
	defer cancel()

	tx, err := r.db.Begin(ctx)
	if err != nil {
		tracing.TraceError(span, err)
		r.log.WarnMsg("CreateUser.Begin", err)
		return 0, err
	}

	defer func(ctx context.Context, t pgx.Tx) {
		err = transaction.Finish(ctx, t, err)
	}(ctx, tx)

	id, err := r.createUserTx(ctx, tx, user)
	if err != nil {
		tracing.TraceError(span, err)
		r.log.WarnMsg("CreateUser.createUserTx", err)
		return 0, err
	}

	return id, nil

}

func (r *repository) createUserTx(ctx context.Context, tx pgx.Tx, user *models.User) (uint32, error) {
	r.log.Info("repository.createUserTx")

	sqlGen, args, err := r.genSQL.Insert(dao.UserTableName).
		Columns(dao.UserInsert...).
		Values(user.Phone, user.Email).
		Suffix("Returning user_id").
		ToSql()
	if err != nil {
		r.log.WarnMsg("createUserTx.genSQL", err)
		return 0, nil
	}

	var id uint32
	if err = tx.QueryRow(ctx, sqlGen, args...).Scan(&id); err != nil {
		r.log.WarnMsg("createUserTx.QueryRow", err)
		return 0, nil
	}

	return id, err
}

func (r *repository) FinOneUser(ctx context.Context, user *models.User) (*models.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "repository.FindOneUser")
	defer span.Finish()

	ctx, cancel := context.WithTimeout(ctx, r.options.Timeout)
	defer cancel()

	tx, err := r.db.Begin(ctx)
	if err != nil {
		tracing.TraceError(span, err)
		r.log.WarnMsg("CreateUser.Begin", err)
		return nil, err
	}

	defer func(ctx context.Context, t pgx.Tx) {
		err = transaction.Finish(ctx, t, err)
	}(ctx, tx)

	user, err = r.findOneUserTx(ctx, tx, user)
	if err != nil {
		tracing.TraceError(span, err)
		return nil, err
	}

	return user, nil
}

func (r *repository) findOneUserTx(ctx context.Context, tx pgx.Tx, user *models.User) (*models.User, error) {
	r.log.Info("repository.findOneUserTx")

	builder := r.genSQL.Select(dao.UsersSelect...).From(dao.UserTableName)

	if user.Id != 0 {
		builder = builder.Where(squirrel.And{squirrel.Eq{"user_id": user.Id}})
	}

	if user.Email != "" {
		builder = builder.Where(squirrel.And{squirrel.Eq{"email": user.Email}})
	}

	if user.Phone != "" {
		builder = builder.Where(squirrel.And{squirrel.Eq{"phone": user.Phone}})
	}

	sqlGen, args, err := builder.ToSql()
	if err != nil {
		r.log.WarnMsg("findUserOneTx.genSQL", err)
		return nil, err
	}

	var daoUser dao.User

	rows, err := tx.Query(ctx, sqlGen, args...)
	if err != nil {
		r.log.WarnMsg("findUserOneTx.Query", err)
		return nil, err
	}

	for rows.Next() {
		if err = rows.Scan(&daoUser.Id, &daoUser.Phone, &daoUser.Email); err != nil {
			r.log.WarnMsg("findUserOneTx.Scan", err)
			return nil, err
		}
	}

	return r.userFromDaoToDomain(&daoUser), nil
}

func (r *repository) GetSubUsers(ctx context.Context, productId uint32) ([]*models.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "repository.GetSubUsers")
	defer span.Finish()

	r.log.Info("repository.GetSubUsers")

	tx, err := r.db.Begin(ctx)
	if err != nil {
		tracing.TraceError(span, err)
		r.log.WarnMsg("GetSubUsers.Begin", err)
		return nil, err
	}

	defer func(ctx context.Context, t pgx.Tx) {
		err = transaction.Finish(ctx, t, err)
	}(ctx, tx)

	users, err := r.getSubUsersTx(ctx, tx, productId)
	if err != nil {
		tracing.TraceError(span, err)
		return nil, err
	}
	return users, nil

}

func (r *repository) getSubUsersTx(ctx context.Context, tx pgx.Tx, productId uint32) ([]*models.User, error) {
	r.log.Info("repository.getSubUsersTx")

	sqlGen, args, err := r.genSQL.Select("users.user_id", "phone", "email").From(dao.UserTableName).
		LeftJoin(fmt.Sprintf("%s on %s.user_id = %s.user_id AND %s.product_id =%d", dao.UserSubProductsTableName,
			dao.UserTableName, dao.UserSubProductsTableName, dao.UserSubProductsTableName, productId)).ToSql()
	if err != nil {
		r.log.WarnMsg("getSubUsersTx.genSQL", err)
		return nil, err
	}
	fmt.Println(sqlGen)
	rows, err := tx.Query(ctx, sqlGen, args...)
	if err != nil {
		r.log.WarnMsg("getSubUsersTx.Query", err)
		return nil, err
	}

	daoUsers := make([]*dao.User, 0, 3)
	err = pgxscan.ScanAll(&daoUsers, rows)
	if err != nil {
		r.log.WarnMsg("getSubUsersTx.ScanAll", err)
		return nil, err
	}

	return r.usersFromDaoToDomain(daoUsers), nil
}
