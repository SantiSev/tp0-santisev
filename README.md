# TP0 | Parte 2: Repaso de Comunicaciones | Ejercicio 8

Esta documentación sirve como referencia sobre el funcionamiento del código y las decisiones tomadas para resolver los ejercicios.

# Decisiones Tomadas

Para esta seccion, se modifca el server para manejar multiples clientes en paralelo

> Todo el código está redactado en inglés, con excepción de algunos logs específicos que permanecen en español para garantizar la compatibilidad con los tests proporcionados.

# Cambios en la Arquitectura del Servidor

## Server

Se modifica la arquitectura del servidor para manejar múltiples clientes en paralelo utilizando multiprocessing, permitiendo procesamiento concurrente de apuestas.

### Cambios Implementados

- **Procesamiento paralelo**: Cada cliente se maneja en un proceso separado usando `multiprocessing.Process`
- **File locking**: Implementa `Lock()` para sincronizar el acceso al archivo de apuestas entre procesos
- **Gestión de procesos**: Mantiene una lista de procesos activos para control y cleanup
- **Función externa**: `handle_client()` se ejecuta en proceso separado para aislamiento

**Características del flujo:**

- **Creación inmediata**: Cada cliente conectado genera un proceso instantáneamente
- **Ejecución paralela**: Los procesos corren simultáneamente sin bloquearse entre sí
- **Función externa**: `handle_client()` se ejecuta fuera del contexto del servidor principal
- **Sincronización final**: `join()` garantiza que todos los procesos terminen antes del sorteo, en caso de haber un error con alguno de los clientes, se desconecta gracefully el cliente del server
- **Cleanup automático**: Los procesos se liberan automáticamente al completarse


## Session

Se modifica el manejo de sesiones para procesar múltiples clientes de forma concurrente y gestionar desconexiones robustamente durante el procesamiento paralelo.

### Cambio en el Modelo de Manejo de Errores

**Implementación anterior:**
- `client_session.begin()` retornaba un booleano indicando éxito/fallo
- El servidor manejaba la desconexión directamente a través del `ClientManager`
> Esto se debe a que en el ejercicio 7 el procesamiento es secuencial y es facil hacer la limpieza de esta forma

**Implementación actual:**
- `client_session.begin()` propaga excepciones directamente (`raise exception`)
- El proceso externo `handle_client()` captura la excepción y termina con `exit(1)`
- Los procesos exitosos terminan con `exit(0)`

### Ventajas del Nuevo Diseño

**Detección automática de fallos:**
Al ejecutar `join()` en todos los procesos, el servidor puede identificar automáticamente qué clientes fallaron mediante sus códigos de salida. Los procesos con `exitcode != 0` indican clientes que experimentaron errores o desconexiones.

**Cleanup selectivo:**
El `ClientManager` desconecta únicamente a los clientes fallidos, manteniendo activos aquellos que completaron exitosamente su procesamiento. Esto permite que la lotería proceda normalmente con los participantes válidos.


### Flujo de Procesamiento

```
1. Múltiples procesos ejecutan clientes en paralelo
2. Cada proceso termina con exit(0) (éxito) o exit(1) (fallo)
3. El servidor ejecuta join() y examina códigos de salida
4. Clientes fallidos se remueven del ClientManager
5. La lotería procede solo con clientes exitosos
6. Resultados se envían únicamente a participantes válidos
```

**Resultado:** Sistema robusto que garantiza la entrega de resultados de lotería a todos los clientes que completaron exitosamente el procesamiento de sus apuestas, sin verse afectado por desconexiones o fallos de otros participantes.


## Business
Se modifica `LotteryService` para manejar acceso concurrente al archivo de apuestas mediante sincronización con locks.

**Problema identificado:**
- Múltiples procesos intentan escribir simultáneamente en `bets.csv`
- Sin sincronización, esto puede causar corrupción de datos o escrituras perdidas
- Race conditions al acceder al recurso compartido (archivo)


# Cómo Ejecutar

1. generar un archivo .yaml de docker-compose mediante la funcion

```bash
./generar-compose.sh docker-compose-dev.yaml 2
```

2. **Limpieza inicial**: Ejecutar `make docker-compose-down` para asegurar un inicio limpio
3. **Inicio de contenedores**: Ejecutar `make docker-compose-up` para iniciar los contenedores de servidor y cliente
4. **Visualización de logs**: Ejecutar `make docker-compose-logs` para ver los resultados y outputs del servidor y clientes
5. **Verificación de estado**: Ejecutar `docker ps -a` para confirmar que los contenedores finalizaron con exit status 0

## Script de Automatización

> **Alternativa conveniente:** Se incluye el script `run_local_test.sh` que automatiza los primeros 3 comandos y genera un archivo `logs.txt` con el output de `make docker-compose-logs` para visualización offline.

### Uso del script:

```bash
./run_local_test.sh
```

Este script ejecuta automáticamente toda la secuencia de testing y guarda los logs en un archivo para análisis posterior.
