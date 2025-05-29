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

### 🧪 Ejemplo básico para leer JSON

```C
#include <stdio.h>
#include "JSON.h"

int main() {
    char* json = "{\"nombre\":\"Juan\", \"edad\":30, \"direccion\": {\"pais\":\"Villa Lactea\",\"departamento\":\"Tierra\"} }";
    
    // Analizar JSON
    JsonResult resultado = ParseJSON(json);
    
    if (resultado.is_valid) {
        printf("JSON válido: %s\n", resultado.value);
    } else {
        printf("Error: %s\n", resultado.error);
    }
    
    // Obtener valores
    JsonResult nombre = GetJSONValue(json, "nombre");
    JsonResult pais = GetJSONValueByPath(json, "direccion.pais");
    
    // Mostrar valores sin comillas
    printf("Nombre: %s\n", nombre.value);
    printf("País: %s\n", pais.value);
    
    // Liberar memoria
    FreeJsonResult(&resultado);
    FreeJsonResult(&nombre);
    FreeJsonResult(&pais);
    
    return 0;
}
```

---

### 🧪 Ejemplo para escribir JSON

```C
#include <stdio.h>
#include "JSON.h"

int main() {
    // 1. Crear un objeto JSON vacío
    JsonResult empty_json = CreateEmptyJSON();
    printf("JSON vacío: %s\n", empty_json.value);
    FreeJsonResult(&empty_json);

    // 2. Crear un objeto JSON con propiedades básicas
    JsonResult person = CreateEmptyJSON();
    person = AddStringToJSON(person.value, "name", "Juan Pérez");
    person = AddNumberToJSON(person.value, "age", 30);
    person = AddBooleanToJSON(person.value, "is_student", 0); // 0 = false
    
    printf("\nPersona básica:\n%s\n", person.value);

    // 3. Crear una dirección como JSON y añadirla a la persona
    JsonResult address = CreateEmptyJSON();
    address = AddStringToJSON(address.value, "street", "Calle Principal 123");
    address = AddStringToJSON(address.value, "city", "Ciudad Ejemplo");
    address = AddStringToJSON(address.value, "country", "España");
    
    person = AddJSONToJSON(person.value, "address", address.value);
    FreeJsonResult(&address);

    // 4. Crear un array de hobbies y añadirlo
    JsonResult hobbies = CreateEmptyArray();
    hobbies = AddItemToArray(hobbies.value, "\"fútbol\"");
    hobbies = AddItemToArray(hobbies.value, "\"lectura\"");
    hobbies = AddItemToArray(hobbies.value, "\"programación\"");
    
    person = AddJSONToJSON(person.value, "hobbies", hobbies.value);
    FreeJsonResult(&hobbies);

    // 5. Formatear el JSON para mejor legibilidad
    JsonResult pretty_person = PrettyPrintJSON(person.value);
    printf("\nPersona con formato bonito:\n%s\n", pretty_person.value);
    FreeJsonResult(&pretty_person);

    // 6. Modificar el JSON existente
    person = AddNumberToJSON(person.value, "age", 31); // Actualizar edad
    person = AddStringToJSON(person.value, "email", "juan@example.com");
    
    printf("\nPersona actualizada:\n%s\n", person.value);

    // 7. Eliminar una propiedad
    person = RemoveKeyFromJSON(person.value, "is_student");
    printf("\nPersona sin is_student:\n%s\n", person.value);

    // 8. Crear otro JSON para combinar
    JsonResult work_info = CreateEmptyJSON();
    work_info = AddStringToJSON(work_info.value, "company", "Tech Solutions");
    work_info = AddStringToJSON(work_info.value, "position", "Desarrollador");
    
    // Combinar con el JSON de persona
    person = MergeJSON(person.value, work_info.value);
    printf("\nPersona con info laboral:\n%s\n", person.value);
    FreeJsonResult(&work_info);

    // 9. Verificar si es JSON válido
    int is_valid = IsValidJSON(person.value);
    printf("\n¿JSON válido? %s\n", is_valid ? "Sí" : "No");

    // Liberar memoria
    FreeJsonResult(&person);

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
            
            printf("  Tipo documento: %s\n", documento.value);
            printf("  Número: %s\n", numero.value);
                        
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
