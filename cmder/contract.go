package cmder

import (
	"golang.org/x/exp/constraints"
)

type ArgType interface {
	constraints.Integer | ~bool | ~string
}

type Arg[T ArgType] interface {
	Value() T
}

type SliceArg[T ArgType] interface {
	Value() []T
}
type FlagType interface {
	~bool | ~string | ~float64 | ~int | ~int64 | ~uint | ~uint64
}

type Optional[T any] interface {
	IsSet() bool
	Value() T
	ValOrDefault(a T) T
}
type Flag[T FlagType] interface{ Value() Optional[T] }
type RequiredFlag[T FlagType] interface{ Value() T }
