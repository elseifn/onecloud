package printutils

import (
	"reflect"
	"strings"
	"fmt"

	"github.com/yunionio/jsonutils"
	"github.com/yunionio/pkg/utils"
	"github.com/yunionio/onecloud/pkg/mcclient/modules"
)

func getter2json(obj interface{}) jsonutils.JSONObject {
	jsonDict := jsonutils.NewDict()

	objValue := reflect.ValueOf(obj)
	objType := reflect.TypeOf(obj)

	// log.Debugf("getter2json %d", objValue.NumMethod())

	for i := 0; i < objValue.NumMethod(); i += 1 {
		methodValue := objValue.Method(i)
		method := objType.Method(i)
		methodName := method.Name
		methodType := methodValue.Type()

		if strings.HasPrefix(methodName, "Get") && methodType.NumIn() == 0 && methodType.NumOut() == 1 {
			fieldName := utils.CamelSplit(methodName[3:], "_")
			out := methodValue.Call([]reflect.Value{})
			if len(out) == 1 {
				jsonDict.Add(jsonutils.Marshal(out[0].Interface()), fieldName)
			}
		}
	}

	return jsonDict
}

func PrintGetterList(data interface{}, columns []string) {
	dataValue := reflect.ValueOf(data)
	if dataValue.Kind() != reflect.Slice {
		fmt.Println("Invalid list data")
		return
	}
	jsonList := make([]jsonutils.JSONObject, dataValue.Len())
	for i := 0; i < dataValue.Len(); i += 1 {
		jsonList[i] = getter2json(dataValue.Index(i).Interface())
	}
	list := &modules.ListResult{
		Data:   jsonList,
		Total:  dataValue.Len(),
		Limit:  0,
		Offset: 0,
	}
	PrintJSONList(list, columns)
}

func PrintGetterObject(obj interface{}) {
	PrintJSONObject(getter2json(obj))
}