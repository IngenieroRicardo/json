package main

/*
#include <stdlib.h>
#include <string.h>

typedef struct {
    char* value;
    int is_valid;
    char* error;
} JsonResult;

typedef struct {
    char** items;
    int count;
    int is_valid;
    char* error;
} JsonArrayResult;
*/
import "C"
import (
	"encoding/json"
	"unsafe"
	"fmt"
	"strconv"
	"strings"
)

// ... (las funciones anteriores ParseJSON, GetJSONValue, etc. se mantienen igual)

//export GetArrayLength
func GetArrayLength(jsonStr *C.char) C.JsonResult {
    goStr := C.GoString(jsonStr)
    var result C.JsonResult

    // Verificar si es un array
    if len(goStr) == 0 || goStr[0] != '[' {
        result.is_valid = 0
        result.error = C.CString("not a JSON array")
        return result
    }

    var arr []interface{}
    err := json.Unmarshal([]byte(goStr), &arr)
    if err != nil {
        result.is_valid = 0
        result.error = C.CString(err.Error())
        return result
    }

    // Convertir el length a string
    lengthStr := strconv.Itoa(len(arr))
    
    result.is_valid = 1
    result.value = C.CString(lengthStr)
    result.error = nil
    return result
}

//export GetArrayItem
func GetArrayItem(jsonStr *C.char, index int) C.JsonResult {
	goStr := C.GoString(jsonStr)
	var result C.JsonResult

	var arr []interface{}
	err := json.Unmarshal([]byte(goStr), &arr)
	if err != nil {
		result.is_valid = 0
		result.error = C.CString(err.Error())
		return result
	}

	if index < 0 || index >= len(arr) {
		result.is_valid = 0
		result.error = C.CString("index out of bounds")
		return result
	}

	itemBytes, err := json.Marshal(arr[index])
	if err != nil {
		result.is_valid = 0
		result.error = C.CString(err.Error())
		return result
	}

	result.is_valid = 1
	result.value = C.CString(string(itemBytes))
	result.error = nil
	return result
}

//export FreeJsonResult
func FreeJsonResult(result *C.JsonResult) {
	if result.value != nil {
		C.free(unsafe.Pointer(result.value))
	}
	if result.error != nil {
		C.free(unsafe.Pointer(result.error))
	}
}








//export ParseJSON
func ParseJSON(jsonStr *C.char) C.JsonResult {
	goStr := C.GoString(jsonStr)
	var result C.JsonResult

	var data interface{}
	err := json.Unmarshal([]byte(goStr), &data)
	if err != nil {
		result.is_valid = 0
		result.error = C.CString(err.Error())
		return result
	}

	jsonBytes, err := json.Marshal(data)
	if err != nil {
		result.is_valid = 0
		result.error = C.CString(err.Error())
		return result
	}

	result.is_valid = 1
	result.value = C.CString(string(jsonBytes))
	result.error = nil
	return result
}

//export GetJSONValue
func GetJSONValue(jsonStr *C.char, key *C.char) C.JsonResult {
	goJsonStr := C.GoString(jsonStr)
	goKey := C.GoString(key)
	var result C.JsonResult

	var data map[string]interface{}
	err := json.Unmarshal([]byte(goJsonStr), &data)
	if err != nil {
		result.is_valid = 0
		result.error = C.CString(err.Error())
		return result
	}

	value, exists := data[goKey]
	if !exists {
		result.is_valid = 0
		result.error = C.CString("key not found")
		return result
	}

	jsonBytes, err := json.Marshal(value)
	if err != nil {
		result.is_valid = 0
		result.error = C.CString(err.Error())
		return result
	}

	result.is_valid = 1
	result.value = C.CString(string(jsonBytes))
	result.error = nil
	return result
}

//export GetJSONValueByPath
func GetJSONValueByPath(jsonStr *C.char, path *C.char) C.JsonResult {
	goJsonStr := C.GoString(jsonStr)
	goPath := C.GoString(path)
	var result C.JsonResult

	var data interface{}
	err := json.Unmarshal([]byte(goJsonStr), &data)
	if err != nil {
		result.is_valid = 0
		result.error = C.CString(err.Error())
		return result
	}

	current := data
	pathParts := strings.Split(goPath, ".")
	for _, part := range pathParts {
		if part == "" {
			continue
		}

		switch v := current.(type) {
		case map[string]interface{}:
			val, exists := v[part]
			if !exists {
				result.is_valid = 0
				result.error = C.CString(fmt.Sprintf("path '%s' not found", part))
				return result
			}
			current = val
		case []interface{}:
			index, err := strconv.Atoi(part)
			if err != nil || index < 0 || index >= len(v) {
				result.is_valid = 0
				result.error = C.CString(fmt.Sprintf("invalid array index '%s'", part))
				return result
			}
			current = v[index]
		default:
			result.is_valid = 0
			result.error = C.CString(fmt.Sprintf("cannot traverse path '%s'", part))
			return result
		}
	}

	jsonBytes, err := json.Marshal(current)
	if err != nil {
		result.is_valid = 0
		result.error = C.CString(err.Error())
		return result
	}

	result.is_valid = 1
	result.value = C.CString(string(jsonBytes))
	result.error = nil
	return result
}



func main() {}