package utils

import (
	"database/sql"
	"net/http"
)

type ServiceError struct {
	Code    int
	Message any
}

func ValidateDataRow[T any](data *T, err error) (valueData *T, errData *ServiceError) {
	errData = checkError(err)

	return data, errData
}

func ValidateDataRows[T any](data []*T, err error) (valueData []*T, errData *ServiceError) {
	errData = checkError(err)

	return data, errData
}

func ValidatePrimitiveValue[T any](data T, err error) (valueData T, errData *ServiceError) {
	errData = checkError(err)

	return data, errData
}

func checkError(err error) *ServiceError {
	var errData *ServiceError

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			errData = &ServiceError{
				Code:    http.StatusNotFound,
				Message: err.Error(),
			}
		default:
			errData = &ServiceError{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			}
		}
	}

	return errData
}
