# Integración de Authentik con PicoShare

## Descripción

PicoShare ahora soporta autenticación a través de Authentik, permitiendo:
- Single Sign-On (SSO) empresarial
- Gestión centralizada de usuarios
- Soporte para múltiples proveedores de identidad
- Compatibilidad con reverse proxy

## Configuración

### 1. Variables de Entorno

```bash
# Cambiar el proveedor de autenticación
PS_AUTH_PROVIDER=authentik

# Configuración de Authentik
PS_AUTHENTIK_URL=https://auth.yourdomain.com
PS_AUTHENTIK_CLIENT_ID=picoshare
PS_AUTHENTIK_CLIENT_SECRET=your-client-secret
PS_AUTHENTIK_REDIRECT_URI=http://localhost:4001/auth/callback

# Opcional: Modo Reverse Proxy
PS_AUTHENTIK_REVERSE_PROXY=true
PS_AUTHENTIK_TRUSTED_PROXIES=127.0.0.1,10.0.0.1
```

### 2. Configuración en Authentik

#### Crear Aplicación OAuth2/OIDC:

1. Ir a **Applications > Applications**
2. Crear nueva aplicación:
   - **Name**: PicoShare
   - **Slug**: picoshare
   - **Provider**: Crear nuevo OAuth2/OpenID Provider

3. Configurar Provider:
   - **Client Type**: Confidential
   - **Client ID**: picoshare
   - **Client Secret**: (generar automáticamente)
   - **Redirect URIs**: `http://localhost:4001/auth/callback`
   - **Scopes**: `openid profile email`

#### Para Reverse Proxy (Recomendado):

1. Ir a **Applications > Outposts**
2. Crear nuevo Outpost de tipo **Proxy**
3. Configurar:
   - **External Host**: `http://localhost:4001`
   - **Internal Host**: `http://localhost:4001`
   - **Forward Auth**: Habilitado

## Modos de Operación

### Modo OAuth2 Directo
- PicoShare maneja el flujo OAuth2 completo
- Redirecciona a Authentik para login
- Maneja tokens JWT internamente

### Modo Reverse Proxy (Recomendado)
- Authentik proxy maneja toda la autenticación
- PicoShare confía en headers HTTP
- Headers esperados:
  - `X-Authentik-Username`
  - `X-Authentik-Email`
  - `X-Authentik-Groups`

## Ejemplo Docker Compose

```yaml
version: '3.8'
services:
  picoshare:
    build: .
    ports:
      - "4001:4001"
    environment:
      - PS_AUTH_PROVIDER=authentik
      - PS_AUTHENTIK_URL=https://auth.example.com
      - PS_AUTHENTIK_CLIENT_ID=picoshare
      - PS_AUTHENTIK_CLIENT_SECRET=your-secret
      - PS_AUTHENTIK_REDIRECT_URI=http://localhost:4001/auth/callback
      - PS_AUTHENTIK_REVERSE_PROXY=true
      - PS_AUTHENTIK_TRUSTED_PROXIES=172.18.0.1
    volumes:
      - ./data:/app/data
```

## Compatibilidad

La integración mantiene **compatibilidad completa** con el sistema de autenticación por secreto compartido existente. Para continuar usando el método anterior:

```bash
PS_AUTH_PROVIDER=shared_secret
PS_SHARED_SECRET=your-existing-secret
```

## Flujo de Autenticación

### OAuth2 Mode:
1. Usuario accede a PicoShare
2. Redirección a Authentik para login
3. Usuario se autentica en Authentik
4. Callback con código de autorización
5. PicoShare intercambia código por token
6. Token almacenado en cookie segura

### Reverse Proxy Mode:
1. Authentik proxy intercepta todas las requests
2. Autentica usuario si es necesario
3. Inyecta headers con información del usuario
4. PicoShare valida headers de proxy confiable
5. Extrae información del usuario de headers

## Arquitectura

```
handlers/auth/
├── provider.go           # Interface común y MultiProviderAuthenticator
├── authentik/
│   └── authentik.go     # Implementación Authentik (OAuth2 + Reverse Proxy)
└── shared_secret/
    ├── shared_secret.go # Implementación original
    └── adapter.go       # Adaptador para nueva interfaz
```