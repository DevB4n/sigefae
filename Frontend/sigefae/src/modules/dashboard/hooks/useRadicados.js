import { useState, useEffect } from "react";
import { API } from "../constants/api";
import { isFinalState } from "../helpers/formatters";

export function useRadicados(obtenerToken, activeTab, userId) {
  const [radicados, setRadicados] = useState([]);
  const [selectedRadicadoId, setSelectedRadicadoId] = useState(null);
  const [radicadoDetail, setRadicadoDetail] = useState(null);
  const [loadingRadicados, setLoadingRadicados] = useState(false);
  const [searchRadicados, setSearchRadicados] = useState("");
  const [sortRadicados, setSortRadicados] = useState("fecha_desc");
  const [tareasPorRadicado, setTareasPorRadicado] = useState({});

  useEffect(() => {
    if (activeTab !== "radicados") return;
    setLoadingRadicados(true);
    fetch(`${API}/documentoradicado`, { headers: { Authorization: `Bearer ${obtenerToken()}` } })
      .then((res) => res.json())
      .then(async (data) => {
        if (!Array.isArray(data)) return;
        setRadicados(data);
        const tareasMap = {};
        await Promise.all(data.map(async (rad) => {
          if (isFinalState(rad.estado_posesion)) return;
          try {
            const res = await fetch(`${API}/documentoradicado/${rad.id}/tareas`, { headers: { Authorization: `Bearer ${obtenerToken()}` } });
            const tareas = await res.json();
            const tareaActiva = (Array.isArray(tareas) ? tareas : []).find(t => t.estado?.nombre === "En Proceso");
            if (tareaActiva) tareasMap[rad.id] = tareaActiva;
          } catch (e) {}
        }));
        setTareasPorRadicado(tareasMap);
      })
      .catch((err) => console.error(err))
      .finally(() => setLoadingRadicados(false));
  }, [activeTab, obtenerToken]);

  useEffect(() => {
    if (!selectedRadicadoId) return;
    setRadicadoDetail(null);
    fetch(`${API}/documentoradicado/${selectedRadicadoId}`, { headers: { Authorization: `Bearer ${obtenerToken()}` } })
      .then((res) => res.json())
      .then(async (data) => {
        if (data?.id) {
          if (data.ruta?.id && !data.ruta?.area) {
            try {
              const rutasRes = await fetch(`${API}/rutas`, { headers: { Authorization: `Bearer ${obtenerToken()}` } });
              const rutasList = await rutasRes.json();
              const rutaCompleta = (Array.isArray(rutasList) ? rutasList : []).find(r => r.id === data.ruta.id);
              if (rutaCompleta) data.ruta.area = rutaCompleta.area;
            } catch (e) {}
          }
          setRadicadoDetail(data);
        }
      })
      .catch((err) => console.error(err));
  }, [selectedRadicadoId, obtenerToken]);

  const getFilteredRadicados = () => {
    let result = [...radicados];
    if (searchRadicados) {
      const q = searchRadicados.toLowerCase();
      result = result.filter(r =>
        (r.numero_radicado?.toLowerCase().includes(q)) ||
        (r.documento_comercial?.numero_documento?.toLowerCase().includes(q)) ||
        (r.documento_comercial?.proveedor?.razon_social?.toLowerCase().includes(q))
      );
    }
    result.sort((a, b) => {
      if (sortRadicados === 'fecha_desc') return new Date(b.fecha_radicacion) - new Date(a.fecha_radicacion);
      if (sortRadicados === 'fecha_asc') return new Date(a.fecha_radicacion) - new Date(b.fecha_radicacion);
      if (sortRadicados === 'estado') return (a.estado_posesion || '').localeCompare(b.estado_posesion || '');
      return 0;
    });
    return result;
  };

  return {
    radicados, setRadicados, selectedRadicadoId, setSelectedRadicadoId, radicadoDetail, setRadicadoDetail,
    loadingRadicados, searchRadicados, setSearchRadicados, sortRadicados, setSortRadicados,
    tareasPorRadicado, getFilteredRadicados
  };
}