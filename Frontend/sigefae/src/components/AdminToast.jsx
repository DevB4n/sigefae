import { useEffect, useState } from "react";
import "./notificaciones.css";

// Recibe esAdmin como prop para evitar leer localStorage en render
export default function AdminToast({ esAdmin }) {
  const [toasts, setToasts] = useState([]);

  useEffect(() => {
    if (!esAdmin) return; // solo Superadministrador recibe estos toasts

    const handler = (e) => {
      const n = e.detail;
      const id = `toast-${Date.now()}-${Math.random()}`;
      setToasts((t) => [{ id, mensaje: n.mensaje, documento_radicado_id: n.documento_radicado_id }, ...t.slice(0, 4)]);
      // auto-remover después de 6s
      setTimeout(() => {
        setToasts((t) => t.filter(x => x.id !== id));
      }, 6000);
    };

    window.addEventListener('notificacion-nueva', handler);
    return () => window.removeEventListener('notificacion-nueva', handler);
  }, [esAdmin]);

  const handleClick = (docId) => {
    if (!docId) return;
    window.dispatchEvent(new CustomEvent('navegar-a-radicado', { detail: docId }));
  };

  const handleClose = (e, id) => {
    e.stopPropagation();
    setToasts((t) => t.filter(x => x.id !== id));
  };

  if (toasts.length === 0) return null;

  return (
    <div className="admin-toast-wrap">
      {toasts.map(t => (
        <div key={t.id} className="admin-toast" onClick={() => handleClick(t.documento_radicado_id)}>
          <div className="admin-toast-title">
            <i className="fa-solid fa-bell" style={{ marginRight: 6 }}></i>
            Nueva Notificación
          </div>
          <div className="admin-toast-body">{t.mensaje}</div>
          <button
            onClick={(e) => handleClose(e, t.id)}
            style={{ position: 'absolute', top: 6, right: 8, background: 'none', border: 'none', color: 'inherit', cursor: 'pointer', fontSize: '0.85em', opacity: 0.7 }}
            title="Cerrar"
          >
            <i className="fa-solid fa-xmark"></i>
          </button>
        </div>
      ))}
    </div>
  );
}
