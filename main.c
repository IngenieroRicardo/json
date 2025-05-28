#include <stdio.h>
#include <stdlib.h>
#include "SQL2JSON.h"
#include "JSON.h"

void imprimir_sin_comillas(char* str) {
    if (str == NULL) return;
    
    size_t len = strlen(str);
    if (len >= 2 && str[0] == '"' && str[len-1] == '"') {
        printf("%.*s", (int)(len-2), str+1);
    } else {
        printf("%s", str);
    }
}

void mostrar_elemento_json(char* json_str, int indice) {
    JsonResult parseado = ParseJSON(json_str);
    if (!parseado.is_valid) {
        printf("Error al parsear elemento: %s\n", parseado.error);
        FreeJsonResult(&parseado);
        return;
    }
    FreeJsonResult(&parseado);

    JsonArrayResult claves = GetJSONKeys(json_str);
    if (!claves.is_valid) {
        printf("Error al obtener claves: %s\n", claves.error);
        FreeJsonArrayResult(&claves);
        return;
    }

    printf("Elemento %d:\n", indice + 1);
    
    for (int i = 0; i < claves.count; i++) {
        char* clave = claves.items[i];
        JsonResult valor = GetJSONValue(json_str, clave);
        
        printf("  %s: ", clave);
        imprimir_sin_comillas(valor.value);
        printf("\n");
        
        FreeJsonResult(&valor);
    }
    printf("\n");
    
    FreeJsonArrayResult(&claves);
}

int main() {
    char* conexion = "root:123456@tcp(192.100.1.210:3306)/chat";
    char* query = "select usuarios, mensajes from mensajeria;";
    
    char* json = SQLrun(conexion, query, 0, 0);
    if (json == NULL) {
        printf("Error al ejecutar la consulta SQL\n");
        return 1;
    }
    printf("Resultado JSON:\n%s\n\n", json);

    JsonResult parseado = ParseJSON(json);
    if (!parseado.is_valid) {
        printf("Error al parsear JSON: %s\n", parseado.error);
        FreeJsonResult(&parseado);
        FreeString(json);
        return 1;
    }
    FreeJsonResult(&parseado);
    
    int es_array = (json[0] == '[');
    
    if (es_array) {
        JsonResult longitud = GetArrayLength(json);
        if (!longitud.is_valid) {
            printf("Error al obtener longitud: %s\n", longitud.error);
            FreeJsonResult(&longitud);
            FreeString(json);
            return 1;
        }
        
        int num_elementos = atoi(longitud.value);
        FreeJsonResult(&longitud);
        
        printf("Total de elementos: %d\n\n", num_elementos);
        
        for (int i = 0; i < num_elementos; i++) {
            JsonResult elemento = GetArrayItem(json, i);
            if (!elemento.is_valid) {
                printf("Error al obtener elemento %d: %s\n", i, elemento.error);
                continue;
            }
            
            mostrar_elemento_json(elemento.value, i);
            
            FreeJsonResult(&elemento);
        }
    } else {
        mostrar_elemento_json(json, 0);
    }

    FreeString(json);
    
    return 0;
}
