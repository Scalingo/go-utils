package publicgen

import (
	"fmt"
	"reflect"
	"strings"
)

var typeStruct = map[string]string{
	"bson.ObjectId": "string",
}

func listFields(st reflect.Type) []field {
	fields := make([]field, 0)
	for i := 0; i < st.NumField(); i++ {
		curField := st.Field(i)

		if curField.Anonymous {
			newFields := listFields(st.Field(i).Type)
			fields = append(fields, newFields...)
			continue
		}

		name := curField.Name
		json := curField.Tag.Get("json")
		typ := curField.Type
		pointer := false
		if curField.Type.Kind() == reflect.Ptr {
			typ = curField.Type.Elem()
			pointer = true
		}
		pkg := typ.PkgPath()
		typePrefix := ""

		if pkg != "" {
			path := strings.Split(pkg, "/")
			typePrefix = path[len(path)-1] + "."
		}
		typStr := typ.Name()
		fullType := fmt.Sprintf("%s%s", typePrefix, typStr)

		if newName, ok := typeStruct[fullType]; ok {
			fullType = newName
			typePrefix = ""
			pkg = ""
		}

		fields = append(fields, field{
			Name:    name,
			Type:    fullType,
			JSONTag: json,
			PkgPath: pkg,
			Pointer: pointer,
		})
	}
	return fields
}
