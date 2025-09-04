# TP0 | Parte 1: Introducción a Docker | Ejercicio 2

Esta documentación sirve como referencia sobre el funcionamiento del código y las decisiones tomadas para resolver los ejercicios.

# Decisiones Tomadas
El objetivo de este ejercicio es permitir que los cambios en los archivos de configuración sean efectivos sin necesidad de reconstruir las imágenes de Docker.


Este ejercicio se resolvio mediante un branch desde la rama ej1 y realizando una breve modificacion al script generar-compose.sh y simplemente agregandole un volumen al servicio del cliente que en si agrega al cliente

Seccion Modificada:

```bash
for i in $(seq 1 "$AMOUNT_CLIENTS"); do
    cat >> docker-compose-dev.yaml << EOF

  client$i:
    container_name: client$i
    image: client:latest
    entrypoint: /client
    environment:
      - CLI_ID=$i
    networks:
      - testing_net
    volumes:   < --- [ESTO FUE AGREGADO]
      - ./client/config.yaml:/config.yaml
    depends_on:
      - server
EOF
done
```
### ¿Por qué agregar un volumen?
Sin el volumen, el archivo config.yaml queda copiado estáticamente dentro de la imagen durante el proceso de build mediante instrucciones como COPY. Esto significa que cualquier modificación al archivo de configuración en el host requiere:

1. Reconstruir completamente la imagen (docker build)
2. Recrear y reiniciar todos los contenedores

Con el volumen, el contenedor accede directamente al archivo desde el sistema de archivos del host, creando un enlace dinámico que refleja los cambios inmediatamente sin necesidad de reconstruir la imagen.


# Como Ejectuar

El script se encuentra en la raíz del proyecto y recibirá por parámetro el nombre del archivo de salida y la cantidad de clientes esperados:

Ejemplo de uso:
```bash
./generar-compose.sh [DOCKER-COMPOSE.YAML] [CANT_CLIENTES]
```

Ejemplo de uso:
```bash
./generar-compose.sh docker-compose-dev.yaml 5
```


