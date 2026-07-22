// Package db contiene los modelos GORM generados a partir del diagrama
// entidad-relación de Sigefae, junto con la conexión a la base de datos.
//
// Los modelos están organizados en varios archivos por dominio/entidad
// (ver catalogos.go, geografia.go, usuario.go, proveedor.go, receptor.go,
// correo.go, archivo.go, factura.go, pago.go, documento_radicado.go)
// en lugar de un único schema.go, para que sea más fácil de mantener y leer.
//
// Convenciones usadas:
//   - Todas las PK se llaman "id" y son de tipo uint (INTEGER en el diagrama).
//   - Los campos FK se representan como <nombre>ID uint, con su tag
//     `gorm:"column:..."` respetando el nombre exacto mostrado en el diagrama.
//   - Se agregan structs de asociación (belongs to) con su propio tag
//     `foreignKey`/`references` para poder hacer Preload() cuando se necesite.
//   - TIMESTAMP y DATETIME -> time.Time
//   - VARCHAR             -> string
//   - INTEGER              -> uint / int (uint cuando es PK/FK, int cuando es
//     un campo numérico simple como "peso" o "orden")
//   - BIGINT               -> int64
//   - DOUBLE                -> float64
//   - ENUM (orientacion_sello_recibido) -> string, se documenta el dominio
//     esperado en un comentario porque el diagrama no lista los valores.
//
// Cada struct define su propio TableName() para que el nombre de la tabla en
// la base de datos coincida EXACTAMENTE con el nombre mostrado en el diagrama
// (evita que el pluralizador de GORM cambie nombres como "pais" -> "pais" de
// forma inconsistente).
package db
