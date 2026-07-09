# SIGEFAE Backend

Backend de **SIGEFAE**, encargado de la recepción, procesamiento y administración automática de documentos recibidos mediante Microsoft Graph.

---

# Objetivos

- Conectarse a Microsoft Graph.
- Leer periódicamente los correos de una cuenta específica.
- Detectar nuevos correos con documentos.
- Descargar y almacenar los archivos adjuntos.
- Procesar el contenido (XML, PDF, ZIP, etc.).
- Registrar únicamente los metadatos en la base de datos.
- Servir la información al frontend mediante una API REST.

---

# Arquitectura

```
Microsoft Graph
        │
        ▼
 Autenticación OAuth2
        │
        ▼
 Obtener correos
        │
        ▼
¿Correo procesado?
   │         │
  Sí         No
  │          │
 Ignorar     ▼
      Descargar adjuntos
             │
             ▼
      Almacenamiento
             │
             ▼
 Procesamiento interno
             │
             ▼
 Base de datos (metadatos)
```

---

# Estructura del proyecto

```
backend/

│
├── cmd/
│   └── main.go
│
├── internal/
│   ├── env/
│   └── graph/
│
├── downloads/
│
├── go.mod
├── go.sum
└── .env
```

---

# Variables de entorno

```
GRAPH_CLIENT_ID=
GRAPH_CLIENT_SECRET=
GRAPH_TENANT_ID=
GRAPH_USER_EMAIL=
```

---

# Microsoft Graph

Actualmente el proyecto implementa:

- Autenticación mediante Client Credentials.
- Renovación automática del Access Token.
- Consulta de correos.
- Consulta de adjuntos.
- Descarga de archivos.

---

# Flujo actual

```
Inicio

↓

Cargar variables de entorno

↓

Solicitar Access Token

↓

Consultar correos

↓

Para cada correo

    ¿Tiene adjuntos?

        Sí

            Obtener adjuntos

            Descargar archivos

        No

            Continuar

↓

Fin
```

---

# Próximas funcionalidades

## Persistencia

- PostgreSQL
- Registro de correos procesados
- Registro de documentos
- Registro de empresas
- Registro de usuarios

---

## Procesamiento

- Descomprimir ZIP
- Leer XML
- Leer PDF
- Validaciones
- Clasificación automática

---

## API

- Login
- Gestión de usuarios
- Empresas
- Documentos
- Descarga de archivos
- Auditoría

---

## Seguridad

- JWT
- Roles
- Permisos
- Auditoría de acciones

---

# Almacenamiento

Los archivos **no se almacenarán en la base de datos**.

Se almacenarán directamente en el sistema de archivos de la VM.

Ejemplo:

```
storage/

    Empresa A/

        2026/

            Julio/

                Factura.pdf
                Factura.xml

    Empresa B/

        2026/

            Agosto/

                Documento.pdf
```

La base de datos únicamente almacenará información como:

- id
- nombre
- ruta
- empresa
- fecha
- hash
- message_id de Microsoft Graph

---

# Tecnologías

- Go
- Microsoft Graph
- OAuth2
- PostgreSQL (pendiente)
- JWT (pendiente)

---

# Estado del proyecto

## Implementado

- [x] Configuración mediante `.env`
- [x] Cliente Microsoft Graph
- [x] OAuth2 Client Credentials
- [x] Renovación automática del Token
- [x] Lectura de correos
- [x] Lectura de adjuntos
- [x] Descarga de archivos

## En desarrollo

- [ ] Evitar reprocesar correos
- [ ] Base de datos
- [ ] API REST
- [ ] Procesamiento de ZIP
- [ ] Extracción XML
- [ ] Extracción PDF
- [ ] Gestión de usuarios
- [ ] Autenticación JWT
- [ ] Panel administrativo

---

# Licencia

Proyecto privado desarrollado para **SIGEFAE**.