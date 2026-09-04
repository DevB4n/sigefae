export default function RenderNormasReparto({
  normasRepartoRadicado, radicadoId, readOnly = false,
  esAdmin, puedeGestionarRecurso, openNormaModal, handleEliminarNorma
}) {
  const totalPct = (normasRepartoRadicado || []).reduce((s, n) => s + (parseFloat(n.porcentaje) || 0), 0);
  const estaCompleto = Math.abs(totalPct - 100) < 0.01;

  return (
    <div className="doc-section">
      <h4 style={{ margin: 0, display: 'flex', alignItems: 'center', gap: 8, justifyContent: 'space-between' }}>
        <span style={{ display: 'inline-flex', alignItems: 'center', gap: 8 }}><i className="fa-solid fa-chart-pie"></i> Normas de Reparto ({(normasRepartoRadicado || []).length})</span>
        {!readOnly && (
          <button className="doc-btn doc-btn-secondary" onClick={() => openNormaModal(radicadoId)} style={{ padding: "6px 12px", fontSize: "0.8em" }}>
            <i className="fa-solid fa-plus"></i> Agregar Norma
          </button>
        )}
      </h4>

      {(normasRepartoRadicado || []).length === 0 ? (
        <p style={{ color: "#6b7280", fontSize: "0.9em" }}>No hay normas de reparto asignadas.</p>
      ) : (
        <>
          <table className="doc-items-table">
            <thead>
              <tr>
                <th>Código</th>
                <th>Nombre</th>
                <th>Proyecto</th>
                <th>Sede</th>
                <th>Área</th>
                <th style={{ textAlign: "right" }}>%</th>
                {!readOnly && <th style={{ width: 90, textAlign: "center" }}>Acciones</th>}
              </tr>
            </thead>
            <tbody>
              {normasRepartoRadicado.map((n) => (
                <tr key={n.id}>
                  <td><strong>{n.norma_reparto?.codigo || n.codigo}</strong></td>
                  <td>{n.norma_reparto?.nombre || n.nombre}</td>
                  <td>{n.norma_reparto?.proyecto || n.proyecto || "—"}</td>
                  <td>{n.norma_reparto?.sucursal || n.sucursal}</td>
                  <td>{n.norma_reparto?.departamento || n.departamento}</td>
                  <td style={{ textAlign: "right", fontWeight: 700 }}>{parseFloat(n.porcentaje).toFixed(2)}%</td>
                  {!readOnly && (
                    <td style={{ textAlign: "center" }}>
                      {esAdmin || puedeGestionarRecurso(n.creado_por_id || n.creado_por?.id) ? (
                        <div style={{ display: "flex", gap: 6, justifyContent: "center" }}>
                          <button className="btn-icon btn-edit" onClick={() => openNormaModal(radicadoId, n)} title="Editar">
                            <i className="fa-solid fa-pen"></i>
                          </button>
                          <button className="btn-icon btn-toggle" onClick={() => handleEliminarNorma(n, radicadoId)} title="Eliminar">
                            <i className="fa-solid fa-xmark"></i>
                          </button>
                        </div>
                      ) : null}
                    </td>
                  )}
                </tr>
              ))}
            </tbody>
          </table>

          <div style={{
            marginTop: 10,
            padding: "10px 14px",
            background: estaCompleto ? "#d1fae5" : "#fef3c7",
            borderRadius: 8,
            display: "flex",
            justifyContent: "space-between",
            alignItems: "center",
            fontSize: "0.9em",
            fontWeight: 700,
            color: estaCompleto ? "#065f46" : "#92400e",
            border: `1px solid ${estaCompleto ? "#a7f3d0" : "#fde68a"}`
          }}>
            <span><i className={`fa-solid ${estaCompleto ? "fa-check-circle" : "fa-triangle-exclamation"}`} style={{ marginRight: 6 }}></i> Total asignado</span>
            <span>{totalPct.toFixed(2)} % {estaCompleto ? "(Completo)" : `(Faltan ${(100 - totalPct).toFixed(2)} %)`}</span>
          </div>
        </>
      )}
    </div>
  );
}
