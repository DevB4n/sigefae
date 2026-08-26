import { useState, useEffect } from "react";
import { API } from "../constants/api";
import { isFinalState } from "../helpers/formatters";

export function useTareas(obtenerToken, activeTab, userId) {
  const [misTareas, setMisTareas] = useState([]);
  const [misTareasCompletadas, setMisTareasCompletadas] = useState([]);
  const [selectedTareaId, setSelectedTareaId] = useState(null);
  const [tareaDetail, setTareaDetail] = useState(null);
  const [loadingTareas, setLoadingTareas] = useState(false);
  const [tareasSubTab, setTareasSubTab] = useState("activas");
  const [searchTareas, setSearchTareas] = useState("");
  const [sortTareas, setSortTareas] = useState("fecha_desc");

  useEffect(() => {
    if (activeTab !== "tareas") return;
    setLoadingTareas(true);
    const t = `?_t=${new Date().getTime()}`;
    fetch(`${API}/documentoradicado${t}`, { headers: { Authorization: `Bearer ${obtenerToken()}` } })
      .then((res) => res.json())
      .then((data) => {
        if (Array.isArray(data)) {
          setMisTareas(data.filter(r => r.usuario_actual_id === userId && !isFinalState(r.estado_posesion)));
          setMisTareasCompletadas(data.filter(r => isFinalState(r.estado_posesion)));
        }
      })
      .catch((err) => console.error(err))
      .finally(() => setLoadingTareas(false));
  }, [activeTab, userId, obtenerToken]);

  useEffect(() => {
    if (!selectedTareaId) return;
    setTareaDetail(null);
    const t = `?_t=${new Date().getTime()}`;
    fetch(`${API}/documentoradicado/${selectedTareaId}${t}`, { headers: { Authorization: `Bearer ${obtenerToken()}` } })
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
          setTareaDetail(data);
        }
      })
      .catch((err) => console.error(err));
  }, [selectedTareaId, obtenerToken]);

  const getFilteredTareas = (baseList) => {
    let result = [...baseList];
    if (searchTareas) {
      const q = searchTareas.toLowerCase();
      result = result.filter(t => (t.numero_radicado?.toLowerCase().includes(q)) || (t.documento_comercial?.numero_documento?.toLowerCase().includes(q)));
    }
    result.sort((a, b) => {
      if (sortTareas === 'fecha_desc') return new Date(b.fecha_creacion) - new Date(a.fecha_creacion);
      if (sortTareas === 'fecha_asc') return new Date(a.fecha_creacion) - new Date(b.fecha_creacion);
      if (sortTareas === 'estado') return (a.estado_posesion || '').localeCompare(b.estado_posesion || '');
      return 0;
    });
    return result;
  };

  return {
    misTareas, setMisTareas, misTareasCompletadas, setMisTareasCompletadas, selectedTareaId, setSelectedTareaId,
    tareaDetail, setTareaDetail, loadingTareas, tareasSubTab, setTareasSubTab,
    searchTareas, setSearchTareas, sortTareas, setSortTareas, getFilteredTareas
  };
}