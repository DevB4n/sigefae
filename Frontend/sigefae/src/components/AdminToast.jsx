import { useEffect, useState } from "react";
import "./notificaciones.css";

export default function AdminToast() {
  const [toasts, setToasts] = useState([]);

  useEffect(() => {
    const handler = (e) => {
      const rol = localStorage.getItem('rol') || '';
      if (rol !== 'Superadministrador') return; // only admins
      const n = e.detail;
      const id = `toast-${Date.now()}`;
      setToasts((t) => [{ id, mensaje: n.mensaje, documento_radicado_id: n.documento_radicado_id }, ...t]);
      // auto remove after 6s
      setTimeout(() => {
        setToasts((t) => t.filter(x => x.id !== id));
      }, 6000);
    };
    window.addEventListener('notificacion-nueva', handler);
    return () => window.removeEventListener('notificacion-nueva', handler);
  }, []);

  const handleClick = (docId) => {
    if (!docId) return;
    window.dispatchEvent(new CustomEvent('navegar-a-radicado', { detail: docId }));
  };

  if (toasts.length === 0) return null;

  return (
    <div className="admin-toast-wrap">
      {toasts.map(t => (
        <div key={t.id} className="admin-toast" onClick={() => handleClick(t.documento_radicado_id)}>
          <div className="admin-toast-title">Nueva Notificación</div>
          <div className="admin-toast-body">{t.mensaje}</div>
        </div>
      ))}
    </div>
  );
}
