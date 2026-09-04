import { catalogoConfig } from "../constants/catalogoConfig";

export default function RenderCatalogos({
  catalogoActivo, setCatalogoActivo, catalogoItems, catalogoLoading,
  showCatalogoForm, setShowCatalogoForm, catalogoEditing, catalogoForm,
  tiposPagoCatalogo, areasCatalogo, rutasCatalogo, usuariosCatalogo, monedasCatalogo,
  cfg, openCatalogoCreate, openCatalogoEdit, handleCatalogoFormChange,
  handleCatalogoSubmit, handleToggleCatalogoStatus
}) {
  const renderFormField = (field) => {
    const value = catalogoForm[field] || "";
    if (field === "tipo_pago_id") return (
      <div className="modal-field" key={field}><label>Tipo de Pago <span className="required">*</span></label>
        <select name={field} value={value} onChange={handleCatalogoFormChange} className="doc-input"><option value="">Seleccione...</option>{tiposPagoCatalogo.map(tp => <option key={tp.id} value={tp.id}>{tp.nombre}</option>)}</select>
      </div>
    );
    if (field === "area_id") return (
      <div className="modal-field" key={field}><label>Área <span className="required">*</span></label>
        <select name={field} value={value} onChange={handleCatalogoFormChange} className="doc-input"><option value="">Seleccione...</option>{areasCatalogo.map(a => <option key={a.id} value={a.id}>{a.nombre}</option>)}</select>
      </div>
    );
    if (field === "ruta_id") return (
      <div className="modal-field" key={field}><label>Ruta <span className="required">*</span></label>
        <select name={field} value={value} onChange={handleCatalogoFormChange} className="doc-input"><option value="">Seleccione...</option>{rutasCatalogo.map(r => <option key={r.id} value={r.id}>{r.nombre}</option>)}</select>
      </div>
    );
    if (field === "usuario_id" || field === "usuario_aprobador_id") return (
      <div className="modal-field" key={field}><label>{field === "usuario_id" ? "Usuario Responsable" : "Usuario Aprobador"} <span className="required">*</span></label>
        <select name={field} value={value} onChange={handleCatalogoFormChange} className="doc-input"><option value="">Seleccione...</option>{usuariosCatalogo.map(u => <option key={u.id} value={u.id}>{u.nombre}</option>)}</select>
      </div>
    );
    if (field === "orden") return <div className="modal-field" key={field} style={{ maxWidth: 120 }}><label>Orden <span className="required">*</span></label><input type="number" name={field} value={value} onChange={handleCatalogoFormChange} min="1" className="doc-input" /></div>;
    if (field === "monto_minimo_smmlv") return <div className="modal-field" key={field} style={{ maxWidth: 160 }}><label>Mínimo SMMLV <span className="required">*</span></label><input type="number" step="0.01" name={field} value={value} onChange={handleCatalogoFormChange} min="0" className="doc-input" /></div>;
    if (field === "monto_maximo_smmlv") return <div className="modal-field" key={field} style={{ maxWidth: 160 }}><label>Máximo SMMLV <span className="required">*</span></label><input type="number" step="0.01" name={field} value={value} onChange={handleCatalogoFormChange} min="0" title="Dejar en 0 para sin límite" className="doc-input" /></div>;
    if (field === "ano") return <div className="modal-field" key={field} style={{ maxWidth: 120 }}><label>Año <span className="required">*</span></label><input type="number" name={field} value={value} onChange={handleCatalogoFormChange} min="2000" className="doc-input" /></div>;
    if (field === "valor") return <div className="modal-field" key={field} style={{ maxWidth: 200 }}><label>Valor ($) <span className="required">*</span></label><input type="number" name={field} value={value} onChange={handleCatalogoFormChange} min="0" className="doc-input" /></div>;

    if (field === "posicion_insercion") return (
      <div className="modal-field" key={field}><label>Posición Inserción <span className="required">*</span></label>
        <select name={field} value={value} onChange={handleCatalogoFormChange} className="doc-input"><option value="">Seleccione...</option><option value="PRIMERO">Al inicio</option><option value="ANTES_FINAL">Antes del final</option><option value="ULTIMO">Al final</option></select>
      </div>
    );
    if (field === "prioridad") return <div className="modal-field" key={field} style={{ maxWidth: 120 }}><label>Prioridad <span className="required">*</span></label><input type="number" name={field} value={value} onChange={handleCatalogoFormChange} min="0" placeholder="0" title="Menor número = se ejecuta primero dentro de su posición" className="doc-input" /></div>;
    if (field === "codigo") return <div className="modal-field" key={field} style={{ flex: 1 }}><label>Código <span className="required">*</span></label><input type="text" name={field} value={value} onChange={handleCatalogoFormChange} placeholder="Ej: BU101" className="doc-input" /></div>;
    if (field === "proyecto") return <div className="modal-field" key={field} style={{ flex: 1 }}><label>Proyecto <span className="required">*</span></label><input type="text" name={field} value={value} onChange={handleCatalogoFormChange} placeholder="Ej: PROY-001" className="doc-input" required /></div>;
    if (field === "descripcion") return <div className="modal-field" key={field} style={{ width: "100%" }}><label>Descripción</label><textarea name={field} value={value} onChange={handleCatalogoFormChange} placeholder="Descripción opcional..." className="doc-input" rows="2" /></div>;
    if (field === "sucursal") {
      const sucursales = ["BUCARAMANGA", "MALAMBO", "CUCUTA", "CB", "CIENAGA DE ORO", "GENERAL"];
      return <div className="modal-field" key={field}><label>Sucursal <span className="required">*</span></label><select name={field} value={value} onChange={handleCatalogoFormChange} className="doc-input"><option value="">Seleccione...</option>{sucursales.map(s => <option key={s} value={s}>{s}</option>)}</select></div>;
    }
    if (field === "departamento") {
      const depts = ["ADMON", "VENTAS", "PRODUCCION"];
      return <div className="modal-field" key={field}><label>Departamento <span className="required">*</span></label><select name={field} value={value} onChange={handleCatalogoFormChange} className="doc-input"><option value="">Seleccione...</option>{depts.map(d => <option key={d} value={d}>{d}</option>)}</select></div>;
    }
    if (field === "tipo") return <div className="modal-field" key={field}><label>Tipo</label><select name={field} value={value} onChange={handleCatalogoFormChange} className="doc-input"><option value="">Ninguno</option><option value="Servicio">Servicio</option><option value="Compra">Compra</option></select></div>;
    if (field === "tarifa_iva") return <div className="modal-field" key={field} style={{ maxWidth: 120 }}><label>Tarifa IVA</label><select name={field} value={value} onChange={handleCatalogoFormChange} className="doc-input"><option value="">Ninguna</option><option value="19%">19%</option><option value="5%">5%</option><option value="0%">0%</option></select></div>;
    return <div className="modal-field" key={field} style={{ flex: 1 }}><label>Nombre <span className="required">*</span></label><input type="text" name={field} value={value} onChange={handleCatalogoFormChange} placeholder={`Nombre del ${cfg.label.toLowerCase()}`} className="doc-input" /></div>;
  };

  const getColumnLabel = (field) => {
    const map = { nombre: "Nombre", codigo: "Código", proyecto: "Proyecto", descripcion: "Descripción", sucursal: "Sucursal", departamento: "Depto", tipo: "Tipo", tarifa_iva: "Tarifa IVA", tipo_pago: "Tipo de Pago", area: "Área", ruta: "Ruta", orden: "Orden", usuario: "Usuario", prioridad: "Prioridad", usuario_aprobador_id: "Aprobador", moneda_id: "Moneda", monto_minimo: "Monto Mínimo", monto_minimo_smmlv: "Mín. SMMLV", monto_maximo_smmlv: "Máx. SMMLV", posicion_insercion: "Posición", ano: "Año", valor: "Valor" };
    return map[field] || field;
  };

  const getItemDisplay = (item, field) => {
    if (field === "tipo_pago_id") return item.tipo_pago || "—";
    if (field === "area_id") return item.area || "—";
    if (field === "ruta_id") return item.ruta || "—";
    if (field === "usuario_id") return item.usuario || "—";
    if (field === "usuario_aprobador_id") return item.usuario_aprobador?.nombre || item.usuario_aprobador_id || "—";
    if (field === "moneda_id") return item.moneda ? `${item.moneda.nombre} (${item.moneda.codigo})` : item.moneda_id || "—";
    if (field === "monto_minimo_smmlv" || field === "monto_maximo_smmlv") return item[field] !== undefined && item[field] !== null ? `${item[field]} SMMLV` : "0 SMMLV";
    if (field === "valor") return item[field] !== undefined && item[field] !== null ? `$${Number(item[field]).toLocaleString()}` : "—";
    if ((field === "tipo" || field === "tarifa_iva") && !item[field]) return "—";
    return item[field] !== undefined ? item[field] : "—";
  };

  return (
    <div className="catalogos-container">
      <div className="catalogo-tabs">
        {Object.entries(catalogoConfig).map(([key, c]) => (
          <button key={key} className={catalogoActivo === key ? "active" : ""} onClick={() => setCatalogoActivo(key)}>{c.label}</button>
        ))}
      </div>
      <div className="catalogo-header">
        <h3>{cfg.label}</h3>
        <button className="doc-btn doc-btn-primary" onClick={openCatalogoCreate}><i className="fa-solid fa-plus"></i> Nuevo {cfg.label}</button>
      </div>
      {showCatalogoForm && (
        <div className="modal-overlay">
          <div className="modal-content" style={{ maxWidth: "700px" }}>
            <div className="modal-header">
              <h3><i className="fa-solid fa-sliders"></i> {catalogoEditing ? "Editar" : "Crear"} {cfg.label}</h3>
              <button className="modal-close" onClick={() => setShowCatalogoForm(false)}><i className="fa-solid fa-xmark"></i></button>
            </div>
            <div className="modal-body">
              <div className="catalogo-form-row" style={{ display: "flex", flexWrap: "wrap", gap: "16px" }}>
                {cfg.fields.map(renderFormField)}
              </div>
            </div>
            <div className="modal-footer">
              <button className="doc-btn doc-btn-secondary" onClick={() => setShowCatalogoForm(false)}>Cancelar</button>
              <button className="doc-btn doc-btn-primary" onClick={handleCatalogoSubmit}><i className="fa-solid fa-floppy-disk"></i> Guardar</button>
            </div>
          </div>
        </div>
      )}
      {catalogoLoading ? <p>Cargando...</p> : catalogoItems.length === 0 ? <p style={{ color: "#6b7280" }}>No hay registros.</p> : (
        <table className="catalogo-table">
          <thead><tr><th>ID</th>{cfg.fields.filter(f => f !== "tipo_pago_id").map(f => <th key={f}>{getColumnLabel(f === "area_id" ? "area" : f === "ruta_id" ? "ruta" : f === "usuario_id" ? "usuario" : f)}</th>)}{cfg.fields.includes("tipo_pago_id") && <th>Tipo de Pago</th>}<th>Estado</th><th style={{ width: 120 }}>Acciones</th></tr></thead>
          <tbody>
            {catalogoItems.map(item => (
              <tr key={item.id}>
                <td>{item.id}</td>
                {cfg.fields.filter(f => f !== "tipo_pago_id").map(f => <td key={f}>{getItemDisplay(item, f)}</td>)}
                {cfg.fields.includes("tipo_pago_id") && <td>{item.tipo_pago || "—"}</td>}
                <td><span className={`status-badge ${item.activo ? "radicado" : "doc-pendiente"}`}>{item.activo ? "Activo" : "Inactivo"}</span></td>
                <td>
                  <div className="catalogo-actions">
                    <button className="btn-icon btn-edit" onClick={() => openCatalogoEdit(item)} title="Editar"><i className="fa-solid fa-pen"></i></button>
                    <button className={`btn-icon btn-toggle ${item.activo ? "active" : ""}`} onClick={() => handleToggleCatalogoStatus(item)} title={item.activo ? "Desactivar" : "Activar"}><i className={`fa-solid ${item.activo ? "fa-check" : "fa-xmark"}`}></i></button>
                  </div>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </div>
  );
}