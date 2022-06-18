package common

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
)

var (
	ErrNotFound = errors.New("not found")
)

type AppError struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
	RootErr    error  `json:"-"`
	ErrorKey   string `json:"errorKey"`
	Log        string `json:"log"`
}

func FailOnError(err error, msg string) {
	if err != nil {
		log.Printf("%s: %s", msg, err)
	}
}

func NewErrorResponse(rootErr error, message string, log, key string) *AppError {
	return &AppError{
		StatusCode: http.StatusBadRequest,
		Message:    message,
		RootErr:    rootErr,
		ErrorKey:   key,
		Log:        log,
	}
}

func NewUnauthorized(rootErr error, message string, key string) *AppError {
	return &AppError{
		StatusCode: http.StatusUnauthorized,
		Message:    message,
		RootErr:    rootErr,
		ErrorKey:   key,
	}
}

func NewFullErrorResponse(statusCode int, rootErr error, message string, log, key string) *AppError {
	return &AppError{
		StatusCode: statusCode,
		Message:    message,
		RootErr:    rootErr,
		ErrorKey:   key,
		Log:        log,
	}
}

func NewCustomError(rootErr error, message string, key string) *AppError {
	if rootErr != nil {
		return NewErrorResponse(rootErr, message, rootErr.Error(), key)
	}
	return NewErrorResponse(errors.New(message), message, message, key)
}

func (e *AppError) RootError() error {
	if err, ok := e.RootErr.(*AppError); ok {
		return err.RootError()
	}
	return e.RootErr
}

func (e *AppError) Error() string {
	return e.RootError().Error()
}

func ErrorDB(err error) *AppError {
	return NewFullErrorResponse(
		http.StatusInternalServerError,
		err,
		"Something went wrong with DB",
		err.Error(),
		"DB_ERROR",
	)
}

func ErrorInternal(err error) *AppError {
	return NewFullErrorResponse(
		http.StatusInternalServerError,
		err,
		"Something went wrong with server",
		err.Error(),
		"INTERNAL_ERROR",
	)
}

func ErrorNotFound(entity string, err error) *AppError {
	return NewErrorResponse(
		err,
		fmt.Sprintf("Cannot found %s", strings.ToLower(entity)),
		err.Error(),
		"NOT_FOUND",
	)
}

func ErrorInvalidRequest(entity string, err error) *AppError {
	return NewErrorResponse(
		err,
		"Invalid request",
		err.Error(),
		"INVALID_REQUEST",
	)
}

func ErrorCannotListEntity(entity string, err error) *AppError {
	return NewCustomError(
		err,
		fmt.Sprintf("Cannot list %s", strings.ToLower(entity)),
		fmt.Sprintf("ERROR_CANNOT_LIST_%s", strings.ToUpper(entity)),
	)
}

func ErrorCannotFoundEntity(entity string, err error) *AppError {
	return NewCustomError(
		err,
		fmt.Sprintf("Cannot found %s", strings.ToLower(entity)),
		fmt.Sprintf("ERROR_CANNOT_FOUND_%s", strings.ToUpper(entity)),
	)
}

func ErrorCannotUpdateEntity(entity string, err error) *AppError {
	return NewCustomError(
		err,
		fmt.Sprintf("Cannot update %s", strings.ToLower(entity)),
		fmt.Sprintf("ERROR_CANNOT_UPDATE_%s", strings.ToUpper(entity)),
	)
}

func ErrorCannotCreateEntity(entity string, err error) *AppError {
	return NewCustomError(
		err,
		fmt.Sprintf("Cannot create %s", strings.ToLower(entity)),
		fmt.Sprintf("ERROR_CANNOT_CREATE_%s", strings.ToUpper(entity)),
	)
}

func ErrorCannotDeleteEntity(entity string, err error) *AppError {
	return NewCustomError(
		err,
		fmt.Sprintf("Cannot delete %s", strings.ToLower(entity)),
		fmt.Sprintf("ERROR_CANNOT_DELETE_%s", strings.ToUpper(entity)),
	)
}

func ErrorEntityExisted(entity string, err error) *AppError {
	return NewCustomError(
		err,
		fmt.Sprintf("%s already exists", strings.ToLower(entity)),
		fmt.Sprintf("%s_EXISTS", strings.ToUpper(entity)),
	)
}

func ErrorEntityDeleted(entity string, err error) *AppError {
	return NewCustomError(
		err,
		fmt.Sprintf("%s has been deleted", strings.ToLower(entity)),
		fmt.Sprintf("%s_DELETED", strings.ToUpper(entity)),
	)
}

func ErrorUnauthorized() *AppError {
	return NewUnauthorized(
		errors.New("Unauthorized"),
		"Unauthorized",
		"UNAUTHORIZED",
	)
}
