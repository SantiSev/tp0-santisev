# TP0 | Documentacion Index

## Ejercicio 1:
> [Documentacion](https://github.com/SantiSev/tp0-santisev/blob/ej1/README.md)

### Resumen:
- Como funciona el script de `./generar-compose.sh`
- Demostracion de como usar el script

## Ejercicio 2:
> [Documentacion](https://github.com/SantiSev/tp0-santisev/blob/ej2/README.md)

### Resumen:
- Modificacion del script `./generar-compose.sh`
- Demostracion de como usar el script

## Ejercicio 3:
> [Documentacion](https://github.com/SantiSev/tp0-santisev/blob/ej3/README.md)

### Resumen:
- Como funciona el script de `./validar-echo-server.sh`
- Demostracion de como usar el script

## Ejercicio 4:
> [Documentacion](https://github.com/SantiSev/tp0-santisev/blob/ej4/README.md)

### Resumen:
- Implementación de `signals` en el servidor para gestionar un apagado ordenado (graceful shutdown)
- Implementación de `signals` en el cliente para gestionar un apagado ordenado (graceful shutdown)

## Ejercicio 5:
> [Documentacion](https://github.com/SantiSev/tp0-santisev/blob/ej5/README.md)

### Resumen:
- Explicacion de la arquitectura base del cliente y servidor
- Explicación de cómo el cliente envía una única apuesta al servidor y cómo ambos finalizan de forma ordenada (graceful shutdown)
- Descripción del manejo eficiente de short reads y short writes
- modificaciones al script `./generar-compose.sh`
- Demostracion de como correr el programa

## Ejercicio 6:
> [Documentacion](https://github.com/SantiSev/tp0-santisev/blob/ej6/README.md)

### Resumen:
- Explicación de las modificaciones en el cliente y el servidor para que el cliente envíe múltiples apuestas en batches y el servidor las almacene correctamente
- Cambios en el script `./generar-compose.sh`
- Demostracion de como correr el programa

## Ejercicio 7:
> [Documentacion](https://github.com/SantiSev/tp0-santisev/blob/ej7/README.md)

### Resumen:
- Explicación de las modificaciones en el cliente y el servidor para que, una vez procesadas todas las apuestas de todas las agencias, se inicie la lotería y se envíen los resultados a cada cliente, mostrando también todos los resultados en el servidor
- Cambios en el script `./generar-compose.sh`
- Demostracion de como correr el programa

## Ejercicio 8:
> [Documentacion](https://github.com/SantiSev/tp0-santisev/blob/ej8/README.md)

### Resumen:
- Agregando paralelismo al servidor para poder processar multiples clientes en paralela en ves de procesarlo sequencialemnte como ocurria en el ej6
- Explicacion de como se manejan las secciones criticas
