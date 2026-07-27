import { useState, useEffect } from "react";
import "./dashboard.css";
import "@fortawesome/fontawesome-free/css/all.min.css";
import logo from "../../assets/login/logo.png";
import { obtenerToken } from "../auth/token.js";

const API = "http://localhost:8080/api";

export default function ProcesosLogistica() {
  const [isSidebarOpen, setIsSidebarOpen] = useState(false);
  const [activeTab, setActiveTab] = useState("welcome");


    // ── Catálogos Admin ──
  const [catalogoActivo, setCatalogoActivo] = useState("tipo-radicacion");
  const [catalogoItems, setCatalogoItems] = useState([]);
  const [catalogoLoading, setCatalogoLoading] = useState(false);
  const [showCatalogoForm, setShowCatalogoForm] = useState(false);
  const [catalogoEditing, setCatalogoEditing] = useState(null);
  const [catalogoForm, setCatalogoForm] = useState({ nombre: "", tipo_pago_id: "" });
  const [tiposPagoCatalogo, setTiposPagoCatalogo] = useState([]);

  // ── Correos ──
  const [correos, setCorreos] = useState([]);
  const [selectedCorreoId, setSelectedCorreoId] = useState(null);
  const [correoDetail, setCorreoDetail] = useState(null);
  const [loading, setLoading] = useState(false);

  // ── Documentos Pendientes ──
  const [documentos, setDocumentos] = useState([]);
  const [selectedDocId, setSelectedDocId] = useState(null);
  const [docDetail, setDocDetail] = useState(null);
  const [loadingDocs, setLoadingDocs] = useState(false);

  // ── Edición de Documento ──
  const [isEditing, setIsEditing] = useState(false);
  const [editForm, setEditForm] = useState({});
  const [saving, setSaving] = useState(false);

  // ── Catálogos ──
  const [areas, setAreas] = useState([]);
  const [tiposRadicacion, setTiposRadicacion] = useState([]);
  const [rutas, setRutas] = useState([]);
  const [metodosPago, setMetodosPago] = useState([]);

  // ── Modal de Radicación ──
  const [showRadicarModal, setShowRadicarModal] = useState(false);
  const [radicarDocId, setRadicarDocId] = useState(null);
  const [radicarForm, setRadicarForm] = useState({
    tipo_radicacion_id: "",
    ruta_id: "",
    metodo_pago_id: "",
    numero_radicado: "",
  });
  const [radicando, setRadicando] = useState(false);

  // ── Radicados ──
  const [radicados, setRadicados] = useState([]);
  const [selectedRadicadoId, setSelectedRadicadoId] = useState(null);
  const [radicadoDetail, setRadicadoDetail] = useState(null);
  const [loadingRadicados, setLoadingRadicados] = useState(false);
  
  const catalogoConfig = {
    "tipo-radicacion": { endpoint: "tipo-radicacion", label: "Tipo de Radicación", hasTipoPago: false, updateMethod: "PUT" },
    "tipos-pago":      { endpoint: "tipos-pago",      label: "Tipo de Pago",      hasTipoPago: false, updateMethod: "PATCH" },
    "metodos-pago":    { endpoint: "metodos-pago",      label: "Método de Pago",    hasTipoPago: true,  updateMethod: "PATCH" },
  };

    const loadCatalogo = async (tipo) => {
    const cfg = catalogoConfig[tipo];
    setCatalogoLoading(true);
    try {
      const res = await fetch(`${API}/${cfg.endpoint}`, {
        headers: { Authorization: `Bearer ${obtenerToken()}` },
      });
      const data = await res.json();
      setCatalogoItems(Array.isArray(data) ? data : []);
    } catch (err) {
      console.error(err);
      setCatalogoItems([]);
    } finally {
      setCatalogoLoading(false);
    }
  };

  const openCatalogoCreate = () => {
    setCatalogoEditing(null);
    setCatalogoForm({ nombre: "", tipo_pago_id: "" });
    setShowCatalogoForm(true);
  };

  const openCatalogoEdit = (item) => {
    setCatalogoEditing(item);
    setCatalogoForm({
      nombre: item.nombre || "",
      tipo_pago_id: item.tipo_pago_id ? String(item.tipo_pago_id) : "",
    });
    setShowCatalogoForm(true);
  };

  const handleCatalogoFormChange = (e) => {
    const { name, value } = e.target;
    setCatalogoForm(prev => ({ ...prev, [name]: value }));
  };

  const handleCatalogoSubmit = async () => {
    if (!catalogoForm.nombre.trim()) {
      alert("El nombre es obligatorio");
      return;
    }

    const cfg = catalogoConfig[catalogoActivo];
    const isEdit = !!catalogoEditing;
    const url = isEdit ? `${API}/${cfg.endpoint}/${catalogoEditing.id}` : `${API}/${cfg.endpoint}`;
    const method = isEdit ? cfg.updateMethod : "POST";

    const body = { nombre: catalogoForm.nombre.trim() };
    if (cfg.hasTipoPago) {
      if (!catalogoForm.tipo_pago_id) {
        alert("Debe seleccionar un tipo de pago");
        return;
      }
      body.tipo_pago_id = parseInt(catalogoForm.tipo_pago_id);
    }

    try {
      const res = await fetch(url, {
        method,
        headers: { "Content-Type": "application/json", Authorization: `Bearer ${obtenerToken()}` },
        body: JSON.stringify(body),
      });
      if (!res.ok) {
        const errData = await res.json();
        throw new Error(errData.error || "Error guardando");
      }
      setShowCatalogoForm(false);
      loadCatalogo(catalogoActivo);
    } catch (err) {
      alert("Error: " + err.message);
    }
  };

  const handleToggleCatalogoStatus = async (item) => {
    const cfg = catalogoConfig[catalogoActivo];
    try {
      const res = await fetch(`${API}/${cfg.endpoint}/${item.id}/activo`, {
        method: "PATCH",
        headers: { "Content-Type": "application/json", Authorization: `Bearer ${obtenerToken()}` },
        body: JSON.stringify({ activo: !item.activo }),
      });
      if (!res.ok) throw new Error("Error cambiando estado");
      loadCatalogo(catalogoActivo);
    } catch (err) {
      alert(err.message);
    }
  };

  // Cargar lista de correos
  useEffect(() => {
    if (activeTab === "correos") {
      setLoading(true);
      fetch(`${API}/correo`, { headers: { Authorization: `Bearer ${obtenerToken()}` } })
        .then((res) => res.json())
        .then((data) => { if (Array.isArray(data)) setCorreos(data); })
        .catch((err) => console.error(err))
        .finally(() => setLoading(false));
    }
  }, [activeTab]);


    // Cargar catálogo activo
  useEffect(() => {
    if (activeTab === "catalogos") {
      loadCatalogo(catalogoActivo);
    }
  }, [activeTab, catalogoActivo]);

  // Cargar tipos de pago cuando se necesiten (para métodos de pago)
  useEffect(() => {
    if (activeTab === "catalogos" && catalogoActivo === "metodos-pago") {
      fetch(`${API}/tipos-pago`, { headers: { Authorization: `Bearer ${obtenerToken()}` } })
        .then(r => r.json())
        .then(data => setTiposPagoCatalogo(Array.isArray(data) ? data : []))
        .catch(err => console.error(err));
    }
  }, [activeTab, catalogoActivo]);

  // Cargar detalle del correo
  useEffect(() => {
    if (selectedCorreoId) {
      setCorreoDetail(null);
      fetch(`${API}/correo/${selectedCorreoId}`, { headers: { Authorization: `Bearer ${obtenerToken()}` } })
        .then((res) => res.json())
        .then((data) => { if (data && data.id) setCorreoDetail(data); })
        .catch((err) => console.error(err));
    }
  }, [selectedCorreoId]);

  // Cargar documentos pendientes
  useEffect(() => {
    if (activeTab === "documentos") {
      setLoadingDocs(true);
      fetch(`${API}/documentocomercial/pendientes`, { headers: { Authorization: `Bearer ${obtenerToken()}` } })
        .then((res) => res.json())
        .then((data) => { if (Array.isArray(data)) setDocumentos(data); })
        .catch((err) => console.error(err))
        .finally(() => setLoadingDocs(false));
    }
  }, [activeTab]);

  // Cargar detalle del documento
  useEffect(() => {
    if (selectedDocId) {
      setDocDetail(null);
      setIsEditing(false);
      fetch(`${API}/documentocomercial/${selectedDocId}`, { headers: { Authorization: `Bearer ${obtenerToken()}` } })
        .then((res) => res.json())
        .then((data) => {
          if (data && data.id) {
            setDocDetail(data);
            setEditForm({
              orden_compra: data.orden_compra || "",
              id_area: data.id_area || "",
              asunto: data.asunto || "",
              fecha_vencimiento: data.fecha_vencimiento ? data.fecha_vencimiento.split("T")[0] : "",
              orientacion_sello_recibido: data.orientacion_sello_recibido || "",
              numero_folios: data.numero_folios || 0,
            });
          }
        })
        .catch((err) => console.error(err));
    }
  }, [selectedDocId]);

  // Cargar catálogos cuando entra en modo edición o radicación
  useEffect(() => {
    if (isEditing || showRadicarModal) {
      const headers = { Authorization: `Bearer ${obtenerToken()}` };
      Promise.all([
        fetch(`${API}/areas`, { headers }).then(r => r.json()),
        fetch(`${API}/tipo-radicacion`, { headers }).then(r => r.json()),
        fetch(`${API}/rutas`, { headers }).then(r => r.json()),
        fetch(`${API}/metodos-pago`, { headers }).then(r => r.json()),
      ])
        .then(([a, tr, r, mp]) => {
          setAreas(Array.isArray(a) ? a : []);
          setTiposRadicacion(Array.isArray(tr) ? tr : []);
          setRutas(Array.isArray(r) ? r : []);
          setMetodosPago(Array.isArray(mp) ? mp : []);
        })
        .catch(err => console.error("Error cargando catálogos:", err));
    }
  }, [isEditing, showRadicarModal]);

  // Cargar radicados
  useEffect(() => {
    if (activeTab === "radicados") {
      setLoadingRadicados(true);
      fetch(`${API}/documentoradicado`, { headers: { Authorization: `Bearer ${obtenerToken()}` } })
        .then((res) => res.json())
        .then((data) => { if (Array.isArray(data)) setRadicados(data); })
        .catch((err) => console.error(err))
        .finally(() => setLoadingRadicados(false));
    }
  }, [activeTab]);

  // Cargar detalle de radicado
  useEffect(() => {
    if (selectedRadicadoId) {
      setRadicadoDetail(null);
      fetch(`${API}/documentoradicado/${selectedRadicadoId}`, { headers: { Authorization: `Bearer ${obtenerToken()}` } })
        .then((res) => res.json())
        .then((data) => { if (data && data.id) setRadicadoDetail(data); })
        .catch((err) => console.error(err));
    }
  }, [selectedRadicadoId]);

  const handleVerArchivo = (filename) => {
    if (!correoDetail) return;
    const url = `${API}/storage/mails/${correoDetail.id_mensaje}/${filename}`;
    window.open(url, "_blank");
  };

  const formatCurrency = (val) => {
    return new Intl.NumberFormat("es-CO", { style: "currency", currency: "COP", minimumFractionDigits: 0 }).format(val || 0);
  };

  const handleInputChange = (e) => {
    const { name, value, type } = e.target;
    setEditForm(prev => ({
      ...prev,
      [name]: type === "number" ? (value === "" ? 0 : parseFloat(value)) : value
    }));
  };

  const handleSave = async () => {
    if (!selectedDocId) return;
    setSaving(true);

    const payload = {
      orden_compra: editForm.orden_compra || "",
      id_area: parseInt(editForm.id_area) || 0,
      asunto: editForm.asunto || "",
      fecha_vencimiento: editForm.fecha_vencimiento ? new Date(editForm.fecha_vencimiento).toISOString() : null,
      orientacion_sello_recibido: editForm.orientacion_sello_recibido || "",
      numero_folios: parseInt(editForm.numero_folios) || 0,
    };

    try {
      const res = await fetch(`${API}/documentocomercial/${selectedDocId}`, {
        method: "PUT",
        headers: { "Content-Type": "application/json", Authorization: `Bearer ${obtenerToken()}` },
        body: JSON.stringify(payload),
      });
      if (!res.ok) {
        const errData = await res.json();
        throw new Error(errData.error || "Error guardando");
      }
      const updated = await res.json();
      setDocDetail(updated);
      setIsEditing(false);
      setDocumentos(prev => prev.map(d => d.id === updated.id ? { ...d, total: updated.total, numero_documento: updated.numero_documento } : d));
    } catch (err) {
      alert("Error al guardar: " + err.message);
    } finally {
      setSaving(false);
    }
  };

  // ── Modal de Radicación ──
  const openRadicarModal = (docId) => {
    setRadicarDocId(docId);
    setRadicarForm({ tipo_radicacion_id: "", ruta_id: "", metodo_pago_id: "", numero_radicado: "" });
    setShowRadicarModal(true);
  };

  const closeRadicarModal = () => {
    setShowRadicarModal(false);
    setRadicarDocId(null);
  };

  const handleRadicarChange = (e) => {
    const { name, value } = e.target;
    setRadicarForm(prev => ({ ...prev, [name]: value }));
  };

  const handleRadicarSubmit = async () => {
    if (!radicarDocId) return;

    if (!radicarForm.tipo_radicacion_id || !radicarForm.ruta_id || !radicarForm.metodo_pago_id) {
      alert("Debe seleccionar tipo de radicación, ruta y método de pago.");
      return;
    }

    setRadicando(true);
    const payload = {
      documento_comercial_id: radicarDocId,
      tipo_radicacion_id: parseInt(radicarForm.tipo_radicacion_id),
      ruta_id: parseInt(radicarForm.ruta_id),
      metodo_pago_id: parseInt(radicarForm.metodo_pago_id),
      numero_radicado: radicarForm.numero_radicado || undefined,
    };

    try {
      const res = await fetch(`${API}/documentoradicado`, {
        method: "POST",
        headers: { "Content-Type": "application/json", Authorization: `Bearer ${obtenerToken()}` },
        body: JSON.stringify(payload),
      });
      if (!res.ok) {
        const errData = await res.json();
        throw new Error(errData.error || "Error al radicar");
      }

      closeRadicarModal();
      // Refrescar pendientes (el documento ya no debe aparecer)
      setSelectedDocId(null);
      setDocDetail(null);
      setDocumentos(prev => prev.filter(d => d.id !== radicarDocId));
      // Cambiar a radicados
      setActiveTab("radicados");
    } catch (err) {
      alert("Error al radicar: " + err.message);
    } finally {
      setRadicando(false);
    }
  };

  const getTabInfo = () => {
    switch (activeTab) {
      case "correos": return { icon: "fa-solid fa-envelope", title: "Recepción de Correos", subtitle: "Gestiona las facturas electrónicas recibidas" };
      case "documentos": return { icon: "fa-solid fa-file-invoice", title: "Documentos Pendientes", subtitle: "Revisa, completa y aprueba documentos para radicación" };
      case "radicados": return { icon: "fa-solid fa-stamp", title: "Documentos Radicados", subtitle: "Consulta el estado de los documentos radicados" };
      case "catalogos": return { icon: "fa-solid fa-sliders", title: "Catálogos del Sistema", subtitle: "Gestiona tipos de radicación, pagos y métodos" };
      default: return { icon: "fa-solid fa-house", title: "Procesos administrativos", subtitle: "Selecciona un formato del menú lateral" };
    }
  };

  const tabInfo = getTabInfo();

  const renderField = (label, name, type = "text", options = null) => {
    if (!isEditing) {
      let displayValue = editForm[name];
      if (type === "select" && options) {
        const opt = options.find(o => o.id == displayValue);
        displayValue = opt ? opt.nombre : (displayValue || "—");
      }
      if (type === "date" && displayValue) {
        displayValue = new Date(displayValue).toLocaleDateString();
      }
      return (
        <div className="doc-field" key={name}>
          <label>{label}</label>
          <span>{displayValue || "—"}</span>
        </div>
      );
    }

    if (type === "select" && options) {
      return (
        <div className="doc-field" key={name}>
          <label>{label}</label>
          <select name={name} value={editForm[name] || ""} onChange={handleInputChange} className="doc-input">
            <option value="">Seleccione...</option>
            {options.map(opt => (
              <option key={opt.id} value={opt.id}>{opt.nombre}</option>
            ))}
          </select>
        </div>
      );
    }

    return (
      <div className="doc-field" key={name}>
        <label>{label}</label>
        <input type={type} name={name} value={editForm[name] || ""} onChange={handleInputChange} className="doc-input" />
      </div>
    );
  };

  const renderWelcome = () => (
    <div className="content-body">
      <div className="welcome-wrap">
        <div className="welcome-icon"><i className="fa-regular fa-user"></i></div>
        <h2>Bienvenido, Usuario</h2>
        <p>En la barra a tu izquierda encontraras todos los procesos de SIGEFAE.</p>
      </div>
    </div>
  );

  const renderCorreos = () => (
    <div className="correos-container">
      <div className="correos-list">
        <h3>Bandeja de Entrada</h3>
        {loading ? <p>Cargando correos...</p> : correos.length === 0 ? <p>No hay correos registrados.</p> : (
          correos.map((c) => (
            <div key={c.id} className={`correo-item ${selectedCorreoId === c.id ? "active" : ""}`} onClick={() => setSelectedCorreoId(c.id)}>
              <div className="correo-item-header"><strong>{c.de}</strong><span className="correo-date">{new Date(c.fecha_recepcion).toLocaleDateString()}</span></div>
              <div className="correo-item-subject">{c.asunto}</div>
              <div className="correo-item-status"><span className={`status-badge ${c.estado_correo?.nombre.toLowerCase().replace(" ", "-")}`}>{c.estado_correo?.nombre || "Sin Estado"}</span></div>
            </div>
          ))
        )}
      </div>
      <div className="correo-detail">
        {!selectedCorreoId ? (
          <div className="correo-empty"><i className="fa-regular fa-envelope-open"></i><p>Selecciona un correo para leerlo</p></div>
        ) : !correoDetail ? <p style={{ padding: "20px" }}>Cargando detalle...</p> : (
          <div className="correo-detail-content">
            <div className="correo-header"><h2>{correoDetail.asunto}</h2><div className="correo-meta"><p><strong>De:</strong> {correoDetail.de}</p><p><strong>Para:</strong> {correoDetail.para}</p><p><strong>Fecha:</strong> {new Date(correoDetail.fecha_recepcion).toLocaleString()}</p></div></div>
            <div className="correo-body"><iframe title="Cuerpo del correo" srcDoc={correoDetail.cuerpo} className="correo-iframe" /></div>
            <div className="correo-attachments">
              <h4>Archivos Adjuntos ({correoDetail.archivos ? correoDetail.archivos.length : 0})</h4>
              <div className="attachments-grid">
                {correoDetail.archivos && correoDetail.archivos.map((file) => {
                  const ext = file.split('.').pop().toLowerCase();
                  let icon = "fa-file", btnClass = "btn-default", actionText = "Ver / Descargar";
                  if (ext === "pdf") { icon = "fa-file-pdf"; btnClass = "btn-pdf"; actionText = "Ver PDF"; }
                  if (ext === "xml") { icon = "fa-file-code"; btnClass = "btn-xml"; actionText = "Ver XML"; }
                  if (ext === "zip") { icon = "fa-file-zipper"; btnClass = "btn-zip"; actionText = "Descargar ZIP"; }
                  if (ext === "eml") { icon = "fa-envelope-open-text"; btnClass = "btn-eml"; actionText = "Ver Original"; }
                  return (
                    <div key={file} className="attachment-card">
                      <i className={`fa-solid ${icon}`}></i>
                      <span className="attachment-name" title={file}>{file}</span>
                      <button className={`attachment-btn ${btnClass}`} onClick={() => handleVerArchivo(file)}>{actionText}</button>
                    </div>
                  );
                })}
              </div>
            </div>
          </div>
        )}
      </div>
    </div>
  );

  const renderDocumentos = () => (
    <div className="correos-container">
      <div className="correos-list">
        <h3>Pendientes de Revisión ({documentos.length})</h3>
        {loadingDocs ? <p style={{ padding: "15px" }}>Cargando documentos...</p> : documentos.length === 0 ? (
          <p style={{ padding: "15px", color: "#6b7280" }}>No hay documentos pendientes.</p>
        ) : (
          documentos.map((doc) => (
            <div key={doc.id} className={`correo-item ${selectedDocId === doc.id ? "active" : ""}`} onClick={() => setSelectedDocId(doc.id)}>
              <div className="correo-item-header"><strong>{doc.proveedor?.razon_social || "Sin proveedor"}</strong><span className="correo-date">{new Date(doc.fecha_documento).toLocaleDateString()}</span></div>
              <div className="correo-item-subject">{doc.tipo} - {doc.numero_documento}</div>
              <div className="correo-item-status"><span className="status-badge doc-pendiente">{formatCurrency(doc.total)}</span></div>
            </div>
          ))
        )}
      </div>

      <div className="correo-detail">
        {!selectedDocId ? (
          <div className="correo-empty"><i className="fa-solid fa-file-invoice"></i><p>Selecciona un documento para revisarlo</p></div>
        ) : !docDetail ? <p style={{ padding: "20px" }}>Cargando detalle...</p> : (
          <div className="doc-detail-content">
            <div className="doc-header">
              <div className="doc-header-top"><h2>{docDetail.tipo} #{docDetail.numero_documento}</h2><span className="doc-total">{formatCurrency(docDetail.total)}</span></div>
              {docDetail.cufe && <p className="doc-cufe"><strong>CUFE:</strong> {docDetail.cufe}</p>}
            </div>

            <div className="doc-body">
              <div className="doc-section">
                <h4><i className="fa-solid fa-circle-info"></i> Información General</h4>
                <div className="doc-grid">
                  <div className="doc-field"><label>Tipo</label><span>{docDetail.tipo}</span></div>
                  <div className="doc-field"><label>Número Documento</label><span>{docDetail.numero_documento}</span></div>
                  <div className="doc-field"><label>Fecha Emisión</label><span>{new Date(docDetail.fecha_documento).toLocaleDateString()}</span></div>
                  <div className="doc-field">
                    <label>Fecha Vencimiento</label>
                    {isEditing ? <input type="date" name="fecha_vencimiento" value={editForm.fecha_vencimiento || ""} onChange={handleInputChange} className="doc-input" /> : <span>{docDetail.fecha_vencimiento ? new Date(docDetail.fecha_vencimiento).toLocaleDateString() : "—"}</span>}
                  </div>
                  <div className="doc-field">
                    <label>Área</label>
                    {isEditing ? <select name="id_area" value={editForm.id_area || ""} onChange={handleInputChange} className="doc-input"><option value="">Seleccione...</option>{areas.map(a => <option key={a.id} value={a.id}>{a.nombre}</option>)}</select> : <span>{docDetail.area?.nombre || "—"}</span>}
                  </div>
                  <div className="doc-field"><label>Moneda</label><span>{docDetail.moneda?.nombre || "—"}</span></div>
                  <div className="doc-field">
                    <label>Orden de Compra</label>
                    {isEditing ? <input type="text" name="orden_compra" value={editForm.orden_compra || ""} onChange={handleInputChange} className="doc-input" /> : <span>{docDetail.orden_compra || "—"}</span>}
                  </div>
                  <div className="doc-field">
                    <label>Asunto</label>
                    {isEditing ? <input type="text" name="asunto" value={editForm.asunto || ""} onChange={handleInputChange} className="doc-input" /> : <span>{docDetail.asunto || "—"}</span>}
                  </div>
                  <div className="doc-field">
                    <label>Orientación Sello</label>
                    {isEditing ? <select name="orientacion_sello_recibido" value={editForm.orientacion_sello_recibido || ""} onChange={handleInputChange} className="doc-input"><option value="">Seleccione...</option><option value="HORIZONTAL">HORIZONTAL</option><option value="VERTICAL">VERTICAL</option></select> : <span>{docDetail.orientacion_sello_recibido || "No definida"}</span>}
                  </div>
                  <div className="doc-field">
                    <label>Número Folios</label>
                    {isEditing ? <input type="number" name="numero_folios" value={editForm.numero_folios || 0} onChange={handleInputChange} className="doc-input" min="0" /> : <span>{docDetail.numero_folios || "—"}</span>}
                  </div>
                </div>
              </div>

              <div className="doc-section">
                <h4><i className="fa-solid fa-building"></i> Proveedor</h4>
                <div className="doc-grid">
                  <div className="doc-field"><label>Razón Social</label><span>{docDetail.proveedor?.razon_social || "—"}</span></div>
                  <div className="doc-field"><label>NIT</label><span>{docDetail.proveedor?.numero_documento || "—"}</span></div>
                </div>
              </div>

              <div className="doc-section">
                <h4><i className="fa-solid fa-user-tie"></i> Receptor</h4>
                <div className="doc-grid">
                  <div className="doc-field"><label>Nombre</label><span>{docDetail.receptor?.nombre || "—"}</span></div>
                  <div className="doc-field"><label>NIT</label><span>{docDetail.receptor?.numero_documento || "—"}</span></div>
                </div>
              </div>

              <div className="doc-section">
                <h4><i className="fa-solid fa-calculator"></i> Resumen Financiero</h4>
                <div className="doc-totals">
                  <div className="doc-total-row"><span>Subtotal</span><span>{formatCurrency(docDetail.subtotal)}</span></div>
                  <div className="doc-total-row"><span>IVA</span><span>{formatCurrency(docDetail.iva)}</span></div>
                  <div className="doc-total-row total-final"><span>Total</span><span>{formatCurrency(docDetail.total)}</span></div>
                </div>
              </div>

              {docDetail.detalles && docDetail.detalles.length > 0 && (
                <div className="doc-section">
                  <h4><i className="fa-solid fa-list"></i> Detalle de Ítems ({docDetail.detalles.length})</h4>
                  <table className="doc-items-table">
                    <thead><tr><th>Descripción</th><th>Cantidad</th><th>Valor Unit.</th><th>Total</th></tr></thead>
                    <tbody>{docDetail.detalles.map((item) => (
                      <tr key={item.id}><td>{item.descripcion}</td><td>{item.cantidad}</td><td>{formatCurrency(item.valor_unitario)}</td><td>{formatCurrency(item.total)}</td></tr>
                    ))}</tbody>
                  </table>
                </div>
              )}

              {docDetail.correo && (
                <div className="doc-section">
                  <h4><i className="fa-solid fa-envelope"></i> Correo de Origen</h4>
                  <p style={{ fontSize: "0.85em", color: "#6b7280" }}><strong>Asunto:</strong> {docDetail.correo.asunto}</p>
                </div>
              )}
            </div>

            <div className="doc-actions">
              {isEditing ? (
                <>
                  <button className="doc-btn doc-btn-secondary" onClick={() => setIsEditing(false)} disabled={saving}><i className="fa-solid fa-xmark"></i> Cancelar</button>
                  <button className="doc-btn doc-btn-primary" onClick={handleSave} disabled={saving}><i className="fa-solid fa-floppy-disk"></i> {saving ? "Guardando..." : "Guardar Cambios"}</button>
                </>
              ) : (
                <>
                  <button className="doc-btn doc-btn-secondary" onClick={() => setIsEditing(true)}><i className="fa-solid fa-pen"></i> Completar Campos</button>
                  <button className="doc-btn doc-btn-primary" onClick={() => openRadicarModal(docDetail.id)}><i className="fa-solid fa-stamp"></i> Aprobar para Radicación</button>
                </>
              )}
            </div>
          </div>
        )}
      </div>

      {/* ── Modal de Radicación ── */}
      {showRadicarModal && (
        <div className="modal-overlay" onClick={closeRadicarModal}>
          <div className="modal-content" onClick={(e) => e.stopPropagation()}>
            <div className="modal-header">
              <h3><i className="fa-solid fa-stamp"></i> Radicar Documento</h3>
              <button className="modal-close" onClick={closeRadicarModal}><i className="fa-solid fa-xmark"></i></button>
            </div>
            <div className="modal-body">
              <div className="modal-field">
                <label>Tipo de Radicación <span className="required">*</span></label>
                <select name="tipo_radicacion_id" value={radicarForm.tipo_radicacion_id} onChange={handleRadicarChange} className="doc-input">
                  <option value="">Seleccione...</option>
                  {tiposRadicacion.map(tr => <option key={tr.id} value={tr.id}>{tr.nombre}</option>)}
                </select>
              </div>
              <div className="modal-field">
                <label>Ruta <span className="required">*</span></label>
                <select name="ruta_id" value={radicarForm.ruta_id} onChange={handleRadicarChange} className="doc-input">
                  <option value="">Seleccione...</option>
                  {rutas.map(r => <option key={r.id} value={r.id}>{r.nombre}</option>)}
                </select>
              </div>
              <div className="modal-field">
                <label>Método de Pago <span className="required">*</span></label>
                <select name="metodo_pago_id" value={radicarForm.metodo_pago_id} onChange={handleRadicarChange} className="doc-input">
                  <option value="">Seleccione...</option>
                  {metodosPago.map(mp => <option key={mp.id} value={mp.id}>{mp.nombre}</option>)}
                </select>
              </div>
              <div className="modal-field">
                <label>Número de Radicado <small>(opcional, se autogenera si está vacío)</small></label>
                <input type="text" name="numero_radicado" value={radicarForm.numero_radicado} onChange={handleRadicarChange} className="doc-input" placeholder="Ej: RAD-2026-00001" />
              </div>
            </div>
            <div className="modal-footer">
              <button className="doc-btn doc-btn-secondary" onClick={closeRadicarModal} disabled={radicando}>Cancelar</button>
              <button className="doc-btn doc-btn-primary" onClick={handleRadicarSubmit} disabled={radicando}>
                <i className="fa-solid fa-stamp"></i> {radicando ? "Radicando..." : "Confirmar Radicación"}
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );

  const renderRadicados = () => (
    <div className="correos-container">
      <div className="correos-list">
        <h3>Documentos Radicados ({radicados.length})</h3>
        {loadingRadicados ? <p style={{ padding: "15px" }}>Cargando...</p> : radicados.length === 0 ? (
          <p style={{ padding: "15px", color: "#6b7280" }}>No hay documentos radicados.</p>
        ) : (
          radicados.map((rad) => (
            <div key={rad.id} className={`correo-item ${selectedRadicadoId === rad.id ? "active" : ""}`} onClick={() => setSelectedRadicadoId(rad.id)}>
              <div className="correo-item-header">
                <strong>{rad.documento_comercial?.numero_documento || "—"}</strong>
                <span className="correo-date">{new Date(rad.fecha_radicacion).toLocaleDateString()}</span>
              </div>
              <div className="correo-item-subject">{rad.documento_comercial?.tipo || "—"} — {rad.numero_radicado}</div>
              <div className="correo-item-status"><span className="status-badge radicado">{rad.estado?.nombre || "Radicado"}</span></div>
            </div>
          ))
        )}
      </div>

      <div className="correo-detail">
        {!selectedRadicadoId ? (
          <div className="correo-empty"><i className="fa-solid fa-stamp"></i><p>Selecciona un radicado para ver su detalle</p></div>
        ) : !radicadoDetail ? <p style={{ padding: "20px" }}>Cargando detalle...</p> : (
          <div className="doc-detail-content">
            <div className="doc-header">
              <div className="doc-header-top">
                <h2>Radicado #{radicadoDetail.numero_radicado}</h2>
                <span className="doc-total">{radicadoDetail.estado?.nombre || "—"}</span>
              </div>
              <p className="doc-cufe"><strong>Documento:</strong> {radicadoDetail.documento_comercial?.tipo} {radicadoDetail.documento_comercial?.numero_documento}</p>
            </div>

            <div className="doc-body">
              <div className="doc-section">
                <h4><i className="fa-solid fa-circle-info"></i> Información del Radicado</h4>
                <div className="doc-grid">
                  <div className="doc-field"><label>Número Radicado</label><span>{radicadoDetail.numero_radicado}</span></div>
                  <div className="doc-field"><label>Fecha Radicación</label><span>{new Date(radicadoDetail.fecha_radicacion).toLocaleString()}</span></div>
                  <div className="doc-field"><label>Tipo Radicación</label><span>{radicadoDetail.tipo_radicacion?.nombre || "—"}</span></div>
                  <div className="doc-field"><label>Ruta</label><span>{radicadoDetail.ruta?.nombre || "—"}</span></div>
                  <div className="doc-field"><label>Método de Pago</label><span>{radicadoDetail.metodo_pago?.nombre || "—"}</span></div>
                  <div className="doc-field"><label>Estado Posesión</label><span>{radicadoDetail.estado_posesion || "—"}</span></div>
                  <div className="doc-field"><label>Paso Actual</label><span>{radicadoDetail.paso_actual?.nombre || "Inicio"}</span></div>
                  <div className="doc-field"><label>Responsable</label><span>{radicadoDetail.usuario_actual?.nombre || "—"}</span></div>
                </div>
              </div>

              {radicadoDetail.documento_comercial && (
                <>
                  <div className="doc-section">
                    <h4><i className="fa-solid fa-file-invoice"></i> Documento Comercial</h4>
                    <div className="doc-grid">
                      <div className="doc-field"><label>Tipo</label><span>{radicadoDetail.documento_comercial.tipo}</span></div>
                      <div className="doc-field"><label>Número</label><span>{radicadoDetail.documento_comercial.numero_documento}</span></div>
                      <div className="doc-field"><label>Proveedor</label><span>{radicadoDetail.documento_comercial.proveedor?.razon_social || "—"}</span></div>
                      <div className="doc-field"><label>Receptor</label><span>{radicadoDetail.documento_comercial.receptor?.nombre || "—"}</span></div>
                    </div>
                  </div>

                  <div className="doc-section">
                    <h4><i className="fa-solid fa-calculator"></i> Valores</h4>
                    <div className="doc-totals">
                      <div className="doc-total-row"><span>Subtotal</span><span>{formatCurrency(radicadoDetail.documento_comercial.subtotal)}</span></div>
                      <div className="doc-total-row"><span>IVA</span><span>{formatCurrency(radicadoDetail.documento_comercial.iva)}</span></div>
                      <div className="doc-total-row total-final"><span>Total</span><span>{formatCurrency(radicadoDetail.documento_comercial.total)}</span></div>
                    </div>
                  </div>
                </>
              )}

              {radicadoDetail.qr && (
                <div className="doc-section">
                  <h4><i className="fa-solid fa-qrcode"></i> Código QR</h4>
                  <p style={{ fontSize: "0.85em", color: "#6b7280", wordBreak: "break-all" }}>{radicadoDetail.qr.url}</p>
                </div>
              )}
            </div>
          </div>
        )}
      </div>
    </div>
  );
    const renderCatalogos = () => {
    const cfg = catalogoConfig[catalogoActivo];

    return (
      <div className="catalogos-container">
        <div className="catalogo-tabs">
          {Object.entries(catalogoConfig).map(([key, c]) => (
            <button
              key={key}
              className={catalogoActivo === key ? "active" : ""}
              onClick={() => setCatalogoActivo(key)}
            >
              {c.label}
            </button>
          ))}
        </div>

        <div className="catalogo-header">
          <h3>{cfg.label}</h3>
          <button className="doc-btn doc-btn-primary" onClick={openCatalogoCreate}>
            <i className="fa-solid fa-plus"></i> Nuevo {cfg.label}
          </button>
        </div>

        {showCatalogoForm && (
          <div className="catalogo-form">
            <h4>{catalogoEditing ? "Editar" : "Crear"} {cfg.label}</h4>
            <div className="catalogo-form-row">
              <div className="modal-field" style={{ flex: 1 }}>
                <label>Nombre <span className="required">*</span></label>
                <input
                  type="text"
                  name="nombre"
                  value={catalogoForm.nombre}
                  onChange={handleCatalogoFormChange}
                  placeholder={`Nombre del ${cfg.label.toLowerCase()}`}
                />
              </div>
              {cfg.hasTipoPago && (
                <div className="modal-field" style={{ flex: 1 }}>
                  <label>Tipo de Pago <span className="required">*</span></label>
                  <select name="tipo_pago_id" value={catalogoForm.tipo_pago_id} onChange={handleCatalogoFormChange}>
                    <option value="">Seleccione...</option>
                    {tiposPagoCatalogo.map(tp => (
                      <option key={tp.id} value={tp.id}>{tp.nombre}</option>
                    ))}
                  </select>
                </div>
              )}
              <div style={{ display: "flex", gap: 8 }}>
                <button className="doc-btn doc-btn-secondary" onClick={() => setShowCatalogoForm(false)}>Cancelar</button>
                <button className="doc-btn doc-btn-primary" onClick={handleCatalogoSubmit}>
                  <i className="fa-solid fa-floppy-disk"></i> Guardar
                </button>
              </div>
            </div>
          </div>
        )}

        {catalogoLoading ? (
          <p>Cargando...</p>
        ) : catalogoItems.length === 0 ? (
          <p style={{ color: "#6b7280" }}>No hay registros.</p>
        ) : (
          <table className="catalogo-table">
            <thead>
              <tr>
                <th>ID</th>
                <th>Nombre</th>
                {cfg.hasTipoPago && <th>Tipo de Pago</th>}
                <th>Estado</th>
                <th style={{ width: 120 }}>Acciones</th>
              </tr>
            </thead>
            <tbody>
              {catalogoItems.map(item => (
                <tr key={item.id}>
                  <td>{item.id}</td>
                  <td>{item.nombre}</td>
                  {cfg.hasTipoPago && <td>{item.tipo_pago || "—"}</td>}
                  <td>
                    <span className={`status-badge ${item.activo ? "radicado" : "doc-pendiente"}`}>
                      {item.activo ? "Activo" : "Inactivo"}
                    </span>
                  </td>
                  <td>
                    <div className="catalogo-actions">
                      <button className="btn-icon btn-edit" onClick={() => openCatalogoEdit(item)} title="Editar">
                        <i className="fa-solid fa-pen"></i>
                      </button>
                      <button
                        className={`btn-icon btn-toggle ${item.activo ? "active" : ""}`}
                        onClick={() => handleToggleCatalogoStatus(item)}
                        title={item.activo ? "Desactivar" : "Activar"}
                      >
                        <i className={`fa-solid ${item.activo ? "fa-check" : "fa-xmark"}`}></i>
                      </button>
                    </div>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </div>
    );
  };

  const renderContent = () => {
    switch (activeTab) {
      case "correos": return renderCorreos();
      case "documentos": return renderDocumentos();
      case "radicados": return renderRadicados();
      case "catalogos": return renderCatalogos();
      default: return renderWelcome();
    }
  };

  return (
    <>
      <header className="top-header">
        <div className="header-left">
          <button className="menu-toggle" title="Menú" onClick={() => setIsSidebarOpen(!isSidebarOpen)}><i className="fas fa-bars" /></button>
          <img src={logo} alt="Logo" className="logo" />
          <div className="header-text"><h1>SIGEFAE</h1><p>Sistema de Gestion de Facturas Electronicas</p></div>
        </div>
      </header>

      <div className="main-container">
        <aside className={`sidebar ${isSidebarOpen ? "open" : ""}`}>
          <div className="sidebar-header"><div className="sidebar-title"><i className="fa-solid fa-gear"></i></div></div>
          <nav className="menu-nav">
            <a href="#" className={`menu-item ${activeTab === "correos" ? "active" : ""}`} onClick={(e) => { e.preventDefault(); setActiveTab("correos"); setIsSidebarOpen(false); }}>
              <div className="item-icon"><i className="fa-solid fa-envelope"></i></div>
              <div className="item-text"><span className="item-nombre">Correos</span></div>
            </a>
            <a href="#" className={`menu-item ${activeTab === "documentos" ? "active" : ""}`} onClick={(e) => { e.preventDefault(); setActiveTab("documentos"); setIsSidebarOpen(false); }}>
              <div className="item-icon"><i className="fa-solid fa-file-invoice"></i></div>
              <div className="item-text"><span className="item-nombre">Documentos</span></div>
            </a>
            <a href="#" className={`menu-item ${activeTab === "radicados" ? "active" : ""}`} onClick={(e) => { e.preventDefault(); setActiveTab("radicados"); setIsSidebarOpen(false); }}>
              <div className="item-icon"><i className="fa-solid fa-stamp"></i></div>
              <div className="item-text"><span className="item-nombre">Radicados</span></div>
            </a>
                        <a 
              href="#" 
              className={`menu-item ${activeTab === "catalogos" ? "active" : ""}`}
              onClick={(e) => { e.preventDefault(); setActiveTab("catalogos"); setIsSidebarOpen(false); }}
            >
              <div className="item-icon">
                <i className="fa-solid fa-sliders"></i>
              </div>
              <div className="item-text">
                <span className="item-nombre">Catálogos</span>
              </div>
            </a>
          </nav>
        </aside>

        <main className="content-area">
          <div className="content-header">
            <div className="content-title">
              <div className="content-icon"><i className={tabInfo.icon}></i></div>
              <div><h2>{tabInfo.title}</h2><p>{tabInfo.subtitle}</p></div>
            </div>
          </div>
          {renderContent()}
        </main>
      </div>
    </>
  );
}