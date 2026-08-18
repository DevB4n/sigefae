import { useState, useEffect } from "react";
import { API } from "../constants/api";

export function useCorreos(obtenerToken, activeTab) {
  const [correos, setCorreos] = useState([]);
  const [selectedCorreoId, setSelectedCorreoId] = useState(null);
  const [correoDetail, setCorreoDetail] = useState(null);
  const [loading, setLoading] = useState(false);
  const [searchCorreos, setSearchCorreos] = useState("");
  const [sortCorreos, setSortCorreos] = useState("fecha_desc");

  useEffect(() => {
    if (activeTab !== "correos") return;
    setLoading(true);
    fetch(`${API}/correo`, { headers: { Authorization: `Bearer ${obtenerToken()}` } })
      .then((res) => res.json())
      .then((data) => { if (Array.isArray(data)) setCorreos(data); })
      .catch((err) => console.error(err))
      .finally(() => setLoading(false));
  }, [activeTab, obtenerToken]);

  useEffect(() => {
    if (!selectedCorreoId) return;
    setCorreoDetail(null);
    fetch(`${API}/correo/${selectedCorreoId}`, { headers: { Authorization: `Bearer ${obtenerToken()}` } })
      .then((res) => res.json())
      .then((data) => { if (data?.id) setCorreoDetail(data); })
      .catch((err) => console.error(err));
  }, [selectedCorreoId, obtenerToken]);

  const handleVerArchivo = (filename) => {
    if (!correoDetail) return;
    window.open(`${API}/storage/mails/${correoDetail.id_mensaje}/${filename}`, "_blank");
  };

  const getFilteredCorreos = () => {
    let result = [...correos];
    if (searchCorreos) {
      const q = searchCorreos.toLowerCase();
      result = result.filter(c => (c.asunto?.toLowerCase().includes(q)) || (c.remitente?.toLowerCase().includes(q)));
    }
    result.sort((a, b) => {
      if (sortCorreos === 'fecha_desc') return new Date(b.fecha_recepcion) - new Date(a.fecha_recepcion);
      if (sortCorreos === 'fecha_asc') return new Date(a.fecha_recepcion) - new Date(b.fecha_recepcion);
      if (sortCorreos === 'estado') return (a.estado_correo?.nombre || '').localeCompare(b.estado_correo?.nombre || '');
      return 0;
    });
    return result;
  };

  return {
    correos, selectedCorreoId, setSelectedCorreoId, correoDetail,
    loading, searchCorreos, setSearchCorreos, sortCorreos, setSortCorreos,
    handleVerArchivo, getFilteredCorreos
  };
}