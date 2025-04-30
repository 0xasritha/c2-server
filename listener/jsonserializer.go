package main

import ()

//type JSONSerializer struct{}
//
//// Scan implements serializer interface
//func (JSONSerializer) Scan(ctx context.Context, field *gorm.Field, dst reflect.Value, dbValue interface{}) (err error) {
//	fieldValue := reflect.New()
//
//	if dbValue != nil {
//		var bytes []byte
//		switch v := dbValue.(type) {
//		case []byte:
//			bytes = v
//		case string:
//			bytes = []byte(v)
//		default:
//			return fmt.Errorf("failed to unmarshal JSONB value: %#v", dbValue)
//		}
//
//		err = json.Unmarshal(bytes, fieldValue.Interface())
//	}
//
//	field.ReflectValueOf(ctx, dst).Set(fieldValue.Elem())
//	return
//}
