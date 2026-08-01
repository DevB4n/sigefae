import { useState, useEffect, useCallback, useRef } from "react";
import { obtenerToken } from "../modules/auth/token.js";
import {
  playNotificacionSound,
  mostrarNotificacionNativa,
  pedirPermisoNotificaciones,
} from "../utils/notificacion.js";

const API = "http://localhost:8080/api";

export function useNotificaciones() {
  const [notificaciones, setNotificaciones] = useState([]);
  const [noLeidas, setNoLeidas] = useState(0);
  const [loading, setLoading] = useState(false);
  const prevCountRef = useRef(0);

  const cargar = useCallback(async () => {
    setLoading(true);
    try {
      const res = await fetch(`${API}/notificacion/mias`, {
        headers: { Authorization: `Bearer ${obtenerToken()}` },
      });
      const data = await res.json();
      if (Array.isArray(data)) {
        const pendientes = data.filter((n) => n.estado === "Pendiente");

        // ── SONIDO: si llegaron nuevas pendientes (no en la carga inicial) ──
        if (pendientes.length > prevCountRef.current && prevCountRef.current > 0) {
          playNotificacionSound();

          const ultima = data[0];
          if (ultima) {
            mostrarNotificacionNativa("SIGEFAE", ultima.mensaje, () => {
              window.dispatchEvent(
                new CustomEvent("navegar-a-radicado", {
                  detail: ultima.documento_radicado_id,
                })
              );
            });
          }
        }

        prevCountRef.current = pendientes.length;
        setNotificaciones(data);
        setNoLeidas(pendientes.length);
      }
    } catch (err) {
      console.error("Error cargando notificaciones:", err);
    } finally {
      setLoading(false);
    }
  }, []);

  const marcarLeida = async (id) => {
    try {
      const res = await fetch(`${API}/notificacion/${id}/leida`, {
        method: "PATCH",
        headers: { Authorization: `Bearer ${obtenerToken()}` },
      });
      if (res.ok) cargar();
    } catch (err) {
      console.error(err);
    }
  };

  useEffect(() => {
    pedirPermisoNotificaciones();
  }, []);

  useEffect(() => {
    cargar();
    const interval = setInterval(cargar, 30000);
    return () => clearInterval(interval);
  }, [cargar]);

  return { notificaciones, noLeidas, loading, marcarLeida, recargar: cargar };
}