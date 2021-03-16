package validate

import (
	"errors"
	"fmt"
	"reflect"
)

const validateTagName = "validate"

func Validate(t interface{}) error {
	e := reflect.ValueOf(t).Elem()
	for i := 0; i < e.NumField(); i++ {
		varName := e.Type().Field(i).Name
		varType := e.Type().Field(i).Type.Name()
		varValue := e.Field(i).Interface()
		tag := e.Type().Field(i).Tag.Get(validateTagName)

		if tag != "required" {
			continue
		}

		switch varType {
		case "int":
			if varValue == 0 {
				return errors.New(fmt.Sprintf("Field %s is required", varName))
			}
		case "string":
			if varValue == "" {
				return errors.New(fmt.Sprintf("Field %s is required", varName))
			}
		}
	}

	return nil
}
