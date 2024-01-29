package base

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

type arg[T ArgType] struct{ v T }

func (a arg[T]) Value() T          { return a.v }
func NewArg[T ArgType](v T) Arg[T] { return arg[T]{v: v} }

type sliceArg[T ArgType] struct{ v []T }

func (a sliceArg[T]) Value() []T               { return a.v }
func NewSliceArg[T ArgType](v []T) SliceArg[T] { return sliceArg[T]{v: v} }
