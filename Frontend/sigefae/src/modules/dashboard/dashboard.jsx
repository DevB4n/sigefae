import { useState, useEffect } from "react";
import "./dashboard.css";
import "@fortawesome/fontawesome-free/css/all.min.css";
import logo from "../../assets/login/logo.png";
import { obtenerToken } from "../auth/token.js";

export default function ProcesosLogistica() {
  const [isSidebarOpen, setIsSidebarOpen] = useState(false);
  const [activeTab, setActiveTab] = useState("welcome");
  const [correos, setCorreos] = useState([]);
  const [selectedCorreoId, setSelectedCorreoId] = useState(null);
  const [correoDetail, setCorreoDetail] = useState(null);
  const [loading, setLoading] = useState(false);

  // Cargar lista de correos
  useEffect(() => {
    if (activeTab === "correos") {
      setLoading(true);
      fetch("http://localhost:8080/api/correo", {
        headers: {
          Authorization: `Bearer ${obtenerToken()}`,
        },
      })
        .then((res) => res.json())
        .then((data) => {
          if (Array.isArray(data)) {
            setCorreos(data);
          } else {
            console.error("Error fetching correos:", data);
          }
        })
        .catch((err) => console.error(err))
        .finally(() => setLoading(false));
    }
  }, [activeTab]);

  // Cargar detalle del correo (con lista de archivos)
  useEffect(() => {
    if (selectedCorreoId) {
      setCorreoDetail(null);
      fetch(`http://localhost:8080/api/correo/${selectedCorreoId}`, {
        headers: {
          Authorization: `Bearer ${obtenerToken()}`,
        },
      })
        .then((res) => res.json())
        .then((data) => {
          if (data && data.id) {
            setCorreoDetail(data);
          }
        })
        .catch((err) => console.error(err));
    }
  }, [selectedCorreoId]);

  const handleVerArchivo = (filename) => {
    if (!correoDetail) return;
    const url = `http://localhost:8080/api/storage/mails/${correoDetail.id_mensaje}/${filename}`;
    window.open(url, "_blank");
  };

  const renderContent = () => {
    if (activeTab === "welcome") {
      return (
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
    }

    if (activeTab === "correos") {
      return (
        <div className="correos-container">
          {/* Panel Izquierdo: Lista de Correos */}
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

          {/* Panel Derecho: Detalle del Correo */}
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
          </nav>
        </aside>

        <main className="content-area">
          <div className="content-header">
            <div className="content-title">
              <div className="content-icon">
                <i className={activeTab === "welcome" ? "fa-solid fa-house" : "fa-solid fa-envelope"}></i>
              </div>
              <div>
                <h2>{activeTab === "welcome" ? "Procesos administrativos" : "Recepción de Correos"}</h2>
                <p>
                  {activeTab === "welcome" 
                    ? "Selecciona un formato del menú lateral" 
                    : "Gestiona las facturas electrónicas recibidas"}
                </p>
              </div>
            </div>
          </div>

          {renderContent()}
        </main>
      </div>
    </>
  );
}