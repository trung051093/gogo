package common

import (
	"context"
	"errors"
	"fmt"

	"github.com/pilagod/gorm-cursor-paginator/v2/cursor"
	"github.com/pilagod/gorm-cursor-paginator/v2/paginator"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type GormEntity interface {
	TableName() string
}

type GormPagePagination interface {
	Value() *PagePagination
	GetLimit() int
	SetLimit(limit int)
	GetTotal() int64
	SetTotal(total int64)
	GetOffset() int
	SetOffset(offset int)
}

type GormCursorPagination interface {
	Value() *CursorPagination
	GetLimit() int
	SetLimit(limit int)
	GetTotal() int64
	SetTotal(total int64)
	GetCursor() cursor.Cursor
	SetCursor(c cursor.Cursor)
	GetPaginator() *paginator.Paginator
}

type Repository[E GormEntity] interface {
	GetDB() *gorm.DB
	ActionName(name string) string
	WithTX(tx *gorm.DB) Repository[E]
	GetDBWithCtx(ctx context.Context) *gorm.DB
	SetPreloadKeys(preloadKeys ...string)
	SetCustomPreloadKeys(preloadKeys map[string]interface{})
	GetPreload(stmt *gorm.DB) *gorm.DB
	HandleQuery(stmt *gorm.DB, q GormQuery) *gorm.DB
	Save(ctx context.Context, entity *E) (*E, error)
	Update(ctx context.Context, entity *E, column string, value interface{}) (*E, error)
	Updates(ctx context.Context, entity *E, updateDto map[string]interface{}) (*E, error)
	Create(ctx context.Context, entity *E) (*E, error)
	CreateList(ctx context.Context, entities []E) ([]E, error)
	DeleteById(ctx context.Context, Id any) error
	DeleteByIds(ctx context.Context, Ids any) error
	GetAll(ctx context.Context) ([]E, error)
	Count(ctx context.Context) (int64, error)
	CountWithCond(ctx context.Context, cond GormQuery) (int64, error)
	FindById(ctx context.Context, id any) (*E, error)
	FindByIds(ctx context.Context, Ids any) ([]E, error)
	FindOneByCond(ctx context.Context, cond GormQuery) (*E, error)
	FindByCond(ctx context.Context, cond GormQuery) ([]E, error)
	FindWithCursorPagination(ctx context.Context, query GormQuery, paging GormCursorPagination) ([]E, GormQuery, GormCursorPagination, error)
	FindWithPagePagination(ctx context.Context, query GormQuery, paging GormPagePagination) ([]E, GormQuery, GormPagePagination, error)
}

type ContextKey string

const (
	TransactionKey ContextKey = "tx_db"
)

type repository[E GormEntity] struct {
	DB          *gorm.DB
	PreloadKeys map[string]interface{} // preload query: selector
}

func NewRepository[E GormEntity](db *gorm.DB) Repository[E] {
	return &repository[E]{
		DB:          db,
		PreloadKeys: make(map[string]interface{}),
	}
}

func WithContext(ctx context.Context, db *gorm.DB) context.Context {
	return context.WithValue(ctx, TransactionKey, db)
}

func CommitTransaction(ctx context.Context) {
	db := ctx.Value(TransactionKey)
	if sessionDB, ok := db.(*gorm.DB); ok {
		_ = sessionDB.Commit()
	}
}

func RollbackTransaction(ctx context.Context) {
	db := ctx.Value(TransactionKey)
	if sessionDB, ok := db.(*gorm.DB); ok {
		_ = sessionDB.Rollback()
	}
}

func (r *repository[E]) WithTX(tx *gorm.DB) Repository[E] {
	return &repository[E]{
		DB:          tx,
		PreloadKeys: r.PreloadKeys,
	}
}

func (r *repository[E]) GetDBWithCtx(ctx context.Context) *gorm.DB {
	txCtx := ctx.Value(TransactionKey)
	if tx, ok := txCtx.(*gorm.DB); ok {
		return tx
	}
	return r.GetDB()
}

func (r *repository[E]) GetDB() *gorm.DB {
	db := r.DB.Session(&gorm.Session{
		NewDB: true,
	})
	return db
}

func (r *repository[E]) ActionName(name string) string {
	var e E
	return fmt.Sprintf("%s.%s", e.TableName(), name)
}

func (r *repository[E]) GetPreload(stmt *gorm.DB) *gorm.DB {
	if len(r.PreloadKeys) == 0 || r.PreloadKeys == nil {
		return stmt
	}
	for key, selector := range r.PreloadKeys {
		if selector != nil {
			stmt = stmt.Preload(key, selector)
		} else {
			stmt = stmt.Preload(key)
		}
	}
	return stmt
}

func (r *repository[E]) SetCustomPreloadKeys(preloadKeys map[string]interface{}) {
	for key, selector := range preloadKeys {
		r.PreloadKeys[key] = selector
	}
}

func (r *repository[E]) SetPreloadKeys(preloadKeys ...string) {
	for _, key := range preloadKeys {
		r.PreloadKeys[key] = nil
	}
}

func (r *repository[E]) Create(ctx context.Context, entity *E) (*E, error) {
	if err := r.GetDBWithCtx(ctx).WithContext(ctx).Model(entity).Create(entity).Error; err != nil {
		return nil, err
	}
	return entity, nil
}

func (r *repository[E]) Update(ctx context.Context, entity *E, column string, value interface{}) (*E, error) {
	if err := r.GetDBWithCtx(ctx).WithContext(ctx).Model(entity).Omit(clause.Associations).Update(column, value).Error; err != nil {
		return nil, err
	}
	return entity, nil
}

func (r *repository[E]) Updates(ctx context.Context, entity *E, updateDto map[string]interface{}) (*E, error) {
	if err := r.GetDBWithCtx(ctx).WithContext(ctx).Model(entity).Omit(clause.Associations).
		Updates(updateDto).Error; err != nil {
		return nil, err
	}
	return entity, nil
}

func (r *repository[E]) Save(ctx context.Context, entity *E) (*E, error) {
	if err := r.GetDBWithCtx(ctx).WithContext(ctx).Model(entity).Omit(clause.Associations).Save(entity).Error; err != nil {
		return nil, err
	}
	return entity, nil
}

func (r *repository[E]) CreateList(ctx context.Context, entities []E) ([]E, error) {
	var entity *E
	if err := r.GetDBWithCtx(ctx).WithContext(ctx).Model(entity).Create(entities).Error; err != nil {
		return nil, err
	}
	return entities, nil
}

func (r *repository[E]) DeleteById(ctx context.Context, id any) error {
	var e E
	if err := r.GetDBWithCtx(ctx).WithContext(ctx).Table(e.TableName()).Delete(e, id).Error; err != nil {
		return err
	}
	return nil
}

func (r *repository[E]) DeleteByIds(ctx context.Context, ids any) error {
	var e E
	if err := r.GetDBWithCtx(ctx).WithContext(ctx).Table(e.TableName()).Delete(e, ids).Error; err != nil {
		return err
	}
	return nil
}

func (r *repository[E]) GetAll(ctx context.Context) ([]E, error) {
	var entities []E
	var e E
	stmt := r.GetDBWithCtx(ctx).WithContext(ctx).Table(e.TableName())
	stmt = r.GetPreload(stmt)
	if err := stmt.Find(&entities).Error; err != nil {
		return nil, err
	}
	return entities, nil
}

func (r *repository[E]) Count(ctx context.Context) (int64, error) {
	var total int64
	var entity E
	stmt := r.GetDBWithCtx(ctx).WithContext(ctx).Model(entity)
	if err := stmt.Count(&total).Error; err != nil {
		return 0, err
	}

	return total, nil
}

func (r *repository[E]) CountWithCond(ctx context.Context, cond GormQuery) (int64, error) {
	var total int64
	var entity E
	stmt := r.GetDBWithCtx(ctx).WithContext(ctx).Model(entity)
	stmt = r.HandleQuery(stmt, cond)
	if err := stmt.Count(&total).Error; err != nil {
		return 0, err
	}

	return total, nil
}

func (r *repository[E]) FindById(ctx context.Context, Id any) (*E, error) {
	var entity E
	stmt := r.GetDBWithCtx(ctx).WithContext(ctx).Model(entity)
	stmt = r.GetPreload(stmt)
	if err := stmt.Where("id", Id).First(&entity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &entity, nil
}

func (r *repository[E]) FindByIds(ctx context.Context, Ids any) ([]E, error) {
	var entities []E
	stmt := r.GetDBWithCtx(ctx).WithContext(ctx).Model(entities)
	stmt = r.GetPreload(stmt)
	if err := stmt.Find(&entities).Where("id IN (?)", Ids).Error; err != nil {
		return nil, err
	}
	return entities, nil
}

func (r *repository[E]) FindOneByCond(ctx context.Context, cond GormQuery) (*E, error) {
	var entity E
	stmt := r.GetDBWithCtx(ctx).WithContext(ctx).Model(entity)
	stmt = r.HandleQuery(stmt, cond)
	stmt = r.GetPreload(stmt)
	if err := stmt.First(&entity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &entity, nil
}

func (r *repository[E]) FindByCond(ctx context.Context, cond GormQuery) ([]E, error) {
	var entities []E
	stmt := r.GetDBWithCtx(ctx).WithContext(ctx).Model(entities)
	stmt = r.HandleQuery(stmt, cond)
	stmt = r.GetPreload(stmt)
	if err := stmt.Find(&entities).Error; err != nil {
		return nil, err
	}
	return entities, nil
}

func (r *repository[E]) FindWithCursorPagination(ctx context.Context, query GormQuery, paging GormCursorPagination) ([]E, GormQuery, GormCursorPagination, error) {
	var entities []E
	var total int64

	stmt := r.GetDBWithCtx(ctx).WithContext(ctx).Model(&entities)

	stmt = r.HandleQuery(stmt, query)

	if err := stmt.Count(&total).Error; err != nil {
		return nil, nil, nil, err
	}

	// preload should be call in there.
	stmt = r.GetPreload(stmt)
	paginator := paging.GetPaginator()
	result, cursor, err := paginator.Paginate(stmt, &entities)
	if err != nil || result.Error != nil {
		return nil, nil, nil, err
	}

	paging.SetTotal(total)
	paging.SetCursor(cursor)
	return entities, query, paging, nil
}

func (r *repository[E]) FindWithPagePagination(ctx context.Context, query GormQuery, paging GormPagePagination) ([]E, GormQuery, GormPagePagination, error) {
	var entities []E
	var total int64

	stmt := r.GetDBWithCtx(ctx).WithContext(ctx).Model(&entities)

	stmt = r.HandleQuery(stmt, query)

	if err := stmt.Count(&total).Error; err != nil {
		return nil, nil, nil, err
	}

	stmt = stmt.Limit(paging.GetLimit()).Offset(paging.GetOffset())

	// preload should be call in there.
	stmt = r.GetPreload(stmt)
	if err := stmt.Error; err != nil {
		return nil, nil, nil, err
	}

	paging.SetTotal(total)
	return entities, query, paging, nil
}
