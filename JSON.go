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
	"strconv"
	"fmt"
	"strings"
	"bytes"
)

// ------------------------- Funciones Optimizadas -------------------------

//export ParseJSON
func ParseJSON(jsonStr *C.char) C.JsonResult {
	goStr := C.GoString(jsonStr)
	var result C.JsonResult

	decoder := json.NewDecoder(bytes.NewReader([]byte(goStr)))
	decoder.UseNumber()

	var data interface{}
	if err := decoder.Decode(&data); err != nil {
		result.is_valid = 0
		result.error = C.CString("Error al analizar JSON: " + err.Error())
		return result
	}

	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false)

	if err := encoder.Encode(data); err != nil {
		result.is_valid = 0
		result.error = C.CString("Error al codificar JSON: " + err.Error())
		return result
	}

	result.is_valid = 1
	result.value = C.CString(strings.TrimSpace(buf.String()))
	result.error = nil
	return result
}

//export GetJSONValue
func GetJSONValue(jsonStr *C.char, key *C.char) C.JsonResult {
	goJsonStr := C.GoString(jsonStr)
	goKey := C.GoString(key)
	var result C.JsonResult

	decoder := json.NewDecoder(bytes.NewReader([]byte(goJsonStr)))
	decoder.UseNumber()

	var data map[string]interface{}
	if err := decoder.Decode(&data); err != nil {
		result.is_valid = 0
		result.error = C.CString("Error al analizar JSON: " + err.Error())
		return result
	}

	value, exists := data[goKey]
	if !exists {
		result.is_valid = 0
		result.error = C.CString(fmt.Sprintf("Clave '%s' no encontrada", goKey))
		return result
	}

	switch v := value.(type) {
	case string:
		result.is_valid = 1
		result.value = C.CString(v)
		result.error = nil
		return result
	case json.Number:
		result.is_valid = 1
		result.value = C.CString(v.String())
		result.error = nil
		return result
	case bool:
		str := "false"
		if v {
			str = "true"
		}
		result.is_valid = 1
		result.value = C.CString(str)
		result.error = nil
		return result
	default:
		var buf bytes.Buffer
		encoder := json.NewEncoder(&buf)
		encoder.SetEscapeHTML(false)
		if err := encoder.Encode(value); err != nil {
			result.is_valid = 0
			result.error = C.CString("Error al codificar valor: " + err.Error())
			return result
		}
		result.is_valid = 1
		result.value = C.CString(strings.TrimSpace(buf.String()))
		result.error = nil
		return result
	}
}

//export GetArrayLength
func GetArrayLength(jsonStr *C.char) C.JsonResult {
	goStr := C.GoString(jsonStr)
	var result C.JsonResult

	if len(goStr) == 0 || goStr[0] != '[' {
		result.is_valid = 0
		result.error = C.CString("No es un arreglo JSON válido")
		return result
	}

	decoder := json.NewDecoder(bytes.NewReader([]byte(goStr)))
	decoder.UseNumber()

	token, err := decoder.Token()
	if err != nil || token != json.Delim('[') {
		result.is_valid = 0
		result.error = C.CString("Arreglo JSON inválido")
		return result
	}

	count := 0
	for decoder.More() {
		var dummy interface{}
		if err := decoder.Decode(&dummy); err != nil {
			result.is_valid = 0
			result.error = C.CString("Error al contar elementos: " + err.Error())
			return result
		}
		count++
	}

	result.is_valid = 1
	result.value = C.CString(strconv.Itoa(count))
	result.error = nil
	return result
}

