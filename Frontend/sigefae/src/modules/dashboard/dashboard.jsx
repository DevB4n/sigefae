import { useState, useEffect } from "react";
import "./dashboard.css";
import "@fortawesome/fontawesome-free/css/all.min.css";
import logo from "../../assets/login/logo.png";
import { obtenerToken } from "../auth/token.js";

const API = "http://localhost:8080/api";

export default function ProcesosLogistica() {
  const [isSidebarOpen, setIsSidebarOpen] = useState(false);
  const [activeTab, setActiveTab] = useState("welcome");

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

  // Cargar lista de correos
  useEffect(() => {
    if (activeTab === "correos") {
      setLoading(true);
      fetch(`${API}/correo`, {
        headers: { Authorization: `Bearer ${obtenerToken()}` },
      })
        .then((res) => res.json())
        .then((data) => { if (Array.isArray(data)) setCorreos(data); })
        .catch((err) => console.error(err))
        .finally(() => setLoading(false));
    }
  }, [activeTab]);

  // Cargar detalle del correo
  useEffect(() => {
    if (selectedCorreoId) {
      setCorreoDetail(null);
      fetch(`${API}/correo/${selectedCorreoId}`, {
        headers: { Authorization: `Bearer ${obtenerToken()}` },
      })
        .then((res) => res.json())
        .then((data) => { if (data && data.id) setCorreoDetail(data); })
        .catch((err) => console.error(err));
    }
  }, [selectedCorreoId]);

  // Cargar documentos pendientes
  useEffect(() => {
    if (activeTab === "documentos") {
      setLoadingDocs(true);
      fetch(`${API}/documentocomercial/pendientes`, {
        headers: { Authorization: `Bearer ${obtenerToken()}` },
      })
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
      fetch(`${API}/documentocomercial/${selectedDocId}`, {
        headers: { Authorization: `Bearer ${obtenerToken()}` },
      })
        .then((res) => res.json())
        .then((data) => { if (data && data.id) setDocDetail(data); })
        .catch((err) => console.error(err));
    }
  }, [selectedDocId]);

  const handleVerArchivo = (filename) => {
    if (!correoDetail) return;
    const url = `${API}/storage/mails/${correoDetail.id_mensaje}/${filename}`;
    window.open(url, "_blank");
  };

  const formatCurrency = (val) => {
    return new Intl.NumberFormat("es-CO", { style: "currency", currency: "COP", minimumFractionDigits: 0 }).format(val || 0);
  };

  // ── Tab de contenido ──
  const getTabInfo = () => {
    switch (activeTab) {
      case "correos":
        return { icon: "fa-solid fa-envelope", title: "Recepción de Correos", subtitle: "Gestiona las facturas electrónicas recibidas" };
      case "documentos":
        return { icon: "fa-solid fa-file-invoice", title: "Documentos Pendientes", subtitle: "Revisa, completa y aprueba documentos para radicación" };
      default:
        return { icon: "fa-solid fa-house", title: "Procesos administrativos", subtitle: "Selecciona un formato del menú lateral" };
    }
  };

  const tabInfo = getTabInfo();

  // ══════════════════════════════════════════
  //  RENDER: Welcome
  // ══════════════════════════════════════════
  const renderWelcome = () => (
    <div className="content-body">
      <div className="welcome-wrap">
        <div className="welcome-icon">
          <i className="fa-regular fa-user"></i>
        </div>
        <h2>Bienvenido, Usuario</h2>
        <p>
          En la barra a tu izquierda encontraras todos los procesos de SIGEFAE.
        </p>
      </div>
    </div>
  );

  // ══════════════════════════════════════════
  //  RENDER: Correos
  // ══════════════════════════════════════════
  const renderCorreos = () => (
    <div className="correos-container">
      <div className="correos-list">
        <h3>Bandeja de Entrada</h3>
        {loading ? (
          <p>Cargando correos...</p>
        ) : correos.length === 0 ? (
          <p>No hay correos registrados.</p>
        ) : (
          correos.map((c) => (
            <div
              key={c.id}
              className={`correo-item ${selectedCorreoId === c.id ? "active" : ""}`}
              onClick={() => setSelectedCorreoId(c.id)}
            >
              <div className="correo-item-header">
                <strong>{c.de}</strong>
                <span className="correo-date">
                  {new Date(c.fecha_recepcion).toLocaleDateString()}
                </span>
              </div>
              <div className="correo-item-subject">{c.asunto}</div>
              <div className="correo-item-status">
                <span className={`status-badge ${c.estado_correo?.nombre.toLowerCase().replace(" ", "-")}`}>
                  {c.estado_correo?.nombre || "Sin Estado"}
                </span>
              </div>
            </div>
          ))
        )}
      </div>

      <div className="correo-detail">
        {!selectedCorreoId ? (
          <div className="correo-empty">
            <i className="fa-regular fa-envelope-open"></i>
            <p>Selecciona un correo para leerlo</p>
          </div>
        ) : !correoDetail ? (
          <p style={{ padding: "20px" }}>Cargando detalle...</p>
        ) : (
          <div className="correo-detail-content">
            <div className="correo-header">
              <h2>{correoDetail.asunto}</h2>
              <div className="correo-meta">
                <p><strong>De:</strong> {correoDetail.de}</p>
                <p><strong>Para:</strong> {correoDetail.para}</p>
                <p><strong>Fecha:</strong> {new Date(correoDetail.fecha_recepcion).toLocaleString()}</p>
              </div>
            </div>

            <div className="correo-body">
              <iframe 
                title="Cuerpo del correo"
                srcDoc={correoDetail.cuerpo}
                className="correo-iframe"
              />
            </div>

            <div className="correo-attachments">
              <h4>Archivos Adjuntos ({correoDetail.archivos ? correoDetail.archivos.length : 0})</h4>
              <div className="attachments-grid">
                {correoDetail.archivos && correoDetail.archivos.map((file) => {
                  const ext = file.split('.').pop().toLowerCase();
                  let icon = "fa-file";
                  let btnClass = "btn-default";
                  let actionText = "Ver / Descargar";

                  if (ext === "pdf") { icon = "fa-file-pdf"; btnClass = "btn-pdf"; actionText = "Ver PDF"; }
                  if (ext === "xml") { icon = "fa-file-code"; btnClass = "btn-xml"; actionText = "Ver XML"; }
                  if (ext === "zip") { icon = "fa-file-zipper"; btnClass = "btn-zip"; actionText = "Descargar ZIP"; }
                  if (ext === "eml") { icon = "fa-envelope-open-text"; btnClass = "btn-eml"; actionText = "Ver Original"; }

                  return (
                    <div key={file} className="attachment-card">
                      <i className={`fa-solid ${icon}`}></i>
                      <span className="attachment-name" title={file}>{file}</span>
                      <button className={`attachment-btn ${btnClass}`} onClick={() => handleVerArchivo(file)}>
                        {actionText}
                      </button>
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

  // ══════════════════════════════════════════
  //  RENDER: Documentos Pendientes
  // ══════════════════════════════════════════
  const renderDocumentos = () => (
    <div className="correos-container">
      {/* Panel Izquierdo: Lista de Documentos */}
      <div className="correos-list">
        <h3>Pendientes de Revisión ({documentos.length})</h3>
        {loadingDocs ? (
          <p style={{ padding: "15px" }}>Cargando documentos...</p>
        ) : documentos.length === 0 ? (
          <p style={{ padding: "15px", color: "#6b7280" }}>No hay documentos pendientes.</p>
        ) : (
          documentos.map((doc) => (
            <div
              key={doc.id}
              className={`correo-item ${selectedDocId === doc.id ? "active" : ""}`}
              onClick={() => setSelectedDocId(doc.id)}
            >
              <div className="correo-item-header">
                <strong>{doc.proveedor?.razon_social || "Sin proveedor"}</strong>
                <span className="correo-date">
                  {new Date(doc.fecha_documento).toLocaleDateString()}
                </span>
              </div>
              <div className="correo-item-subject">
                {doc.tipo} - {doc.numero_documento}
              </div>
              <div className="correo-item-status">
                <span className="status-badge doc-pendiente">
                  {formatCurrency(doc.total)}
                </span>
              </div>
            </div>
          ))
        )}
      </div>

      {/* Panel Derecho: Detalle del Documento */}
      <div className="correo-detail">
        {!selectedDocId ? (
          <div className="correo-empty">
            <i className="fa-solid fa-file-invoice"></i>
            <p>Selecciona un documento para revisarlo</p>
          </div>
        ) : !docDetail ? (
          <p style={{ padding: "20px" }}>Cargando detalle...</p>
        ) : (
          <div className="doc-detail-content">
            {/* Encabezado del documento */}
            <div className="doc-header">
              <div className="doc-header-top">
                <h2>{docDetail.tipo} #{docDetail.numero_documento}</h2>
                <span className="doc-total">{formatCurrency(docDetail.total)}</span>
              </div>
              {docDetail.cufe && (
                <p className="doc-cufe"><strong>CUFE:</strong> {docDetail.cufe}</p>
              )}
            </div>

            {/* Cuerpo scrolleable */}
            <div className="doc-body">
              {/* Información general */}
              <div className="doc-section">
                <h4><i className="fa-solid fa-circle-info"></i> Información General</h4>
                <div className="doc-grid">
                  <div className="doc-field">
                    <label>Fecha Emisión</label>
                    <span>{new Date(docDetail.fecha_documento).toLocaleDateString()}</span>
                  </div>
                  <div className="doc-field">
                    <label>Fecha Vencimiento</label>
                    <span>{docDetail.fecha_vencimiento ? new Date(docDetail.fecha_vencimiento).toLocaleDateString() : "—"}</span>
                  </div>
                  <div className="doc-field">
                    <label>Moneda</label>
                    <span>{docDetail.moneda?.nombre || "—"}</span>
                  </div>
                  <div className="doc-field">
                    <label>Área</label>
                    <span>{docDetail.area?.nombre || "—"}</span>
                  </div>
                  <div className="doc-field">
                    <label>Orden de Compra</label>
                    <span>{docDetail.orden_compra || "—"}</span>
                  </div>
                  <div className="doc-field">
                    <label>Orientación Sello</label>
                    <span>{docDetail.orientacion_sello_recibido || "No definida"}</span>
                  </div>
                </div>
              </div>

              {/* Proveedor */}
              <div className="doc-section">
                <h4><i className="fa-solid fa-building"></i> Proveedor</h4>
                <div className="doc-grid">
                  <div className="doc-field">
                    <label>Razón Social</label>
                    <span>{docDetail.proveedor?.razon_social || "—"}</span>
                  </div>
                  <div className="doc-field">
                    <label>NIT</label>
                    <span>{docDetail.proveedor?.numero_documento || "—"}</span>
                  </div>
                </div>
              </div>

              {/* Receptor */}
              <div className="doc-section">
                <h4><i className="fa-solid fa-user-tie"></i> Receptor</h4>
                <div className="doc-grid">
                  <div className="doc-field">
                    <label>Nombre</label>
                    <span>{docDetail.receptor?.nombre || "—"}</span>
                  </div>
                  <div className="doc-field">
                    <label>NIT</label>
                    <span>{docDetail.receptor?.numero_documento || "—"}</span>
                  </div>
                </div>
              </div>

              {/* Totales */}
              <div className="doc-section">
                <h4><i className="fa-solid fa-calculator"></i> Resumen Financiero</h4>
                <div className="doc-totals">
                  <div className="doc-total-row">
                    <span>Subtotal</span>
                    <span>{formatCurrency(docDetail.subtotal)}</span>
                  </div>
                  <div className="doc-total-row">
                    <span>IVA</span>
                    <span>{formatCurrency(docDetail.iva)}</span>
                  </div>
                  <div className="doc-total-row total-final">
                    <span>Total</span>
                    <span>{formatCurrency(docDetail.total)}</span>
                  </div>
                </div>
              </div>

              {/* Detalle de ítems */}
              {docDetail.detalles && docDetail.detalles.length > 0 && (
                <div className="doc-section">
                  <h4><i className="fa-solid fa-list"></i> Detalle de Ítems ({docDetail.detalles.length})</h4>
                  <table className="doc-items-table">
                    <thead>
                      <tr>
                        <th>Descripción</th>
                        <th>Cantidad</th>
                        <th>Valor Unit.</th>
                        <th>Total</th>
                      </tr>
                    </thead>
                    <tbody>
                      {docDetail.detalles.map((item) => (
                        <tr key={item.id}>
                          <td>{item.descripcion}</td>
                          <td>{item.cantidad}</td>
                          <td>{formatCurrency(item.valor_unitario)}</td>
                          <td>{formatCurrency(item.total)}</td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              )}

              {/* Correo origen */}
              {docDetail.correo && (
                <div className="doc-section">
                  <h4><i className="fa-solid fa-envelope"></i> Correo de Origen</h4>
                  <p style={{ fontSize: "0.85em", color: "#6b7280" }}>
                    <strong>Asunto:</strong> {docDetail.correo.asunto}
                  </p>
                </div>
              )}
            </div>

            {/* Footer con acciones */}
            <div className="doc-actions">
              <button className="doc-btn doc-btn-secondary" title="Funcionalidad próximamente">
                <i className="fa-solid fa-pen"></i> Completar Campos
              </button>
              <button className="doc-btn doc-btn-primary" title="Funcionalidad próximamente">
                <i className="fa-solid fa-check"></i> Aprobar para Radicación
              </button>
            </div>
          </div>
        )}
      </div>
    </div>
  );

  const renderContent = () => {
    switch (activeTab) {
      case "correos": return renderCorreos();
      case "documentos": return renderDocumentos();
      default: return renderWelcome();
    }
  };

  return (
    <>
      <header className="top-header">
        <div className="header-left">
          <button
            className="menu-toggle"
            title="Menú"
            onClick={() => setIsSidebarOpen(!isSidebarOpen)}
          >
            <i className="fas fa-bars" />
          </button>
          <img src={logo} alt="Logo" className="logo" />
          <div className="header-text">
            <h1>SIGEFAE</h1>
            <p>Sistema de Gestion de Facturas Electronicas</p>
          </div>
        </div>
      </header>

      <div className="main-container">
        <aside className={`sidebar ${isSidebarOpen ? "open" : ""}`}>
          <div className="sidebar-header">
            <div className="sidebar-title">
              <i className="fa-solid fa-gear"></i>
            </div>
          </div>

          <nav className="menu-nav">
            <a 
              href="#" 
              className={`menu-item ${activeTab === "correos" ? "active" : ""}`}
              onClick={(e) => { e.preventDefault(); setActiveTab("correos"); setIsSidebarOpen(false); }}
            >
              <div className="item-icon">
                <i className="fa-solid fa-envelope"></i>
              </div>
              <div className="item-text">
                <span className="item-nombre">Correos</span>
              </div>
            </a>
            <a 
              href="#" 
              className={`menu-item ${activeTab === "documentos" ? "active" : ""}`}
              onClick={(e) => { e.preventDefault(); setActiveTab("documentos"); setIsSidebarOpen(false); }}
            >
              <div className="item-icon">
                <i className="fa-solid fa-file-invoice"></i>
              </div>
              <div className="item-text">
                <span className="item-nombre">Documentos</span>
              </div>
            </a>
          </nav>
        </aside>

        <main className="content-area">
          <div className="content-header">
            <div className="content-title">
              <div className="content-icon">
                <i className={tabInfo.icon}></i>
              </div>
              <div>
                <h2>{tabInfo.title}</h2>
                <p>{tabInfo.subtitle}</p>
              </div>
            </div>
          </div>

          {renderContent()}
        </main>
      </div>
    </>
  );
}