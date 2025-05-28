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

/*JSON de prueba:
[
    { 
        "documento":"pasaporte",
        "numero":"B00000001"
    },
    {
        "documento":"pasaporte",
        "numero":"B00000002"
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
    printf("El JSON válido\n");
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
    printf("El array contiene %d elementos\n", array_length);


    
    // 3. Procesar cada elemento si la longitud es correcta
    if (array_length > 0) {
        for (int i = 0; i < array_length; i++) {
            printf("\nElemento %d:\n", i+1);
            
            JsonResult item_result = GetArrayItem(json_data, i);
            if (!item_result.is_valid) {
                printf("Error al obtener elemento: ");
                printf("Error: %s\n", length_result.error);
                FreeJsonResult(&item_result);
                continue;
            }
            
            printf("Contenido JSON: %s\n", item_result.value);
            
            // Extraer valores específicos
            JsonResult documento = GetJSONValue(item_result.value, "documento");
            JsonResult numero = GetJSONValue(item_result.value, "numero");
            
            printf("  data: %s\n", value(documento));
            printf("  numero: %s\n", value(numero));
                        
            FreeJsonResult(&documento);
            FreeJsonResult(&numero);
            FreeJsonResult(&item_result);
        }
    } else {
        printf("\nError: El array parece vacío o no se pudo determinar su longitud\n");
    }
    printf("\n");
    
    return 0;
}





/*
char* simple_test = "[1,2,3]";
JsonResult test = GetArrayLength(simple_test);
print_result(test);
*/

/*
#include <stdio.h>
#include <stdlib.h>
#include "JSON.h"

void print_result(JsonResult result) {
    if (result.is_valid) {
        printf("Value: %s\n", result.value);
    } else {
        printf("Error: %s\n", result.error);
    }
}

int main() {
    // JSON complejo de ejemplo
    char* complex_json = "{"
        "\"person\": {"
            "\"name\": \"John\","
            "\"age\": 30,"
            "\"address\": {"
                "\"street\": \"123 Main St\","
                "\"city\": \"New York\""
            "},"
            "\"hobbies\": [\"reading\", \"swimming\", \"coding\"]"
        "},"
        "\"active\": true"
    "}";
    
    // 1. Parsear el JSON completo
    JsonResult parsed = ParseJSON(complex_json);
    printf("Parsed JSON:\n");
    print_result(parsed);
    FreeJsonResult(&parsed);
    
    // 2. Obtener un objeto anidado
    JsonResult person = GetJSONValue(complex_json, "person");
    printf("\nPerson object:\n");
    print_result(person);
    
    // 3. Obtener un valor del objeto anidado
    JsonResult name = GetJSONValue(person.value, "name");
    printf("\nName:\n");
    print_result(name);
    
    // 4. Obtener un array
    JsonResult hobbies = GetJSONValue(person.value, "hobbies");
    printf("\nHobbies array:\n");
    print_result(hobbies);
    
    // 5. Obtener un valor por path
    JsonResult street = GetJSONValueByPath(complex_json, "person.address.street");
    printf("\nStreet (by path):\n");
    print_result(street);
    
    // 6. Obtener elemento de array por path
    JsonResult first_hobby = GetJSONValueByPath(complex_json, "person.hobbies.0");
    printf("\nFirst hobby (by path):\n");
    print_result(first_hobby);
    
    // 7. Manejo de errores
    JsonResult invalid = GetJSONValueByPath(complex_json, "person.invalid.key");
    printf("\nInvalid path:\n");
    print_result(invalid);
    
    // Liberar memoria
    FreeJsonResult(&person);
    FreeJsonResult(&name);
    FreeJsonResult(&hobbies);
    FreeJsonResult(&street);
    FreeJsonResult(&first_hobby);
    FreeJsonResult(&invalid);
    
    return 0;
}*/
