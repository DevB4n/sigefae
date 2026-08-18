import { useState } from "react";
import { generarExpedientePDF } from "../../../utils/expedientePdf.js";

export function usePdfExpediente() {
  const [generandoPdf, setGenerandoPdf] = useState(false);

  const handleDescargarExpediente = async (radicado, tareasFlujo, historialTrazabilidad, API) => {
    if (!radicado) return;
    try {
      setGenerandoPdf(true);
      const anexosUrls = (radicado.archivos || [])
        .filter(a => a.extension?.toLowerCase() === 'pdf' || a.nombre?.toLowerCase().endsWith('.pdf'))
        .map(a => `${API}/archivo/${a.id}/download?download=1`);
      await generarExpedientePDF(radicado, tareasFlujo, historialTrazabilidad, anexosUrls);
    } catch (err) {
      console.error("Error al generar expediente", err);
      alert("Hubo un error al generar el expediente PDF.");
    } finally { setGenerandoPdf(false); }
  };

  return { generandoPdf, handleDescargarExpediente };
}