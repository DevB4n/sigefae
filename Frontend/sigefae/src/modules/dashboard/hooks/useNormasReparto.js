import { useState } from "react";
import { API } from "../constants/api";

export function useNormasReparto(obtenerToken, esAdmin, puedeGestionarRecurso, normasRepartoRadicado, setNormasRepartoRadicado) {
  const [showNormaModal, setShowNormaModal] = useState(false);
  const [normaEditandoId, setNormaEditandoId] = useState(null);
  const [normaFormDetalle, setNormaFormDetalle] = useState({ norma_reparto_id: "", porcentaje: "", proyecto: "", descripcion: "" });
  const [normaModalRadicadoId, setNormaModalRadicadoId] = useState(null);
  const [normaFiltroSede, setNormaFiltroSede] = useState("");
  const [normaFiltroArea, setNormaFiltroArea] = useState("");
  const [normasRepartoCatalogo, setNormasRepartoCatalogo] = useState([]);

  const sedesDisponibles = ["BUCARAMANGA", "MALAMBO", "CUCUTA", "CB", "CIENAGA DE ORO", "GENERAL"];
  const areasDisponibles = ["ADMON", "VENTAS", "PRODUCCION"];

  const openNormaModal = (radicadoId, normaExistente = null) => {
    setNormaModalRadicadoId(radicadoId);
    if (normaExistente) {
      setNormaEditandoId(normaExistente.id);
      setNormaFormDetalle({
        norma_reparto_id: String(normaExistente.norma_reparto_id || normaExistente.norma_reparto?.id || ""),
        porcentaje: String(normaExistente.porcentaje || ""),
        proyecto: normaExistente.proyecto || "",
        descripcion: normaExistente.descripcion || ""
      });
      const norma = normasRepartoCatalogo.find(n => n.id === (normaExistente.norma_reparto_id || normaExistente.norma_reparto?.id));
      if (norma) { setNormaFiltroSede(norma.sucursal || ""); setNormaFiltroArea(norma.departamento || ""); }
    } else {
      setNormaEditandoId(null); setNormaFormDetalle({ norma_reparto_id: "", porcentaje: "", proyecto: "", descripcion: "" });
      setNormaFiltroSede(""); setNormaFiltroArea("");
    }
    if (normasRepartoCatalogo.length === 0) {
      fetch(`${API}/normas-reparto?activo=true&_t=${new Date().getTime()}`, { headers: { Authorization: `Bearer ${obtenerToken()}` } })
        .then(r => r.json()).then(d => setNormasRepartoCatalogo(Array.isArray(d) ? d : []));
    }
    setShowNormaModal(true);
  };

  const handleGuardarNormaDetalle = async (radicadoId) => {
    if (!normaFormDetalle.norma_reparto_id || normaFormDetalle.porcentaje === "") { alert("Seleccione una norma y el porcentaje"); return; }
    if (!normaFormDetalle.proyecto || !normaFormDetalle.proyecto.trim()) { alert("El campo Proyecto es obligatorio"); return; }
    try {
      // Norma editada / nueva
      const payloadNorma = {
        norma_reparto_id: parseInt(normaFormDetalle.norma_reparto_id),
        porcentaje: parseFloat(normaFormDetalle.porcentaje),
        proyecto: normaFormDetalle.proyecto.trim(),
        descripcion: normaFormDetalle.descripcion?.trim() || ""
      };

      // Construir lista final preservando datos de las normas no editadas
      let nuevasNormas;
      if (normaEditandoId) {
        nuevasNormas = normasRepartoRadicado.map(n => {
          const nId = n.id;
          if (nId === normaEditandoId) {
            // Reemplazar con datos del form
            return payloadNorma;
          }
          // Preservar datos existentes exactamente
          return {
            norma_reparto_id: n.norma_reparto_id || n.norma_reparto?.id,
            porcentaje: parseFloat(n.porcentaje),
            proyecto: n.proyecto || "",
            descripcion: n.descripcion || ""
          };
        });
        // Si no encontró por id, agregar al final (fallback)
        const encontrado = normasRepartoRadicado.some(n => n.id === normaEditandoId);
        if (!encontrado) nuevasNormas.push(payloadNorma);
      } else {
        // Nueva norma: verificar duplicado
        const yaExiste = normasRepartoRadicado.find(n =>
          (n.norma_reparto_id || n.norma_reparto?.id) === payloadNorma.norma_reparto_id
        );
        if (yaExiste) { alert("Esta norma ya fue agregada"); return; }
        nuevasNormas = [
          ...normasRepartoRadicado.map(n => ({
            norma_reparto_id: n.norma_reparto_id || n.norma_reparto?.id,
            porcentaje: parseFloat(n.porcentaje),
            proyecto: n.proyecto || "",
            descripcion: n.descripcion || ""
          })),
          payloadNorma
        ];
      }

      const res = await fetch(`${API}/documentoradicado/${radicadoId}/normas-reparto`, {
        method: "POST", headers: { "Content-Type": "application/json", Authorization: `Bearer ${obtenerToken()}` },
        body: JSON.stringify({ normas: nuevasNormas })
      });
      if (!res.ok) { const err = await res.json(); throw new Error(err.error || "Error guardando norma"); }
      alert(normaEditandoId ? "Norma actualizada correctamente" : "Norma agregada correctamente");
      setShowNormaModal(false); setNormaEditandoId(null);
      const nrRes = await fetch(`${API}/documentoradicado/${radicadoId}/normas-reparto?_t=${new Date().getTime()}`, { headers: { Authorization: `Bearer ${obtenerToken()}` } });
      const nrData = await nrRes.json();
      setNormasRepartoRadicado(Array.isArray(nrData) ? nrData : []);
    } catch (err) { alert("Error: " + err.message); }
  };

  const handleEliminarNorma = async (asignacion, radicadoId) => {
    const norma = asignacion || {};
    const creadorId = Number(norma.creado_por_id || norma.creado_por?.id || 0);
    if (!esAdmin && !puedeGestionarRecurso(creadorId)) { alert("No tienes permisos para eliminar esta norma de reparto."); return; }
    if (!confirm("¿Está seguro de eliminar esta norma de reparto?")) return;
    try {
      const normasRestantes = normasRepartoRadicado.filter(n => n.id !== norma.id).map(n => ({
        norma_reparto_id: n.norma_reparto_id || n.norma_reparto?.id,
        porcentaje: parseFloat(n.porcentaje),
        proyecto: n.proyecto || "",
        descripcion: n.descripcion || ""
      }));
      const res = await fetch(`${API}/documentoradicado/${radicadoId}/normas-reparto`, {
        method: "POST", headers: { "Content-Type": "application/json", Authorization: `Bearer ${obtenerToken()}` },
        body: JSON.stringify({ normas: normasRestantes })
      });
      if (!res.ok) { const err = await res.json(); throw new Error(err.error || "Error eliminando norma"); }
      alert("Norma eliminada");
      const nrRes = await fetch(`${API}/documentoradicado/${radicadoId}/normas-reparto?_t=${new Date().getTime()}`, { headers: { Authorization: `Bearer ${obtenerToken()}` } });
      const nrData = await nrRes.json();
      setNormasRepartoRadicado(Array.isArray(nrData) ? nrData : []);
    } catch (err) { alert("Error: " + err.message); }
  };

  return {
    showNormaModal, setShowNormaModal, normaEditandoId, setNormaEditandoId,
    normaFormDetalle, setNormaFormDetalle, normaModalRadicadoId,
    normaFiltroSede, setNormaFiltroSede, normaFiltroArea, setNormaFiltroArea,
    normasRepartoCatalogo, setNormasRepartoCatalogo,
    sedesDisponibles, areasDisponibles,
    openNormaModal, handleGuardarNormaDetalle, handleEliminarNorma
  };
}