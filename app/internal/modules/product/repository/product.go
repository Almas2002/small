package repository

import (
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"
	"github.com/opentracing/opentracing-go"
	"small/internal/models"
	"small/internal/modules/product/repository/dao"
	userDao "small/internal/modules/user/repository/dao"
	"small/pkg/tools/transaction"
	"small/pkg/tracing"
)

func (r *repository) SaveProduct(ctx context.Context, product *models.Product) (uint32, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "repository.SaveProduct")
	defer span.Finish()

	r.log.Info("repository.SaveProduct")

	ctx, cancel := context.WithTimeout(ctx, r.options.Timeout)
	defer cancel()

	tx, err := r.db.Begin(ctx)
	if err != nil {
		tracing.TraceError(span, err)
		r.log.WarnMsg("CreateCategory.Begin", err)
		return 0, err
	}

	defer func(ctx context.Context, t pgx.Tx) {
		err = transaction.Finish(ctx, t, err)
	}(ctx, tx)

	id, err := r.saveProductTx(ctx, tx, product)
	if err != nil {
		tracing.TraceError(span, err)
		r.log.WarnMsg("CreateCategory.saveProductTx", err)
	}

	return id, nil
}

func (r *repository) saveProductTx(ctx context.Context, tx pgx.Tx, product *models.Product) (uint32, error) {
	r.log.Info("repository.saveProductTx")

	sqlGen, args, err := r.genSQL.Insert(dao.ProductTableName).
		Columns(dao.ProductInsert...).
		Values(product.Title, product.Price).
		Suffix("returning product_id").
		ToSql()

	if err != nil {
		r.log.WarnMsg("saveProductTx.genSQl", err)
		return 0, err
	}

	var id uint32
	if err = tx.QueryRow(ctx, sqlGen, args...).Scan(&id); err != nil {
		r.log.WarnMsg("saveProductTx.genSQl", err)
		return 0, err
	}

	return id, err

}

func (r *repository) FindOneProduct(ctx context.Context, product *models.Product) (*models.Product, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "repository.SaveProduct")
	defer span.Finish()

	r.log.Info("repository.SaveProduct")

	ctx, cancel := context.WithTimeout(ctx, r.options.Timeout)
	defer cancel()

	tx, err := r.db.Begin(ctx)
	if err != nil {
		tracing.TraceError(span, err)
		r.log.WarnMsg("CreateCategory.Begin", err)
		return nil, err
	}

	defer func(ctx context.Context, t pgx.Tx) {
		err = transaction.Finish(ctx, t, err)
	}(ctx, tx)

	domainProduct, err := r.findOneProductTx(ctx, tx, product)
	if err != nil {
		tracing.TraceError(span, err)
		r.log.WarnMsg("CreateCategory.saveProductTx", err)
	}
	return domainProduct, nil
}

func (r *repository) findOneProductTx(ctx context.Context, tx pgx.Tx, product *models.Product) (*models.Product, error) {
	r.log.Info("repository.findOneProductTx")

	builder := r.genSQL.Select(dao.ProductsSelect...).
		From(dao.ProductTableName)

	if product.Id != 0 {
		builder = builder.Where(squirrel.And{squirrel.Eq{"product_id": product.Id}})
	}

	if product.Title != "" {
		builder = builder.Where(squirrel.And{squirrel.Eq{"title": product.Title}})
	}

	sqlGen, args, err := builder.ToSql()
	if err != nil {
		r.log.WarnMsg("findOneProductTx.genSQl", err)
		return nil, err
	}

	var daoProduct dao.Product
	rows, err := tx.Query(ctx, sqlGen, args...)
	if err != nil {
		r.log.WarnMsg("findOneProductTx.Query", err)
		return nil, err
	}

	for rows.Next() {
		if err = rows.Scan(&daoProduct.Id, &daoProduct.Title, &daoProduct.Price); err != nil {
			r.log.WarnMsg("findOneProductTx.Scan", err)
			return nil, err
		}
	}

	return r.productFromDaoToDomain(&daoProduct), err
}

