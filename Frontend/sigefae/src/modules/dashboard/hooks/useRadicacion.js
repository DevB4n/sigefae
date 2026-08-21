import { useState, useEffect } from "react";
import { API } from "../constants/api";

export function useRadicacion(obtenerToken, userId, setDocumentos, setActiveTab, setSelectedDocId, setDocDetail) {
  const [showRadicarModal, setShowRadicarModal] = useState(false);
  const [radicarDocId, setRadicarDocId] = useState(null);
  const [radicarForm, setRadicarForm] = useState({ tipo_radicacion_id: "", ruta_id: "", metodo_pago_id: "", numero_radicado: "", normas_reparto: [] });
  const [radicando, setRadicando] = useState(false);
  const [normasRepartoAutoMsg, setNormasRepartoAutoMsg] = useState("");

  const [tiposRadicacion, setTiposRadicacion] = useState([]);
  const [rutas, setRutas] = useState([]);
  const [metodosPago, setMetodosPago] = useState([]);
  const [normasRepartoCatalogo, setNormasRepartoCatalogo] = useState([]);

  const [normaFiltroSede, setNormaFiltroSede] = useState("");
  const [normaFiltroArea, setNormaFiltroArea] = useState("");
  const [normaSeleccionadaId, setNormaSeleccionadaId] = useState("");
  const [normaPorcentajeInput, setNormaPorcentajeInput] = useState("");
  const [normaValorInput, setNormaValorInput] = useState("");
  const [subtotalDoc, setSubtotalDoc] = useState(0);

  const sedesDisponibles = ["BUCARAMANGA", "MALAMBO", "CUCUTA", "CB", "CIENAGA DE ORO", "GENERAL"];
  const areasDisponibles = ["ADMON", "VENTAS", "PRODUCCION"];

  useEffect(() => {
    if (!showRadicarModal) return;
    const headers = { Authorization: `Bearer ${obtenerToken()}` };
    Promise.all([
      fetch(`${API}/tipo-radicacion`, { headers }).then(r => r.json()),
      fetch(`${API}/rutas`, { headers }).then(r => r.json()),
      fetch(`${API}/metodos-pago`, { headers }).then(r => r.json()),
      fetch(`${API}/normas-reparto?activo=true`, { headers }).then(r => r.json()),
    ]).then(([tr, r, mp, nr]) => {
      setTiposRadicacion(Array.isArray(tr) ? tr : []);
      setRutas(Array.isArray(r) ? r : []);
      setMetodosPago(Array.isArray(mp) ? mp : []);
      setNormasRepartoCatalogo(Array.isArray(nr) ? nr : []);
    }).catch(err => console.error("Error cargando catálogos:", err));
  }, [showRadicarModal, obtenerToken]);

  const openRadicarModal = (docId, subtotal = 0) => {
    setRadicarDocId(docId);
    setSubtotalDoc(parseFloat(subtotal) || 0);
    setRadicarForm({ tipo_radicacion_id: "", ruta_id: "", metodo_pago_id: "", numero_radicado: "", normas_reparto: [] });
    setNormasRepartoAutoMsg(""); setNormaFiltroSede(""); setNormaFiltroArea("");
    setNormaSeleccionadaId(""); setNormaPorcentajeInput(""); setNormaValorInput("");
    setShowRadicarModal(true);
  };

  const handleRadicarChange = (e) => {
    const { name, value } = e.target;
    setRadicarForm(prev => ({ ...prev, [name]: value }));
  };

  const handleAgregarNormaModal = () => {
    if (!normaSeleccionadaId) { alert("Selecciona una norma"); return; }
    if (!normaPorcentajeInput && !normaValorInput) { alert("Ingresa un porcentaje o un valor"); return; }
    
    let pct = parseFloat(normaPorcentajeInput);
    let val = parseFloat(normaValorInput);
    
    // Auto-calcular el que falte
    if (subtotalDoc > 0) {
      if (normaValorInput && !normaPorcentajeInput) {
        pct = (val / subtotalDoc) * 100;
      } else if (normaPorcentajeInput && !normaValorInput) {
        val = (pct / 100) * subtotalDoc;
      }
    }
    
    if (isNaN(pct) || pct <= 0) { alert("Porcentaje inválido"); return; }
    if (pct > 100) { alert("El porcentaje no puede superar el 100%"); return; }
    if (subtotalDoc > 0 && val > subtotalDoc) { alert("El valor no puede superar el subtotal"); return; }
    
    const yaExiste = radicarForm.normas_reparto.find(n => n.norma_reparto_id === normaSeleccionadaId);
    if (yaExiste) { alert("Esta norma ya fue agregada"); return; }
    setRadicarForm(prev => ({ ...prev, normas_reparto: [...prev.normas_reparto, { norma_reparto_id: normaSeleccionadaId, porcentaje: pct.toFixed(2), valor: val.toFixed(2) }] }));
    setNormaSeleccionadaId(""); setNormaPorcentajeInput(""); setNormaValorInput("");
  };

  const handleNormaRepartoChange = (index, field, value) => {
    setRadicarForm(prev => {
      const updated = [...prev.normas_reparto];
      updated[index] = { ...updated[index], [field]: value };
      return { ...prev, normas_reparto: updated };
    });
  };

  const handleRemoveNormaReparto = (index) => {
    setRadicarForm(prev => ({ ...prev, normas_reparto: prev.normas_reparto.filter((_, i) => i !== index) }));
  };

  const totalPorcentajeNormas = (radicarForm.normas_reparto || []).reduce((sum, n) => sum + (parseFloat(n.porcentaje) || 0), 0);

  const normasFiltradas = normasRepartoCatalogo.filter(n => {
    if (normaFiltroSede && n.sucursal !== normaFiltroSede) return false;
    if (normaFiltroArea && n.departamento !== normaFiltroArea) return false;
    return true;
  });

  const handleRadicarSubmit = async () => {
    if (!radicarDocId) return;
    if (!radicarForm.tipo_radicacion_id || !radicarForm.ruta_id || !radicarForm.metodo_pago_id) {
      alert("Debe seleccionar tipo de radicación, ruta y método de pago."); return;
    }
    setRadicando(true);
    const payload = {
      documento_comercial_id: radicarDocId,
      tipo_radicacion_id: parseInt(radicarForm.tipo_radicacion_id),
      ruta_id: parseInt(radicarForm.ruta_id),
      metodo_pago_id: parseInt(radicarForm.metodo_pago_id),
      numero_radicado: radicarForm.numero_radicado?.trim() || "",
      normas_reparto: (radicarForm.normas_reparto || []).filter(n => n.norma_reparto_id && n.porcentaje).map(n => ({ norma_reparto_id: parseInt(n.norma_reparto_id), porcentaje: parseFloat(n.porcentaje) })),
    };
    try {
      const res = await fetch(`${API}/documentoradicado`, {
        method: "POST", headers: { "Content-Type": "application/json", Authorization: `Bearer ${obtenerToken()}` },
        body: JSON.stringify(payload),
      });
      if (!res.ok) { const errData = await res.json(); throw new Error(errData.error || "Error al radicar"); }
      const creado = await res.json();
      if (creado?.usuario_actual?.id && creado.usuario_actual.id !== userId) {
        await fetch(`${API}/notificacion`, {
          method: "POST", headers: { "Content-Type": "application/json", Authorization: `Bearer ${obtenerToken()}` },
          body: JSON.stringify({ usuario_id: creado.usuario_actual.id, documento_radicado_id: creado.id, mensaje: `Nuevo radicado #${creado.numero_radicado} requiere tu revisión`, estado: "Pendiente", tipo: "Asignacion", fecha_creacion: new Date().toISOString() })
        });
      }
      setShowRadicarModal(false); setRadicarDocId(null);
      setSelectedDocId(null); setDocDetail(null);
      setDocumentos(prev => prev.filter(d => d.id !== radicarDocId));
      setActiveTab("radicados");
    } catch (err) { alert("Error al radicar: " + err.message); }
    finally { setRadicando(false); }
  };

  return {
    showRadicarModal, setShowRadicarModal, radicarDocId, setRadicarDocId,
    radicarForm, setRadicarForm, radicando, normasRepartoAutoMsg,
    tiposRadicacion, rutas, metodosPago, normasRepartoCatalogo,
    normaFiltroSede, setNormaFiltroSede, normaFiltroArea, setNormaFiltroArea,
    normaSeleccionadaId, setNormaSeleccionadaId, normaPorcentajeInput, setNormaPorcentajeInput,
    normaValorInput, setNormaValorInput, subtotalDoc,
    sedesDisponibles, areasDisponibles, normasFiltradas, totalPorcentajeNormas,
    openRadicarModal, handleRadicarChange, handleAgregarNormaModal,
    handleNormaRepartoChange, handleRemoveNormaReparto, handleRadicarSubmit
  };
}