import { useState, useEffect } from "react";
import { API } from "../constants/api";

export function useSaia(obtenerToken) {
  const [saiaModalOpen, setSaiaModalOpen] = useState(false);
  const [saiaRadicado, setSaiaRadicado] = useState(null);
  const [saiaActiveTab, setSaiaActiveTab] = useState("info");
  const [saiaAnexoIdx, setSaiaAnexoIdx] = useState(0);
  const [saiaPdfUrl, setSaiaPdfUrl] = useState(null);

  // PDF blob para visor SAIA
  useEffect(() => {
    if (!saiaModalOpen || !saiaRadicado) {
      if (saiaPdfUrl) { URL.revokeObjectURL(saiaPdfUrl); setSaiaPdfUrl(null); }
      return;
    }
    const archivosPdf = (saiaRadicado.archivos || []).filter(a => a.extension?.toLowerCase() === 'pdf' || a.nombre?.toLowerCase().endsWith('.pdf'));
    const anexoActual = archivosPdf[saiaAnexoIdx];
    if (!anexoActual) {
      if (saiaPdfUrl) { URL.revokeObjectURL(saiaPdfUrl); setSaiaPdfUrl(null); }
      return;
    }
    let cancelled = false;
    fetch(`${API}/archivo/${anexoActual.id}/download?t=${new Date().getTime()}`, { headers: { Authorization: `Bearer ${obtenerToken()}` } })
      .then(res => { if (!res.ok) throw new Error("No se pudo cargar el PDF"); return res.blob(); })
      .then(blob => { if (cancelled) return; const url = URL.createObjectURL(blob); setSaiaPdfUrl(prev => { if (prev) URL.revokeObjectURL(prev); return url; }); })
      .catch(() => { if (!cancelled) setSaiaPdfUrl(null); });
    return () => { cancelled = true; };
  }, [saiaModalOpen, saiaRadicado, saiaAnexoIdx, obtenerToken]);

  const openSaia = async (radicadoBase, fromTab, setSelectedTareaId, setSelectedRadicadoId, recargarFlujo, recargarTrazabilidad, recargarNormas, recargarComentarios) => {
    if (saiaPdfUrl) { URL.revokeObjectURL(saiaPdfUrl); setSaiaPdfUrl(null); }
    try {
      const res = await fetch(`${API}/documentoradicado/${radicadoBase.id}`, { headers: { Authorization: `Bearer ${obtenerToken()}` } });
      const rad = await res.json();
      if (rad?.id) {
        if (rad.ruta?.id && !rad.ruta?.area) {
          try {
            const rutasRes = await fetch(`${API}/rutas`, { headers: { Authorization: `Bearer ${obtenerToken()}` } });
            const rutasList = await rutasRes.json();
            const rutaCompleta = (Array.isArray(rutasList) ? rutasList : []).find(r => r.id === rad.ruta.id);
            if (rutaCompleta) rad.ruta.area = rutaCompleta.area;
          } catch (e) {}
        }
        setSaiaRadicado(rad);
        setSaiaAnexoIdx(0); setSaiaActiveTab("info");
        if (fromTab === "tareas") { setSelectedTareaId(rad.id); setSelectedRadicadoId(null); }
        else { setSelectedRadicadoId(rad.id); setSelectedTareaId(null); }
        await recargarFlujo(rad.id); await recargarTrazabilidad(rad.id); await recargarNormas(rad.id); await recargarComentarios(rad.id);
        setSaiaModalOpen(true);
      }
    } catch (err) { console.error(err); alert("Error cargando detalle del radicado"); }
  };

  return {
    saiaModalOpen, setSaiaModalOpen, saiaRadicado, setSaiaRadicado,
    saiaActiveTab, setSaiaActiveTab, saiaAnexoIdx, setSaiaAnexoIdx, saiaPdfUrl, setSaiaPdfUrl,
    openSaia
  };
}