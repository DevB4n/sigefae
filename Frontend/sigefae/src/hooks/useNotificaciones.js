// useNotificaciones.js
import { useState, useEffect, useCallback, useRef } from "react";
import { obtenerToken } from "../modules/auth/token.js";
import {
  playNotificacionSound,
  mostrarNotificacionNativa,
  pedirPermisoNotificaciones,
  inicializarAudio,
} from "../utils/notificacion.js";

const API = "http://localhost:8080/api";

export function useNotificaciones() {
  const [notificaciones, setNotificaciones] = useState([]);
  const [noLeidas, setNoLeidas] = useState(0);
  const [loading, setLoading] = useState(false);
  const idsVistosRef = useRef(new Set());
  const primeraCargaRef = useRef(true);

  const cargar = useCallback(async () => {
    setLoading(true);
    try {
      const res = await fetch(`${API}/notificacion/mias`, {
        headers: { Authorization: `Bearer ${obtenerToken()}` },
      });
      const data = await res.json();
      if (Array.isArray(data)) {
        const pendientes = data.filter((n) => n.estado === "Pendiente");

        if (!primeraCargaRef.current) {
          const nuevas = pendientes.filter((n) => !idsVistosRef.current.has(n.id));
          if (nuevas.length > 0) {
            playNotificacionSound();
            const ultima = nuevas[0];
            mostrarNotificacionNativa("SIGEFAE", ultima.mensaje, () => {
              window.dispatchEvent(
                new CustomEvent("navegar-a-radicado", {
                  detail: ultima.documento_radicado_id,
                })
              );
            });
          }
        }

        idsVistosRef.current = new Set(pendientes.map((n) => n.id));
        primeraCargaRef.current = false;
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
    inicializarAudio();
  }, []);

  useEffect(() => {
    cargar();
    const interval = setInterval(cargar, 30000);
    return () => clearInterval(interval);
  }, [cargar]);

  return { notificaciones, noLeidas, loading, marcarLeida, recargar: cargar };
}