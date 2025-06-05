# JSON
Una biblioteca ligera para manipular JSON en C.  
Compilada usando: `go build -o JSON.dll -buildmode=c-shared JSON.go`

---

### 📥 Descargar la librería

| Linux | Windows |
| --- | --- |
| `wget https://github.com/IngenieroRicardo/JSON/releases/download/1.0/JSON.so` | `Invoke-WebRequest https://github.com/IngenieroRicardo/JSON/releases/download/1.0/JSON.dll -OutFile ./JSON.dll` |
| `wget https://github.com/IngenieroRicardo/JSON/releases/download/1.0/JSON.h` | `Invoke-WebRequest https://github.com/IngenieroRicardo/JSON/releases/download/1.0/JSON.h -OutFile ./JSON.h` |

---

### 🛠️ Compilar

| Linux | Windows |
| --- | --- |
| `gcc -o main.bin main.c ./JSON.so` | `gcc -o main.exe main.c ./JSON.dll` |
| `x86_64-w64-mingw32-gcc -o main.exe main.c ./JSON.dll` |  |

---

### 🧪 Ejemplo básico para leer JSON

```C
#include <stdio.h>
#include "JSON.h"

int main() {
    char* json = "{\"nombre\":\"Juan\", \"edad\":30, \"direccion\": {\"pais\":\"Villa Lactea\",\"departamento\":\"Tierra\"}, \"documentos\": [\"B00000001\",\"00000000-1\"], \"foto\":\"iVBORw0KGgoAAAANSUhEUgAAAAgAAAAICAIAAABLbSncAAAAAXNSR0IArs4c6QAAAARnQU1BAACxjwv8YQUAAAAJcEhZcwAADsMAAA7DAcdvqGQAAAArSURBVBhXY/iPA0AlGBgwGFAKlwQmAKrAIgcVRZODCsI5cAAVgVDo4P9/AHe4m2U/OJCWAAAAAElFTkSuQmCC\" }";
    
    // Analizar JSON
    JsonResult resultado = ParseJSON(json);
    
    if (resultado.is_valid) {
        printf("JSON válido: %s\n", resultado.value);
    } else {
        printf("Error: %s\n", resultado.error);
        FreeJsonResult(resultado);
        return 1;
    }
    
    // Obtener valores
    JsonResult nombre = GetJSONValue(json, "nombre");
    JsonResult pais = GetJSONValueByPath(json, "direccion.pais");
    JsonResult documento1 = GetJSONValueByPath(json, "documentos.0");
    
    // Mostrar valores sin comillas
    printf("Nombre: %s\n", nombre.value);
    printf("País: %s\n", pais.value);
    printf("Primer Documento: %s\n", documento1.value);
    
    // Liberar memoria
    FreeJsonResult(resultado);
    FreeJsonResult(nombre);
    FreeJsonResult(pais);
    FreeJsonResult(documento1);
    
    return 0;
}
```

---

### 🧪 Ejemplo para escribir, editar y eliminar JSON

```C
#include <stdio.h>
#include "JSON.h"

int main() {
    // 1. Crear un objeto JSON vacío
    JsonResult json_vacio = CreateEmptyJSON();
    printf("JSON vacío: %s\n", json_vacio.value);
    FreeJsonResult(json_vacio);

    // 2. Crear un objeto JSON con datos básicos de persona
    JsonResult persona = CreateEmptyJSON();
    persona = AddStringToJSON(persona.value, "nombre", "Juan Pérez");
    persona = AddNumberToJSON(persona.value, "edad", 30);
    persona = AddBooleanToJSON(persona.value, "es_estudiante", 0); // 0 = falso
    
    printf("\nPersona básica:\n%s\n", persona.value);

    // 3. Crear una dirección como JSON y añadirla a la persona
    JsonResult direccion = CreateEmptyJSON();
    direccion = AddStringToJSON(direccion.value, "calle", "Calle Principal 123");
    direccion = AddStringToJSON(direccion.value, "ciudad", "Ciudad Ejemplo");
    direccion = AddStringToJSON(direccion.value, "pais", "España");
    
    persona = AddJSONToJSON(persona.value, "direccion", direccion.value);
    FreeJsonResult(direccion);

    // 4. Crear un array de pasatiempos y añadirlo
    JsonResult pasatiempos = CreateEmptyArray();
    pasatiempos = AddItemToArray(pasatiempos.value, "\"fútbol\"");
    pasatiempos = AddItemToArray(pasatiempos.value, "\"lectura\"");
    pasatiempos = AddItemToArray(pasatiempos.value, "\"programación\"");
    
    persona = AddJSONToJSON(persona.value, "pasatiempos", pasatiempos.value);
    FreeJsonResult(pasatiempos);

    // 5. Modificar el JSON existente
    persona = AddNumberToJSON(persona.value, "edad", 31); // Actualizar edad
    persona = AddStringToJSON(persona.value, "correo", "juan@ejemplo.com");
    
    printf("\nPersona actualizada:\n%s\n", persona.value);

    // 6. Eliminar una propiedad
    persona = RemoveKeyFromJSON(persona.value, "es_estudiante");
    printf("\nPersona sin 'es_estudiante':\n%s\n", persona.value);

    // 7. Crear otro JSON con información laboral
    JsonResult info_laboral = CreateEmptyJSON();
    info_laboral = AddStringToJSON(info_laboral.value, "empresa", "Soluciones Tecnológicas");
    info_laboral = AddStringToJSON(info_laboral.value, "puesto", "Desarrollador");
    
    // Combinar con el JSON de persona
    persona = MergeJSON(persona.value, info_laboral.value);
    printf("\nPersona con información laboral:\n%s\n", persona.value);
    FreeJsonResult(info_laboral);

    // 8. Verificar si el JSON es válido
    int es_valido = IsValidJSON(persona.value);
    printf("\n¿JSON válido? %s\n", es_valido ? "Sí" : "No");

    // Liberar memoria
    FreeJsonResult(persona);

    return 0;
}
```

