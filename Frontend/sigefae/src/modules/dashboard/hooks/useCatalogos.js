import { useState, useEffect } from "react";
import { API } from "../constants/api";
import { catalogoConfig } from "../constants/catalogoConfig";

export function useCatalogos(obtenerToken, activeTab) {
  const [catalogoActivo, setCatalogoActivo] = useState("tipo-radicacion");
  const [catalogoItems, setCatalogoItems] = useState([]);
  const [catalogoLoading, setCatalogoLoading] = useState(false);
  const [showCatalogoForm, setShowCatalogoForm] = useState(false);
  const [catalogoEditing, setCatalogoEditing] = useState(null);
  const [catalogoForm, setCatalogoForm] = useState({});

  const [tiposPagoCatalogo, setTiposPagoCatalogo] = useState([]);
  const [areasCatalogo, setAreasCatalogo] = useState([]);
  const [rutasCatalogo, setRutasCatalogo] = useState([]);
  const [usuariosCatalogo, setUsuariosCatalogo] = useState([]);
  const [monedasCatalogo, setMonedasCatalogo] = useState([]);

  const cfg = catalogoConfig[catalogoActivo];

  const loadCatalogo = async (tipo) => {
    const c = catalogoConfig[tipo];
    setCatalogoLoading(true);
    try {
      const res = await fetch(`${API}/${c.endpoint}`, { headers: { Authorization: `Bearer ${obtenerToken()}` } });
      const data = await res.json();
      setCatalogoItems(Array.isArray(data) ? data : []);
    } catch (err) { console.error(err); setCatalogoItems([]); }
    finally { setCatalogoLoading(false); }
  };

  useEffect(() => { if (activeTab === "catalogos") loadCatalogo(catalogoActivo); }, [activeTab, catalogoActivo]);
  useEffect(() => {
    if (activeTab !== "catalogos") return;
    const headers = { Authorization: `Bearer ${obtenerToken()}` };
    if (catalogoActivo === "reglas-monto") fetch(`${API}/monedas`, { headers }).then(r => r.ok ? r.json() : []).then(d => setMonedasCatalogo(Array.isArray(d) ? d : [])).catch(() => setMonedasCatalogo([]));
    if (catalogoActivo === "metodos-pago") fetch(`${API}/tipos-pago`, { headers }).then(r => r.json()).then(d => setTiposPagoCatalogo(Array.isArray(d) ? d : []));
    if (catalogoActivo === "rutas" || catalogoActivo === "pasos-ruta") fetch(`${API}/areas`, { headers }).then(r => r.json()).then(d => setAreasCatalogo(Array.isArray(d) ? d : []));
    if (catalogoActivo === "pasos-ruta") fetch(`${API}/rutas`, { headers }).then(r => r.json()).then(d => setRutasCatalogo(Array.isArray(d) ? d : [])).catch(err => console.error(err));
    if (catalogoActivo === "pasos-ruta" || catalogoActivo === "reglas-monto") fetch(`${API}/usuarios`, { headers }).then(r => r.json()).then(d => setUsuariosCatalogo(Array.isArray(d) ? d : []));
  }, [activeTab, catalogoActivo, obtenerToken]);

  const openCatalogoCreate = () => {
    setCatalogoEditing(null);
    const empty = {};
    cfg.fields.forEach(f => empty[f] = f === "orden" ? 1 : "");
    setCatalogoForm(empty);
    setShowCatalogoForm(true);
  };

  const openCatalogoEdit = (item) => {
    setCatalogoEditing(item);
    const values = {};
    cfg.fields.forEach(f => { values[f] = item[f] !== undefined ? String(item[f]) : ""; });
    setCatalogoForm(values);
    setShowCatalogoForm(true);
  };

  const handleCatalogoFormChange = (e) => {
    const { name, value, type } = e.target;
    setCatalogoForm(prev => ({ ...prev, [name]: type === "number" ? (value === "" ? 0 : parseInt(value)) : value }));
  };

  const handleCatalogoSubmit = async () => {
    const isEdit = !!catalogoEditing;
    const url = isEdit ? `${API}/${cfg.endpoint}/${catalogoEditing.id}` : `${API}/${cfg.endpoint}`;
    const method = isEdit ? cfg.method : "POST";
    const body = {};
    cfg.fields.forEach(f => {
      if (f === "orden" || f.includes("_id") || f === "monto_minimo") body[f] = parseFloat(catalogoForm[f]) || 0;
      else body[f] = catalogoForm[f]?.trim() || "";
    });
    if (!body.nombre && cfg.fields.includes("nombre")) { alert("El nombre es obligatorio"); return; }
    if (cfg.fields.includes("area_id") && !body.area_id) { alert("Debe seleccionar un área"); return; }
    if (cfg.fields.includes("ruta_id") && !body.ruta_id) { alert("Debe seleccionar una ruta"); return; }
    if (cfg.fields.includes("usuario_id") && !body.usuario_id) { alert("Debe seleccionar un usuario"); return; }

    try {
      const res = await fetch(url, { method, headers: { "Content-Type": "application/json", Authorization: `Bearer ${obtenerToken()}` }, body: JSON.stringify(body) });
      if (!res.ok) { const errData = await res.json(); throw new Error(errData.error || "Error guardando"); }
      setShowCatalogoForm(false); loadCatalogo(catalogoActivo);
    } catch (err) { alert("Error: " + err.message); }
  };

  const handleToggleCatalogoStatus = async (item) => {
    try {
      const res = await fetch(`${API}/${cfg.endpoint}/${item.id}/activo`, {
        method: "PATCH", headers: { "Content-Type": "application/json", Authorization: `Bearer ${obtenerToken()}` },
        body: JSON.stringify({ activo: !item.activo }),
      });
      if (!res.ok) throw new Error("Error cambiando estado");
      loadCatalogo(catalogoActivo);
    } catch (err) { alert(err.message); }
  };

  return {
    catalogoActivo, setCatalogoActivo, catalogoItems, catalogoLoading,
    showCatalogoForm, setShowCatalogoForm, catalogoEditing, catalogoForm,
    tiposPagoCatalogo, areasCatalogo, rutasCatalogo, usuariosCatalogo, monedasCatalogo,
    cfg, openCatalogoCreate, openCatalogoEdit, handleCatalogoFormChange,
    handleCatalogoSubmit, handleToggleCatalogoStatus
  };
}