//export GetArrayItem
func GetArrayItem(jsonStr *C.char, index int) C.JsonResult {
	goStr := C.GoString(jsonStr)
	var result C.JsonResult

	decoder := json.NewDecoder(bytes.NewReader([]byte(goStr)))
	decoder.UseNumber()

	token, err := decoder.Token()
	if err != nil || token != json.Delim('[') {
		result.is_valid = 0
		result.error = C.CString("Arreglo JSON inválido")
		return result
	}

	currentIndex := 0
	for decoder.More() {
		if currentIndex == index {
			var item interface{}
			if err := decoder.Decode(&item); err != nil {
				result.is_valid = 0
				result.error = C.CString("Error al obtener elemento: " + err.Error())
				return result
			}

			var buf bytes.Buffer
			encoder := json.NewEncoder(&buf)
			encoder.SetEscapeHTML(false)
			if err := encoder.Encode(item); err != nil {
				result.is_valid = 0
				result.error = C.CString("Error al codificar elemento: " + err.Error())
				return result
			}

			result.is_valid = 1
			result.value = C.CString(strings.TrimSpace(buf.String()))
			result.error = nil
			return result
		}

		// Saltar este elemento
		var dummy interface{}
		if err := decoder.Decode(&dummy); err != nil {
			result.is_valid = 0
			result.error = C.CString("Error al saltar elemento: " + err.Error())
			return result
		}
		currentIndex++
	}

	result.is_valid = 0
	result.error = C.CString("Índice fuera de rango")
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

//export GetJSONKeys
func GetJSONKeys(jsonStr *C.char) C.JsonArrayResult {
	goJsonStr := C.GoString(jsonStr)
	var result C.JsonArrayResult

	decoder := json.NewDecoder(bytes.NewReader([]byte(goJsonStr)))
	decoder.UseNumber()

	var data map[string]interface{}
	if err := decoder.Decode(&data); err != nil {
		result.is_valid = 0
		result.error = C.CString("Error al analizar JSON: " + err.Error())
		return result
	}

	keys := make([]string, 0, len(data))
	for key := range data {
		keys = append(keys, key)
	}

	cArray := C.malloc(C.size_t(len(keys)) * C.size_t(unsafe.Sizeof(uintptr(0))))
	cKeys := (*[1<<30 - 1]*C.char)(cArray)

	for i, key := range keys {
		cKeys[i] = C.CString(key)
	}

	result.is_valid = 1
	result.items = (**C.char)(cArray)
	result.count = C.int(len(keys))
	result.error = nil
	return result
}

//export FreeJsonArrayResult
func FreeJsonArrayResult(result *C.JsonArrayResult) {
	if result.items != nil {
		cKeys := (*[1<<30]*C.char)(unsafe.Pointer(result.items))[:result.count:result.count]
		for i := 0; i < int(result.count); i++ {
			C.free(unsafe.Pointer(cKeys[i]))
		}
		C.free(unsafe.Pointer(result.items))
	}
	if result.error != nil {
		C.free(unsafe.Pointer(result.error))
	}
}

//export GetJSONValueByPath
func GetJSONValueByPath(jsonStr *C.char, path *C.char) C.JsonResult {
	goJsonStr := C.GoString(jsonStr)
	goPath := C.GoString(path)
	var result C.JsonResult

	decoder := json.NewDecoder(bytes.NewReader([]byte(goJsonStr)))
	decoder.UseNumber()

	var data interface{}
	if err := decoder.Decode(&data); err != nil {
		result.is_valid = 0
		result.error = C.CString("Error al analizar JSON: " + err.Error())
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
				result.error = C.CString(fmt.Sprintf("Ruta '%s' no encontrada", part))
				return result
			}
			current = val
		case []interface{}:
			index, err := strconv.Atoi(part)
			if err != nil || index < 0 || index >= len(v) {
				result.is_valid = 0
				result.error = C.CString(fmt.Sprintf("Índice de arreglo inválido '%s'", part))
				return result
			}
			current = v[index]
		default:
			result.is_valid = 0
			result.error = C.CString(fmt.Sprintf("No se puede navegar por la ruta '%s'", part))
			return result
		}
	}

	switch v := current.(type) {
	case string:
		result.is_valid = 1
		result.value = C.CString(v)
		result.error = nil
		return result
	case json.Number:
		result.is_valid = 1
		result.value = C.CString(v.String())
		result.error = nil
		return result
	case bool:
		str := "false"
		if v {
			str = "true"
		}
		result.is_valid = 1
		result.value = C.CString(str)
		result.error = nil
		return result
	case nil:
		result.is_valid = 1
		result.value = C.CString("null")
		result.error = nil
		return result
	default:
		var buf bytes.Buffer
		encoder := json.NewEncoder(&buf)
		encoder.SetEscapeHTML(false)
		if err := encoder.Encode(current); err != nil {
			result.is_valid = 0
			result.error = C.CString("Error al codificar valor: " + err.Error())
			return result
		}
		result.is_valid = 1
		result.value = C.CString(strings.TrimSpace(buf.String()))
		result.error = nil
		return result
	}
}

