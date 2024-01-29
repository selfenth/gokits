package urfave

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/stoewer/go-strcase"
	"github.com/urfave/cli/v2"
	"github.com/xdeeper/gokits/cmder/base"
)

type Runable interface {
	Run(ctx context.Context) error
}

type urfaveRunable struct{ app *cli.App }

func (a urfaveRunable) Run(ctx context.Context) error { return a.app.RunContext(ctx, os.Args) }

func Build(a any) Runable {
	app := cli.NewApp()
	rv := reflect.Indirect(reflect.ValueOf(a))
	app.Name = strcase.SnakeCase(rv.Type().Name())

	c := cli.Command{}
	buildCommand(app.Name, rv, &c)

	app.Usage = c.Usage
	app.UsageText = c.UsageText
	app.ArgsUsage = c.ArgsUsage
	app.Action = c.Action
	app.Description = c.Description
	app.Flags = c.Flags
	app.Commands = c.Subcommands

	return urfaveRunable{app: app}
}

func buildCommand(name string, rv reflect.Value, c *cli.Command) error {
	cmd := urfaveCommand{rt: rv.Type()}
	if desc := rv.FieldByName(base.C_Description); desc.IsValid() {
		c.Description, _ = desc.Interface().(string)
	}
	c.Name = name
	for i := 0; i < rv.NumField(); i++ {
		field := rv.Type().Field(i)
		if field.Tag.Get("cmd") != "" {
			fmt.Printf("get sub command:%s\n", field.Tag.Get("cmd"))
			subCmd := cli.Command{}
			if err := buildCommand(field.Tag.Get("cmd"), rv.Field(i), &subCmd); err != nil {
				return err
			}
			c.Subcommands = append(c.Subcommands, &subCmd)
		}
		if field.Type.PkgPath() != "github.com/xdeeper/gokits/cmder" {
			continue
		}
		switch {
		case strings.HasPrefix(field.Type.Name(), "Arg["):
			c.ArgsUsage = field.Tag.Get(base.C_Usage)
			cmd.setters = append(cmd.setters, fieldSetter{
				index: i,
				name:  field.Name,
				r:     getArgsSetter(field.Type),
			})
			if len(c.ArgsUsage) > 0 {
				c.ArgsUsage = fmt.Sprintf("[ arguments(single value): %s ]", c.ArgsUsage)
			}
		case strings.HasPrefix(field.Type.Name(), "SliceArg["):
			c.ArgsUsage = field.Tag.Get(base.C_Usage)
			cmd.setters = append(cmd.setters, fieldSetter{
				index: i,
				name:  field.Name,
				r:     getArgSliceSetter(field.Type),
			})
			if len(c.ArgsUsage) > 0 {
				c.ArgsUsage = fmt.Sprintf("[ arguments(slice): %s ]", c.ArgsUsage)
			}
		case strings.HasPrefix(field.Type.Name(), "RequiredFlag["):
			a := parseFlag(field, true)
			cmd.setters = append(cmd.setters, fieldSetter{
				index: i,
				r:     a,
			})
		case strings.HasPrefix(field.Type.Name(), "Flag["):
			a := parseFlag(field, false)
			cmd.setters = append(cmd.setters, fieldSetter{
				index: i,
				name:  field.Name,
				r:     a,
			})
		}
	}
	if _, ok := rv.Interface().(Runable); ok {
		c.Action = cmd.Action
	}
	return nil
}

type fieldSetter struct {
	index int
	name  string
	r     iValReader
}

func (s *fieldSetter) set(ctx *cli.Context, a reflect.Value) error {
	v, err := s.r.read(ctx)
	if err != nil {
		return err
	}
	rv := reflect.Indirect(a)
	fv := rv.Field(s.index)
	switch {
	case !fv.CanSet():
		return fmt.Errorf("cannot set field:%d-%s", s.index, rv.Type().Field(s.index).Name)
	case !reflect.TypeOf(v).AssignableTo(fv.Type()):
		return fmt.Errorf("cannot assign value(%T) to field of type:%T", v, fv.Type())
	default:
		rv.Field(s.index).Set(reflect.ValueOf(v))
		return nil
	}
}

type urfaveCommand struct {
	rt      reflect.Type
	setters []fieldSetter
}

func (c urfaveCommand) Action(ctx *cli.Context) error {
	a := reflect.New(c.rt)
	for _, s := range c.setters {
		if err := s.set(ctx, a); err != nil {
			return err
		}
	}
	return a.Interface().(Runable).Run(ctx.Context)
}
