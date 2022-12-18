package domain

import "context"

type FindWith int

const (
	FindWithId FindWith = iota
	FindWithRelationId
	FindWithName
)

type (
	// ICreateRepository - Create
	ICreateRepository[T any] interface {
		Create(ctx context.Context, param *T) error
	}

	// IReadAllRepository - ReadAllNoCondition
	IReadAllRepository[T any] interface {
		All(ctx context.Context) (data []*T, err error)
	}

	// IReadAllWhereInRepository - ReadAllWhereIn
	IReadAllWhereInRepository[T any] interface {
		AllWhereIn(ctx context.Context, id []int) (data []*T, err error)
	}

	// IReadOneRepository - ReadOne/Show
	IReadOneRepository[T any] interface {
		Find(ctx context.Context, key FindWith, val any) (data *T, err error)
	}

	// IUpdateRepository - Update
	IUpdateRepository[T any] interface {
		Update(ctx context.Context, param *T) error
	}

	// IDeleteRepository - Delete
	IDeleteRepository[T any] interface {
		Delete(ctx context.Context, param *T) error
	}
)
