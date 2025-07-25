# PicoShare - Instrucciones para Claude

## Workflow de Desarrollo

### Después de cada modificación de código:

1. **Construir imagen Docker:**
   ```bash
   docker build -t picoshare .
   ```
   
   **IMPORTANTE:** Si el `docker build` se demora más de 2-3 minutos:
   - Cancelar el proceso (Ctrl+C)
   - Volver a ejecutar el mismo comando
   - NO intentar soluciones alternativas como builds locales
   - Repetir hasta que el build se complete exitosamente

2. **Ejecutar la nueva imagen:**
   ```bash
   docker run -p 4001:4001 picoshare
   ```

3. **Verificar funcionamiento:**
   - Acceder a http://10.10.10.202:4001 (usar IP en lugar de localhost)
   - Probar la funcionalidad modificada
   - Verificar que el modo oscuro funciona correctamente

4. **Pruebas automatizadas con Playwright:**
   - Usar el servidor MCP de Playwright para validar cambios automáticamente
   - Especialmente importante para cambios de UI, autenticación, y flujos de navegación
   - Ejecutar tests antes de considerar los cambios completos

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

### Configuración de autenticación:

- **Authentik Integration**: El proyecto soporta autenticación con Authentik OAuth2
- **Variables de entorno**: Configurar en `.env` el tipo de autenticación (`PS_AUTH_PROVIDER`)
- **Modo Shared Secret**: Para desarrollo simple con `PS_SHARED_SECRET`
- **Modo Authentik**: Para SSO empresarial con configuración OAuth2 completa

#### Problemas resueltos de autenticación OAuth2:

1. **Error 404 en URLs `/auth/{token}`**: 
   - **Problema**: El callback OAuth2 redirigía a URLs inexistentes como `/auth/wuJ5O1Tp5bVwOZWMPtGSty2OOdd2rGa2du9Bo4b-pQ0=`
   - **Causa**: En `authentik.go:244-248`, se usaba el parámetro `state` como URL de redirect
   - **Solución**: Cambiar redirect a URL fija `"/"` en lugar de usar `state` variable

2. **Bucle infinito de autenticación**:
   - **Problema**: Después de autenticación exitosa en Authentik, el usuario seguía siendo redirigido a auth
   - **Causa**: Cookie con `SameSite: http.SameSiteStrictMode` no persistía después del redirect OAuth2
   - **Solución**: Cambiar a `SameSite: http.SameSiteLaxMode` en `authentik.go:247`
   - **Archivo modificado**: `handlers/auth/authentik/authentik.go`

3. **Configuración requerida en Authentik**:
   - **Redirect URI exacta**: `http://10.10.10.202:4001/auth/callback`
   - **Scopes**: `openid profile email`
   - **Client Type**: Confidential

#### Usuario de prueba para validación Authentik:
- **Usuario**: claude-debug
- **Contraseña**: Inicio$1

### Notas importantes:

- Siempre probar cambios en Docker antes de confirmar
- El modo oscuro usa variables CSS para facilitar mantenimiento
- Las transiciones están configuradas para cambios suaves de tema
- El toggle del tema guarda preferencia en localStorage
- **IMPORTANTE:** Usar siempre la IP 10.10.10.202 en lugar de localhost para pruebas
- **Testing**: Siempre validar cambios con Playwright MCP antes de completar tareas