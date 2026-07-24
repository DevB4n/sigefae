import { useState } from "react";
import "./dashboard.css";
import "@fortawesome/fontawesome-free/css/all.min.css";
import logo from "../../assets/login/logo.png";

export default function ProcesosLogistica() {
  const [isSidebarOpen, setIsSidebarOpen] = useState(false);

  return (
    <>
      {/* ── HEADER ─────────────────────────────────────────────────── */}
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

      {/* ── LAYOUT ─────────────────────────────────────────────────── */}
      <div className="main-container">
        {/* SIDEBAR */}
        <aside className={`sidebar ${isSidebarOpen ? "open" : ""}`}>
          <div className="sidebar-header">
            <div className="sidebar-title">
              <i class="fa-solid fa-gear"></i>
            </div>
          </div>

          <nav className="menu-nav"></nav>
        </aside>

        {/* ÁREA DE CONTENIDO */}
        <main className="content-area">
          <div className="content-header">
            <div className="content-title">
              <div className="content-icon">
                <i className="fa-solid fa-file"></i>
              </div>
              <div>
                <h2>Procesos administrativos</h2>
                <p>Selecciona un formato del menú lateral</p>
              </div>
            </div>
          </div>

          <div className="content-body">
            <div className="welcome-wrap">
              <div className="welcome-icon">
                <i className="fa-regular fa-user"></i>
              </div>
              <h2>Bienvenido, Usuario</h2>
              <p>
                En la barra a tu izquierda encontraras todos los procesos de SIGEFAE
              </p>
            </div>
          </div>
        </main>
      </div>
    </>
  );
}