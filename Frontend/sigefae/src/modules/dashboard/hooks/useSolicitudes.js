import { useState, useEffect } from "react";
import { API } from "../constants/api";

export function useSolicitudes(obtenerToken, activeTab, esAdmin) {
  const [solicitudes, setSolicitudes] = useState([]);
  const [loadingSolicitudes, setLoadingSolicitudes] = useState(false);

  useEffect(() => {
    if (activeTab !== "solicitudes") return;
    setLoadingSolicitudes(true);
    const endpoint = esAdmin ? `${API}/solicitud-rechazo` : `${API}/solicitud-rechazo/mias`;
    fetch(endpoint, { headers: { Authorization: `Bearer ${obtenerToken()}` } })
      .then(r => r.json())
      .then(data => setSolicitudes(Array.isArray(data) ? data : []))
      .catch(err => { console.error(err); setSolicitudes([]); })
      .finally(() => setLoadingSolicitudes(false));
  }, [activeTab, esAdmin, obtenerToken]);

  const handleDecidir = async (id, accept) => {
    if (!confirm(accept ? 'Aceptar solicitud y marcar como rechazado?' : 'Rechazar solicitud (no marcar documento)?')) return;
    try {
      const res = await fetch(`${API}/solicitud-rechazo/${id}/decidir`, {
        method: 'POST', headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${obtenerToken()}` },
        body: JSON.stringify({ accept, mensaje: '' })
      });
      if (!res.ok) throw new Error((await res.json()).error || 'Error');
      alert(accept ? 'Solicitud aceptada' : 'Solicitud rechazada');
      setSolicitudes(prev => prev.filter(x => x.id !== id));
    } catch (err) { alert('Error: ' + err.message); }
  };

  return { solicitudes, loadingSolicitudes, handleDecidir };
}