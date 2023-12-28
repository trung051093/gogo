package common

import (
	"context"

	"github.com/go-playground/validator"
	"gorm.io/gorm"
)

type Service[E GormEntity, R any] interface {
	Create(ctx context.Context, dto CreateDto[E]) (*E, error)
	CreateList(ctx context.Context, dto CreateListDto[E]) ([]E, error)
	UpdateByID(ctx context.Context, id any, dto UpdateDto[E]) (*E, error)
	Updates(ctx context.Context, entity *E, dto UpdateDto[E]) (*E, error)
	Save(ctx context.Context, entity *E) (*E, error)
	DeleteById(ctx context.Context, Id any) error
	DeleteByIds(ctx context.Context, Ids any) error
	SearchCursorPaging(ctx context.Context, dto SearchCursorPagingDto) ([]E, GormQuery, GormCursorPagination, error)
	Search(ctx context.Context, dto SearchDto) ([]E, error)
	FindById(ctx context.Context, id any) (*E, error)
	FindByIds(ctx context.Context, ids any) ([]E, error)
	GetAll(ctx context.Context) ([]E, error)
	GetRepository() R
}

type service[E GormEntity, R Repository[E]] struct {
	repository R
	validator  *validator.Validate
}

func NewService[E GormEntity, R Repository[E]](repo R) Service[E, R] {
	validator := validator.New()
	return &service[E, R]{repo, validator}
}

func (s *service[E, R]) WithTX(trxHandle *gorm.DB) Service[E, R] {
	s.repository = s.repository.WithTX(trxHandle).(R)
	return s
}

func (s *service[E, R]) GetRepository() R {
	return s.repository
}

func (s *service[E, R]) validateDto(dto ValidateDto) error {
	return dto.Validate()
}

func (s *service[E, R]) Create(ctx context.Context, dto CreateDto[E]) (*E, error) {
	if err := s.validateDto(dto); err != nil {
		return nil, err
	}
	entity := dto.ToEntity()
	return s.repository.Create(ctx, entity)
}

func (s *service[E, R]) CreateList(ctx context.Context, dto CreateListDto[E]) ([]E, error) {
	if err := s.validateDto(dto); err != nil {
		return nil, err
	}
	entities := dto.ToEntities()
	return s.repository.CreateList(ctx, entities)
}

func (s *service[E, R]) UpdateByID(ctx context.Context, id any, dto UpdateDto[E]) (*E, error) {
	if err := s.validateDto(dto); err != nil {
		return nil, err
	}
	currentEntity, err := s.FindById(ctx, id)
	if err != nil {
		panic(err)
	}
	return s.repository.Save(ctx, dto.ToEntity(currentEntity))
}

func (s *service[E, R]) Updates(ctx context.Context, entity *E, dto UpdateDto[E]) (*E, error) {
	if err := s.validateDto(dto); err != nil {
		return nil, err
	}
	return s.repository.Updates(ctx, entity, dto.ToMapInterface())
}

func (s *service[E, R]) Save(ctx context.Context, entity *E) (*E, error) {
	return s.repository.Save(ctx, entity)
}

func (s *service[E, R]) DeleteById(ctx context.Context, id any) error {
	return s.repository.DeleteById(ctx, id)
}

func (s *service[E, R]) DeleteByIds(ctx context.Context, ids any) error {
	return s.repository.DeleteByIds(ctx, ids)
}

func (s *service[E, R]) SearchCursorPaging(ctx context.Context, dto SearchCursorPagingDto) ([]E, GormQuery, GormCursorPagination, error) {
	if err := s.validateDto(dto); err != nil {
		return nil, nil, nil, err
	}
	paging := dto.ToPagination()
	query := dto.ToQuery()
	return s.repository.FindWithCursorPagination(ctx, query, paging)
}

func (s *service[E, R]) Search(ctx context.Context, dto SearchDto) ([]E, error) {
	if err := s.validateDto(dto); err != nil {
		return nil, err
	}
	query := dto.ToQuery()
	return s.repository.FindByCond(ctx, query)
}

func (s *service[E, R]) FindById(ctx context.Context, id any) (*E, error) {
	return s.repository.FindById(ctx, id)
}

func (s *service[E, R]) FindByIds(ctx context.Context, ids any) ([]E, error) {
	return s.repository.FindByIds(ctx, ids)
}

func (s *service[E, R]) GetAll(ctx context.Context) ([]E, error) {
	return s.repository.GetAll(ctx)
}