func (r *repository) UpdateProduct(ctx context.Context, id uint32, updateFn func(old *models.Product) *models.Product) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "repository.UpdateProduct")
	defer span.Finish()

	r.log.Info("repository.UpdateProduct")

	ctx, cancel := context.WithTimeout(ctx, r.options.Timeout)
	defer cancel()

	tx, err := r.db.Begin(ctx)
	if err != nil {
		tracing.TraceError(span, err)
		r.log.WarnMsg("UpdateProduct.Begin", err)
		return err
	}

	defer func(ctx context.Context, t pgx.Tx) {
		err = transaction.Finish(ctx, t, err)
	}(ctx, tx)

	oldProduct, err := r.findOneProductTx(ctx, tx, &models.Product{Id: id})
	if err != nil {
		tracing.TraceError(span, err)
		r.log.WarnMsg("UpdateProduct.Old.findOneProductTx", err)
		return err
	}

	newProduct := updateFn(oldProduct)
	if err = r.updateProductTx(ctx, tx, newProduct); err != nil {
		tracing.TraceError(span, err)
		r.log.WarnMsg("UpdateProduct.updateProductTx", err)
		return err
	}

	return nil

}

func (r *repository) updateProductTx(ctx context.Context, tx pgx.Tx, product *models.Product) error {
	r.log.Info("repository.updateProductTx")

	sqlGen, args, err := r.genSQL.Update(dao.ProductTableName).
		Set("title", product.Title).
		Set("price", product.Price).Where(squirrel.And{squirrel.Eq{
		"product_id": product.Id,
	}}).ToSql()
	if err != nil {
		r.log.WarnMsg("updateProductTx.genSQL", err)
		return err
	}

	tag, err := tx.Exec(ctx, sqlGen, args...)
	if err != nil {
		r.log.WarnMsg("updateProductTx.Exec", err)
		return err
	}
	r.log.Debug(fmt.Sprintf("count updated products: %d", tag.RowsAffected()))
	return nil
}

func (r *repository) UserSubToProduct(ctx context.Context, userId uint32, productId uint32) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "repository.UserSubToProduct")

	r.log.Info("repository.UserSubToProduct")

	tx, err := r.db.Begin(ctx)
	if err != nil {
		tracing.TraceError(span, err)
		r.log.WarnMsg("UserSubToProduct.Begin", err)
		return err
	}

	defer func(ctx context.Context, t pgx.Tx) {
		err = transaction.Finish(ctx, t, err)
	}(ctx, tx)

	if err = r.userSubToProductTx(ctx, tx, userId, productId); err != nil {
		tracing.TraceError(span, err)
		r.log.WarnMsg("UserSubToProduct.userSubToProductTx", err)
		return err
	}
	return nil
}

func (r *repository) userSubToProductTx(ctx context.Context, tx pgx.Tx, userId uint32, productId uint32) error {
	r.log.Info("repository.userSubToProductTx")

	sqlGen, args, err := r.genSQL.Insert(userDao.UserSubProductsTableName).
		Columns(userDao.UserSubToProductInsert...).
		Values(userId, productId).
		ToSql()

	tag, err := tx.Exec(ctx, sqlGen, args...)
	if err != nil {
		r.log.WarnMsg("userSubToProductTx.Exec", err)
		return err
	}

	r.log.Debug(fmt.Sprintf("created subs: %d", tag.RowsAffected()))

	return nil
}

func (r *repository) UserUnsubToProduct(ctx context.Context, userId uint32, productId uint32) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "repository.UserUnsubToProduct")
	defer span.Finish()

	tx, err := r.db.Begin(ctx)
	if err != nil {
		tracing.TraceError(span, err)
		r.log.WarnMsg("UserUnsubToProduct.Begin", err)
		return err
	}

	defer func(ctx context.Context, t pgx.Tx) {
		err = transaction.Finish(ctx, t, err)
	}(ctx, tx)

	if err = r.userUnsubToProductTx(ctx, tx, userId, productId); err != nil {
		tracing.TraceError(span, err)
		r.log.WarnMsg("UserUnsubToProduct.userUnsubToProductTx", err)
		return err
	}
	return nil
}

func (r *repository) userUnsubToProductTx(ctx context.Context, tx pgx.Tx, userId uint32, productId uint32) error {
	r.log.Info("repository.userUnsubToProductTx")

	sqlGen, args, err := r.genSQL.Delete(userDao.UserSubProductsTableName).
		Where(squirrel.And{squirrel.Eq{"user_id": userId}, squirrel.Eq{"product_id": productId}}).
		ToSql()

	tag, err := tx.Exec(ctx, sqlGen, args...)
	if err != nil {
		r.log.WarnMsg("userUnsubToProductTx.Exec", err)
		return err
	}

	r.log.Debug(fmt.Sprintf("deleted subs: %d", tag.RowsAffected()))

	return nil
}
