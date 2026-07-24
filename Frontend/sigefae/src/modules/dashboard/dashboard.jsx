import "./dashboard.css";

export default function ProcesosLogistica() {
  return (
    <>
      {/* ── HEADER ─────────────────────────────────────────────────── */}
      <header className="top-header">
        <div className="header-left">
          <button className="menu-toggle" title="Menú">
            <i className="fas fa-bars" />
          </button>
          <img src="" alt="Logo" className="logo" />
          <div className="header-text">
            <h1>SIGEFAE</h1>
            <p>Sistema de Gestion de Facturas Electronicas</p>
          </div>
        </div>
      </header>

      {/* ── LAYOUT ─────────────────────────────────────────────────── */}
      <div className="main-container">
        {/* SIDEBAR */}
        <aside className="sidebar">
          <div className="sidebar-header">
            <div className="sidebar-title">
              <i className="fas fa-clipboard-list" />
              
            </div>
          </div>

          <nav className="menu-nav">
          </nav>
        </aside>

        {/* ÁREA DE CONTENIDO */}
        <main className="content-area">
          <div className="content-header">
            <div className="content-title">
              <div className="content-icon">
                <i className="fas fa-truck" />
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
                <i className="fas fa-boxes" />
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