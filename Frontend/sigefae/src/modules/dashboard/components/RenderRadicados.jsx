import { isFinalState, formatCurrency } from "../helpers/formatters";

export default function RenderRadicados({
  radicados, selectedRadicadoId, loadingRadicados, searchRadicados, setSearchRadicados,
  sortRadicados, setSortRadicados, tareasPorRadicado, getFilteredRadicados, openSaia
}) {
  const filteredList = getFilteredRadicados();
  return (
    <div className="fullwidth-list-container">
      <div className="fullwidth-list-header">
        <h3><i className="fa-solid fa-stamp"></i> Documentos Radicados ({filteredList.length})</h3>
        <div className="list-controls">
          <input type="text" placeholder="Buscar radicado, doc, proveedor..." value={searchRadicados} onChange={e => setSearchRadicados(e.target.value)} />
          <select value={sortRadicados} onChange={e => setSortRadicados(e.target.value)}>
            <option value="fecha_desc">Más recientes</option>
            <option value="fecha_asc">Más antiguos</option>
            <option value="estado">Por Estado</option>
          </select>
        </div>
      </div>
      <div className="fullwidth-list-body">
        {loadingRadicados ? <p style={{ padding: "15px" }}>Cargando...</p> : filteredList.length === 0 ? (
          <div className="correo-empty"><i className="fa-solid fa-inbox"></i><p>No hay documentos radicados.</p></div>
        ) : (
          <div className="card-grid">
            {filteredList.map((rad) => (
              <div key={rad.id} className={`rad-card ${selectedRadicadoId === rad.id ? "active" : ""}`} onClick={() => openSaia(rad, "radicados")}>
                <div className="rad-card-top">
                  <span className="rad-card-tipo">{rad.documento_comercial?.tipo || "—"}</span>
                  <span className={`status-badge ${rad.estado_posesion === "Completado" ? "radicado" : rad.estado_posesion === "EnProceso" ? "doc-pendiente" : (rad.estado_posesion === "Devuelto" || rad.estado_posesion === "Rechazado") ? "doc-rechazado" : ""}`}>
                    {rad.estado_posesion === "Completado" ? "Finalizado" : rad.estado_posesion === "EnProceso" ? "En Proceso" : (rad.estado_posesion === "Devuelto" || rad.estado_posesion === "Rechazado") ? "Rechazado" : rad.estado_posesion === "Libre" ? "Libre" : "En espera"}
                  </span>
                </div>
                <div className="rad-card-numero"><i className="fa-solid fa-hashtag"></i> {rad.numero_radicado}</div>
                <div className="rad-card-doc">{rad.documento_comercial?.numero_documento || "Sin documento"}</div>
                <div className="rad-card-bottom">
                  <span className="rad-card-date"><i className="fa-regular fa-calendar"></i> {new Date(rad.fecha_radicacion).toLocaleDateString()}</span>
                  {(!isFinalState(rad.estado_posesion)) && rad.usuario_actual?.nombre && (
                    <span className="rad-card-user" title={`${rad.usuario_actual.nombre} — ${tareasPorRadicado[rad.id]?.descripcion || rad.paso_actual?.nombre || "Sin paso"}`}>
                      <i className="fa-solid fa-user"></i> {rad.usuario_actual.nombre}
                    </span>
                  )}
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}