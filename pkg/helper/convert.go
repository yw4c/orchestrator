package helper

import "reflect"

func StructAssign(target interface{}, original interface{}) {
	tValue := reflect.ValueOf(target).Elem()
	oValue := reflect.ValueOf(original).Elem()
	oTypeOfT := oValue.Type()
	for i := 0; i < oValue.NumField(); i++ {
		name := oTypeOfT.Field(i).Name
		if ok := tValue.FieldByName(name).IsValid(); ok {
			if ok := tValue.FieldByName(name).Type() == oValue.FieldByName(name).Type(); ok {
				tValue.FieldByName(name).Set(reflect.ValueOf(oValue.Field(i).Interface()))
			}
		}
	}
}