//export GetArrayItems
func GetArrayItems(jsonStr *C.char) C.JsonArrayResult {
	goStr := C.GoString(jsonStr)
	var result C.JsonArrayResult

	if len(goStr) == 0 || goStr[0] != '[' {
		result.is_valid = 0
		result.error = C.CString("No es un arreglo JSON válido")
		return result
	}

	decoder := json.NewDecoder(bytes.NewReader([]byte(goStr)))
	decoder.UseNumber()

	token, err := decoder.Token()
	if err != nil || token != json.Delim('[') {
		result.is_valid = 0
		result.error = C.CString("Arreglo JSON inválido")
		return result
	}

	var items []string
	for decoder.More() {
		var item interface{}
		if err := decoder.Decode(&item); err != nil {
			// Liberar memoria ya asignada si hay error
			for _, str := range items {
				C.free(unsafe.Pointer(C.CString(str)))
			}
			result.is_valid = 0
			result.error = C.CString("Error al obtener elementos: " + err.Error())
			return result
		}

		var buf bytes.Buffer
		encoder := json.NewEncoder(&buf)
		encoder.SetEscapeHTML(false)
		if err := encoder.Encode(item); err != nil {
			for _, str := range items {
				C.free(unsafe.Pointer(C.CString(str)))
			}
			result.is_valid = 0
			result.error = C.CString("Error al codificar elemento: " + err.Error())
			return result
		}

		items = append(items, strings.TrimSpace(buf.String()))
	}

	cArray := C.malloc(C.size_t(len(items)) * C.size_t(unsafe.Sizeof(uintptr(0))))
	cItems := (*[1<<30 - 1]*C.char)(cArray)

	for i, item := range items {
		cItems[i] = C.CString(item)
	}

	result.is_valid = 1
	result.items = (**C.char)(cArray)
	result.count = C.int(len(items))
	result.error = nil
	return result
}

// ------------------------- Funciones de Construcción Optimizadas -------------------------

//export CreateEmptyJSON
func CreateEmptyJSON() C.JsonResult {
	var result C.JsonResult
	result.is_valid = 1
	result.value = C.CString("{}")
	result.error = nil
	return result
}

//export CreateEmptyArray
func CreateEmptyArray() C.JsonResult {
	var result C.JsonResult
	result.is_valid = 1
	result.value = C.CString("[]")
	result.error = nil
	return result
}

//export AddStringToJSON
func AddStringToJSON(jsonStr *C.char, key *C.char, value *C.char) C.JsonResult {
	goJsonStr := C.GoString(jsonStr)
	goKey := C.GoString(key)
	goValue := C.GoString(value)
	var result C.JsonResult

	decoder := json.NewDecoder(bytes.NewReader([]byte(goJsonStr)))
	decoder.UseNumber()

	var data map[string]interface{}
	if err := decoder.Decode(&data); err != nil {
		result.is_valid = 0
		result.error = C.CString("Error al analizar JSON: " + err.Error())
		return result
	}

	data[goKey] = goValue

	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(data); err != nil {
		result.is_valid = 0
		result.error = C.CString("Error al codificar JSON: " + err.Error())
		return result
	}

	result.is_valid = 1
	result.value = C.CString(strings.TrimSpace(buf.String()))
	result.error = nil
	return result
}

//export AddNumberToJSON
func AddNumberToJSON(jsonStr *C.char, key *C.char, value float64) C.JsonResult {
	goJsonStr := C.GoString(jsonStr)
	goKey := C.GoString(key)
	var result C.JsonResult

	decoder := json.NewDecoder(bytes.NewReader([]byte(goJsonStr)))
	decoder.UseNumber()

	var data map[string]interface{}
	if err := decoder.Decode(&data); err != nil {
		result.is_valid = 0
		result.error = C.CString("Error al analizar JSON: " + err.Error())
		return result
	}

	data[goKey] = value

	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(data); err != nil {
		result.is_valid = 0
		result.error = C.CString("Error al codificar JSON: " + err.Error())
		return result
	}

	result.is_valid = 1
	result.value = C.CString(strings.TrimSpace(buf.String()))
	result.error = nil
	return result
}

