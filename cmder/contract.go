package cmder

import "time"

type ArgType interface {
	bool | string | time.Duration | float32 | float64 |
		int | int8 | int16 | int32 | int64 |
		uint | uint8 | uint16 | uint32 | uint64
}

type Arg[T ArgType] interface {
	Value() T
}

type SliceArg[T ArgType] interface {
	Value() []T
}
type FlagType interface {
	bool | string | time.Duration | float64 |
		int | int64 | uint | uint64
}

type Optional[T any] interface {
	IsSet() bool
	Value() T
	ValOrDefault(a T) T
}
type Flag[T FlagType] interface{ Value() Optional[T] }
type RequiredFlag[T FlagType] interface{ Value() T }
