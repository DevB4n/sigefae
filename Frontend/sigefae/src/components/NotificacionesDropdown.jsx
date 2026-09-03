import { useState, useRef, useEffect } from "react";
import "./notificaciones.css";

// Recibe notificaciones, noLeidas, marcarLeida y recargar como props
// desde el componente padre (que tiene la única instancia de useNotificaciones).
// Esto evita que se creen múltiples instancias del hook con pollers duplicados.
export default function NotificacionesDropdown({ notificaciones = [], noLeidas = 0, marcarLeida, recargar, onNavigate }) {
  const [open, setOpen] = useState(false);
  const dropdownRef = useRef(null);

  useEffect(() => {
    const handleClickOutside = (e) => {
      if (dropdownRef.current && !dropdownRef.current.contains(e.target)) {
        setOpen(false);
      }
    };
    document.addEventListener("mousedown", handleClickOutside);
    return () => document.removeEventListener("mousedown", handleClickOutside);
  }, []);

  const handleClick = (n) => {
    if (n.estado === "Pendiente" && marcarLeida) marcarLeida(n.id);
    setOpen(false);
    if (n.documento_radicado_id && onNavigate) {
      onNavigate(n.documento_radicado_id);
    }
  };

  return (
    <div className="notif-dropdown" ref={dropdownRef}>
      <button className="notif-bell" onClick={() => setOpen(!open)} title="Notificaciones">
        <i className="fa-solid fa-bell"></i>
        {noLeidas > 0 && (
          <span className="notif-badge">{noLeidas > 99 ? "99+" : noLeidas}</span>
        )}
      </button>

      {open && (
        <div className="notif-menu">
          <div className="notif-header">
            <strong>Notificaciones</strong>
            <button className="notif-refresh" onClick={recargar} title="Recargar">
              <i className="fa-solid fa-rotate-right"></i>
            </button>
          </div>
          <div className="notif-list">
            {notificaciones.length === 0 ? (
              <p className="notif-empty">No tienes notificaciones</p>
            ) : (
              notificaciones.map((n) => (
                <div
                  key={n.id}
                  className={`notif-item ${n.estado === "Pendiente" ? "no-leida" : ""}`}
                  onClick={() => handleClick(n)}
                >
                  <div className="notif-icon">
                    {n.tipo === "Asignacion" && <i className="fa-solid fa-user-check"></i>}
                    {n.tipo === "Recordatorio" && <i className="fa-solid fa-clock"></i>}
                    {n.tipo === "Sistema" && <i className="fa-solid fa-gear"></i>}
                    {n.tipo === "Rechazo" && <i className="fa-solid fa-ban"></i>}
                    {n.tipo === "Devolucion" && <i className="fa-solid fa-rotate-left"></i>}
                    {n.tipo === "Finalizado" && <i className="fa-solid fa-circle-check"></i>}
                  </div>
                  <div className="notif-content">
                    <p className="notif-msg">{n.mensaje}</p>
                    <span className="notif-time">
                      {new Date(n.fecha_creacion).toLocaleString()}
                    </span>
                  </div>
                  {n.estado === "Pendiente" && <span className="notif-dot"></span>}
                </div>
              ))
            )}
          </div>
        </div>
      )}
    </div>
  );
}