package urfave

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/urfave/cli/v2"
	"github.com/xdeeper/gokits/cmder/base"
)

type argsParser[T base.ArgType] struct {
	conv func(string) (T, error)
}

func (s argsParser[T]) read(ctx *cli.Context) (any, error) {
	switch ctx.Args().Len() {
	case 1:
		v, err := s.conv(strings.TrimSpace(ctx.Args().Get(0)))
		if err != nil {
			return nil, err
		}
		return base.NewArg[T](v), nil
	case 0:
		return nil, fmt.Errorf("no arg given...")
	default:
		return nil, fmt.Errorf("too many(%d) args given...", ctx.Args().Len())
	}
}

type argsSliceParser[T base.ArgType] struct {
	conv func(string) (T, error)
}

func (s argsSliceParser[T]) read(ctx *cli.Context) (any, error) {
	switch ctx.Args().Len() {
	case 0:
		return nil, fmt.Errorf("no arg given...")
	default:
		slice := make([]T, ctx.Args().Len())
		var err error
		for i, v := range ctx.Args().Slice() {
			slice[i], err = s.conv(strings.TrimSpace(v))
			if err != nil {
				var a T
				return nil, fmt.Errorf("failed conv(%s) to type %T", v, a)
			}
		}
		return base.NewSliceArg[T](slice), nil
	}
}

type iValReader interface {
	read(*cli.Context) (any, error)
}

func getArgSliceSetter(rt reflect.Type) iValReader {
	switch {
	case typeOfArgT[bool](rt):
		return argsSliceParser[bool]{conv: strconv.ParseBool}
	case typeOfArgT[string](rt):
		return argsSliceParser[string]{conv: func(s string) (string, error) { return s, nil }}
	case typeOfArgT[time.Duration](rt):
		return argsSliceParser[time.Duration]{conv: time.ParseDuration}
	case typeOfArgT[float32](rt):
		return argsSliceParser[float32]{conv: convNum[float32]()}
	case typeOfArgT[float64](rt):
		return argsSliceParser[float64]{conv: convNum[float64]()}
	case typeOfArgT[int](rt):
		return argsSliceParser[int]{conv: convNum[int]()}
	case typeOfArgT[int8](rt):
		return argsSliceParser[int8]{conv: convNum[int8]()}
	case typeOfArgT[int16](rt):
		return argsSliceParser[int16]{conv: convNum[int16]()}
	case typeOfArgT[int32](rt):
		return argsSliceParser[int32]{conv: convNum[int32]()}
	case typeOfArgT[int64](rt):
		return argsSliceParser[int64]{conv: convNum[int64]()}
	case typeOfArgT[uint](rt):
		return argsSliceParser[uint]{conv: convNum[uint]()}
	case typeOfArgT[uint8](rt):
		return argsSliceParser[uint8]{conv: convNum[uint8]()}
	case typeOfArgT[uint32](rt):
		return argsSliceParser[uint32]{conv: convNum[uint32]()}
	case typeOfArgT[uint64](rt):
		return argsSliceParser[uint64]{conv: convNum[uint64]()}
	default:
		panic(fmt.Errorf("cannot get argsParser of type:%s", rt))
	}
}

func getArgsSetter(rt reflect.Type) iValReader {
	switch {
	case typeOfArgT[bool](rt):
		return argsParser[bool]{conv: strconv.ParseBool}
	case typeOfArgT[string](rt):
		return argsParser[string]{conv: func(s string) (string, error) { return s, nil }}
	case typeOfArgT[time.Duration](rt):
		return argsParser[time.Duration]{conv: time.ParseDuration}
	case typeOfArgT[float32](rt):
		return argsParser[float32]{conv: convNum[float32]()}
	case typeOfArgT[float64](rt):
		return argsParser[float64]{conv: convNum[float64]()}
	case typeOfArgT[int](rt):
		return argsParser[int]{conv: convNum[int]()}
	case typeOfArgT[int8](rt):
		return argsParser[int8]{conv: convNum[int8]()}
	case typeOfArgT[int16](rt):
		return argsParser[int16]{conv: convNum[int16]()}
	case typeOfArgT[int32](rt):
		return argsParser[int32]{conv: convNum[int32]()}
	case typeOfArgT[int64](rt):
		return argsParser[int64]{conv: convNum[int64]()}
	case typeOfArgT[uint](rt):
		return argsParser[uint]{conv: convNum[uint]()}
	case typeOfArgT[uint8](rt):
		return argsParser[uint8]{conv: convNum[uint8]()}
	case typeOfArgT[uint32](rt):
		return argsParser[uint32]{conv: convNum[uint32]()}
	case typeOfArgT[uint64](rt):
		return argsParser[uint64]{conv: convNum[uint64]()}
	default:
		panic(fmt.Errorf("cannot get argsParser of type:%s", rt))
	}
}

func convNum[T interface {
	int | int8 | int16 | int32 | int64 |
		uint | uint8 | uint16 | uint32 | uint64 |
		float32 | float64
}]() func(s string) (T, error) {
	var t T
	rt := reflect.TypeOf(t)
	if strings.HasPrefix(rt.Name(), "int") {
		nStr := strings.TrimPrefix(rt.Name(), "int")
		n, _ := strconv.ParseInt(nStr, 10, 32)
		return func(s string) (T, error) {
			if v, err := strconv.ParseInt(s, 10, int(n)); err != nil {
				return 0, err
			} else {
				return T(v), nil
			}
		}
	}
	switch {
	case strings.HasPrefix(rt.Name(), "int"):
		nStr := strings.TrimPrefix(rt.Name(), "int")
		n, _ := strconv.ParseInt(nStr, 10, 32)
		return func(s string) (T, error) {
			if v, err := strconv.ParseInt(s, 10, int(n)); err != nil {
				return 0, err
			} else {
				return T(v), nil
			}
		}
	case strings.HasPrefix(rt.Name(), "uint"):
		nStr := strings.TrimPrefix(rt.Name(), "uint")
		n, _ := strconv.ParseInt(nStr, 10, 32)
		return func(s string) (T, error) {
			if v, err := strconv.ParseUint(s, 10, int(n)); err != nil {
				return 0, err
			} else {
				return T(v), nil
			}
		}
	case rt.Name() == "float32":
		return func(s string) (T, error) {
			if v, err := strconv.ParseFloat(s, 32); err != nil {
				return 0, err
			} else {
				return T(v), nil
			}
		}
	default:
		return func(s string) (T, error) {
			if v, err := strconv.ParseFloat(s, 64); err != nil {
				return 0, err
			} else {
				return T(v), nil
			}
		}
	}
}

func typeOfArgT[T base.ArgType](rt reflect.Type) bool {
	return reflect.TypeOf((*base.Arg[T])(nil)).Elem().Name() == rt.Name()
}
