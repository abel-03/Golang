//go:build !solution

package structtags

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

func Unpack(req *http.Request, ptr interface{}) error {
	if err := req.ParseForm(); err != nil {
		return err
	}

	v := reflect.ValueOf(ptr).Elem()
	typeOfPtr := v.Type()

	fields := make(map[string]reflect.Value)
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := typeOfPtr.Field(i)
		tag := fieldType.Tag.Get("http")
		if tag == "" {
			tag = strings.ToLower(fieldType.Name)
		}
		fields[tag] = field
	}

	for name, values := range req.Form {
		f, ok := fields[name]
		if !ok {
			continue
		}

		for _, value := range values {
			if f.Kind() == reflect.Slice {
				elem := reflect.New(f.Type().Elem()).Elem()
				if err := populate(elem, value); err != nil {
					return fmt.Errorf("%s: %v", name, err)
				}
				f.Set(reflect.Append(f, elem))
			} else {
				if err := populate(f, value); err != nil {
					return fmt.Errorf("%s: %v", name, err)
				}
			}
		}
	}
	return nil
}

func populate(v reflect.Value, value string) error {
	switch v.Kind() {
	case reflect.String:
		v.SetString(value)

	case reflect.Int:
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		v.SetInt(i)

	case reflect.Bool:
		b, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		v.SetBool(b)

	default:
		return fmt.Errorf("unsupported kind %s", v.Type())
	}
	return nil
}
