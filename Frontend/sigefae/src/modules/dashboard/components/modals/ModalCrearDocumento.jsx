import { formatCurrency } from "../../helpers/formatters";

export default function ModalCrearDocumento({
  showCrearDocModal, setShowCrearDocModal, docForm, creandoDoc,
  proveedoresCatalogo, receptoresCatalogo, areasCatalogo, monedasCatalogo,
  handleDocFormChange, handleAddDetalle, handleDetalleChange, handleRemoveDetalle,
  handleCrearDocumentoSubmit, openCrearProveedorModal
}) {
  if (!showCrearDocModal) return null;
  return (
    <div className="modal-overlay" onClick={() => setShowCrearDocModal(false)}>
      <div className="modal-content" onClick={(e) => e.stopPropagation()} style={{ maxWidth: 800, maxHeight: "90vh", overflowY: "auto" }}>
        <div className="modal-header">
          <h3><i className="fa-solid fa-file-circle-plus"></i> Crear Documento Manual</h3>
          <button className="modal-close" onClick={() => setShowCrearDocModal(false)}><i className="fa-solid fa-xmark"></i></button>
        </div>
        <div className="modal-body">
          <div style={{ display: "grid", gridTemplateColumns: "1fr 1fr", gap: 12 }}>
            <div className="modal-field">
              <label>Tipo <span className="required">*</span></label>
              <select name="tipo" value={docForm.tipo} onChange={handleDocFormChange} className="doc-input">
                <option value="FACTURA_FISICA">Factura Física</option>
                <option value="CUENTA_COBRO">Cuenta de Cobro</option>
                <option value="NOTA_CREDITO">Nota Crédito</option>
                <option value="NOTA_DEBITO">Nota Débito</option>
              </select>
            </div>
            <div className="modal-field">
              <label>Número Documento <span className="required">*</span></label>
              <input type="text" name="numero_documento" value={docForm.numero_documento} onChange={handleDocFormChange} className="doc-input" placeholder="Ej: F001-1234" />
            </div>
            <div className="modal-field">
              <label style={{ display: "flex", justifyContent: "space-between", alignItems: "center" }}>
                <span>Proveedor <span className="required">*</span></span>
                <button type="button" className="btn-icon btn-edit" onClick={openCrearProveedorModal} title="Nuevo Proveedor" style={{ width: 24, height: 24, fontSize: "0.8em" }}><i className="fa-solid fa-plus"></i></button>
              </label>
              <select name="id_proveedor" value={docForm.id_proveedor} onChange={handleDocFormChange} className="doc-input">
                <option value="">Seleccione...</option>{proveedoresCatalogo.map(p => <option key={p.id} value={p.id}>{p.razon_social} — {p.numero_documento}</option>)}
              </select>
            </div>
            <div className="modal-field">
              <label>Receptor <span className="required">*</span></label>
              <select name="id_receptor" value={docForm.id_receptor} onChange={handleDocFormChange} className="doc-input">
                <option value="">Seleccione...</option>{receptoresCatalogo.map(r => <option key={r.id} value={r.id}>{r.nombre} — {r.numero_documento}</option>)}
              </select>
            </div>
            <div className="modal-field">
              <label>Área <span className="required">*</span></label>
              <select name="id_area" value={docForm.id_area} onChange={handleDocFormChange} className="doc-input">
                <option value="">Seleccione...</option>{areasCatalogo.map(a => <option key={a.id} value={a.id}>{a.nombre}</option>)}
              </select>
            </div>
            <div className="modal-field">
              <label>Moneda <span className="required">*</span></label>
              <select name="moneda_id" value={docForm.moneda_id} onChange={handleDocFormChange} className="doc-input">
                <option value="">Seleccione...</option>{monedasCatalogo.map(m => <option key={m.id} value={m.id}>{m.nombre} ({m.codigo})</option>)}
              </select>
            </div>
            <div className="modal-field">
              <label>Fecha Emisión <span className="required">*</span></label>
              <input type="date" name="fecha_documento" value={docForm.fecha_documento} onChange={handleDocFormChange} className="doc-input" />
            </div>
            <div className="modal-field" style={{ gridColumn: "1 / -1" }}>
              <label>Asunto / Observación</label>
              <input type="text" name="asunto" value={docForm.asunto} onChange={handleDocFormChange} className="doc-input" placeholder="Descripción general del documento" />
            </div>
          </div>
          <div className="doc-section" style={{ marginTop: 20 }}>
            <div style={{ display: "flex", justifyContent: "space-between", alignItems: "center", marginBottom: 10 }}>
              <h4><i className="fa-solid fa-list"></i> Detalle de Ítems</h4>
              <button className="doc-btn doc-btn-secondary" onClick={handleAddDetalle} type="button"><i className="fa-solid fa-plus"></i> Agregar Ítem</button>
            </div>
            {docForm.detalles.length === 0 ? <p style={{ color: "#6b7280", fontSize: "0.9em" }}>Sin ítems. Agregue al menos uno.</p> : (
              <table className="doc-items-table">
                <thead><tr><th>Descripción</th><th style={{ width: 90 }}>Cant.</th><th style={{ width: 130 }}>Valor Unit.</th><th style={{ width: 110 }}>IVA Unit.</th><th style={{ width: 130 }}>Total Línea</th><th style={{ width: 40 }}></th></tr></thead>
                <tbody>
                  {docForm.detalles.map((d, idx) => (
                    <tr key={idx}>
                      <td><input type="text" value={d.descripcion} onChange={(e) => handleDetalleChange(idx, "descripcion", e.target.value)} className="doc-input" placeholder="Descripción" /></td>
                      <td><input type="number" value={d.cantidad} onChange={(e) => handleDetalleChange(idx, "cantidad", e.target.value)} className="doc-input" min="0" step="0.01" /></td>
                      <td><input type="number" value={d.valor_unitario} onChange={(e) => handleDetalleChange(idx, "valor_unitario", e.target.value)} className="doc-input" min="0" /></td>
                      <td><input type="number" value={d.iva_unitario} onChange={(e) => handleDetalleChange(idx, "iva_unitario", e.target.value)} className="doc-input" min="0" /></td>
                      <td style={{ textAlign: "right", fontWeight: 600 }}>{formatCurrency(d.total)}</td>
                      <td><button className="btn-icon btn-toggle" onClick={() => handleRemoveDetalle(idx)} title="Eliminar"><i className="fa-solid fa-xmark"></i></button></td>
                    </tr>
                  ))}
                </tbody>
              </table>
            )}
          </div>
          <div className="doc-section" style={{ marginTop: 16, background: "#f9fafb", padding: 12, borderRadius: 8 }}>
            <div className="doc-totals">
              <div className="doc-total-row"><span>Subtotal</span><span>{formatCurrency(docForm.subtotal)}</span></div>
              <div className="doc-total-row"><span>IVA</span><span>{formatCurrency(docForm.iva)}</span></div>
              <div className="doc-total-row total-final"><span>Total Documento</span><span>{formatCurrency(docForm.total)}</span></div>
            </div>
          </div>
        </div>
        <div className="modal-footer">
          <button className="doc-btn doc-btn-secondary" onClick={() => setShowCrearDocModal(false)} disabled={creandoDoc}>Cancelar</button>
          <button className="doc-btn doc-btn-primary" onClick={handleCrearDocumentoSubmit} disabled={creandoDoc}>
            <i className="fa-solid fa-file-circle-plus"></i> {creandoDoc ? "Creando..." : "Crear Documento"}
          </button>
        </div>
      </div>
    </div>
  );
}