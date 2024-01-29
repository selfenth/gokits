package base

import "time"

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
type flag[T FlagType] struct{ val Optional[T] }

func (a flag[T]) Value() Optional[T]  { return a.val }
func NewFlag[T FlagType](a T) Flag[T] { return flag[T]{val: optional[T]{val: a, set: true}} }
func NotsetFlag[T FlagType]() Flag[T] { return flag[T]{val: optional[T]{}} }

type RequiredFlag[T FlagType] interface{ Value() T }
type requiredFlag[T FlagType] struct{ val T }

func (a requiredFlag[T]) Value() T                    { return a.val }
func NewRequiredFlag[T FlagType](a T) RequiredFlag[T] { return requiredFlag[T]{val: a} }

type optional[T any] struct {
	val T
	set bool
}

func (a optional[T]) IsSet() bool { return a.set }
func (a optional[T]) Value() T    { return a.val }
func (a optional[T]) ValOrDefault(defVal T) T {
	if a.set {
		return a.val
	}
	return defVal
}
