# PicoShare - Instrucciones para Claude

## Workflow de Desarrollo

### Después de cada modificación de código:

1. **Construir imagen Docker:**
   ```bash
   docker build -t picoshare .
   ```

2. **Ejecutar la nueva imagen:**
   ```bash
   docker run -p 4001:4001 picoshare
   ```

3. **Verificar funcionamiento:**
   - Acceder a http://10.10.10.202:4001 (usar IP en lugar de localhost)
   - Probar la funcionalidad modificada
   - Verificar que el modo oscuro funciona correctamente

### Comandos útiles:

- **Detener contenedor actual:**
  ```bash
  docker stop $(docker ps -q --filter ancestor=picoshare)
  ```

- **Ver logs del contenedor:**
  ```bash
  docker logs $(docker ps -q --filter ancestor=picoshare)
  ```

- **Ejecutar en modo desarrollo con hot reload:**
  ```bash
  docker run -p 4001:4001 -v $(pwd):/app picoshare
  ```

### Estructura del proyecto:

- `handlers/static/css/style.css` - Estilos CSS principales con variables para modo oscuro/claro
- `handlers/static/js/theme-toggle.js` - JavaScript para alternar tema
- `handlers/templates/` - Plantillas HTML de Go
- `Dockerfile` - Configuración de contenedor Docker

### Notas importantes:

- Siempre probar cambios en Docker antes de confirmar
- El modo oscuro usa variables CSS para facilitar mantenimiento
- Las transiciones están configuradas para cambios suaves de tema
- El toggle del tema guarda preferencia en localStorage
- **IMPORTANTE:** Usar siempre la IP 10.10.10.202 en lugar de localhost para pruebas