---

### 🧪 Ejemplo avanzado para leer JSON

```C
#include <stdio.h>
#include "JSON.h"

int main() {
    /* Ejemplo de JSON para pruebas:
    [
        { 
            "documento": "pasaporte",
            "numero": "B00000001"
        },
        {
            "documento": "pasaporte",
            "numero": "B00000002"
        }
    ]
    */
    char* json_data = "[{\"documento\":\"pasaporte\",\"numero\":\"B00000001\"},{\"documento\":\"pasaporte\",\"numero\":\"B00000002\"}]";
    printf("Procesando JSON completo:\n%s\n", json_data);

    // 1. Verificar parsing básico
    JsonResult parse_test = ParseJSON(json_data);
    printf("\nTest de parsing: ");
    if (!parse_test.is_valid) {
        printf("Error: %s\n", parse_test.error);
        FreeJsonResult(parse_test);
        return 1;
    }
    printf("JSON válido\n");
    FreeJsonResult(parse_test);

    // 2. Obtener longitud del array
    JsonResult length_result = GetArrayLength(json_data);
    printf("\nLongitud del array: ");
    if (!length_result.is_valid) {
        printf("Error: %s\n", length_result.error);
        FreeJsonResult(length_result);
        return 1;
    } 
    int array_length = atoi(length_result.value);
    FreeJsonResult(length_result);
    printf("El array contiene %d elemento(s)\n", array_length);

    // 3. Procesar cada elemento si la longitud es correcta
    if (array_length > 0) {
        for (int i = 0; i < array_length; i++) {
            printf("\nElemento %d:\n", i+1);
            
            JsonResult item_result = GetArrayItem(json_data, i);
            if (!item_result.is_valid) {
                printf("Error al obtener elemento: %s\n", item_result.error);
                FreeJsonResult(item_result);
                continue;
            }
            
            printf("Contenido JSON: %s\n", item_result.value);
            
            // Extraer valores específicos
            JsonResult documento = GetJSONValue(item_result.value, "documento");
            JsonResult numero = GetJSONValue(item_result.value, "numero");
            
            printf("  Tipo documento: %s\n", documento.value);
            printf("  Número: %s\n", numero.value);
                        
            FreeJsonResult(documento);
            FreeJsonResult(numero);
            FreeJsonResult(item_result);
        }
    } else {
        printf("\nError: El array está vacío o no se pudo determinar su longitud\n");
    }
    
    printf("\nProcesamiento completado\n");
    return 0;
}
```

---



## 📚 Documentación de la API

### Funciones Principales

#### Manejo Básico de JSON
- `JsonResult ParseJSON(char* jsonStr)`: Analiza una cadena JSON
- `int IsValidJSON(char* json_str)`: Verifica si una cadena es JSON válido

#### Obtención de Valores
- `JsonResult GetJSONValue(char* json_str, char* key)`: Obtiene valor por clave
- `JsonResult GetJSONValueByPath(char* json_str, char* path)`: Obtiene valor por ruta
- `JsonResult GetArrayLength(char* json_str)`: Obtiene longitud de array
- `JsonResult GetArrayItem(char* json_str, int index)`: Obtiene elemento de array

#### Construcción/Modificación
- `JsonResult CreateEmptyJSON()`: Crea objeto JSON vacío
- `JsonResult CreateEmptyArray()`: Crea array JSON vacío
- `JsonResult AddStringToJSON(char* json_str, char* key, char* value)`
- `JsonResult AddNumberToJSON(char* json_str, char* key, double value)`
- `JsonResult AddBooleanToJSON(char* json_str, char* key, int value)`
- `JsonResult AddJSONToJSON(char* parent_json, char* key, char* child_json)`
- `JsonResult AddItemToArray(char* json_array, char* item)`
- `JsonResult RemoveKeyFromJSON(char* json_str, char* key)`
- `JsonResult RemoveItemFromArray(char* json_array, int index)`
- `JsonResult MergeJSON(char* json1, char* json2)`: Combina dos JSONs

#### Utilidades
- `void FreeJsonResult(JsonResult result)`: Libera memoria de resultados
- `void FreeJsonArrayResult(JsonArrayResult result)`: Libera memoria de arrays

### Estructuras
```c
typedef struct {
    char* value;      // Valor obtenido
    int is_valid;     // 1 si es válido, 0 si hay error
    char* error;      // Mensaje de error (si lo hay)
} JsonResult;

typedef struct {
    char** items;     // Array de elementos
    int count;        // Número de elementos
    int is_valid;     // 1 si es válido, 0 si hay error
    char* error;      // Mensaje de error (si lo hay)
} JsonArrayResult;
```
