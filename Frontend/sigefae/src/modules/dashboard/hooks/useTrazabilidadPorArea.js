import { useState, useEffect } from "react";
import { API } from "../constants/api";

export function useTrazabilidadPorArea(obtenerToken) {
  const [areas, setAreas] = useState([]);
  const [areaSeleccionada, setAreaSeleccionada] = useState("");
  const [fechaDesde, setFechaDesde] = useState("");
  const [fechaHasta, setFechaHasta] = useState("");
  const [radicados, setRadicados] = useState([]);
  const [seleccionados, setSeleccionados] = useState(new Set());
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    fetch(`${API}/areas`, { headers: { Authorization: `Bearer ${obtenerToken()}` } })
      .then(r => r.json())
      .then(data => setAreas(Array.isArray(data) ? data : []))
      .catch(err => console.error(err));
  }, [obtenerToken]);

  const buscar = async () => {
    if (!areaSeleccionada || !fechaDesde || !fechaHasta) {
      alert("Seleccione un área y el rango de fechas");
      return;
    }
    setLoading(true);
    try {
      const url = `${API}/trazabilidad/por-area?area_id=${areaSeleccionada}&fecha_desde=${fechaDesde}&fecha_hasta=${fechaHasta}`;
      const res = await fetch(url, { headers: { Authorization: `Bearer ${obtenerToken()}` } });
      const data = await res.json();
      if (!res.ok) throw new Error(data.error || "Error del servidor");
      setRadicados(Array.isArray(data) ? data : []);
      setSeleccionados(new Set());
    } catch (err) {
      alert("Error: " + err.message);
    } finally {
      setLoading(false);
    }
  };

  const toggleSeleccion = (id) => {
    setSeleccionados(prev => {
      const next = new Set(prev);
      if (next.has(id)) next.delete(id);
      else next.add(id);
      return next;
    });
  };

  const toggleTodos = () => {
    if (seleccionados.size === radicados.length && radicados.length > 0) {
      setSeleccionados(new Set());
    } else {
      setSeleccionados(new Set(radicados.map(r => r.id)));
    }
  };

  return {
    areas, areaSeleccionada, setAreaSeleccionada,
    fechaDesde, setFechaDesde, fechaHasta, setFechaHasta,
    radicados, loading, seleccionados,
    buscar, toggleSeleccion, toggleTodos, setSeleccionados
  };
}