export default function ModalRadicar({
  showRadicarModal, setShowRadicarModal, radicarForm, radicando,
  normasRepartoAutoMsg, tiposRadicacion, rutas, metodosPago,
  normasRepartoCatalogo,
  normaFiltroSede, setNormaFiltroSede, normaFiltroArea, setNormaFiltroArea,
  normaSeleccionadaId, setNormaSeleccionadaId, normaPorcentajeInput, setNormaPorcentajeInput,
  sedesDisponibles, areasDisponibles, normasFiltradas, totalPorcentajeNormas,
  handleRadicarChange, handleAgregarNormaModal, handleNormaRepartoChange, handleRemoveNormaReparto, handleRadicarSubmit
}) {
  if (!showRadicarModal) return null;
  return (
    <div className="modal-overlay" onClick={() => setShowRadicarModal(false)}>
      <div className="modal-content" onClick={(e) => e.stopPropagation()}>
        <div className="modal-header">
          <h3><i className="fa-solid fa-stamp"></i> Radicar Documento</h3>
          <button className="modal-close" onClick={() => setShowRadicarModal(false)}><i className="fa-solid fa-xmark"></i></button>
        </div>
        <div className="modal-body">
          <div className="modal-field">
            <label>Tipo de Radicación <span className="required">*</span></label>
            <select name="tipo_radicacion_id" value={radicarForm.tipo_radicacion_id} onChange={handleRadicarChange} className="doc-input">
              <option value="">Seleccione...</option>{tiposRadicacion.map(tr => <option key={tr.id} value={tr.id}>{tr.nombre}</option>)}
            </select>
          </div>
          <div className="modal-field">
            <label>Ruta <span className="required">*</span></label>
            <select name="ruta_id" value={radicarForm.ruta_id} onChange={handleRadicarChange} className="doc-input">
              <option value="">Seleccione...</option>{rutas.map(r => <option key={r.id} value={r.id}>{r.nombre}</option>)}
            </select>
          </div>
          <div className="modal-field">
            <label>Método de Pago <span className="required">*</span></label>
            <select name="metodo_pago_id" value={radicarForm.metodo_pago_id} onChange={handleRadicarChange} className="doc-input">
              <option value="">Seleccione...</option>{metodosPago.map(mp => <option key={mp.id} value={mp.id}>{mp.nombre}</option>)}
            </select>
          </div>
          {normasRepartoAutoMsg && (
            <div className="modal-field" style={{ padding: "10px 12px", borderRadius: 8, background: "#fef3c7", border: "1px solid #f59e0b", color: "#92400e", fontSize: "0.92em", fontWeight: 600 }}>{normasRepartoAutoMsg}</div>
          )}
          <div className="modal-field">
            <label>Número de Radicado <small>(opcional, se autogenera si está vacío)</small></label>
            <input type="text" name="numero_radicado" value={radicarForm.numero_radicado} onChange={handleRadicarChange} className="doc-input" placeholder="Ej: RAD-2026-00001" />
          </div>
          <div className="modal-field">
            <label style={{ display: "flex", justifyContent: "space-between", alignItems: "center" }}>
              <span><i className="fa-solid fa-chart-pie"></i> Normas de Reparto</span>
              <span style={{ fontSize: "0.8em", color: totalPorcentajeNormas === 100 ? "#059669" : totalPorcentajeNormas > 0 ? "#dc2626" : "#6b7280", fontWeight: 700 }}>Total: {totalPorcentajeNormas.toFixed(2)}%</span>
            </label>
            <div style={{ display: "flex", gap: 8, marginBottom: 10 }}>
              <select className="doc-input" style={{ flex: 1 }} value={normaFiltroSede} onChange={(e) => { setNormaFiltroSede(e.target.value); setNormaFiltroArea(""); setNormaSeleccionadaId(""); }}>
                <option value="">Todas las sedes</option>{sedesDisponibles.map(s => <option key={s} value={s}>{s}</option>)}
              </select>
              <select className="doc-input" style={{ flex: 1 }} value={normaFiltroArea} onChange={(e) => { setNormaFiltroArea(e.target.value); setNormaSeleccionadaId(""); }} disabled={!normaFiltroSede}>
                <option value="">Todas las áreas</option>{areasDisponibles.map(a => <option key={a} value={a}>{a}</option>)}
              </select>
            </div>
            <div style={{ display: "flex", gap: 8, marginBottom: 10, alignItems: "center" }}>
              <select className="doc-input" style={{ flex: 2 }} value={normaSeleccionadaId} onChange={(e) => setNormaSeleccionadaId(e.target.value)} disabled={!normaFiltroSede || !normaFiltroArea}>
                <option value="">Seleccione norma...</option>{normasFiltradas.map(n => <option key={n.id} value={String(n.id)}>{n.codigo} — {n.nombre}</option>)}
              </select>
              <input type="number" min="0" max="100" step="0.01" placeholder="%" className="doc-input" style={{ flex: 0.6, textAlign: "right" }} value={normaPorcentajeInput} onChange={(e) => setNormaPorcentajeInput(e.target.value)} />
              <button className="doc-btn doc-btn-secondary" onClick={handleAgregarNormaModal} disabled={!normaSeleccionadaId || !normaPorcentajeInput} style={{ padding: "6px 12px", fontSize: "0.85em" }}>
                <i className="fa-solid fa-plus"></i> Agregar
              </button>
            </div>
            {(radicarForm.normas_reparto || []).length === 0 ? (
              <p style={{ fontSize: "0.85em", color: "#6b7280" }}>No se han asignado normas de reparto.</p>
            ) : (
              <div style={{ display: "flex", flexDirection: "column", gap: 6, marginBottom: 10 }}>
                {radicarForm.normas_reparto.map((norma, idx) => {
                  const info = normasRepartoCatalogo.find(n => String(n.id) === norma.norma_reparto_id);
                  return (
                    <div key={idx} style={{ display: "flex", gap: 8, alignItems: "center", padding: "8px 10px", background: "#f9fafb", borderRadius: 6, border: "1px solid #e5e7eb" }}>
                      <span style={{ flex: 2, fontSize: "0.85em", fontWeight: 600, color: "#1f2937" }}>{info ? `${info.codigo} — ${info.nombre}` : "Norma desconocida"}</span>
                      <span style={{ flex: 1, fontSize: "0.8em", color: "#6b7280" }}>{info ? `${info.sucursal} / ${info.departamento}` : ""}</span>
                      <input type="number" min="0" max="100" step="0.01" value={norma.porcentaje} onChange={(e) => handleNormaRepartoChange(idx, "porcentaje", e.target.value)} className="doc-input" style={{ width: 70, textAlign: "right", fontSize: "0.85em" }} />
                      <span style={{ fontSize: "0.85em", fontWeight: 700 }}>%</span>
                      <button className="btn-icon btn-toggle" onClick={() => handleRemoveNormaReparto(idx)} title="Quitar" style={{ width: 28, height: 28, flexShrink: 0 }}><i className="fa-solid fa-xmark"></i></button>
                    </div>
                  );
                })}
              </div>
            )}
          </div>
        </div>
        <div className="modal-footer">
          <button className="doc-btn doc-btn-secondary" onClick={() => setShowRadicarModal(false)} disabled={radicando}>Cancelar</button>
          <button className="doc-btn doc-btn-primary" onClick={handleRadicarSubmit} disabled={radicando}>
            <i className="fa-solid fa-stamp"></i> {radicando ? "Radicando..." : "Radicar Documento"}
          </button>
        </div>
      </div>
    </div>
  );
}