//export AddBooleanToJSON
func AddBooleanToJSON(jsonStr *C.char, key *C.char, value C.int) C.JsonResult {
	goJsonStr := C.GoString(jsonStr)
	goKey := C.GoString(key)
	goValue := value != 0
	var result C.JsonResult

	decoder := json.NewDecoder(bytes.NewReader([]byte(goJsonStr)))
	decoder.UseNumber()

	var data map[string]interface{}
	if err := decoder.Decode(&data); err != nil {
		result.is_valid = 0
		result.error = C.CString("Error al analizar JSON: " + err.Error())
		return result
	}

	data[goKey] = goValue

	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(data); err != nil {
		result.is_valid = 0
		result.error = C.CString("Error al codificar JSON: " + err.Error())
		return result
	}

	result.is_valid = 1
	result.value = C.CString(strings.TrimSpace(buf.String()))
	result.error = nil
	return result
}

//export AddJSONToJSON
func AddJSONToJSON(parentJson *C.char, key *C.char, childJson *C.char) C.JsonResult {
	goParent := C.GoString(parentJson)
	goKey := C.GoString(key)
	goChild := C.GoString(childJson)
	var result C.JsonResult

	decoder := json.NewDecoder(bytes.NewReader([]byte(goParent)))
	decoder.UseNumber()

	var parentData map[string]interface{}
	if err := decoder.Decode(&parentData); err != nil {
		result.is_valid = 0
		result.error = C.CString("Error al analizar JSON padre: " + err.Error())
		return result
	}

	childDecoder := json.NewDecoder(bytes.NewReader([]byte(goChild)))
	childDecoder.UseNumber()

	var childData interface{}
	if err := childDecoder.Decode(&childData); err != nil {
		result.is_valid = 0
		result.error = C.CString("Error al analizar JSON hijo: " + err.Error())
		return result
	}

	parentData[goKey] = childData

	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(parentData); err != nil {
		result.is_valid = 0
		result.error = C.CString("Error al codificar JSON combinado: " + err.Error())
		return result
	}

	result.is_valid = 1
	result.value = C.CString(strings.TrimSpace(buf.String()))
	result.error = nil
	return result
}

//export AddItemToArray
func AddItemToArray(jsonArray *C.char, item *C.char) C.JsonResult {
	goArray := C.GoString(jsonArray)
	goItem := C.GoString(item)
	var result C.JsonResult

	decoder := json.NewDecoder(bytes.NewReader([]byte(goArray)))
	decoder.UseNumber()

	var arrayData []interface{}
	if err := decoder.Decode(&arrayData); err != nil {
		result.is_valid = 0
		result.error = C.CString("Error al analizar arreglo JSON: " + err.Error())
		return result
	}

	itemDecoder := json.NewDecoder(bytes.NewReader([]byte(goItem)))
	itemDecoder.UseNumber()

	var itemData interface{}
	if err := itemDecoder.Decode(&itemData); err != nil {
		result.is_valid = 0
		result.error = C.CString("Error al analizar elemento: " + err.Error())
		return result
	}

	arrayData = append(arrayData, itemData)

	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(arrayData); err != nil {
		result.is_valid = 0
		result.error = C.CString("Error al codificar arreglo actualizado: " + err.Error())
		return result
	}

	result.is_valid = 1
	result.value = C.CString(strings.TrimSpace(buf.String()))
	result.error = nil
	return result
}

