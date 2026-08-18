export default function ModalCrearProveedor({
  showCrearProveedorModal, setShowCrearProveedorModal, proveedorForm,
  creandoProveedor, tiposDocumentoCatalogo, handleProveedorFormChange, handleCrearProveedorSubmit
}) {
  if (!showCrearProveedorModal) return null;
  return (
    <div className="modal-overlay" onClick={() => setShowCrearProveedorModal(false)}>
      <div className="modal-content" onClick={(e) => e.stopPropagation()} style={{ maxWidth: 500 }}>
        <div className="modal-header">
          <h3><i className="fa-solid fa-building"></i> Nuevo Proveedor</h3>
          <button className="modal-close" onClick={() => setShowCrearProveedorModal(false)}><i className="fa-solid fa-xmark"></i></button>
        </div>
        <div className="modal-body">
          <div className="modal-field">
            <label>Razón Social <span className="required">*</span></label>
            <input type="text" name="razon_social" value={proveedorForm.razon_social} onChange={handleProveedorFormChange} className="doc-input" placeholder="Ej: INVERCOMER DEL CARIBE SAS" />
          </div>
          <div style={{ display: "grid", gridTemplateColumns: "1fr 1fr", gap: 12 }}>
            <div className="modal-field">
              <label>Tipo Documento <span className="required">*</span></label>
              <select name="tipo_documento_id" value={proveedorForm.tipo_documento_id} onChange={handleProveedorFormChange} className="doc-input">
                <option value="">Seleccione...</option>{tiposDocumentoCatalogo.map(t => <option key={t.id} value={t.id}>{t.nombre}</option>)}
              </select>
            </div>
            <div className="modal-field">
              <label>Número Documento <span className="required">*</span></label>
              <input type="text" name="numero_documento" value={proveedorForm.numero_documento} onChange={handleProveedorFormChange} className="doc-input" placeholder="Ej: 900383385" />
            </div>
          </div>
          <div style={{ display: "grid", gridTemplateColumns: "1fr 1fr", gap: 12 }}>
            <div className="modal-field">
              <label>Email</label>
              <input type="email" name="email" value={proveedorForm.email} onChange={handleProveedorFormChange} className="doc-input" placeholder="opcional" />
            </div>
            <div className="modal-field">
              <label>Teléfono</label>
              <input type="text" name="telefono" value={proveedorForm.telefono} onChange={handleProveedorFormChange} className="doc-input" placeholder="opcional" />
            </div>
          </div>
        </div>
        <div className="modal-footer">
          <button className="doc-btn doc-btn-secondary" onClick={() => setShowCrearProveedorModal(false)} disabled={creandoProveedor}>Cancelar</button>
          <button className="doc-btn doc-btn-primary" onClick={handleCrearProveedorSubmit} disabled={creandoProveedor}>
            <i className="fa-solid fa-floppy-disk"></i> {creandoProveedor ? "Guardando..." : "Guardar Proveedor"}
          </button>
        </div>
      </div>
    </div>
  );
}