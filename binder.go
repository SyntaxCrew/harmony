package harmony

import (
	"encoding/json"
	"errors"
	"net/http"
	"reflect"
	"strconv"
)

type (
	Binder interface {
		Bind(ctx Context, dest any) error
		BindJSON(ctx Context, dest any) error
	}

	binder struct{}
)

func newBinder() Binder {
	return &binder{}
}

func (b *binder) BindJSON(ctx Context, dest any) error {
	return json.NewDecoder(ctx.Request().Body).Decode(dest)
}

func (b *binder) Bind(ctx Context, dest any) (err error) {
	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(error); ok {
				err = e
			}
		}
	}()
	values := reflect.ValueOf(dest)
	if values.Kind() != reflect.Pointer {
		return errors.New("harmony: binder: failed to bind due to dest is not a pointer")
	}
	fields := reflect.TypeOf(dest).Elem()

	num := fields.NumField()
	for i := 0; i < num; i++ {

		field := fields.Field(i)
		value := values.Elem().Field(i)

		val, ok := b.getRequestValueFromTag(ctx, field)
		if !ok || val == "" {
			continue
		}

		err = b.setFieldValue(value, val)
		if err != nil {
			return err
		}
	}

	switch ctx.Request().Method {
	case http.MethodGet:
		return nil
	default:
		return b.BindJSON(ctx, dest)
	}
}

func (b *binder) getRequestValueFromTag(ctx Context, field reflect.StructField) (string, bool) {
	var (
		v  string
		ok bool
	)

	if v, ok = field.Tag.Lookup("query"); ok {
		val := ctx.QueryString(v)
		return val, true
	}

	if v, ok = field.Tag.Lookup("path"); ok {
		val := ctx.PathParam(v)
		return val, true
	}

	return "", false
}

func (b *binder) setFieldValue(value reflect.Value, val string) error {
	if !value.CanSet() {
		return nil
	}

	switch value.Kind() {
	case reflect.String:
		value.SetString(val)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return err
		}
		value.SetInt(v)
	case reflect.Float32, reflect.Float64:
		v, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return err
		}
		value.SetFloat(v)
	case reflect.Bool:
		v, err := strconv.ParseBool(val)
		if err != nil {
			return err
		}
		value.SetBool(v)
	}
	return nil
}
