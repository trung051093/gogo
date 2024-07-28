package common

type ValidateDto interface {
	Validate() error
}

type CreateListDto[E GormEntity] interface {
	ValidateDto
	ToEntities() []*E
}

type CreateDto[E GormEntity] interface {
	ValidateDto
	ToEntity() *E
}

type UpdateDto[E GormEntity] interface {
	ValidateDto
	ToEntity(currentEntity *E) *E
	ToMapInterface() map[string]interface{}
}

type DeleteDto[E GormEntity] interface {
	ValidateDto
	ToEntity() *E
}

type SearchDto interface {
	ValidateDto
	ToQuery() GormQuery
}

type SearchCursorPagingDto interface {
	ValidateDto
	SearchDto
	ToPagination() *CursorPagination
}
