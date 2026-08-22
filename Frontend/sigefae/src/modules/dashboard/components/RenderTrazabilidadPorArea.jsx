import "./RenderTrazabilidadPorArea.css";

export default function RenderTrazabilidadPorArea({
  areas, areaSeleccionada, setAreaSeleccionada,
  fechaDesde, setFechaDesde, fechaHasta, setFechaHasta,
  radicados, loading, seleccionados,
  buscar, toggleSeleccion, toggleTodos
}) {
  const areaNombre = areas.find(a => String(a.id) === String(areaSeleccionada))?.nombre || "Área";

  const handleDescargar = () => {
    const seleccionadosList = radicados.filter(r => seleccionados.has(r.id));
    if (seleccionadosList.length === 0) {
      alert("Seleccione al menos un radicado para descargar");
      return;
    }
    import("../../../utils/trazabilidadPorAreaPdf.js").then(mod => {
      mod.generarTrazabilidadPorAreaPDF(areaNombre, fechaDesde, fechaHasta, seleccionadosList);
    });
  };

  const estadoColor = (estado) => {
    const map = {
      "enproceso": "#3498db", "en proceso": "#3498db", "completado": "#27ae60",
      "rechazado": "#e74c3c", "devuelto": "#f39c12", "libre": "#95a5a6",
      "tomado": "#9b59b6", "pendiente": "#f39c12",
    };
    return map[estado?.toLowerCase()] || "#7f8c8d";
  };

  const formatFechaCorta = (fecha) =>
    fecha ? new Date(fecha).toLocaleDateString("es-CO", { day: "2-digit", month: "short", year: "numeric" }) : "—";

  const formatFechaLarga = (fecha) =>
    fecha ? new Date(fecha).toLocaleString("es-CO", { day: "2-digit", month: "2-digit", year: "numeric", hour: "2-digit", minute: "2-digit" }) : "—";

  return (
    <div className="trazabilidad-area-container">
      
      {/* ── Filtros ── */}
      <div className="tza-filters-card">
        <div className="tza-filters-header">
          <div className="tza-filters-icon"><i className="fa-solid fa-route" /></div>
          <div>
            <h3 className="tza-filters-title">Trazabilidad por Área</h3>
            <p className="tza-filters-subtitle">Filtra por área y rango de fechas para generar reportes</p>
          </div>
        </div>

        <div className="tza-filters-row">
          <div className="tza-form-group" style={{ flex: 1, minWidth: 220 }}>
            <label>Área</label>
            <select value={areaSeleccionada} onChange={e => setAreaSeleccionada(e.target.value)}>
              <option value="">Seleccione un área...</option>
              {areas.map(a => <option key={a.id} value={a.id}>{a.nombre}</option>)}
            </select>
          </div>

          <div className="tza-form-group">
            <label>Desde</label>
            <input type="date" value={fechaDesde} onChange={e => setFechaDesde(e.target.value)} />
          </div>

          <div className="tza-form-group">
            <label>Hasta</label>
            <input type="date" value={fechaHasta} onChange={e => setFechaHasta(e.target.value)} />
          </div>

          <button className="tza-btn tza-btn-primary" onClick={buscar} disabled={loading}>
            <i className={loading ? "fas fa-spinner fa-spin" : "fas fa-search"} />
            {loading ? "Buscando..." : "Buscar"}
          </button>

          {radicados.length > 0 && (
            <button className="tza-btn tza-btn-success" onClick={handleDescargar}>
              <i className="fas fa-file-pdf" /> Descargar seleccionadas
            </button>
          )}
        </div>
      </div>

      {/* ── Resumen ── */}
      {radicados.length > 0 && (
        <div className="tza-summary">
          <div className="tza-summary-badges">
            <span className="tza-badge tza-badge-default">
              {radicados.length} radicado{radicados.length !== 1 ? "s" : ""}
            </span>
            <span className={`tza-badge ${seleccionados.size > 0 ? "tza-badge-active" : "tza-badge-default"}`}>
              {seleccionados.size} seleccionado{seleccionados.size !== 1 ? "s" : ""}
            </span>
          </div>

          <label className="tza-select-all">
            <input
              type="checkbox"
              checked={seleccionados.size === radicados.length && radicados.length > 0}
              onChange={toggleTodos}
            />
            Seleccionar todos
          </label>
        </div>
      )}

      {/* ── Estado vacío ── */}
      {radicados.length === 0 && !loading && (
        <div className="tza-empty">
          <i className="fas fa-inbox" />
          <h4>Sin resultados</h4>
          <p>Seleccione un área y un rango de fechas, luego pulse <strong>Buscar</strong>.</p>
        </div>
      )}

      {/* ── Lista ── */}
      {radicados.length > 0 && (
        <div className="tza-radicados-list">
          {radicados.map(rad => {
            const checked = seleccionados.has(rad.id);
            return (
              <div key={rad.id} className={`tza-card ${checked ? "tza-card-selected" : ""}`}>
                
                <div className={`tza-card-header ${checked ? "tza-card-header-selected" : ""}`}>
                  <input
                    type="checkbox"
                    className="tza-card-checkbox"
                    checked={checked}
                    onChange={() => toggleSeleccion(rad.id)}
                  />

                  <div className="tza-card-info">
                    <div className="tza-card-title-row">
                      <span className="tza-card-numero">{rad.numero_radicado}</span>
                      <span className="tza-card-estado" style={{ background: estadoColor(rad.estado_posesion) }}>
                        {rad.estado_posesion}
                      </span>
                    </div>
                    <div className="tza-card-proveedor">
                      <i className="fas fa-building" />
                      <span>{rad.proveedor?.razon_social || "Sin proveedor"}</span>
                      {rad.proveedor?.numero_documento && (
                        <span className="tza-card-nit">• NIT {rad.proveedor.numero_documento}</span>
                      )}
                    </div>
                  </div>

                  <div className="tza-card-meta">
                    <div className="tza-card-fecha">
                      <i className="far fa-calendar-alt" /> {formatFechaCorta(rad.fecha_radicacion)}
                    </div>
                    <div className="tza-card-eventos">
                      {rad.trazabilidad?.length || 0} evento{(rad.trazabilidad?.length || 0) !== 1 ? "s" : ""}
                    </div>
                  </div>
                </div>

                <div className="tza-table-scroll">
                  {rad.trazabilidad && rad.trazabilidad.length > 0 ? (
                    <table className="tza-table">
                      <thead>
                        <tr>
                          <th style={{ width: 150 }}>Fecha</th>
                          <th style={{ width: 140 }}>Acción</th>
                          <th>Descripción</th>
                          <th style={{ width: 130 }}>Usuario</th>
                        </tr>
                      </thead>
                      <tbody>
                        {rad.trazabilidad.map((t, i) => (
                          <tr key={i}>
                            <td className="tza-table-fecha">{formatFechaLarga(t.fecha)}</td>
                            <td><span className="tza-table-accion">{t.accion}</span></td>
                            <td className="tza-table-desc">
                              {t.descripcion || <span className="tza-table-desc-empty">Sin descripción</span>}
                            </td>
                            <td className="tza-table-usuario">
                              <i className="far fa-user" />{t.usuario_nombre || "Sistema"}
                            </td>
                          </tr>
                        ))}
                      </tbody>
                    </table>
                  ) : (
                    <div className="tza-table-no-data">
                      <i className="fas fa-info-circle" /> Sin registros de trazabilidad para este radicado.
                    </div>
                  )}
                </div>
              </div>
            );
          })}
        </div>
      )}
    </div>
  );
}