//export RemoveKeyFromJSON
func RemoveKeyFromJSON(jsonStr *C.char, key *C.char) C.JsonResult {
	goJsonStr := C.GoString(jsonStr)
	goKey := C.GoString(key)
	var result C.JsonResult

	decoder := json.NewDecoder(bytes.NewReader([]byte(goJsonStr)))
	decoder.UseNumber()

	var data map[string]interface{}
	if err := decoder.Decode(&data); err != nil {
		result.is_valid = 0
		result.error = C.CString("Error al analizar JSON: " + err.Error())
		return result
	}

	if _, exists := data[goKey]; !exists {
		result.is_valid = 0
		result.error = C.CString(fmt.Sprintf("Clave '%s' no encontrada", goKey))
		return result
	}

	delete(data, goKey)

	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(data); err != nil {
		result.is_valid = 0
		result.error = C.CString("Error al codificar JSON actualizado: " + err.Error())
		return result
	}

	result.is_valid = 1
	result.value = C.CString(strings.TrimSpace(buf.String()))
	result.error = nil
	return result
}

//export RemoveItemFromArray
func RemoveItemFromArray(jsonArray *C.char, index C.int) C.JsonResult {
	goArray := C.GoString(jsonArray)
	goIndex := int(index)
	var result C.JsonResult

	decoder := json.NewDecoder(bytes.NewReader([]byte(goArray)))
	decoder.UseNumber()

	var arrayData []interface{}
	if err := decoder.Decode(&arrayData); err != nil {
		result.is_valid = 0
		result.error = C.CString("Error al analizar arreglo JSON: " + err.Error())
		return result
	}

	if goIndex < 0 || goIndex >= len(arrayData) {
		result.is_valid = 0
		result.error = C.CString("Índice fuera de rango")
		return result
	}

	arrayData = append(arrayData[:goIndex], arrayData[goIndex+1:]...)

	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(arrayData); err != nil {
		result.is_valid = 0
		result.error = C.CString("Error al codificar arreglo actualizado: " + err.Error())
		return result
	}

	result.is_valid = 1
	result.value = C.CString(strings.TrimSpace(buf.String()))
	result.error = nil
	return result
}

//export PrettyPrintJSON
func PrettyPrintJSON(jsonStr *C.char) C.JsonResult {
	goJsonStr := C.GoString(jsonStr)
	var result C.JsonResult

	decoder := json.NewDecoder(bytes.NewReader([]byte(goJsonStr)))
	decoder.UseNumber()

	var data interface{}
	if err := decoder.Decode(&data); err != nil {
		result.is_valid = 0
		result.error = C.CString("Error al analizar JSON: " + err.Error())
		return result
	}

	jsonBytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		result.is_valid = 0
		result.error = C.CString("Error al formatear JSON: " + err.Error())
		return result
	}

	result.is_valid = 1
	result.value = C.CString(string(jsonBytes))
	result.error = nil
	return result
}

//export MergeJSON
func MergeJSON(json1 *C.char, json2 *C.char) C.JsonResult {
	goJson1 := C.GoString(json1)
	goJson2 := C.GoString(json2)
	var result C.JsonResult

	decoder1 := json.NewDecoder(bytes.NewReader([]byte(goJson1)))
	decoder1.UseNumber()

	var data1 map[string]interface{}
	if err := decoder1.Decode(&data1); err != nil {
		result.is_valid = 0
		result.error = C.CString("Error al analizar primer JSON: " + err.Error())
		return result
	}

	decoder2 := json.NewDecoder(bytes.NewReader([]byte(goJson2)))
	decoder2.UseNumber()

	var data2 map[string]interface{}
	if err := decoder2.Decode(&data2); err != nil {
		result.is_valid = 0
		result.error = C.CString("Error al analizar segundo JSON: " + err.Error())
		return result
	}

	for key, value := range data2 {
		data1[key] = value
	}

	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(data1); err != nil {
		result.is_valid = 0
		result.error = C.CString("Error al codificar JSON combinado: " + err.Error())
		return result
	}

	result.is_valid = 1
	result.value = C.CString(strings.TrimSpace(buf.String()))
	result.error = nil
	return result
}

//export IsValidJSON
func IsValidJSON(jsonStr *C.char) C.int {
	goJsonStr := C.GoString(jsonStr)
	
	decoder := json.NewDecoder(bytes.NewReader([]byte(goJsonStr)))
	decoder.UseNumber()
	
	var dummy interface{}
	if err := decoder.Decode(&dummy); err != nil {
		return 0
	}
	return 1
}

func main() {}
