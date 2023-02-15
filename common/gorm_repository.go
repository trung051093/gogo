package common

import (
	"context"
	"errors"

	"github.com/pilagod/gorm-cursor-paginator/v2/cursor"
	"github.com/pilagod/gorm-cursor-paginator/v2/paginator"
	"gorm.io/gorm"
)

type GormEntity interface {
	TableName() string
}

type GormPagePagination interface {
	GetLimit() int
	SetLimit(limit int)
	GetTotal() int64
	SetTotal(total int64)
	GetOffset() int
	SetOffset(offset int)
}

type GormCursorPagination interface {
	GetLimit() int
	SetLimit(limit int)
	GetTotal() int64
	SetTotal(total int64)
	GetCursor() cursor.Cursor
	SetCursor(c cursor.Cursor)
	GetPaginator() *paginator.Paginator
}

type GormFilter interface {
	GetOrders() string
	SetOrders(order string) string
	GetFieldSelect() string
}

type Repository[E GormEntity] interface {
	GetDB() *gorm.DB
	SetPreloadKeys(preloadKeys ...string)
	SetCustomPreloadKeys(preloadKeys map[string]interface{})
	GetPreload(stmt *gorm.DB) *gorm.DB
	Save(ctx context.Context, entity *E) (*E, error)
	CreateList(ctx context.Context, entities []E) ([]E, error)
	Delete(ctx context.Context, entity *E) (*E, error)
	GetAll(ctx context.Context) ([]E, error)
	Count(ctx context.Context) (int64, error)
	CountWithCondition(ctx context.Context, cond map[string]interface{}) (int64, error)
	FindByID(ctx context.Context, id any) (*E, error)
	FindOneByCond(ctx context.Context, cond map[string]interface{}) (*E, error)
	FindByCond(ctx context.Context, cond map[string]interface{}) ([]E, error)
	FindWithCursorPagination(ctx context.Context, cond map[string]interface{}, filter GormFilter, paging GormCursorPagination) ([]E, GormFilter, GormCursorPagination, error)
	FindWithPagePagination(ctx context.Context, cond map[string]interface{}, filter GormFilter, paging GormPagePagination) ([]E, GormFilter, GormPagePagination, error)
}

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

func (r *repository[E]) WithTX(tx *gorm.DB) Repository[E] {
	return &repository[E]{
		DB:          tx,
		PreloadKeys: r.PreloadKeys,
	}
}

func (r *repository[E]) GetDB() *gorm.DB {
	return r.DB
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

func (r *repository[E]) Save(ctx context.Context, entity *E) (*E, error) {
	if err := r.DB.WithContext(ctx).Model(entity).Save(&entity).Error; err != nil {
		return nil, err
	}
	return entity, nil
}

func (r *repository[E]) CreateList(ctx context.Context, entities []E) ([]E, error) {
	var e E
	if err := r.DB.WithContext(ctx).Table(e.TableName()).Create(&entities).Error; err != nil {
		return nil, err
	}
	return entities, nil
}

func (r *repository[E]) Delete(ctx context.Context, entity *E) (*E, error) {
	if err := r.DB.WithContext(ctx).Model(entity).Delete(&entity).Error; err != nil {
		return nil, err
	}
	return entity, nil
}

func (r *repository[E]) GetAll(ctx context.Context) ([]E, error) {
	var entities []E
	var e E
	stmt := r.DB.WithContext(ctx).Table(e.TableName())
	stmt = r.GetPreload(stmt)
	if err := stmt.Find(&entities).Error; err != nil {
		return nil, err
	}
	return entities, nil
}

func (r *repository[E]) Count(ctx context.Context) (int64, error) {
	var total int64
	var entity E
	stmt := r.DB.WithContext(ctx).Model(entity)
	if err := stmt.Count(&total).Error; err != nil {
		return 0, err
	}

	return total, nil
}

func (r *repository[E]) CountWithCondition(ctx context.Context, cond map[string]interface{}) (int64, error) {
	var total int64
	var entity E
	stmt := r.DB.WithContext(ctx).Model(entity).Where(cond)
	if err := stmt.Count(&total).Error; err != nil {
		return 0, err
	}

	return total, nil
}

func (r *repository[E]) FindByID(ctx context.Context, Id any) (*E, error) {
	var entity E
	stmt := r.DB.WithContext(ctx).Model(entity)
	stmt = r.GetPreload(stmt)
	if err := stmt.First(&entity, Id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &entity, nil
}

func (r *repository[E]) FindOneByCond(ctx context.Context, cond map[string]interface{}) (*E, error) {
	var entity E
	stmt := r.DB.WithContext(ctx).Model(entity).Where(cond)
	stmt = r.GetPreload(stmt)
	if err := stmt.First(&entity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &entity, nil
}

func (r *repository[E]) FindByCond(ctx context.Context, cond map[string]interface{}) ([]E, error) {
	var entities []E
	stmt := r.DB.WithContext(ctx).Model(entities).Where(cond)
	stmt = r.GetPreload(stmt)
	if err := stmt.Find(&entities).Error; err != nil {
		return nil, err
	}
	return entities, nil
}

func (r *repository[E]) FindWithCursorPagination(ctx context.Context, cond map[string]interface{}, filter GormFilter, paging GormCursorPagination) ([]E, GormFilter, GormCursorPagination, error) {
	var entities []E
	var total int64

	stmt := r.DB.WithContext(ctx).Model(&entities).Where(cond)

	if err := stmt.Count(&total).Error; err != nil {
		return nil, nil, nil, err
	}

	if filter != nil {
		if filter.GetOrders() != "" {
			stmt = stmt.Order(filter.GetOrders())
		}

		if filter.GetFieldSelect() != "" {
			stmt = stmt.Select(filter.GetFieldSelect())
		}
	}

	// preload should be call in there.
	stmt = r.GetPreload(stmt)
	paginator := paging.GetPaginator()
	result, cursor, err := paginator.Paginate(stmt, &entities)
	if err != nil || result.Error != nil {
		return nil, nil, nil, err
	}

	if err := stmt.Find(&entities).Error; err != nil {
		return nil, nil, nil, err
	}

	paging.SetTotal(total)
	paging.SetCursor(cursor)
	return entities, filter, paging, nil
}

func (r *repository[E]) FindWithPagePagination(ctx context.Context, cond map[string]interface{}, filter GormFilter, paging GormPagePagination) ([]E, GormFilter, GormPagePagination, error) {
	var entities []E
	var total int64

	stmt := r.DB.WithContext(ctx).Model(&entities).Where(cond)

	if err := stmt.Count(&total).Error; err != nil {
		return nil, nil, nil, err
	}

	if filter != nil {
		if filter.GetOrders() != "" {
			stmt = stmt.Order(filter.GetOrders())
		}

		if filter.GetFieldSelect() != "" {
			stmt = stmt.Select(filter.GetFieldSelect())
		}
	}

	stmt = stmt.Limit(paging.GetLimit()).Offset(paging.GetOffset())

	// preload should be call in there.
	stmt = r.GetPreload(stmt)
	if err := stmt.Error; err != nil {
		return nil, nil, nil, err
	}

	paging.SetTotal(total)
	return entities, filter, paging, nil
}
