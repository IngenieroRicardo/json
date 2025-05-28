# JSON
Una biblioteca ligera para leer y procesar JSON en C.  
Compilada usando: `go build -o JSON.dll -buildmode=c-shared JSON.go`

---

### 📥 Descargar la librería

| Linux | Windows |
| --- | --- |
| `wget https://raw.githubusercontent.com/IngenieroRicardo/JSON/refs/heads/main/JSON.so` | `Invoke-WebRequest https://raw.githubusercontent.com/IngenieroRicardo/JSON/refs/heads/main/JSON.dll -OutFile ./JSON.dll` |
| `wget https://raw.githubusercontent.com/IngenieroRicardo/JSON/refs/heads/main/JSON.h` | `Invoke-WebRequest https://raw.githubusercontent.com/IngenieroRicardo/JSON/refs/heads/main/JSON.h -OutFile ./JSON.h` |

---

### 🛠️ Compilar

| Linux | Windows |
| --- | --- |
| `gcc -o main.bin main.c ./JSON.so` | `gcc -o main.exe main.c ./JSON.dll` |
| `x86_64-w64-mingw32-gcc -o main.exe main.c ./JSON.dll` |  |

---

### 🧪 Ejemplo básico

```C
#include <stdio.h>
#include <stdlib.h>
#include "JSON.h"


int main() {
    char* json = "{\"nombre\":\"Juan\", \"edad\":30}";
    
    // Analizar JSON
    JsonResult resultado = ParseJSON(json);
    
    if (resultado.is_valid) {
        printf("JSON válido: %s\n", resultado.value);
    } else {
        printf("Error: %s\n", resultado.error);
    }
    
    // Obtener valor específico
    JsonResult nombre = GetJSONValue(json, "nombre");
    printf("Nombre: %s\n", nombre.value);
    
    // Liberar memoria
    FreeJsonResult(&resultado);
    FreeJsonResult(&nombre);
    
    return 0;
}
```



### 🧪 Ejemplo avanzado 1

```C
#include <stdio.h>
#include <stdlib.h>
#include "JSON.h"

void imprimir_resultado(JsonResult resultado) {
    if (resultado.is_valid) {
        printf("Valor: %s\n", resultado.value);
    } else {
        printf("Error: %s\n", resultado.error);
    }
}

int main() {
    // JSON complejo de ejemplo
    char* json_complejo = "{"
        "\"persona\": {"
            "\"nombre\": \"Juan\","
            "\"edad\": 30,"
            "\"direccion\": {"
                "\"calle\": \"123 Calle Principal\","
                "\"ciudad\": \"Ciudad de México\""
            "},"
            "\"pasatiempos\": [\"leer\", \"nadar\", \"programar\"]"
        "},"
        "\"activo\": true"
    "}";
    
    // 1. Analizar (parsear) el JSON completo
    JsonResult parseado = ParseJSON(json_complejo);
    printf("JSON parseado:\n");
    imprimir_resultado(parseado);
    FreeJsonResult(&parseado);
    
    // 2. Obtener un objeto anidado
    JsonResult persona = GetJSONValue(json_complejo, "persona");
    printf("\nObjeto persona:\n");
    imprimir_resultado(persona);
    
    // 3. Obtener un valor del objeto anidado
    JsonResult nombre = GetJSONValue(persona.value, "nombre");
    printf("\nNombre:\n");
    imprimir_resultado(nombre);
    
    // 4. Obtener un array
    JsonResult pasatiempos = GetJSONValue(persona.value, "pasatiempos");
    printf("\nArray de pasatiempos:\n");
    imprimir_resultado(pasatiempos);
    
    // 5. Obtener un valor por ruta (path)
    JsonResult calle = GetJSONValueByPath(json_complejo, "persona.direccion.calle");
    printf("\nCalle (por ruta):\n");
    imprimir_resultado(calle);
    
    // 6. Obtener elemento de array por ruta
    JsonResult primer_pasatiempo = GetJSONValueByPath(json_complejo, "persona.pasatiempos.0");
    printf("\nPrimer pasatiempo (por ruta):\n");
    imprimir_resultado(primer_pasatiempo);
    
    // 7. Manejo de errores
    JsonResult invalido = GetJSONValueByPath(json_complejo, "persona.clave.invalida");
    printf("\nRuta inválida:\n");
    imprimir_resultado(invalido);
    
    // Liberar memoria
    FreeJsonResult(&persona);
    FreeJsonResult(&nombre);
    FreeJsonResult(&pasatiempos);
    FreeJsonResult(&calle);
    FreeJsonResult(&primer_pasatiempo);
    FreeJsonResult(&invalido);
    
    return 0;
}
```

---

### 🧪 Ejemplo avanzado 2

```C
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include "JSON.h"

char* value(JsonResult result) {
    if (result.is_valid) {
        return result.value;
    } else {
        return "";
    }
}

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
        FreeJsonResult(&parse_test);
        return 1;
    }
    printf("JSON válido\n");
    FreeJsonResult(&parse_test);

    // 2. Obtener longitud del array
    JsonResult length_result = GetArrayLength(json_data);
    printf("\nLongitud del array: ");
    if (!length_result.is_valid) {
        printf("Error: %s\n", length_result.error);
        FreeJsonResult(&length_result);
        return 1;
    } 
    int array_length = atoi(length_result.value);
    FreeJsonResult(&length_result);
    printf("El array contiene %d elemento(s)\n", array_length);

    // 3. Procesar cada elemento si la longitud es correcta
    if (array_length > 0) {
        for (int i = 0; i < array_length; i++) {
            printf("\nElemento %d:\n", i+1);
            
            JsonResult item_result = GetArrayItem(json_data, i);
            if (!item_result.is_valid) {
                printf("Error al obtener elemento: %s\n", item_result.error);
                FreeJsonResult(&item_result);
                continue;
            }
            
            printf("Contenido JSON: %s\n", item_result.value);
            
            // Extraer valores específicos
            JsonResult documento = GetJSONValue(item_result.value, "documento");
            JsonResult numero = GetJSONValue(item_result.value, "numero");
            
            printf("  Tipo documento: %s\n", value(documento));
            printf("  Número: %s\n", value(numero));
                        
            FreeJsonResult(&documento);
            FreeJsonResult(&numero);
            FreeJsonResult(&item_result);
        }
    } else {
        printf("\nError: El array está vacío o no se pudo determinar su longitud\n");
    }
    
    printf("\nProcesamiento completado\n");
    return 0;
}
```

## Características

- ✅ Analizar y validar cadenas JSON
- 🔍 Extraer valores por clave o ruta JSON
- 📦 Manejar arrays JSON (obtener longitud, acceder a elementos)
- 🚀 Interfaz compatible con C para integración con otros lenguajes
- 🧠 Operaciones seguras en memoria con limpieza adecuada
- 📝 Manejo completo de errores


---
