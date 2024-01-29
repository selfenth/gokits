package urfave

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/stoewer/go-strcase"
	"github.com/urfave/cli/v2"
	"github.com/xdeeper/gokits/cmder/base"
)

var (
	_ iValReader = iflag{}
)

type iflag struct {
	name     string
	required bool
	rt       reflect.Type
	a        cli.Flag
}

// read implements iValReader.
func (a iflag) read(ctx *cli.Context) (any, error) {
	if a.required {
		if !ctx.IsSet(a.name) {
			return nil, fmt.Errorf("flag:%s not set", a.name)
		}
		switch tv := ctx.Value(a.name).(type) {
		case string:
			if len(tv) == 0 {
				return nil, fmt.Errorf("flag:%s not set", a.name)
			}
			return base.NewRequiredFlag(tv), nil
		case bool:
			return base.NewRequiredFlag(tv), nil
		case int:
			return base.NewRequiredFlag(tv), nil
		case int64:
			return base.NewRequiredFlag(tv), nil
		case uint:
			return base.NewRequiredFlag(tv), nil
		case uint64:
			return base.NewRequiredFlag(tv), nil
		case time.Duration:
			return base.NewRequiredFlag(tv), nil
		case float64:
			return base.NewRequiredFlag(tv), nil
		default:
			return nil, fmt.Errorf("type(%T) of flag:%s is not supported", ctx.Value(a.name), a.name)
		}
	}
	if !ctx.IsSet(a.name) {
		return a.getNotSet(), nil
	}
	switch tv := ctx.Value(a.name).(type) {
	case string:
		if len(tv) == 0 {
			return nil, fmt.Errorf("flag:%s not set", a.name)
		}
		return base.NewFlag(tv), nil
	case bool:
		return base.NewFlag(tv), nil
	case int:
		return base.NewFlag(tv), nil
	case int64:
		return base.NewFlag(tv), nil
	case uint:
		return base.NewFlag(tv), nil
	case uint64:
		return base.NewFlag(tv), nil
	case time.Duration:
		return base.NewFlag(tv), nil
	case float64:
		return base.NewFlag(tv), nil
	default:
		return nil, fmt.Errorf("type(%T) of flag:%s is not supported", ctx.Value(a.name), a.name)
	}
}

func (a iflag) getNotSet() any {
	switch {
	case typeOfFlagT[bool](a.rt):
		return base.NotsetFlag[bool]()
	case typeOfFlagT[string](a.rt):
		return base.NotsetFlag[string]()
	case typeOfFlagT[time.Duration](a.rt):
		return base.NotsetFlag[time.Duration]()
	case typeOfFlagT[float64](a.rt):
		return base.NotsetFlag[float64]()
	case typeOfFlagT[int](a.rt):
		return base.NotsetFlag[int]()
	case typeOfFlagT[int64](a.rt):
		return base.NotsetFlag[int64]()
	case typeOfFlagT[uint](a.rt):
		return base.NotsetFlag[uint]()
	case typeOfFlagT[uint64](a.rt):
		return base.NotsetFlag[uint64]()
	default:
		panic(fmt.Errorf("cannot get argsParser of type:%s", a.rt.String()))
	}
}

func parseFlag(f reflect.StructField, required bool) iflag {
	var (
		hidden bool
	)
	a := iflag{name: strcase.KebabCase(f.Name), required: required, rt: f.Type}
	flagTag := f.Tag.Get(base.C_Flag)
	if len(flagTag) > 0 {
		attrs := strings.Split(flagTag, ",")
		a.name = attrs[0]
		for _, s := range attrs[1:] {
			switch s {
			case base.C_Hidden:
				hidden = true
			}
		}
	}
	usage := f.Tag.Get(base.C_Usage)
	defaultVal := f.Tag.Get(base.C_Default)

	switch {
	case typeOfFlagT[bool](f.Type):
	case typeOfRequiredFlagT[bool](f.Type):
		a.a = &cli.BoolFlag{Name: a.name, Usage: usage, DefaultText: defaultVal, Hidden: hidden, Required: required}
	case typeOfFlagT[string](f.Type):
	case typeOfRequiredFlagT[string](f.Type):
		a.a = &cli.StringFlag{Name: a.name, Usage: usage, DefaultText: defaultVal, Hidden: hidden, Required: required}
	case typeOfFlagT[time.Duration](f.Type):
	case typeOfRequiredFlagT[time.Duration](f.Type):
		a.a = &cli.DurationFlag{Name: a.name, Usage: usage, DefaultText: defaultVal, Hidden: hidden, Required: required}
	case typeOfFlagT[float64](f.Type):
	case typeOfRequiredFlagT[float64](f.Type):
		a.a = &cli.Float64Flag{Name: a.name, Usage: usage, DefaultText: defaultVal, Hidden: hidden, Required: required}
	case typeOfFlagT[int](f.Type):
	case typeOfRequiredFlagT[int](f.Type):
		a.a = &cli.IntFlag{Name: a.name, Usage: usage, DefaultText: defaultVal, Hidden: hidden, Required: required}
	case typeOfFlagT[int64](f.Type):
	case typeOfRequiredFlagT[int64](f.Type):
		a.a = &cli.Int64Flag{Name: a.name, Usage: usage, DefaultText: defaultVal, Hidden: hidden, Required: required}
	case typeOfFlagT[uint](f.Type):
	case typeOfRequiredFlagT[uint](f.Type):
		a.a = &cli.UintFlag{Name: a.name, Usage: usage, DefaultText: defaultVal, Hidden: hidden, Required: required}
	case typeOfFlagT[uint64](f.Type):
	case typeOfRequiredFlagT[uint64](f.Type):
		a.a = &cli.Uint64Flag{Name: a.name, Usage: usage, DefaultText: defaultVal, Hidden: hidden, Required: required}
	default:
		panic(fmt.Errorf("cannot get argsParser of type:%s", f.Type.String()))
	}
	return a
}

func typeOfFlagT[T base.FlagType](rt reflect.Type) bool {
	return reflect.TypeOf((*base.Flag[T])(nil)).Elem().Name() == rt.Name()
}
func typeOfRequiredFlagT[T base.FlagType](rt reflect.Type) bool {
	return reflect.TypeOf((*base.RequiredFlag[T])(nil)).Elem().Name() == rt.Name()
}
