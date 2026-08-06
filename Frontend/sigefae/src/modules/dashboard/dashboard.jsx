import { useState, useEffect } from "react";
import "./dashboard.css";
import "@fortawesome/fontawesome-free/css/all.min.css";
import logo from "../../assets/login/logo.png";
import { obtenerToken } from "../auth/token.js";
import PdfEditor from "../../components/PdfEditor";
import NotificacionesDropdown from "../../components/NotificacionesDropdown";
import { generarExpedientePDF } from "../../utils/expedientePdf";

const API = "http://localhost:8080/api";

// ── Helpers para leer datos del usuario logueado ──
const obtenerRol = () => localStorage.getItem("rol") || "";
const obtenerUserId = () => parseInt(localStorage.getItem("user_id")) || 0;

export default function ProcesosLogistica() {
  const [isSidebarOpen, setIsSidebarOpen] = useState(false);
  const userRol = obtenerRol();
  const userId = obtenerUserId();
   const esAdmin = userRol === "Superadministrador";
  const esUsuario = !esAdmin;
 
  
  const handleSubirAnexo = async (e, radicadoId) => {
  const file = e.target.files[0];
  if (!file) return;

  const formData = new FormData();
  formData.append("file", file);

    try {
      const res = await fetch(`${API}/documentoradicado/${radicadoId}/anexos`, {
        method: "POST",
        headers: { Authorization: `Bearer ${obtenerToken()}` },
        body: formData,
      });

      if (!res.ok) {
        const err = await res.json();
        throw new Error(err.error || "Error subiendo archivo");
      }

      alert("Archivo subido correctamente");

      // Refrescar según la pestaña activa
      if (activeTab === "tareas") {
        setSelectedTareaId(null);
        setTimeout(() => setSelectedTareaId(radicadoId), 10);
      } else if (activeTab === "radicados") {
        setSelectedRadicadoId(null);
        setTimeout(() => setSelectedRadicadoId(radicadoId), 10);
      }
    } catch (err) {
      alert("Error: " + err.message);
    }

    e.target.value = "";
  };

  // Si es aprobador, arranca directo en "tareas"
  const [activeTab, setActiveTab] = useState(esAdmin? "welcome" : "tareas");

  //pdf
  const [pdfEditor, setPdfEditor] = useState({ open: false, archivoId: null, radicadoId: null });

    // ── Crear Documento Manual ──
  const [showCrearDocModal, setShowCrearDocModal] = useState(false);
  const [docForm, setDocForm] = useState({
    tipo: "FACTURA_FISICA",
    numero_documento: "",
    id_proveedor: "",
    id_receptor: "",
    id_area: "",
    tipo_factura_id: "",
    asunto: "",
    fecha_documento: new Date().toISOString().split("T")[0],
    moneda_id: "",
    subtotal: 0,
    iva: 0,
    total: 0,
    detalles: []
  });
  const [proveedoresCatalogo, setProveedoresCatalogo] = useState([]);
  const [receptoresCatalogo, setReceptoresCatalogo] = useState([]);
  const [creandoDoc, setCreandoDoc] = useState(false);

  //normas de reparto
    // ── Filtros jerárquicos para normas en el modal de radicación ──
  const [normaFiltroSede, setNormaFiltroSede] = useState("");
  const [normaFiltroArea, setNormaFiltroArea] = useState("");
  const [normaSeleccionadaId, setNormaSeleccionadaId] = useState("");
  const [normaPorcentajeInput, setNormaPorcentajeInput] = useState("");
  const [normasRepartoCatalogo, setNormasRepartoCatalogo] = useState([]);
  const [normasRepartoRadicado, setNormasRepartoRadicado] = useState([]);

  // ── Catálogos Admin ──
  const [catalogoActivo, setCatalogoActivo] = useState("tipo-radicacion");
  const [catalogoItems, setCatalogoItems] = useState([]);
  const [catalogoLoading, setCatalogoLoading] = useState(false);
  const [showCatalogoForm, setShowCatalogoForm] = useState(false);
  const [catalogoEditing, setCatalogoEditing] = useState(null);
  const [catalogoForm, setCatalogoForm] = useState({});

  // Catálogos auxiliares para dropdowns
  const [tiposPagoCatalogo, setTiposPagoCatalogo] = useState([]);
  const [areasCatalogo, setAreasCatalogo] = useState([]);
  const [rutasCatalogo, setRutasCatalogo] = useState([]);
  const [usuariosCatalogo, setUsuariosCatalogo] = useState([]);

  // ── Correos ──
  const [correos, setCorreos] = useState([]);
  const [selectedCorreoId, setSelectedCorreoId] = useState(null);
  const [correoDetail, setCorreoDetail] = useState(null);
  const [loading, setLoading] = useState(false);

  // ── Documentos Pendientes ──
  const [documentos, setDocumentos] = useState([]);
  const [selectedDocId, setSelectedDocId] = useState(null);
  const [docDetail, setDocDetail] = useState(null);
  const [loadingDocs, setLoadingDocs] = useState(false);

  // ── Catálogos ──
  const [tiposRadicacion, setTiposRadicacion] = useState([]);
  const [rutas, setRutas] = useState([]);
  const [metodosPago, setMetodosPago] = useState([]);
  const [monedasCatalogo, setMonedasCatalogo] = useState([]);

  // ── Modal de Radicación ──
  const [showRadicarModal, setShowRadicarModal] = useState(false);
  const [radicarDocId, setRadicarDocId] = useState(null);
  const [radicarForm, setRadicarForm] = useState({
    tipo_radicacion_id: "",
    ruta_id: "",
    metodo_pago_id: "",
    numero_radicado: "",
    normas_reparto: [],
  });
  const [radicando, setRadicando] = useState(false);

  // ── Radicados ──
  const [radicados, setRadicados] = useState([]);
  const [selectedRadicadoId, setSelectedRadicadoId] = useState(null);
  const [radicadoDetail, setRadicadoDetail] = useState(null);
  const [loadingRadicados, setLoadingRadicados] = useState(false);

  // ── Mis Tareas (filtradas para el aprobador) ──
  const [misTareas, setMisTareas] = useState([]);
  const [selectedTareaId, setSelectedTareaId] = useState(null);
  const [tareaDetail, setTareaDetail] = useState(null);
  const [loadingTareas, setLoadingTareas] = useState(false);
  const [tareasSubTab, setTareasSubTab] = useState("activas"); // "activas" | "completadas"
  const [misTareasCompletadas, setMisTareasCompletadas] = useState([]);

    // ── Flujo de tareas del radicado ──
  const [tareasFlujo, setTareasFlujo] = useState([]);

  // ── Devolución de Tareas ──
  const [showDevolverModal, setShowDevolverModal] = useState(false);
  const [devolverForm, setDevolverForm] = useState({
    tarea_destino_id: "",
    observacion: "",
    retorno_directo: true
  });
  const [devolviendo, setDevolviendo] = useState(false);

  const [comentarios, setComentarios] = useState([]);
  const [nuevoComentario, setNuevoComentario] = useState("");
  const [enviandoComentario, setEnviandoComentario] = useState(false);

  // ── Trazabilidad ──
  const [historialTrazabilidad, setHistorialTrazabilidad] = useState([]);

  // ── Expediente PDF ──
  const [generandoPdf, setGenerandoPdf] = useState(false);
  
  
  const totalPorcentajeNormas = (radicarForm.normas_reparto || []).reduce(
    (sum, n) => sum + (parseFloat(n.porcentaje) || 0), 0
  );


  // Cargar tareas e historial de trazabilidad cuando se selecciona un radicado o tarea
  useEffect(() => {
    const radicadoId = selectedRadicadoId || selectedTareaId;
    if (radicadoId) {
      // Cargar tareas
      fetch(`${API}/documentoradicado/${radicadoId}/tareas`, {
        headers: { Authorization: `Bearer ${obtenerToken()}` }
      })
        .then(r => r.json())
        .then(data => setTareasFlujo(Array.isArray(data) ? data : []))
        .catch(err => console.error(err));

      // Cargar trazabilidad
      fetch(`${API}/trazabilidad?documento_radicado_id=${radicadoId}`, {
        headers: { Authorization: `Bearer ${obtenerToken()}` }
      })
        .then(r => r.json())
        .then(data => setHistorialTrazabilidad(Array.isArray(data) ? data : []))
        .catch(err => console.error(err));
            // ── CARGAR NORMAS DE REPARTO ──
      fetch(`${API}/documentoradicado/${radicadoId}/normas-reparto`, {
        headers: { Authorization: `Bearer ${obtenerToken()}` }
      })
        .then(r => r.json())
        .then(data => setNormasRepartoRadicado(Array.isArray(data) ? data : []))
        .catch(() => setNormasRepartoRadicado([]));
      }
     else {
      setTareasFlujo([]);
      setHistorialTrazabilidad([]);
      setNormasRepartoRadicado([]); 
    }
  }, [selectedRadicadoId, selectedTareaId]);

  const handleDescargarExpediente = async (radicado) => {
    if (!radicado) return;
    try {
      setGenerandoPdf(true);
      // Collect PDF urls from archivos
      const anexosUrls = (radicado.archivos || [])
        .filter(a => a.extension?.toLowerCase() === 'pdf' || a.nombre?.toLowerCase().endsWith('.pdf'))
        .map(a => `${API}/archivo/${a.id}/download?download=1`);
      
      await generarExpedientePDF(radicado, tareasFlujo, historialTrazabilidad, anexosUrls);
    } catch (err) {
      console.error("Error al generar expediente", err);
      alert("Hubo un error al generar el expediente PDF.");
    } finally {
      setGenerandoPdf(false);
    }
  };

    // Cargar catálogos necesarios para crear documento manual
  useEffect(() => {
    if (activeTab !== "documentos") return;
    const headers = { Authorization: `Bearer ${obtenerToken()}` };

    fetch(`${API}/proveedor`, { headers })
      .then(r => r.json())
      .then(d => setProveedoresCatalogo(Array.isArray(d) ? d : []))
      .catch(() => setProveedoresCatalogo([]));

    fetch(`${API}/receptor`, { headers })
      .then(r => r.json())
      .then(d => setReceptoresCatalogo(Array.isArray(d) ? d : []))
      .catch(() => setReceptoresCatalogo([]));

    if (areasCatalogo.length === 0) {
      fetch(`${API}/areas`, { headers })
        .then(r => r.json())
        .then(d => setAreasCatalogo(Array.isArray(d) ? d : []))
        .catch(() => {});
    }
    if (monedasCatalogo.length === 0) {
      fetch(`${API}/monedas`, { headers })
        .then(r => r.json())
        .then(d => setMonedasCatalogo(Array.isArray(d) ? d : []))
        .catch(() => {});
    }
  }, [activeTab]);

    // ── Cargar comentarios del radicado ──
  useEffect(() => {
    const radicadoId = activeTab === "tareas" ? selectedTareaId :
                       activeTab === "radicados" ? selectedRadicadoId : null;
    if (!radicadoId) {
      setComentarios([]);
      return;
    }

    setComentarios([]);
    let cancelled = false;

    fetch(`${API}/comentario?documento_radicado_id=${radicadoId}`, {
      headers: { Authorization: `Bearer ${obtenerToken()}` }
    })
      .then(r => r.json())
      .then(data => {
        if (!cancelled) setComentarios(Array.isArray(data) ? data : []);
      })
      .catch(err => {
        if (!cancelled) {
          console.error("Error cargando comentarios:", err);
          setComentarios([]);
        }
      });

    return () => { cancelled = true; };
  }, [activeTab, selectedRadicadoId, selectedTareaId]);


    const openCrearDocModal = () => {
    setSelectedDocId(null);
    setDocDetail(null);
    setDocForm({
      tipo: "FACTURA_FISICA",
      numero_documento: "",
      orden_compra: "",
      id_proveedor: "",
      id_receptor: "",
      id_area: "",
      tipo_factura_id: "",
      asunto: "",
      fecha_documento: new Date().toISOString().split("T")[0],
      fecha_vencimiento: "",
      moneda_id: "",
      subtotal: 0,
      iva: 0,
      total: 0,
      numero_folios: 1,
      detalles: []
    });
    setShowCrearDocModal(true);
  };

  const recalcularTotales = (detalles) => {
    const subtotal = detalles.reduce((s, d) => s + (parseFloat(d.valor_unitario) || 0) * (parseFloat(d.cantidad) || 0), 0);
    const iva = detalles.reduce((s, d) => s + (parseFloat(d.iva_unitario) || 0) * (parseFloat(d.cantidad) || 0), 0);
    return { subtotal, iva, total: subtotal + iva };
  };

  const handleDocFormChange = (e) => {
    const { name, value } = e.target;
    setDocForm(prev => ({ ...prev, [name]: value }));
  };

  const handleAddDetalle = () => {
    setDocForm(prev => {
      const detalles = [...prev.detalles, { descripcion: "", cantidad: 1, valor_unitario: 0, iva_unitario: 0, total: 0 }];
      const t = recalcularTotales(detalles);
      return { ...prev, detalles, ...t };
    });
  };

  const handleDetalleChange = (idx, field, value) => {
    setDocForm(prev => {
      const detalles = prev.detalles.map((d, i) => {
        if (i !== idx) return d;
        const upd = { ...d, [field]: field === "descripcion" ? value : parseFloat(value) || 0 };
        upd.total = ((upd.valor_unitario || 0) + (upd.iva_unitario || 0)) * (upd.cantidad || 0);
        return upd;
      });
      const t = recalcularTotales(detalles);
      return { ...prev, detalles, ...t };
    });
  };

  const handleRemoveDetalle = (idx) => {
    setDocForm(prev => {
      const detalles = prev.detalles.filter((_, i) => i !== idx);
      const t = recalcularTotales(detalles);
      return { ...prev, detalles, ...t };
    });
  };

  const handleCrearDocumentoSubmit = async () => {
    if (!docForm.numero_documento || !docForm.id_proveedor || !docForm.id_receptor || !docForm.id_area || !docForm.moneda_id || !docForm.fecha_documento) {
      alert("Complete todos los campos obligatorios (*)");
      return;
    }
    if (docForm.detalles.length === 0) {
      alert("Agregue al menos un ítem al documento");
      return;
    }

    setCreandoDoc(true);
    const payload = {
      tipo: docForm.tipo,
      numero_documento: docForm.numero_documento.trim(),
      orden_compra: "",
      id_proveedor: parseInt(docForm.id_proveedor),
      id_receptor: parseInt(docForm.id_receptor),
      id_area: parseInt(docForm.id_area),
      tipo_factura_id: docForm.tipo_factura_id ? parseInt(docForm.tipo_factura_id) : null,
      asunto: docForm.asunto.trim(),
      fecha_documento: docForm.fecha_documento + "T00:00:00Z",
      fecha_vencimiento: null,
      moneda_id: parseInt(docForm.moneda_id),
      subtotal: parseFloat(docForm.subtotal) || 0,
      iva: parseFloat(docForm.iva) || 0,
      total: parseFloat(docForm.total) || 0,
      numero_folios: 1,
      orientacion_sello_recibido: "",
      activo: true,
      detalles: docForm.detalles.map(d => ({
        descripcion: d.descripcion,
        cantidad: parseFloat(d.cantidad) || 0,
        valor_unitario: parseFloat(d.valor_unitario) || 0,
        iva_unitario: parseFloat(d.iva_unitario) || 0,
        total: parseFloat(d.total) || 0
      }))
    };

    try {
      const res = await fetch(`${API}/documentocomercial`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${obtenerToken()}`
        },
        body: JSON.stringify(payload)
      });
      if (!res.ok) {
        const err = await res.json();
        throw new Error(err.error || "Error creando documento");
      }
      alert("Documento creado correctamente");
      setShowCrearDocModal(false);
      // Refrescar lista de pendientes
      setLoadingDocs(true);
      fetch(`${API}/documentocomercial/pendientes`, { headers: { Authorization: `Bearer ${obtenerToken()}` } })
        .then(r => r.json())
        .then(data => { if (Array.isArray(data)) setDocumentos(data); })
        .finally(() => setLoadingDocs(false));
    } catch (err) {
      alert("Error: " + err.message);
    } finally {
      setCreandoDoc(false);
    }
  };

  const renderFlujoAprobacion = () => (
    <div className="doc-section">
      <h4><i className="fa-solid fa-route"></i> Flujo de Aprobación</h4>
      {tareasFlujo.length === 0 ? (
        <p style={{ color: "#6b7280", fontSize: "0.9em" }}>No hay pasos definidos para este flujo.</p>
      ) : (
        <table className="doc-items-table">
          <thead>
            <tr>
              <th>Usuario</th>
              <th>Correo</th>
              <th>Acción a realizar</th>
              <th>Estado</th>
              <th>Fecha Realización</th>
            </tr>
          </thead>
          <tbody>
            {tareasFlujo.map((t) => (
              <tr key={t.id} style={t.estado?.nombre === "En Proceso" ? { background: "rgba(44,82,130,0.05)" } : {}}>
                <td>{t.usuario_asignado?.nombre || "—"}</td>
                <td>{t.usuario_asignado?.email || "—"}</td>
                <td>{t.descripcion}</td>
                <td>
                  <span className={`status-badge ${
                    t.estado?.nombre === "Completada" ? "radicado" :
                    t.estado?.nombre === "En Proceso" ? "doc-pendiente" :
                    t.estado?.nombre === "Devuelta" ? "doc-rechazado" : ""
                  }`}>
                    {t.estado?.nombre || "Pendiente"}
                  </span>
                </td>
                <td>{t.fecha_finalizacion ? new Date(t.fecha_finalizacion).toLocaleString() : "—"}</td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </div>
  );

  //normas para modal normas de reparto
    const sedesDisponibles = ["BUCARAMANGA", "MALAMBO", "CUCUTA", "CB", "CIENAGA DE ORO", "GENERAL"];
  const areasDisponibles = ["ADMON", "VENTAS", "PRODUCCION"];

  const normasFiltradas = normasRepartoCatalogo.filter(n => {
    if (normaFiltroSede && n.sucursal !== normaFiltroSede) return false;
    if (normaFiltroArea && n.departamento !== normaFiltroArea) return false;
    return true;
  });

  const handleAgregarNormaModal = () => {
    if (!normaSeleccionadaId || !normaPorcentajeInput) {
      alert("Selecciona una norma y escribe el porcentaje");
      return;
    }
    const pct = parseFloat(normaPorcentajeInput);
    if (isNaN(pct) || pct <= 0) {
      alert("Porcentaje inválido");
      return;
    }
    const yaExiste = radicarForm.normas_reparto.find(n => n.norma_reparto_id === normaSeleccionadaId);
    if (yaExiste) {
      alert("Esta norma ya fue agregada");
      return;
    }
    setRadicarForm(prev => ({
      ...prev,
      normas_reparto: [...prev.normas_reparto, {
        norma_reparto_id: normaSeleccionadaId,
        porcentaje: normaPorcentajeInput
      }]
    }));
    setNormaSeleccionadaId("");
    setNormaPorcentajeInput("");
  };
    const handleNormaRepartoChange = (index, field, value) => {
    setRadicarForm(prev => {
      const updated = [...prev.normas_reparto];
      updated[index] = { ...updated[index], [field]: value };
      return { ...prev, normas_reparto: updated };
    });
  };

  const handleRemoveNormaReparto = (index) => {
    setRadicarForm(prev => ({
      ...prev,
      normas_reparto: prev.normas_reparto.filter((_, i) => i !== index)
    }));
  };
  // ── Render Trazabilidad ──
  const renderTrazabilidad = () => (
    <div className="doc-section">
      <h4>
        <i className="fa-solid fa-clock-rotate-left"></i>{" "}
        Historial de Trazabilidad ({historialTrazabilidad.length})
      </h4>
      <div
        style={{
          maxHeight: 250,
          overflowY: "auto",
          marginBottom: 12,
          paddingRight: 4
        }}
      >
        {historialTrazabilidad.length === 0 ? (
          <p style={{ color: "#6b7280", fontSize: "0.9em" }}>
            No hay historial disponible.
          </p>
        ) : (
          <div style={{ display: "flex", flexDirection: "column", gap: 8 }}>
            {historialTrazabilidad.map((t, idx) => (
              <div
                key={idx}
                style={{
                  padding: "10px",
                  background: "#f9fafb",
                  borderLeft: "3px solid var(--pardo-red)",
                  borderRadius: 4
                }}
              >
                <div style={{ display: "flex", justifyContent: "space-between", marginBottom: 4 }}>
                  <strong style={{ color: "var(--pardo-navy)", fontSize: "0.9em" }}>{t.accion}</strong>
                  <span style={{ fontSize: "0.8em", color: "#6b7280" }}>
                    {new Date(t.fecha).toLocaleString()}
                  </span>
                </div>
                <div style={{ fontSize: "0.85em", color: "#4b5563", marginBottom: 4 }}>
                  {t.descripcion}
                </div>
                <div style={{ fontSize: "0.8em", color: "#9ca3af", textAlign: "right" }}>
                  <i className="fa-solid fa-user" style={{ marginRight: 4 }}></i>
                  {t.usuario_nombre || "Sistema"}
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );

  // ── Render comentarios ──
  const renderComentarios = (radicadoId, readOnly = false) => (
  <div className="doc-section">
    <h4>
      <i className="fa-solid fa-comments"></i>{" "}
      Comentarios ({comentarios.length})
    </h4>

    <div
      style={{
        maxHeight: 320,
        overflowY: "auto",
        marginBottom: 12,
        paddingRight: 4
      }}
    >
      {comentarios.length === 0 ? (
        <p style={{ color: "#6b7280", fontSize: "0.9em" }}>
          No hay comentarios.
        </p>
      ) : (
        comentarios.map((c) => (
          <div
            key={c.id}
            style={{
              background: "#f3f4f6",
              borderRadius: 8,
              padding: "10px 12px",
              marginBottom: 8
            }}
          >
            <div
              style={{
                display: "flex",
                justifyContent: "space-between",
                alignItems: "center",
                marginBottom: 4
              }}
            >
              <strong style={{ fontSize: "0.85em", color: "#1f2937" }}>
                {c.usuario_nombre || `Usuario #${c.usuario_id}`}
              </strong>

              <span style={{ fontSize: "0.75em", color: "#6b7280" }}>
                {new Date(c.fecha).toLocaleString()}
              </span>
            </div>

            <p
              style={{
                margin: 0,
                fontSize: "0.9em",
                color: "#374151",
                whiteSpace: "pre-wrap"
              }}
            >
              {c.descripcion}
            </p>
          </div>
        ))
      )}
    </div>

    {readOnly ? (
      <div
        style={{
          padding: "10px 12px",
          borderRadius: 6,
          background: "#f3f4f6",
          color: "#6b7280",
          fontSize: "0.9em"
        }}
      >
        <i className="fa-solid fa-lock"></i>{" "}
        El documento está finalizado. No se pueden agregar más comentarios.
      </div>
    ) : (
      <div style={{ display: "flex", gap: 8, alignItems: "flex-start" }}>
        <textarea
          value={nuevoComentario}
          onChange={(e) => setNuevoComentario(e.target.value)}
          placeholder="Escribe un comentario..."
          rows={2}
          style={{
            flex: 1,
            padding: 8,
            borderRadius: 6,
            border: "1px solid #d1d5db",
            resize: "vertical",
            fontFamily: "inherit",
            fontSize: "0.9em"
          }}
        />

        <button
          className="doc-btn doc-btn-primary"
          onClick={() => handleEnviarComentario(radicadoId)}
          disabled={enviandoComentario || !nuevoComentario.trim()}
          style={{ marginTop: 2 }}
        >
          <i className="fa-solid fa-paper-plane"></i>{" "}
          {enviandoComentario ? "Enviando..." : "Enviar"}
        </button>
      </div>
    )}
  </div>
);

  // ── Completar tarea ──

    // ── Completar tarea ──
  const handleCompletarTarea = async (tareaId, radicadoId) => {
    try {
      const res = await fetch(`${API}/tarea/${tareaId}/completar`, {
        method: "PATCH",
        headers: { Authorization: `Bearer ${obtenerToken()}` }
      });
      if (!res.ok) {
        const errData = await res.json();
        throw new Error(errData.error || "Error completando tarea");
      }

      alert("Tarea completada correctamente");

      const flujoRes = await fetch(`${API}/documentoradicado/${radicadoId}/tareas`, {
        headers: { Authorization: `Bearer ${obtenerToken()}` }
      });
      const flujoData = await flujoRes.json();
      setTareasFlujo(Array.isArray(flujoData) ? flujoData : []);

      const detalleRes = await fetch(`${API}/documentoradicado/${radicadoId}`, {
        headers: { Authorization: `Bearer ${obtenerToken()}` }
      });
      const detalleData = await detalleRes.json();
      if (detalleData && detalleData.id) {
        if (activeTab === "tareas") setTareaDetail(detalleData);
        if (activeTab === "radicados") setRadicadoDetail(detalleData);
      }

      if (activeTab === "tareas") {
        const listaRes = await fetch(`${API}/documentoradicado`, {
          headers: { Authorization: `Bearer ${obtenerToken()}` }
        });
        const listaData = await listaRes.json();
        if (Array.isArray(listaData)) {
          setMisTareas(listaData.filter(r => r.usuario_actual_id === userId && r.estado_posesion !== "Completado"));
        }
      }

      // ── NOTIFICAR al siguiente usuario en el flujo ──
      const siguiente = flujoData.find(t => t.estado?.nombre === "En Proceso");
      if (siguiente?.usuario_asignado?.id && siguiente.usuario_asignado.id !== userId) {
        await fetch(`${API}/notificacion`, {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${obtenerToken()}`
          },
          body: JSON.stringify({
            usuario_id: siguiente.usuario_asignado.id,
            documento_radicado_id: radicadoId,
            mensaje: `Te asignaron el radicado #${detalleData?.numero_radicado || radicadoId} — Paso: ${siguiente.descripcion || 'revisar'}`,
            estado: "Pendiente",
            tipo: "Asignacion",
            fecha_creacion: new Date().toISOString(),
          })
        });
      }

    } catch (err) {
      alert("Error: " + err.message);
    }
  };

  
  
      // ── Devolver Tarea ──
const handleDevolverTarea = async (tareaId, radicadoId) => {
  if (!devolverForm.tarea_destino_id || !devolverForm.observacion.trim()) {
    alert("Debe seleccionar un paso de destino y escribir un motivo de devolución.");
    return;
  }
  setDevolviendo(true);
  try {
    const res = await fetch(`${API}/tarea/${tareaId}/devolver`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${obtenerToken()}`
      },
      body: JSON.stringify({
        tarea_destino_id: parseInt(devolverForm.tarea_destino_id),
        observacion: devolverForm.observacion.trim(),
        retorno_directo: devolverForm.retorno_directo
      })
    });

    if (!res.ok) {
      const errData = await res.json();
      throw new Error(errData.error || "Error al devolver la tarea");
    }

    // 1. Cerrar modal y limpiar formulario
    setShowDevolverModal(false);
    setDevolverForm({ tarea_destino_id: "", observacion: "", retorno_directo: true });

    // 2. Limpiar selección para que el usuario ya no vea el detalle del documento devuelto
    if (activeTab === "tareas") {
      setSelectedTareaId(null);
      setTareaDetail(null);
    } else {
      setSelectedRadicadoId(null);
      setRadicadoDetail(null);
    }
    setTareasFlujo([]);
    setHistorialTrazabilidad([]);
    setComentarios([]);

    // 3. Recargar la lista según la pestaña activa
    const listaRes = await fetch(`${API}/documentoradicado`, {
      headers: { Authorization: `Bearer ${obtenerToken()}` }
    });
    const listaData = await listaRes.json();

    if (Array.isArray(listaData)) {
      if (activeTab === "tareas") {
        const activas = listaData.filter(r => r.usuario_actual_id === userId && r.estado_posesion !== "Completado");
        const completadas = listaData.filter(r => r.estado_posesion === "Completado");
        setMisTareas(activas);
        setMisTareasCompletadas(completadas);
      } else if (activeTab === "radicados") {
        setRadicados(listaData);
      }
    }

    alert("Tarea devuelta correctamente");

  } catch (err) {
    alert("Error: " + err.message);
  } finally {
    setDevolviendo(false);
  }
};

  // ── Enviar comentario ──
    const handleEnviarComentario = async (radicadoId) => {
    if (!nuevoComentario.trim()) return;
    setEnviandoComentario(true);
    try {
      const res = await fetch(`${API}/comentario`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${obtenerToken()}`
        },
        body: JSON.stringify({
          documento_radicado_id: radicadoId,
          usuario_id: userId,
          descripcion: nuevoComentario.trim()
        })
      });
      if (!res.ok) {
        const errData = await res.json();
        throw new Error(errData.error || "Error enviando comentario");
      }

      setNuevoComentario("");
      const listRes = await fetch(`${API}/comentario?documento_radicado_id=${radicadoId}`, {
        headers: { Authorization: `Bearer ${obtenerToken()}` }
      });
      const listData = await listRes.json();
      setComentarios(Array.isArray(listData) ? listData : []);


    } catch (err) {
      alert("Error: " + err.message);
    } finally {
      setEnviandoComentario(false);
    }
  };

  const catalogoConfig = {
    "tipo-radicacion": { endpoint: "tipo-radicacion", label: "Tipo de Radicación", method: "PUT", fields: ["nombre"] },
    "tipos-pago":      { endpoint: "tipos-pago",      label: "Tipo de Pago",      method: "PATCH", fields: ["nombre"] },
    "metodos-pago":    { endpoint: "metodos-pago",    label: "Método de Pago",    method: "PATCH", fields: ["nombre", "tipo_pago_id"] },
    "areas":           { endpoint: "areas",           label: "Área",              method: "PATCH", fields: ["nombre"] },
    "rutas":           { endpoint: "rutas",           label: "Ruta",              method: "PUT",   fields: ["nombre", "area_id"] },
    "pasos-ruta":      { endpoint: "pasos-ruta",      label: "Paso de Ruta",      method: "PUT",   fields: ["ruta_id", "orden", "nombre", "usuario_id"] },
    "reglas-monto":    { endpoint: "regla-monto-ruta",label: "Regla de Monto", method: "PUT", fields: ["monto_minimo", "moneda_id", "posicion_insercion", "usuario_aprobador_id"] },
    "normas-reparto":  {endpoint: "normas-reparto",label: "Norma de Reparto",method: "PUT",fields: ["codigo", "nombre", "sucursal", "departamento", "tipo", "tarifa_iva"]},
  };

  const loadCatalogo = async (tipo) => {
    const cfg = catalogoConfig[tipo];
    setCatalogoLoading(true);
    try {
      const res = await fetch(`${API}/${cfg.endpoint}`, {
        headers: { Authorization: `Bearer ${obtenerToken()}` },
      });
      const data = await res.json();
      setCatalogoItems(Array.isArray(data) ? data : []);
    } catch (err) {
      console.error(err);
      setCatalogoItems([]);
    } finally {
      setCatalogoLoading(false);
    }
  };

  const openCatalogoCreate = () => {
    setCatalogoEditing(null);
    const empty = {};
    catalogoConfig[catalogoActivo].fields.forEach(f => empty[f] = f === "orden" ? 1 : "");
    setCatalogoForm(empty);
    setShowCatalogoForm(true);
  };

  const openCatalogoEdit = (item) => {
    setCatalogoEditing(item);
    const values = {};
    catalogoConfig[catalogoActivo].fields.forEach(f => {
      values[f] = item[f] !== undefined ? String(item[f]) : "";
    });
    setCatalogoForm(values);
    setShowCatalogoForm(true);
  };

  const handleCatalogoFormChange = (e) => {
    const { name, value, type } = e.target;
    setCatalogoForm(prev => ({
      ...prev,
      [name]: type === "number" ? (value === "" ? 0 : parseInt(value)) : value
    }));
  };

  const handleCatalogoSubmit = async () => {
    const cfg = catalogoConfig[catalogoActivo];
    const isEdit = !!catalogoEditing;
    const url = isEdit ? `${API}/${cfg.endpoint}/${catalogoEditing.id}` : `${API}/${cfg.endpoint}`;
    const method = isEdit ? cfg.method : "POST";

    const body = {};
    cfg.fields.forEach(f => {
      if (f === "orden" || f.includes("_id") || f === "monto_minimo") {
        body[f] = parseFloat(catalogoForm[f]) || 0;
      } else {
        body[f] = catalogoForm[f]?.trim() || "";
      }
    });
    if (!body.nombre && cfg.fields.includes("nombre")) {
      alert("El nombre es obligatorio");
      return;
    }
    if (cfg.fields.includes("area_id") && !body.area_id) {
      alert("Debe seleccionar un área");
      return;
    }
    if (cfg.fields.includes("ruta_id") && !body.ruta_id) {
      alert("Debe seleccionar una ruta");
      return;
    }
    if (cfg.fields.includes("usuario_id") && !body.usuario_id) {
      alert("Debe seleccionar un usuario");
      return;
    }

    try {
      const res = await fetch(url, {
        method,
        headers: { "Content-Type": "application/json", Authorization: `Bearer ${obtenerToken()}` },
        body: JSON.stringify(body),
      });
      if (!res.ok) {
        const errData = await res.json();
        throw new Error(errData.error || "Error guardando");
      }
      setShowCatalogoForm(false);
      loadCatalogo(catalogoActivo);
    } catch (err) {
      alert("Error: " + err.message);
    }
  };

  const handleToggleCatalogoStatus = async (item) => {
    const cfg = catalogoConfig[catalogoActivo];
    try {
      const res = await fetch(`${API}/${cfg.endpoint}/${item.id}/activo`, {
        method: "PATCH",
        headers: { "Content-Type": "application/json", Authorization: `Bearer ${obtenerToken()}` },
        body: JSON.stringify({ activo: !item.activo }),
      });
      if (!res.ok) throw new Error("Error cambiando estado");
      loadCatalogo(catalogoActivo);
    } catch (err) {
      alert(err.message);
    }
  };

  // Cargar lista de correos
  useEffect(() => {
    if (activeTab === "correos") {
      setLoading(true);
      fetch(`${API}/correo`, { headers: { Authorization: `Bearer ${obtenerToken()}` } })
        .then((res) => res.json())
        .then((data) => { if (Array.isArray(data)) setCorreos(data); })
        .catch((err) => console.error(err))
        .finally(() => setLoading(false));
    }
  }, [activeTab]);

  // Cargar catálogo activo
  useEffect(() => {
    if (activeTab === "catalogos") {
      loadCatalogo(catalogoActivo);
    }
  }, [activeTab, catalogoActivo]);

  // Cargar catálogos auxiliares según la pestaña activa
    useEffect(() => {
    if (activeTab !== "catalogos") return;
    const headers = { Authorization: `Bearer ${obtenerToken()}` };

    if (catalogoActivo === "reglas-monto") {
      fetch(`${API}/monedas`, { headers })
        .then(r => r.ok ? r.json() : [])
        .then(d => setMonedasCatalogo(Array.isArray(d) ? d : []))
        .catch(() => setMonedasCatalogo([]));
    }
    
    if (catalogoActivo === "metodos-pago") {
      fetch(`${API}/tipos-pago`, { headers }).then(r => r.json()).then(d => setTiposPagoCatalogo(Array.isArray(d) ? d : []));
    }
    if (catalogoActivo === "rutas" || catalogoActivo === "pasos-ruta") {
      fetch(`${API}/areas`, { headers }).then(r => r.json()).then(d => setAreasCatalogo(Array.isArray(d) ? d : []));
    }
    if (catalogoActivo === "pasos-ruta") {
      fetch(`${API}/rutas`, { headers })
        .then(r => r.json())
        .then(d => setRutasCatalogo(Array.isArray(d) ? d : []))
        .catch(err => console.error("Error cargando rutas:", err));
    }
    if (catalogoActivo === "pasos-ruta" || catalogoActivo === "reglas-monto") {
      fetch(`${API}/usuarios`, { headers }).then(r => r.json()).then(d => setUsuariosCatalogo(Array.isArray(d) ? d : []));
    }
  }, [activeTab, catalogoActivo]);
  
  // Cargar tipos de pago cuando se necesiten (para métodos de pago)
  useEffect(() => {
    if (activeTab === "catalogos" && catalogoActivo === "metodos-pago") {
      fetch(`${API}/tipos-pago`, { headers: { Authorization: `Bearer ${obtenerToken()}` } })
        .then(r => r.json())
        .then(data => setTiposPagoCatalogo(Array.isArray(data) ? data : []))
        .catch(err => console.error(err));
    }
  }, [activeTab, catalogoActivo]);

  // Cargar detalle del correo
  useEffect(() => {
    if (selectedCorreoId) {
      setCorreoDetail(null);
      fetch(`${API}/correo/${selectedCorreoId}`, { headers: { Authorization: `Bearer ${obtenerToken()}` } })
        .then((res) => res.json())
        .then((data) => { if (data && data.id) setCorreoDetail(data); })
        .catch((err) => console.error(err));
    }
  }, [selectedCorreoId]);

  // Cargar documentos pendientes
  useEffect(() => {
    if (activeTab === "documentos") {
      setLoadingDocs(true);
      fetch(`${API}/documentocomercial/pendientes`, { headers: { Authorization: `Bearer ${obtenerToken()}` } })
        .then((res) => res.json())
        .then((data) => { if (Array.isArray(data)) setDocumentos(data); })
        .catch((err) => console.error(err))
        .finally(() => setLoadingDocs(false));
    }
  }, [activeTab]);

    // Cargar detalle del documento
  useEffect(() => {
    if (selectedDocId) {
      setDocDetail(null);
      fetch(`${API}/documentocomercial/${selectedDocId}`, { headers: { Authorization: `Bearer ${obtenerToken()}` } })
        .then((res) => res.json())
        .then((data) => {
          if (data && data.id) {
            setDocDetail(data);
          }
        })
        .catch((err) => console.error(err));
    }
  }, [selectedDocId]);


  // Cargar catálogos cuando entra en modo edición o radicación
  useEffect(() => {
    if (showRadicarModal) {
      const headers = { Authorization: `Bearer ${obtenerToken()}` };
      Promise.all([
        fetch(`${API}/tipo-radicacion`, { headers }).then(r => r.json()),
        fetch(`${API}/rutas`, { headers }).then(r => r.json()),
        fetch(`${API}/metodos-pago`, { headers }).then(r => r.json()),
        fetch(`${API}/normas-reparto?activo=true`, { headers }).then(r => r.json()),
      ])
        .then(([tr, r, mp, nr]) => {
          setTiposRadicacion(Array.isArray(tr) ? tr : []);
          setRutas(Array.isArray(r) ? r : []);
          setMetodosPago(Array.isArray(mp) ? mp : []);
          setNormasRepartoCatalogo(Array.isArray(nr) ? nr : []);
        })
        .catch(err => console.error("Error cargando catálogos:", err));
    }
  }, [showRadicarModal]);

  // Cargar radicados
  useEffect(() => {
    if (activeTab === "radicados") {
      setLoadingRadicados(true);
      fetch(`${API}/documentoradicado`, { headers: { Authorization: `Bearer ${obtenerToken()}` } })
        .then((res) => res.json())
        .then((data) => { if (Array.isArray(data)) setRadicados(data); })
        .catch((err) => console.error(err))
        .finally(() => setLoadingRadicados(false));
    }
  }, [activeTab]);

    // Cargar detalle de radicado (con área completa de la ruta)
  useEffect(() => {
    if (selectedRadicadoId) {
      setRadicadoDetail(null);
      fetch(`${API}/documentoradicado/${selectedRadicadoId}`, { headers: { Authorization: `Bearer ${obtenerToken()}` } })
        .then((res) => res.json())
        .then(async (data) => {
          if (data && data.id) {
            if (data.ruta?.id && !data.ruta?.area) {
              try {
                const rutasRes = await fetch(`${API}/rutas`, { headers: { Authorization: `Bearer ${obtenerToken()}` } });
                const rutasList = await rutasRes.json();
                const rutaCompleta = (Array.isArray(rutasList) ? rutasList : []).find(r => r.id === data.ruta.id);
                if (rutaCompleta) {
                  data.ruta.area = rutaCompleta.area;
                }
              } catch (e) {}
            }
            setRadicadoDetail(data);
          }
        })
        .catch((err) => console.error(err));
    }
  }, [selectedRadicadoId]);

  // ── Cargar MIS TAREAS (filtradas por usuario logueado) ──
      useEffect(() => {
    if (activeTab === "tareas") {
      setLoadingTareas(true);
      fetch(`${API}/documentoradicado`, { headers: { Authorization: `Bearer ${obtenerToken()}` } })
        .then((res) => res.json())
        .then((data) => {
          if (Array.isArray(data)) {
            const activas = data.filter(r => r.usuario_actual_id === userId && r.estado_posesion !== "Completado");
            const completadas = data.filter(r => r.estado_posesion === "Completado");
            setMisTareas(activas);
            setMisTareasCompletadas(completadas);
          }
        })
        .catch((err) => console.error(err))
        .finally(() => setLoadingTareas(false));
    }
  }, [activeTab, userId]);

 
    // Cargar detalle de una tarea (con área completa de la ruta)
  useEffect(() => {
    if (selectedTareaId) {
      setTareaDetail(null);
      fetch(`${API}/documentoradicado/${selectedTareaId}`, { headers: { Authorization: `Bearer ${obtenerToken()}` } })
        .then((res) => res.json())
        .then(async (data) => {
          if (data && data.id) {
            if (data.ruta?.id && !data.ruta?.area) {
              try {
                const rutasRes = await fetch(`${API}/rutas`, { headers: { Authorization: `Bearer ${obtenerToken()}` } });
                const rutasList = await rutasRes.json();
                const rutaCompleta = (Array.isArray(rutasList) ? rutasList : []).find(r => r.id === data.ruta.id);
                if (rutaCompleta) {
                  data.ruta.area = rutaCompleta.area;
                }
              } catch (e) {}
            }
            setTareaDetail(data);
          }
        })
        .catch((err) => console.error(err));
    }
  }, [selectedTareaId]);
  const handleVerArchivo = (filename) => {
    if (!correoDetail) return;
    const url = `${API}/storage/mails/${correoDetail.id_mensaje}/${filename}`;
    window.open(url, "_blank");
  };

  const formatCurrency = (val) => {
    return new Intl.NumberFormat("es-CO", { style: "currency", currency: "COP", minimumFractionDigits: 0 }).format(val || 0);
  };

  // ── Modal de Radicación ──
  const openRadicarModal = (docId) => {
    setRadicarDocId(docId);
    setRadicarForm({
      tipo_radicacion_id: "",
      ruta_id: "",
      metodo_pago_id: "",
      numero_radicado: "",
      normas_reparto: []   // ← FALTABA ESTO
    });
    setShowRadicarModal(true);
  };

  const closeRadicarModal = () => {
    setShowRadicarModal(false);
    setRadicarDocId(null);
  };

  const handleRadicarChange = (e) => {
    const { name, value } = e.target;
    setRadicarForm(prev => ({ ...prev, [name]: value }));
  };

    const handleRadicarSubmit = async () => {
    if (!radicarDocId) return;

    if (!radicarForm.tipo_radicacion_id || !radicarForm.ruta_id || !radicarForm.metodo_pago_id) {
      alert("Debe seleccionar tipo de radicación, ruta y método de pago.");
      return;
    }

    setRadicando(true);
        const payload = {
      documento_comercial_id: radicarDocId,
      tipo_radicacion_id: parseInt(radicarForm.tipo_radicacion_id),
      ruta_id: parseInt(radicarForm.ruta_id),
      metodo_pago_id: parseInt(radicarForm.metodo_pago_id),
      numero_radicado: radicarForm.numero_radicado?.trim() || "",
      normas_reparto: (radicarForm.normas_reparto || [])
        .filter(n => n.norma_reparto_id && n.porcentaje)
        .map(n => ({
          norma_reparto_id: parseInt(n.norma_reparto_id),
          porcentaje: parseFloat(n.porcentaje)
        })),
    };

    try {
      const res = await fetch(`${API}/documentoradicado`, {
        method: "POST",
        headers: { "Content-Type": "application/json", Authorization: `Bearer ${obtenerToken()}` },
        body: JSON.stringify(payload),
      });
      if (!res.ok) {
        const errData = await res.json();
        throw new Error(errData.error || "Error al radicar");
      }

      const creado = await res.json();

      // ── NOTIFICAR al primer responsable ──
      if (creado?.usuario_actual?.id && creado.usuario_actual.id !== userId) {
        await fetch(`${API}/notificacion`, {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${obtenerToken()}`
          },
          body: JSON.stringify({
            usuario_id: creado.usuario_actual.id,
            documento_radicado_id: creado.id,
            mensaje: `Nuevo radicado #${creado.numero_radicado} requiere tu revisión`,
            estado: "Pendiente",
            tipo: "Asignacion",
            fecha_creacion: new Date().toISOString(),
          })
        });
      }

      closeRadicarModal();
      setSelectedDocId(null);
      setDocDetail(null);
      setDocumentos(prev => prev.filter(d => d.id !== radicarDocId));
      setActiveTab("radicados");
    } catch (err) {
      alert("Error al radicar: " + err.message);
    } finally {
      setRadicando(false);
    }
  };


    const renderNormasReparto = () => {
    if (normasRepartoRadicado.length === 0) return null;
    return (
      <div className="doc-section">
        <h4><i className="fa-solid fa-chart-pie"></i> Normas de Reparto ({normasRepartoRadicado.length})</h4>
        <table className="doc-items-table">
          <thead><tr><th>Código</th><th>Nombre</th><th>Sede</th><th>Área</th><th style={{ textAlign: "right" }}>%</th></tr></thead>
          <tbody>
            {normasRepartoRadicado.map((n) => (
              <tr key={n.id}>
                <td><strong>{n.norma_reparto?.codigo || n.codigo}</strong></td>
                <td>{n.norma_reparto?.nombre || n.nombre}</td>
                <td>{n.norma_reparto?.sucursal || n.sucursal}</td>
                <td>{n.norma_reparto?.departamento || n.departamento}</td>
                <td style={{ textAlign: "right", fontWeight: 700 }}>{parseFloat(n.porcentaje).toFixed(2)}%</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    );
  };

    // ── Ver anexo en nueva pestaña (preview) ──
  const handleVerAnexo = async (archivoId, nombre) => {
    try {
      const res = await fetch(`${API}/archivo/${archivoId}/download`, {
        headers: { Authorization: `Bearer ${obtenerToken()}` }
      });
      if (!res.ok) throw new Error("Error obteniendo archivo");

      const blob = await res.blob();
      const url = window.URL.createObjectURL(blob);
      window.open(url, "_blank");

      // Liberar memoria después de 1 minuto
      setTimeout(() => window.URL.revokeObjectURL(url), 60000);
    } catch (err) {
      alert("Error: " + err.message);
    }
  };

  // ── Descargar anexo directamente ──
  const handleDescargarAnexo = async (archivoId, nombre) => {
    try {
      const res = await fetch(`${API}/archivo/${archivoId}/download`, {
        headers: { Authorization: `Bearer ${obtenerToken()}` }
      });
      if (!res.ok) throw new Error("Error descargando archivo");

      const blob = await res.blob();
      const url = window.URL.createObjectURL(blob);

      const a = document.createElement("a");
      a.href = url;
      a.download = nombre;
      document.body.appendChild(a);
      a.click();
      document.body.removeChild(a);

      setTimeout(() => window.URL.revokeObjectURL(url), 60000);
    } catch (err) {
      alert("Error: " + err.message);
    }
  };

  // ── Borrar anexo (solo admin) ──
  const handleBorrarAnexo = async (archivoId, radicadoId) => {
    if (!confirm("¿Está seguro de eliminar este archivo?")) return;
    try {
      const res = await fetch(`${API}/archivo/${archivoId}`, {
        method: "DELETE",
        headers: { Authorization: `Bearer ${obtenerToken()}` }
      });
      if (!res.ok) throw new Error("Error eliminando archivo");

      alert("Archivo eliminado");

      // Refrescar detalle según pestaña activa
      if (activeTab === "tareas") {
        setSelectedTareaId(null);
        setTimeout(() => setSelectedTareaId(radicadoId), 10);
      } else if (activeTab === "radicados") {
        setSelectedRadicadoId(null);
        setTimeout(() => setSelectedRadicadoId(radicadoId), 10);
      }
    } catch (err) {
      alert("Error: " + err.message);
    }
  };

    const crearNotificacion = async ({ usuario_id, mensaje, tipo = "Sistema", documento_radicado_id = null }) => {
    try {
      await fetch(`${API}/notificacion`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${obtenerToken()}`,
        },
        body: JSON.stringify({
          usuario_id,
          documento_radicado_id,
          mensaje,
          estado: "Pendiente",
          tipo,
          fecha_creacion: new Date().toISOString(),
        }),
      });
    } catch (err) {
      console.error("Error creando notificación:", err);
    }
  };

  const getTabInfo = () => {
    switch (activeTab) {
      case "correos": return { icon: "fa-solid fa-envelope", title: "Recepción de Correos", subtitle: "Gestiona las facturas electrónicas recibidas" };
      case "documentos": return { icon: "fa-solid fa-file-invoice", title: "Documentos Pendientes", subtitle: "Revisa y aprueba documentos para radicación" };
      case "radicados": return { icon: "fa-solid fa-stamp", title: "Documentos Radicados", subtitle: "Consulta el estado de los documentos radicados" };
      case "catalogos": return { icon: "fa-solid fa-sliders", title: "Catálogos del Sistema", subtitle: "Gestiona tipos de radicación, pagos y métodos" };
      case "tareas": return { icon: "fa-solid fa-clipboard-list", title: "Mis Tareas", subtitle: "Documentos radicados asignados a ti" };
      default: return { icon: "fa-solid fa-house", title: "Procesos administrativos", subtitle: "Selecciona un formato del menú lateral" };
    }
  };

  const tabInfo = getTabInfo();

  const renderWelcome = () => (
    <div className="content-body">
      <div className="welcome-wrap">
        <div className="welcome-icon"><i className="fa-solid fa-user"></i></div>
        <h2>Bienvenido, Usuario</h2>
        <p>En la barra a tu izquierda encontraras todos los procesos de SIGEFAE.</p>
      </div>
    </div>
  );

  const renderCorreos = () => (
    <div className="correos-container">
      <div className="correos-list">
        <h3>Bandeja de Entrada</h3>
        {loading ? <p>Cargando correos...</p> : correos.length === 0 ? <p>No hay correos registrados.</p> : (
          correos.map((c) => (
            <div key={c.id} className={`correo-item ${selectedCorreoId === c.id ? "active" : ""}`} onClick={() => setSelectedCorreoId(c.id)}>
              <div className="correo-item-header"><strong>{c.de}</strong><span className="correo-date">{new Date(c.fecha_recepcion).toLocaleDateString()}</span></div>
              <div className="correo-item-subject">{c.asunto}</div>
              <div className="correo-item-status"><span className={`status-badge ${c.estado_correo?.nombre.toLowerCase().replace(" ", "-")}`}>{c.estado_correo?.nombre || "Sin Estado"}</span></div>
            </div>
          ))
        )}
      </div>
      <div className="correo-detail">
        {!selectedCorreoId ? (
          <div className="correo-empty"><i className="fa-regular fa-envelope-open"></i><p>Selecciona un correo para leerlo</p></div>
        ) : !correoDetail ? <p style={{ padding: "20px" }}>Cargando detalle...</p> : (
          <div className="correo-detail-content">
            <div className="correo-header"><h2>{correoDetail.asunto}</h2><div className="correo-meta"><p><strong>De:</strong> {correoDetail.de}</p><p><strong>Para:</strong> {correoDetail.para}</p><p><strong>Fecha:</strong> {new Date(correoDetail.fecha_recepcion).toLocaleString()}</p></div></div>
            <div className="correo-body"><iframe title="Cuerpo del correo" srcDoc={correoDetail.cuerpo} className="correo-iframe" /></div>
            <div className="correo-attachments">
              <h4>Archivos Adjuntos ({correoDetail.archivos ? correoDetail.archivos.length : 0})</h4>
              <div className="attachments-grid">
                {correoDetail.archivos && correoDetail.archivos.map((file) => {
                  const ext = file.split('.').pop().toLowerCase();
                  let icon = "fa-file", btnClass = "btn-default", actionText = "Ver / Descargar";
                  if (ext === "pdf") { icon = "fa-file-pdf"; btnClass = "btn-pdf"; actionText = "Ver PDF"; }
                  if (ext === "xml") { icon = "fa-file-code"; btnClass = "btn-xml"; actionText = "Ver XML"; }
                  if (ext === "zip") { icon = "fa-file-zipper"; btnClass = "btn-zip"; actionText = "Descargar ZIP"; }
                  if (ext === "eml") { icon = "fa-envelope-open-text"; btnClass = "btn-eml"; actionText = "Ver Original"; }
                  return (
                    <div key={file} className="attachment-card">
                      <i className={`fa-solid ${icon}`}></i>
                      <span className="attachment-name" title={file}>{file}</span>
                      <button className={`attachment-btn ${btnClass}`} onClick={() => handleVerArchivo(file)}>{actionText}</button>
                    </div>
                  );
                })}
              </div>
            </div>
          </div>
        )}
      </div>
    </div>
  );

  const renderDocumentos = () => (
    <div className="correos-container">
      <div className="correos-list">
        <div style={{ display: "flex", justifyContent: "space-between", alignItems: "center", padding: "0 15px" }}>
          <h3>Pendientes de Revisión ({documentos.length})</h3>
          {esAdmin && (
            <button className="doc-btn doc-btn-primary" onClick={openCrearDocModal}>
              <i className="fa-solid fa-plus"></i> Crear Documento Manual
            </button>
          )}
        </div>
        {loadingDocs ? <p style={{ padding: "15px" }}>Cargando documentos...</p> : documentos.length === 0 ? (
          <p style={{ padding: "15px", color: "#6b7280" }}>No hay documentos pendientes.</p>
        ) : (
          documentos.map((doc) => (
            <div key={doc.id} className={`correo-item ${selectedDocId === doc.id ? "active" : ""}`} onClick={() => setSelectedDocId(doc.id)}>
              <div className="correo-item-header"><strong>{doc.proveedor?.razon_social || "Sin proveedor"}</strong><span className="correo-date">{new Date(doc.fecha_documento).toLocaleDateString()}</span></div>
              <div className="correo-item-subject">{doc.tipo} - {doc.numero_documento}</div>
              <div className="correo-item-status"><span className="status-badge doc-pendiente">{formatCurrency(doc.total)}</span></div>
            </div>
          ))
        )}
      </div>

      <div className="correo-detail">
        {!selectedDocId ? (
          <div className="correo-empty"><i className="fa-solid fa-file-invoice"></i><p>Selecciona un documento para revisarlo</p></div>
        ) : !docDetail ? <p style={{ padding: "20px" }}>Cargando detalle...</p> : (
          <div className="doc-detail-content">
            <div className="doc-header">
              <div className="doc-header-top"><h2>{docDetail.tipo} #{docDetail.numero_documento}</h2><span className="doc-total">{formatCurrency(docDetail.total)}</span></div>
              {docDetail.cufe && <p className="doc-cufe"><strong>CUFE:</strong> {docDetail.cufe}</p>}
            </div>

            <div className="doc-body">
              <div className="doc-section">
                <h4><i className="fa-solid fa-circle-info"></i> Información General</h4>
                <div className="doc-grid">
                  <div className="doc-field"><label>Tipo</label><span>{docDetail.tipo}</span></div>
                  <div className="doc-field"><label>Número Documento</label><span>{docDetail.numero_documento}</span></div>
                  <div className="doc-field"><label>Fecha Emisión</label><span>{new Date(docDetail.fecha_documento).toLocaleDateString()}</span></div>
                  <div className="doc-field"><label>Moneda</label><span>{docDetail.moneda?.nombre || "—"}</span></div>
                  <div className="doc-field"><label>Asunto</label><span>{docDetail.correo?.asunto || "—"}</span></div>
                </div>
              </div>

              <div className="doc-section">
                <h4><i className="fa-solid fa-building"></i> Proveedor</h4>
                <div className="doc-grid">
                  <div className="doc-field"><label>Razón Social</label><span>{docDetail.proveedor?.razon_social || "—"}</span></div>
                  <div className="doc-field"><label>NIT</label><span>{docDetail.proveedor?.numero_documento || "—"}</span></div>
                </div>
              </div>

              <div className="doc-section">
                <h4><i className="fa-solid fa-user-tie"></i> Receptor</h4>
                <div className="doc-grid">
                  <div className="doc-field"><label>Nombre</label><span>{docDetail.receptor?.nombre || "—"}</span></div>
                  <div className="doc-field"><label>NIT</label><span>{docDetail.receptor?.numero_documento || "—"}</span></div>
                </div>
              </div>

              <div className="doc-section">
                <h4><i className="fa-solid fa-calculator"></i> Resumen Financiero</h4>
                <div className="doc-totals">
                  <div className="doc-total-row"><span>Subtotal</span><span>{formatCurrency(docDetail.subtotal)}</span></div>
                  <div className="doc-total-row"><span>IVA</span><span>{formatCurrency(docDetail.iva)}</span></div>
                  <div className="doc-total-row total-final"><span>Total</span><span>{formatCurrency(docDetail.total)}</span></div>
                </div>
              </div>

              {docDetail.detalles && docDetail.detalles.length > 0 && (
                <div className="doc-section">
                  <h4><i className="fa-solid fa-list"></i> Detalle de Ítems ({docDetail.detalles.length})</h4>
                  <table className="doc-items-table">
                    <thead><tr><th>Descripción</th><th>Cantidad</th><th>Valor Unit.</th><th>Total</th></tr></thead>
                    <tbody>{docDetail.detalles.map((item) => (
                      <tr key={item.id}><td>{item.descripcion}</td><td>{item.cantidad}</td><td>{formatCurrency(item.valor_unitario)}</td><td>{formatCurrency(item.total)}</td></tr>
                    ))}</tbody>
                  </table>
                </div>
              )}

              {docDetail.correo && (
                <div className="doc-section">
                  <h4><i className="fa-solid fa-envelope"></i> Correo de Origen</h4>
                  <p style={{ fontSize: "0.85em", color: "#6b7280" }}><strong>Asunto:</strong> {docDetail.correo.asunto}</p>
                </div>
              )}
            </div>
            <div className="doc-actions">
              <button className="doc-btn doc-btn-primary" onClick={() => openRadicarModal(docDetail.id)}>
                <i className="fa-solid fa-stamp"></i> Aprobar para Radicación
              </button>
            </div>
          </div>
        )}
      </div>

      {/* ── Modal de Radicación ── */}
      {showRadicarModal && (
        <div className="modal-overlay" onClick={closeRadicarModal}>
          <div className="modal-content" onClick={(e) => e.stopPropagation()}>
            <div className="modal-header">
              <h3><i className="fa-solid fa-stamp"></i> Radicar Documento</h3>
              <button className="modal-close" onClick={closeRadicarModal}><i className="fa-solid fa-xmark"></i></button>
            </div>
            <div className="modal-body">
              <div className="modal-field">
                <label>Tipo de Radicación <span className="required">*</span></label>
                <select name="tipo_radicacion_id" value={radicarForm.tipo_radicacion_id} onChange={handleRadicarChange} className="doc-input">
                  <option value="">Seleccione...</option>
                  {tiposRadicacion.map(tr => <option key={tr.id} value={tr.id}>{tr.nombre}</option>)}
                </select>
              </div>
              <div className="modal-field">
                <label>Ruta <span className="required">*</span></label>
                <select name="ruta_id" value={radicarForm.ruta_id} onChange={handleRadicarChange} className="doc-input">
                  <option value="">Seleccione...</option>
                  {rutas.map(r => <option key={r.id} value={r.id}>{r.nombre}</option>)}
                </select>
              </div>
              <div className="modal-field">
                <label>Método de Pago <span className="required">*</span></label>
                <select name="metodo_pago_id" value={radicarForm.metodo_pago_id} onChange={handleRadicarChange} className="doc-input">
                  <option value="">Seleccione...</option>
                  {metodosPago.map(mp => <option key={mp.id} value={mp.id}>{mp.nombre}</option>)}
                </select>
              </div>
              <div className="modal-field">
                <label>Número de Radicado <small>(opcional, se autogenera si está vacío)</small></label>
                <input type="text" name="numero_radicado" value={radicarForm.numero_radicado} onChange={handleRadicarChange} className="doc-input" placeholder="Ej: RAD-2026-00001" />
              </div>
                            {/* ── NORMAS DE REPARTO JERÁRQUICAS ── */}
              <div className="modal-field">
                <label style={{ display: "flex", justifyContent: "space-between", alignItems: "center" }}>
                  <span><i className="fa-solid fa-chart-pie"></i> Normas de Reparto</span>
                  <span style={{ 
                    fontSize: "0.8em", 
                    color: totalPorcentajeNormas === 100 ? "#059669" : totalPorcentajeNormas > 0 ? "#dc2626" : "#6b7280",
                    fontWeight: 700 
                  }}>
                    Total: {totalPorcentajeNormas.toFixed(2)}%
                  </span>
                </label>

                {/* Filtros */}
                <div style={{ display: "flex", gap: 8, marginBottom: 10 }}>
                  <select
                    className="doc-input"
                    style={{ flex: 1 }}
                    value={normaFiltroSede}
                    onChange={(e) => { setNormaFiltroSede(e.target.value); setNormaFiltroArea(""); setNormaSeleccionadaId(""); }}
                  >
                    <option value="">Todas las sedes</option>
                    {sedesDisponibles.map(s => <option key={s} value={s}>{s}</option>)}
                  </select>
                  <select
                    className="doc-input"
                    style={{ flex: 1 }}
                    value={normaFiltroArea}
                    onChange={(e) => { setNormaFiltroArea(e.target.value); setNormaSeleccionadaId(""); }}
                    disabled={!normaFiltroSede}
                  >
                    <option value="">Todas las áreas</option>
                    {areasDisponibles.map(a => <option key={a} value={a}>{a}</option>)}
                  </select>
                </div>

                {/* Agregar norma */}
                <div style={{ display: "flex", gap: 8, marginBottom: 10, alignItems: "center" }}>
                  <select
                    className="doc-input"
                    style={{ flex: 2 }}
                    value={normaSeleccionadaId}
                    onChange={(e) => setNormaSeleccionadaId(e.target.value)}
                    disabled={!normaFiltroSede || !normaFiltroArea}
                  >
                    <option value="">Seleccione norma...</option>
                    {normasFiltradas.map(n => (
                      <option key={n.id} value={String(n.id)}>
                        {n.codigo} — {n.nombre}
                      </option>
                    ))}
                  </select>
                  <input
                    type="number"
                    min="0"
                    max="100"
                    step="0.01"
                    placeholder="%"
                    className="doc-input"
                    style={{ flex: 0.6, textAlign: "right" }}
                    value={normaPorcentajeInput}
                    onChange={(e) => setNormaPorcentajeInput(e.target.value)}
                  />
                  <button
                    className="doc-btn doc-btn-secondary"
                    onClick={handleAgregarNormaModal}
                    disabled={!normaSeleccionadaId || !normaPorcentajeInput}
                    style={{ padding: "6px 12px", fontSize: "0.85em" }}
                  >
                    <i className="fa-solid fa-plus"></i> Agregar
                  </button>
                </div>

                {/* Lista de normas agregadas */}
                {(radicarForm.normas_reparto || []).length === 0 ? (
                  <p style={{ fontSize: "0.85em", color: "#6b7280" }}>No se han asignado normas de reparto.</p>
                ) : (
                  <div style={{ display: "flex", flexDirection: "column", gap: 6, marginBottom: 10 }}>
                    {radicarForm.normas_reparto.map((norma, idx) => {
                      const info = normasRepartoCatalogo.find(n => String(n.id) === norma.norma_reparto_id);
                      return (
                        <div key={idx} style={{ 
                          display: "flex", 
                          gap: 8, 
                          alignItems: "center", 
                          padding: "8px 10px", 
                          background: "#f9fafb", 
                          borderRadius: 6,
                          border: "1px solid #e5e7eb"
                        }}>
                          <span style={{ flex: 2, fontSize: "0.85em", fontWeight: 600, color: "#1f2937" }}>
                            {info ? `${info.codigo} — ${info.nombre}` : "Norma desconocida"}
                          </span>
                          <span style={{ flex: 1, fontSize: "0.8em", color: "#6b7280" }}>
                            {info ? `${info.sucursal} / ${info.departamento}` : ""}
                          </span>
                          <input
                            type="number"
                            min="0"
                            max="100"
                            step="0.01"
                            value={norma.porcentaje}
                            onChange={(e) => handleNormaRepartoChange(idx, "porcentaje", e.target.value)}
                            className="doc-input"
                            style={{ width: 70, textAlign: "right", fontSize: "0.85em" }}
                          />
                          <span style={{ fontSize: "0.85em", fontWeight: 700 }}>%</span>
                          <button
                            className="btn-icon btn-toggle"
                            onClick={() => handleRemoveNormaReparto(idx)}
                            title="Quitar"
                            style={{ width: 28, height: 28, flexShrink: 0 }}
                          >
                            <i className="fa-solid fa-xmark"></i>
                          </button>
                        </div>
                      );
                    })}
                  </div>
                )}

                {totalPorcentajeNormas > 0 && totalPorcentajeNormas !== 100 && (
                  <p style={{ fontSize: "0.8em", color: "#dc2626", marginTop: 4 }}>
                    <i className="fa-solid fa-triangle-exclamation"></i> La suma de porcentajes debe ser exactamente 100%
                  </p>
                )}
              </div>
            </div>
            <div className="modal-footer">
              <button className="doc-btn doc-btn-secondary" onClick={closeRadicarModal} disabled={radicando}>Cancelar</button>
              <button className="doc-btn doc-btn-primary" onClick={handleRadicarSubmit} disabled={radicando || (radicarForm.normas_reparto.length > 0 && totalPorcentajeNormas !== 100)}>
                <i className="fa-solid fa-stamp"></i> {radicando ? "Radicando..." : "Confirmar Radicación"}
              </button>
            </div>
          </div>
        </div>
      )}

    </div>
  );

  const renderRadicados = () => (
  <div className="correos-container">
    <div className="correos-list">
      <h3>Documentos Radicados ({radicados.length})</h3>
      {loadingRadicados ? <p style={{ padding: "15px" }}>Cargando...</p> : radicados.length === 0 ? (
        <p style={{ padding: "15px", color: "#6b7280" }}>No hay documentos radicados.</p>
      ) : (
        radicados.map((rad) => (
          <div key={rad.id} className={`correo-item ${selectedRadicadoId === rad.id ? "active" : ""}`} onClick={() => setSelectedRadicadoId(rad.id)}>
            <div className="correo-item-header">
              <strong>{rad.documento_comercial?.numero_documento || "—"}</strong>
              <span className="correo-date">{new Date(rad.fecha_radicacion).toLocaleDateString()}</span>
            </div>
            <div className="correo-item-subject">{rad.documento_comercial?.tipo || "—"} — {rad.numero_radicado}</div>
            <div className="correo-item-status">
              <span className={`status-badge ${
                rad.estado_posesion === "Completado" ? "radicado" :
                rad.estado_posesion === "EnProceso" ? "doc-pendiente" :
                rad.estado_posesion === "Devuelto" ? "doc-rechazado" : ""
              }`}>
                {rad.estado_posesion === "Completado" ? "Finalizado" :
                 rad.estado_posesion === "EnProceso" ? "En Proceso" :
                 rad.estado_posesion === "Devuelto" ? "Devuelto" :
                 rad.estado_posesion === "Libre" ? "Libre" : "En espera"}
              </span>
            </div>
          </div>
        ))
      )}
    </div>

    <div className="correo-detail">
      {!selectedRadicadoId ? (
        <div className="correo-empty"><i className="fa-solid fa-stamp"></i><p>Selecciona un radicado para ver su detalle</p></div>
      ) : !radicadoDetail ? <p style={{ padding: "20px" }}>Cargando detalle...</p> : (
        <div className="doc-detail-content">
          <div className="doc-header">
            <div className="doc-header-top">
              <div style={{ display: 'flex', alignItems: 'center', gap: '12px' }}>
                <h2>Radicado #{radicadoDetail.numero_radicado}</h2>
                <span className={`status-badge ${
                  radicadoDetail.estado_posesion === "Completado" ? "radicado" :
                  radicadoDetail.estado_posesion === "EnProceso" ? "doc-pendiente" :
                  radicadoDetail.estado_posesion === "Devuelto" ? "doc-rechazado" : ""
                }`}>
                  {radicadoDetail.estado_posesion === "Completado" ? "Finalizado" :
                   radicadoDetail.estado_posesion === "EnProceso" ? "En Proceso" :
                   radicadoDetail.estado_posesion === "Devuelto" ? "Devuelto" : "Pendiente"}
                </span>
              </div>
              <div style={{ display: 'flex', gap: '10px', alignItems: 'center' }}>
                <span className="doc-total">{
                radicadoDetail.estado_posesion === "Completado" ? "Finalizado" :
                radicadoDetail.estado_posesion === "EnProceso" ? "En Proceso" :
                radicadoDetail.estado_posesion === "Libre" ? "Libre" : "En espera"}
                </span>
                <button 
                  className="btn btn-outline" 
                  onClick={() => handleDescargarExpediente(radicadoDetail)}
                  disabled={generandoPdf}
                >
                  <i className={`fa-solid ${generandoPdf ? 'fa-spinner fa-spin' : 'fa-file-pdf'}`}></i> 
                  {generandoPdf ? " Generando..." : " Expediente"}
                </button>
              </div>
            </div>
            <p className="doc-cufe"><strong>Documento:</strong> {radicadoDetail.documento_comercial?.tipo} {radicadoDetail.documento_comercial?.numero_documento}</p>
          </div>

          <div className="doc-body">
            {/* ── Anexos (CORREGIDO: radicadoDetail en vez de tareaDetail) ── */}
            {radicadoDetail.archivos && radicadoDetail.archivos.length > 0 && (
  <div className="doc-section">
    <h4><i className="fa-solid fa-paperclip"></i> Anexos ({radicadoDetail.archivos.length})</h4>
    <div className="attachments-grid">
      {radicadoDetail.archivos.map((arch) => {
        const ext = arch.extension?.toLowerCase() || "";
        const esPreview = ["pdf", "jpg", "jpeg", "png", "gif", "txt"].includes(ext);
        return (
          <div key={arch.id} className="attachment-card" style={{ position: "relative" }}>
            <i className={`fa-solid fa-file${ext === 'pdf' ? '-pdf' : ext === 'xml' ? '-code' : ''}`}></i>
            <span className="attachment-name" title={arch.nombre}>{arch.nombre}</span>

                  <div style={{ display: "flex", gap: 6, flexWrap: "wrap", justifyContent: "center" }}>
                    {esPreview && (
                    radicadoDetail.estado_posesion === "Completado" ? (
                      <button
                        className="attachment-btn btn-pdf"
                        onClick={() => handleVerAnexo(arch.id, arch.nombre)}
                      >
                        <i className="fa-solid fa-eye"></i> Ver
                      </button>
                    ) : (
                      <button
                        className="attachment-btn btn-pdf"
                        onClick={() =>setPdfEditor({open: true,archivoId: arch.id,radicadoId: radicadoDetail.id})}>
                        <i className="fa-solid fa-pen-to-square"></i> Abrir editor
                      </button>
                    )
                  )}
                    <button className="attachment-btn btn-default" onClick={() => handleDescargarAnexo(arch.id, arch.nombre)}>
                      <i className="fa-solid fa-download"></i> Descargar
                    </button>
                    {userRol === "Superadministrador" && (
                      <button
                        className="attachment-btn btn-toggle"
                        onClick={() => handleBorrarAnexo(arch.id, radicadoDetail.id)}
                        title="Eliminar"
                      >
                        <i className="fa-solid fa-xmark"></i>
                      </button>
                    )}
                  </div>
                </div>
              );
            })}
          </div>
        </div>
      )}

            {/* ── Subir nuevo anexo (CORREGIDO: radicadoDetail.id en vez de tareaDetail.id) ── */}
            <div className="doc-section">
              <h4><i className="fa-solid fa-cloud-arrow-up"></i> Adjuntar archivo</h4>
              <input 
                type="file" 
                id="anexo-input-radicado"
                onChange={(e) => handleSubirAnexo(e, radicadoDetail.id)} 
                style={{ display: 'none' }}
              />
              <button 
                className="doc-btn doc-btn-secondary" 
                onClick={() => document.getElementById('anexo-input-radicado').click()}
              >
                <i className="fa-solid fa-upload"></i> Seleccionar archivo
              </button>
            </div>

            {renderFlujoAprobacion()}

            {renderTrazabilidad()}
            
            {renderNormasReparto()}
            
            {/* ── Comentarios ── */}
            {renderComentarios(radicadoDetail.id,radicadoDetail.estado_posesion === "Completado")}

            <div className="doc-section">
              <h4><i className="fa-solid fa-circle-info"></i> Información del Radicado</h4>
              <div className="doc-grid">
                <div className="doc-field"><label>Número Radicado</label><span>{radicadoDetail.numero_radicado}</span></div>
                <div className="doc-field"><label>Fecha Radicación</label><span>{new Date(radicadoDetail.fecha_radicacion).toLocaleString()}</span></div>
                <div className="doc-field"><label>Tipo Radicación</label><span>{radicadoDetail.tipo_radicacion?.nombre || "—"}</span></div>
                <div className="doc-field"><label>Ruta</label><span>{radicadoDetail.ruta?.nombre || "—"}</span></div>
                <div className="doc-field"><label>Área</label><span>{typeof radicadoDetail.ruta?.area === 'string' ? radicadoDetail.ruta.area : (radicadoDetail.ruta?.area?.nombre || "—")}</span></div>
                <div className="doc-field"><label>Método de Pago</label><span>{radicadoDetail.metodo_pago?.nombre || "—"}</span></div>
                <div className="doc-field"><label>Estado Posesión</label><span>{radicadoDetail.estado_posesion || "—"}</span></div>
                <div className="doc-field"><label>Paso Actual</label><span>{radicadoDetail.paso_actual?.nombre || "Inicio"}</span></div>
                <div className="doc-field"><label>Responsable</label><span>{radicadoDetail.usuario_actual?.nombre || "—"}</span></div>
              </div>
            </div>

            {radicadoDetail.documento_comercial && (
              <>
                <div className="doc-section">
                  <h4><i className="fa-solid fa-file-invoice"></i> Documento Comercial</h4>
                  <div className="doc-grid">
                    <div className="doc-field"><label>Tipo</label><span>{radicadoDetail.documento_comercial.tipo}</span></div>
                    <div className="doc-field"><label>Número</label><span>{radicadoDetail.documento_comercial.numero_documento}</span></div>
                    <div className="doc-field"><label>Proveedor</label><span>{radicadoDetail.documento_comercial.proveedor?.razon_social || "—"}</span></div>
                    <div className="doc-field"><label>Receptor</label><span>{radicadoDetail.documento_comercial.receptor?.nombre || "—"}</span></div>
                  </div>
                </div>

                <div className="doc-section">
                  <h4><i className="fa-solid fa-calculator"></i> Valores</h4>
                  <div className="doc-totals">
                    <div className="doc-total-row"><span>Subtotal</span><span>{formatCurrency(radicadoDetail.documento_comercial.subtotal)}</span></div>
                    <div className="doc-total-row"><span>IVA</span><span>{formatCurrency(radicadoDetail.documento_comercial.iva)}</span></div>
                    <div className="doc-total-row total-final"><span>Total</span><span>{formatCurrency(radicadoDetail.documento_comercial.total)}</span></div>
                  </div>
                </div>
              </>
            )}

            {radicadoDetail.qr?.url && (
              <div className="doc-section" style={{ textAlign: "center" }}>
                <h4><i className="fa-solid fa-qrcode"></i> Código QR del Expediente</h4>
                <img
                  src={`https://api.qrserver.com/v1/create-qr-code/?size=180x180&data=${encodeURIComponent(radicadoDetail.qr.url)}`}
                  alt="Código QR"
                  style={{ margin: "12px auto", display: "block", borderRadius: 8, border: "1px solid #e5e7eb" }}
                />
                <p style={{ fontSize: "0.75em", color: "#6b7280", wordBreak: "break-all", marginTop: 8 }}>
                  {radicadoDetail.qr.url}
                </p>
              </div>
            )}
          </div>
        </div>
      )}
    </div>
  </div>
);

  // ── NUEVO: Panel de Mis Tareas (solo para aprobadores) ──
    const renderTareas = () => {
    const tareasMostradas = tareasSubTab === "activas" ? misTareas : misTareasCompletadas;

    return (
      <div className="correos-container">
        <div className="correos-list">
          <h3>Mis Tareas</h3>

          {/* ── Sub-pestañas ── */}
          <div style={{ display: "flex", borderBottom: "2px solid var(--gray-200)" }}>
            <button
              onClick={() => { setTareasSubTab("activas"); setSelectedTareaId(null); }}
              style={{
                flex: 1,
                padding: "12px",
                border: "none",
                background: tareasSubTab === "activas" ? "white" : "transparent",
                borderBottom: tareasSubTab === "activas" ? "3px solid var(--pardo-yellow)" : "3px solid transparent",
                fontWeight: 600,
                color: tareasSubTab === "activas" ? "var(--pardo-navy)" : "var(--gray-600)",
                cursor: "pointer",
                fontSize: "0.9em"
              }}
            >
              <i className="fa-solid fa-spinner" style={{ marginRight: 6 }}></i>
              En Proceso ({misTareas.length})
            </button>
            <button
              onClick={() => { setTareasSubTab("completadas"); setSelectedTareaId(null); }}
              style={{
                flex: 1,
                padding: "12px",
                border: "none",
                background: tareasSubTab === "completadas" ? "white" : "transparent",
                borderBottom: tareasSubTab === "completadas" ? "3px solid var(--pardo-yellow)" : "3px solid transparent",
                fontWeight: 600,
                color: tareasSubTab === "completadas" ? "var(--pardo-navy)" : "var(--gray-600)",
                cursor: "pointer",
                fontSize: "0.9em"
              }}
            >
              <i className="fa-solid fa-check-circle" style={{ marginRight: 6 }}></i>
              Completadas ({misTareasCompletadas.length})
            </button>
          </div>

          {loadingTareas ? (
            <p style={{ padding: "15px" }}>Cargando...</p>
          ) : tareasMostradas.length === 0 ? (
            <p style={{ padding: "15px", color: "#6b7280" }}>
              {tareasSubTab === "activas"
                ? "No tienes documentos asignados."
                : "No hay documentos finalizados."}
            </p>
          ) : (
            tareasMostradas.map((rad) => (
              <div key={rad.id} className={`correo-item ${selectedTareaId === rad.id ? "active" : ""}`} onClick={() => setSelectedTareaId(rad.id)}>
                <div className="correo-item-header">
                  <strong>{rad.documento_comercial?.numero_documento || "—"}</strong>
                  <span className="correo-date">{new Date(rad.fecha_radicacion).toLocaleDateString()}</span>
                </div>
                <div className="correo-item-subject">{rad.documento_comercial?.tipo || "—"} — {rad.numero_radicado}</div>
                <div className="correo-item-status">
                  <span className={`status-badge ${
                    rad.estado_posesion === "Completado" ? "radicado" :
                    rad.estado_posesion === "EnProceso" ? "doc-pendiente" :
                    rad.estado_posesion === "Devuelto" ? "doc-rechazado" : ""
                  }`}>
                    {rad.estado_posesion === "Completado" ? "Finalizado" :
                     rad.estado_posesion === "EnProceso" ? "En Proceso" :
                     rad.estado_posesion === "Devuelto" ? "Devuelto" : "Sin estado"}
                  </span>
                </div>
              </div>
            ))
          )}
        </div>

        <div className="correo-detail">
          {!selectedTareaId ? (
            <div className="correo-empty"><i className="fa-solid fa-clipboard-list"></i><p>Selecciona una tarea para ver su detalle</p></div>
          ) : !tareaDetail ? (
            <p style={{ padding: "20px" }}>Cargando detalle...</p>
          ) : (
            <div className="doc-detail-content">
              <div className="doc-header">
                <div className="doc-header-top">
                  <div style={{ display: 'flex', alignItems: 'center', gap: '12px' }}>
                    <h2>Radicado #{tareaDetail.numero_radicado}</h2>
                    <span className={`status-badge ${
                      tareaDetail.estado_posesion === "Completado" ? "radicado" :
                      tareaDetail.estado_posesion === "EnProceso" ? "doc-pendiente" :
                      tareaDetail.estado_posesion === "Devuelto" ? "doc-rechazado" : ""
                    }`}>
                      {tareaDetail.estado_posesion === "Completado" ? "Finalizado" :
                       tareaDetail.estado_posesion === "EnProceso" ? "En Proceso" :
                       tareaDetail.estado_posesion === "Devuelto" ? "Devuelto" : "Pendiente"}
                    </span>
                  </div>
                  <div style={{ display: 'flex', gap: '10px', alignItems: 'center' }}>
                    <span className="doc-total">{
                    tareaDetail.estado_posesion === "Completado" ? "Finalizado" :
                    tareaDetail.estado_posesion === "EnProceso" ? "En Proceso" :
                    tareaDetail.estado_posesion === "Libre" ? "Libre" : "En espera"}
                    </span>
                    <button 
                      className="btn btn-outline" 
                      onClick={() => handleDescargarExpediente(tareaDetail)}
                      disabled={generandoPdf}
                    >
                      <i className={`fa-solid ${generandoPdf ? 'fa-spinner fa-spin' : 'fa-file-pdf'}`}></i> 
                      {generandoPdf ? " Generando..." : " Expediente"}
                    </button>
                  </div>
                </div>
                <p className="doc-cufe"><strong>Documento:</strong> {tareaDetail.documento_comercial?.tipo} {tareaDetail.documento_comercial?.numero_documento}</p>
              </div>

              <div className="doc-body">
                {/* ── Anexos ── */}
                {tareaDetail.archivos && tareaDetail.archivos.length > 0 && (
                  <div className="doc-section">
                    <h4><i className="fa-solid fa-paperclip"></i> Anexos ({tareaDetail.archivos.length})</h4>
                    <div className="attachments-grid">
                      {tareaDetail.archivos.map((arch) => {
                        const ext = arch.extension?.toLowerCase() || "";
                        const esPreview = ["pdf", "jpg", "jpeg", "png", "gif", "txt"].includes(ext);
                        return (
                          <div key={arch.id} className="attachment-card" style={{ position: "relative" }}>
                            <i className={`fa-solid fa-file${ext === 'pdf' ? '-pdf' : ext === 'xml' ? '-code' : ''}`}></i>
                            <span className="attachment-name" title={arch.nombre}>{arch.nombre}</span>
                            <div style={{ display: "flex", gap: 6, flexWrap: "wrap", justifyContent: "center" }}>
                              {esPreview && (
                              tareaDetail.estado_posesion === "Completado" ? (
                                <button
                                  className="attachment-btn btn-pdf"
                                  onClick={() => handleVerAnexo(arch.id, arch.nombre)}
                                >
                                  <i className="fa-solid fa-eye"></i> Ver
                                </button>
                              ) : (
                                <button
                                  className="attachment-btn btn-pdf"
                                  onClick={() =>
                                    setPdfEditor({open: true,archivoId: arch.id,radicadoId: tareaDetail.id})}>
                                  <i className="fa-solid fa-pen-to-square"></i> Abrir editor
                                  </button>))}
                              <button className="attachment-btn btn-default" onClick={() => handleDescargarAnexo(arch.id, arch.nombre)}>
                                <i className="fa-solid fa-download"></i> Descargar
                              </button>
                              {userRol === "Superadministrador" && (
                                <button className="attachment-btn btn-toggle" onClick={() => handleBorrarAnexo(arch.id, tareaDetail.id)} title="Eliminar">
                                  <i className="fa-solid fa-xmark"></i>
                                </button>
                              )}
                            </div>
                          </div>
                        );
                      })}
                    </div>
                  </div>
                )}

                {/* ── Subir anexo (solo si está activo) ── */}
                {tareaDetail.estado_posesion !== "Completado" && (
                  <div className="doc-section">
                    <h4><i className="fa-solid fa-cloud-arrow-up"></i> Adjuntar archivo</h4>
                    <input type="file" id="anexo-input" onChange={(e) => handleSubirAnexo(e, tareaDetail.id)} style={{ display: 'none' }} />
                    <button className="doc-btn doc-btn-secondary" onClick={() => document.getElementById('anexo-input').click()}>
                      <i className="fa-solid fa-upload"></i> Seleccionar archivo
                    </button>
                  </div>
                )}

                {renderFlujoAprobacion()}

                {renderTrazabilidad()}
                {renderNormasReparto()}
                {renderComentarios(tareaDetail.id,tareaDetail.estado_posesion === "Completado")}

                <div className="doc-section">
                  <h4><i className="fa-solid fa-circle-info"></i> Información del Radicado</h4>
                  <div className="doc-grid">
                    <div className="doc-field"><label>Número Radicado</label><span>{tareaDetail.numero_radicado}</span></div>
                    <div className="doc-field"><label>Fecha Radicación</label><span>{new Date(tareaDetail.fecha_radicacion).toLocaleString()}</span></div>
                    <div className="doc-field"><label>Tipo Radicación</label><span>{tareaDetail.tipo_radicacion?.nombre || "—"}</span></div>
                    <div className="doc-field"><label>Ruta</label><span>{tareaDetail.ruta?.nombre || "—"}</span></div>
                    <div className="doc-field"><label>Área</label><span>{typeof tareaDetail.ruta?.area === 'string' ? tareaDetail.ruta.area : (tareaDetail.ruta?.area?.nombre || "—")}</span></div>                    <div className="doc-field"><label>Método de Pago</label><span>{tareaDetail.metodo_pago?.nombre || "—"}</span></div>
                    <div className="doc-field"><label>Estado Posesión</label><span>{tareaDetail.estado_posesion || "—"}</span></div>
                    <div className="doc-field"><label>Paso Actual</label><span>{tareaDetail.paso_actual?.nombre || "Inicio"}</span></div>
                    <div className="doc-field"><label>Responsable</label><span>{tareaDetail.usuario_actual?.nombre || "—"}</span></div>
                  </div>
                </div>

                {tareaDetail.documento_comercial && (
                  <>
                    <div className="doc-section">
                      <h4><i className="fa-solid fa-file-invoice"></i> Documento Comercial</h4>
                      <div className="doc-grid">
                        <div className="doc-field"><label>Tipo</label><span>{tareaDetail.documento_comercial.tipo}</span></div>
                        <div className="doc-field"><label>Número</label><span>{tareaDetail.documento_comercial.numero_documento}</span></div>
                        <div className="doc-field"><label>Proveedor</label><span>{tareaDetail.documento_comercial.proveedor?.razon_social || "—"}</span></div>
                        <div className="doc-field"><label>Receptor</label><span>{tareaDetail.documento_comercial.receptor?.nombre || "—"}</span></div>
                      </div>
                    </div>

                    <div className="doc-section">
                      <h4><i className="fa-solid fa-calculator"></i> Valores</h4>
                      <div className="doc-totals">
                        <div className="doc-total-row"><span>Subtotal</span><span>{formatCurrency(tareaDetail.documento_comercial.subtotal)}</span></div>
                        <div className="doc-total-row"><span>IVA</span><span>{formatCurrency(tareaDetail.documento_comercial.iva)}</span></div>
                        <div className="doc-total-row total-final"><span>Total</span><span>{formatCurrency(tareaDetail.documento_comercial.total)}</span></div>
                      </div>
                    </div>
                  </>
                )}

                {tareaDetail.qr?.url && (
                  <div className="doc-section" style={{ textAlign: "center" }}>
                    <h4><i className="fa-solid fa-qrcode"></i> Código QR del Expediente</h4>
                    <img
                      src={`https://api.qrserver.com/v1/create-qr-code/?size=180x180&data=${encodeURIComponent(tareaDetail.qr.url)}`}
                      alt="Código QR"
                      style={{ margin: "12px auto", display: "block", borderRadius: 8, border: "1px solid #e5e7eb" }}
                    />
                    <p style={{ fontSize: "0.75em", color: "#6b7280", wordBreak: "break-all", marginTop: 8 }}>
                      {tareaDetail.qr.url}
                    </p>
                  </div>
                )}
              </div>

              <div className="doc-actions">
                {(() => {
                  const tareaActiva = tareasFlujo.find(t => t.estado?.nombre === "En Proceso");
                  const esResponsable = tareaActiva && (tareaActiva.usuario_asignado_id === userId || tareaDetail.usuario_actual?.id === userId);

                  if (esResponsable && tareaDetail.estado_posesion !== "Completado") {
                    return (
                      <>
                        <button className="doc-btn doc-btn-primary" onClick={() => handleCompletarTarea(tareaActiva.id, tareaDetail.id)}>
                          <i className="fa-solid fa-check"></i> Marcar como Completado
                        </button>
                        <button 
                          className="doc-btn doc-btn-secondary" 
                          onClick={() => {
                            setDevolverForm({ tarea_destino_id: "", observacion: "", retorno_directo: true });
                            setShowDevolverModal(true);
                          }}
                          style={{ borderColor: "var(--pardo-red)", color: "var(--pardo-red)" }}
                        >
                          <i className="fa-solid fa-reply"></i> Devolver
                        </button>
                      </>
                    );
                  }
                  if (tareaDetail.estado_posesion === "Completado") {
                    return (
                      <span className="status-badge radicado" style={{ padding: "10px 16px" }}>
                        <i className="fa-solid fa-check-circle"></i> Proceso Finalizado
                      </span>
                    );
                  }
                  return null;
                })()}
              </div>
            </div>
          )}
        </div>
      </div>
    );
  };

  const renderCatalogos = () => {
    const cfg = catalogoConfig[catalogoActivo];

    const renderFormField = (field) => {
      const value = catalogoForm[field] || "";
      if (field === "tipo_pago_id") {
        return (
          <div className="modal-field" key={field}>
            <label>Tipo de Pago <span className="required">*</span></label>
            <select name={field} value={value} onChange={handleCatalogoFormChange}>
              <option value="">Seleccione...</option>
              {tiposPagoCatalogo.map(tp => <option key={tp.id} value={tp.id}>{tp.nombre}</option>)}
            </select>
          </div>
        );
      }
      if (field === "area_id") {
        return (
          <div className="modal-field" key={field}>
            <label>Área <span className="required">*</span></label>
            <select name={field} value={value} onChange={handleCatalogoFormChange}>
              <option value="">Seleccione...</option>
              {areasCatalogo.map(a => <option key={a.id} value={a.id}>{a.nombre}</option>)}
            </select>
          </div>
        );
      }
      if (field === "ruta_id") {
        return (
          <div className="modal-field" key={field}>
            <label>Ruta <span className="required">*</span></label>
            <select name={field} value={value} onChange={handleCatalogoFormChange}>
              <option value="">Seleccione...</option>
              {rutasCatalogo.map(r => <option key={r.id} value={r.id}>{r.nombre}</option>)}
            </select>
          </div>
        );
      }
      if (field === "usuario_id") {
        return (
          <div className="modal-field" key={field}>
            <label>Usuario Responsable <span className="required">*</span></label>
            <select name={field} value={value} onChange={handleCatalogoFormChange}>
              <option value="">Seleccione...</option>
              {usuariosCatalogo.map(u => <option key={u.id} value={u.id}>{u.nombre}</option>)}
            </select>
          </div>
        );
      }
      if (field === "usuario_aprobador_id") {
        return (
          <div className="modal-field" key={field}>
            <label>Usuario Aprobador <span className="required">*</span></label>
            <select name={field} value={value} onChange={handleCatalogoFormChange}>
              <option value="">Seleccione...</option>
              {usuariosCatalogo.map(u => <option key={u.id} value={u.id}>{u.nombre}</option>)}
            </select>
          </div>
        );
      }
      if (field === "orden") {
        return (
          <div className="modal-field" key={field} style={{ maxWidth: 120 }}>
            <label>Orden <span className="required">*</span></label>
            <input type="number" name={field} value={value} onChange={handleCatalogoFormChange} min="1" />
          </div>
        );
      }
      if (field === "monto_minimo") {
        return (
          <div className="modal-field" key={field} style={{ maxWidth: 160 }}>
            <label>Monto Mínimo <span className="required">*</span></label>
            <input type="number" name={field} value={value} onChange={handleCatalogoFormChange} min="0" />
          </div>
        );
      }
      if (field === "moneda_id") {
        return (
          <div className="modal-field" key={field}>
            <label>Moneda <span className="required">*</span></label>
            <select name={field} value={value} onChange={handleCatalogoFormChange}>
              <option value="">Seleccione...</option>
              {monedasCatalogo.map(m => <option key={m.id} value={m.id}>{m.nombre} ({m.codigo})</option>)}
            </select>
          </div>
        );
      }
      if (field === "posicion_insercion") {
        return (
          <div className="modal-field" key={field}>
            <label>Posición Inserción <span className="required">*</span></label>
            <select name={field} value={value} onChange={handleCatalogoFormChange}>
              <option value="">Seleccione...</option>
              <option value="PRIMERO">Al inicio</option>
              <option value="ANTES_FINAL">Antes del final</option>
              <option value="ULTIMO">Al final</option>
            </select>
          </div>
        );
      }
            if (field === "codigo") {
        return (
          <div className="modal-field" key={field} style={{ flex: 1 }}>
            <label>Código <span className="required">*</span></label>
            <input type="text" name={field} value={value} onChange={handleCatalogoFormChange} placeholder="Ej: BU101" />
          </div>
        );
      }
            if (field === "sucursal") {
        const sucursales = ["BUCARAMANGA", "MALAMBO", "CUCUTA", "CB", "CIENAGA DE ORO", "GENERAL"];
        return (
          <div className="modal-field" key={field}>
            <label>Sucursal <span className="required">*</span></label>
            <select name={field} value={value} onChange={handleCatalogoFormChange}>
              <option value="">Seleccione...</option>
              {sucursales.map(s => <option key={s} value={s}>{s}</option>)}
            </select>
          </div>
        );
      }
      if (field === "departamento") {
        const depts = ["ADMON", "VENTAS", "PRODUCCION"];
        return (
          <div className="modal-field" key={field}>
            <label>Departamento <span className="required">*</span></label>
            <select name={field} value={value} onChange={handleCatalogoFormChange}>
              <option value="">Seleccione...</option>
              {depts.map(d => <option key={d} value={d}>{d}</option>)}
            </select>
          </div>
        );
      }
      if (field === "tipo") {
        return (
          <div className="modal-field" key={field}>
            <label>Tipo</label>
            <select name={field} value={value} onChange={handleCatalogoFormChange}>
              <option value="">Ninguno</option>
              <option value="Servicio">Servicio</option>
              <option value="Compra">Compra</option>
            </select>
          </div>
        );
      }
      if (field === "tarifa_iva") {
        return (
          <div className="modal-field" key={field} style={{ maxWidth: 120 }}>
            <label>Tarifa IVA</label>
            <select name={field} value={value} onChange={handleCatalogoFormChange}>
              <option value="">Ninguna</option>
              <option value="19%">19%</option>
              <option value="5%">5%</option>
              <option value="0%">0%</option>
            </select>
          </div>
        );
      }
      return (
        <div className="modal-field" key={field} style={{ flex: 1 }}>
          <label>Nombre <span className="required">*</span></label>
          <input type="text" name={field} value={value} onChange={handleCatalogoFormChange} placeholder={`Nombre del ${cfg.label.toLowerCase()}`} />
        </div>
      );
    };

        const getColumnLabel = (field) => {
      const map = { 
        nombre: "Nombre", 
        codigo: "Código",
        sucursal: "Sucursal",
        departamento: "Depto",
        tipo: "Tipo",
        tarifa_iva: "Tarifa IVA",
        tipo_pago: "Tipo de Pago", 
        area: "Área", 
        ruta: "Ruta", 
        orden: "Orden", 
        usuario: "Usuario" 
      };
      return map[field] || field;
    };

    const getItemDisplay = (item, field) => {
      if (field === "tipo_pago_id") return item.tipo_pago || "—";
      if (field === "area_id") return item.area || "—";
      if (field === "ruta_id") return item.ruta || "—";
      if (field === "usuario_id") return item.usuario || "—";
      if ((field === "tipo" || field === "tarifa_iva") && !item[field]) return "—";
      return item[field] !== undefined ? item[field] : "—";
    };

    return (
      <div className="catalogos-container">
        <div className="catalogo-tabs">
          {Object.entries(catalogoConfig).map(([key, c]) => (
            <button key={key} className={catalogoActivo === key ? "active" : ""} onClick={() => setCatalogoActivo(key)}>
              {c.label}
            </button>
          ))}
        </div>

        <div className="catalogo-header">
          <h3>{cfg.label}</h3>
          <button className="doc-btn doc-btn-primary" onClick={openCatalogoCreate}>
            <i className="fa-solid fa-plus"></i> Nuevo {cfg.label}
          </button>
        </div>

        {showCatalogoForm && (
          <div className="catalogo-form">
            <h4>{catalogoEditing ? "Editar" : "Crear"} {cfg.label}</h4>
            <div className="catalogo-form-row">
              {cfg.fields.map(renderFormField)}
              <div style={{ display: "flex", gap: 8, alignItems: "flex-end" }}>
                <button className="doc-btn doc-btn-secondary" onClick={() => setShowCatalogoForm(false)}>Cancelar</button>
                <button className="doc-btn doc-btn-primary" onClick={handleCatalogoSubmit}>
                  <i className="fa-solid fa-floppy-disk"></i> Guardar
                </button>
              </div>
            </div>
          </div>
        )}

        {catalogoLoading ? (
          <p>Cargando...</p>
        ) : catalogoItems.length === 0 ? (
          <p style={{ color: "#6b7280" }}>No hay registros.</p>
        ) : (
          <table className="catalogo-table">
            <thead>
              <tr>
                <th>ID</th>
                {cfg.fields.filter(f => f !== "tipo_pago_id").map(f => (
                  <th key={f}>{getColumnLabel(f === "area_id" ? "area" : f === "ruta_id" ? "ruta" : f === "usuario_id" ? "usuario" : f)}</th>
                ))}
                {cfg.fields.includes("tipo_pago_id") && <th>Tipo de Pago</th>}
                <th>Estado</th>
                <th style={{ width: 120 }}>Acciones</th>
              </tr>
            </thead>
            <tbody>
              {catalogoItems.map(item => (
                <tr key={item.id}>
                  <td>{item.id}</td>
                  {cfg.fields.filter(f => f !== "tipo_pago_id").map(f => (
                    <td key={f}>{getItemDisplay(item, f)}</td>
                  ))}
                  {cfg.fields.includes("tipo_pago_id") && <td>{item.tipo_pago || "—"}</td>}
                  <td>
                    <span className={`status-badge ${item.activo ? "radicado" : "doc-pendiente"}`}>
                      {item.activo ? "Activo" : "Inactivo"}
                    </span>
                  </td>
                  <td>
                    <div className="catalogo-actions">
                      <button className="btn-icon btn-edit" onClick={() => openCatalogoEdit(item)} title="Editar">
                        <i className="fa-solid fa-pen"></i>
                      </button>
                      <button className={`btn-icon btn-toggle ${item.activo ? "active" : ""}`} onClick={() => handleToggleCatalogoStatus(item)} title={item.activo ? "Desactivar" : "Activar"}>
                        <i className={`fa-solid ${item.activo ? "fa-check" : "fa-xmark"}`}></i>
                      </button>
                    </div>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </div>
    );
  };

  const renderContent = () => {
    switch (activeTab) {
      case "correos": return renderCorreos();
      case "documentos": return renderDocumentos();
      case "radicados": return renderRadicados();
      case "catalogos": return renderCatalogos();
      case "tareas": return renderTareas();
      default: return renderWelcome();
    }
  };
  const detalleRadicadoActual =
  activeTab === "tareas" ? tareaDetail : radicadoDetail;

const editorPermitido =
  detalleRadicadoActual?.estado_posesion !== "Completado";

  return (
    <>
        <header className="top-header">
        <div className="header-left">
          <button className="menu-toggle" title="Menú" onClick={() => setIsSidebarOpen(!isSidebarOpen)}><i className="fas fa-bars" /></button>
          <img src={logo} alt="Logo" className="logo" />
          <div className="header-text"><h1>SIGEFAE</h1><p>Sistema de Gestion de Facturas Electronicas</p></div>
        </div>
        <div className="header-right">
          <NotificacionesDropdown onNavigate={(id) => {
            setActiveTab("tareas");
            setSelectedTareaId(id);
            setSelectedRadicadoId(null);
          }} />
        </div>
      </header>

      <div className="main-container">
                <aside className={`sidebar ${isSidebarOpen ? "open" : ""}`}>
          <div className="sidebar-header">
            <div className="sidebar-title"><i className="fa-solid fa-gear"></i></div>
          </div>
          <nav className="menu-nav">
            {/* ── Mis Tareas: todos los usuarios (no admin) ── */}
            {esUsuario && (
              <a href="#" className={`menu-item ${activeTab === "tareas" ? "active" : ""}`} onClick={(e) => { e.preventDefault(); setActiveTab("tareas"); setIsSidebarOpen(false); }}>
                <div className="item-icon"><i className="fa-solid fa-clipboard-list"></i></div>
                <div className="item-text"><span className="item-nombre">Mis Tareas</span></div>
              </a>
            )}

            {/* ── Admin: ve el flujo completo de recepción ── */}
            {esAdmin && (
              <>
                <a href="#" className={`menu-item ${activeTab === "correos" ? "active" : ""}`} onClick={(e) => { e.preventDefault(); setActiveTab("correos"); setIsSidebarOpen(false); }}>
                  <div className="item-icon"><i className="fa-solid fa-envelope"></i></div>
                  <div className="item-text"><span className="item-nombre">Correos</span></div>
                </a>
                <a href="#" className={`menu-item ${activeTab === "documentos" ? "active" : ""}`} onClick={(e) => { e.preventDefault(); setActiveTab("documentos"); setIsSidebarOpen(false); }}>
                  <div className="item-icon"><i className="fa-solid fa-file-invoice"></i></div>
                  <div className="item-text"><span className="item-nombre">Documentos</span></div>
                </a>
              </>
            )}

            {/* ── Radicados: visible para todos ── */}
            <a href="#" className={`menu-item ${activeTab === "radicados" ? "active" : ""}`} onClick={(e) => { e.preventDefault(); setActiveTab("radicados"); setIsSidebarOpen(false); }}>
              <div className="item-icon"><i className="fa-solid fa-stamp"></i></div>
              <div className="item-text"><span className="item-nombre">Radicados</span></div>
            </a>

            {/* ── Catálogos: solo admin ── */}
            {esAdmin && (
              <a href="#" className={`menu-item ${activeTab === "catalogos" ? "active" : ""}`} onClick={(e) => { e.preventDefault(); setActiveTab("catalogos"); setIsSidebarOpen(false); }}>
                <div className="item-icon"><i className="fa-solid fa-sliders"></i></div>
                <div className="item-text"><span className="item-nombre">Catálogos</span></div>
              </a>
            )}
          </nav>
        </aside>
        <main className="content-area">
          <div className="content-header">
            <div className="content-title">
              <div className="content-icon"><i className={tabInfo.icon}></i></div>
              <div><h2>{tabInfo.title}</h2><p>{tabInfo.subtitle}</p></div>
            </div>
          </div>
          {renderContent()}
        </main>
      </div>
      {pdfEditor.open && editorPermitido && (
      <PdfEditor
        archivoId={pdfEditor.archivoId}
        radicadoId={pdfEditor.radicadoId}
        onClose={() =>
          setPdfEditor({
            open: false,
            archivoId: null,
            radicadoId: null
          })
        }
        onSaved={() => {
          const radicadoId = pdfEditor.radicadoId;
        
          setPdfEditor({
            open: false,
            archivoId: null,
            radicadoId: null
          });
        
          if (activeTab === "tareas") {
            setSelectedTareaId(null);
            setTimeout(() => setSelectedTareaId(radicadoId), 10);
          } else {
            setSelectedRadicadoId(null);
            setTimeout(() => setSelectedRadicadoId(radicadoId), 10);
          }
        }}
      />
      )}

            {/* ── Modal Crear Documento Manual ── */}
      {showCrearDocModal && (
        <div className="modal-overlay" onClick={() => setShowCrearDocModal(false)}>
          <div className="modal-content" onClick={(e) => e.stopPropagation()} style={{ maxWidth: 800, maxHeight: "90vh", overflowY: "auto" }}>
            <div className="modal-header">
              <h3><i className="fa-solid fa-file-circle-plus"></i> Crear Documento Manual</h3>
              <button className="modal-close" onClick={() => setShowCrearDocModal(false)}><i className="fa-solid fa-xmark"></i></button>
            </div>
            <div className="modal-body">
              <div style={{ display: "grid", gridTemplateColumns: "1fr 1fr", gap: 12 }}>
                <div className="modal-field">
                  <label>Tipo <span className="required">*</span></label>
                  <select name="tipo" value={docForm.tipo} onChange={handleDocFormChange} className="doc-input">
                    <option value="FACTURA_FISICA">Factura Física</option>
                    <option value="CUENTA_COBRO">Cuenta de Cobro</option>
                    <option value="NOTA_CREDITO">Nota Crédito</option>
                    <option value="NOTA_DEBITO">Nota Débito</option>
                  </select>
                </div>
                <div className="modal-field">
                  <label>Número Documento <span className="required">*</span></label>
                  <input type="text" name="numero_documento" value={docForm.numero_documento} onChange={handleDocFormChange} className="doc-input" placeholder="Ej: F001-1234" />
                </div>
                <div className="modal-field">
                  <label>Proveedor <span className="required">*</span></label>
                  <select name="id_proveedor" value={docForm.id_proveedor} onChange={handleDocFormChange} className="doc-input">
                    <option value="">Seleccione...</option>
                    {proveedoresCatalogo.map(p => (
                      <option key={p.id} value={p.id}>{p.razon_social} — {p.numero_documento}</option>
                    ))}
                  </select>
                </div>
                <div className="modal-field">
                  <label>Receptor <span className="required">*</span></label>
                  <select name="id_receptor" value={docForm.id_receptor} onChange={handleDocFormChange} className="doc-input">
                    <option value="">Seleccione...</option>
                    {receptoresCatalogo.map(r => (
                      <option key={r.id} value={r.id}>{r.nombre} — {r.numero_documento}</option>
                    ))}
                  </select>
                </div>
                <div className="modal-field">
                  <label>Área <span className="required">*</span></label>
                  <select name="id_area" value={docForm.id_area} onChange={handleDocFormChange} className="doc-input">
                    <option value="">Seleccione...</option>
                    {areasCatalogo.map(a => (
                      <option key={a.id} value={a.id}>{a.nombre}</option>
                    ))}
                  </select>
                </div>
                <div className="modal-field">
                  <label>Moneda <span className="required">*</span></label>
                  <select name="moneda_id" value={docForm.moneda_id} onChange={handleDocFormChange} className="doc-input">
                    <option value="">Seleccione...</option>
                    {monedasCatalogo.map(m => (
                      <option key={m.id} value={m.id}>{m.nombre} ({m.codigo})</option>
                    ))}
                  </select>
                </div>
                <div className="modal-field">
                  <label>Fecha Emisión <span className="required">*</span></label>
                  <input type="date" name="fecha_documento" value={docForm.fecha_documento} onChange={handleDocFormChange} className="doc-input" />
                </div>
                <div className="modal-field" style={{ gridColumn: "1 / -1" }}>
                  <label>Asunto / Observación</label>
                  <input type="text" name="asunto" value={docForm.asunto} onChange={handleDocFormChange} className="doc-input" placeholder="Descripción general del documento" />
                </div>
              </div>

              {/* ── Detalles ── */}
              <div className="doc-section" style={{ marginTop: 20 }}>
                <div style={{ display: "flex", justifyContent: "space-between", alignItems: "center", marginBottom: 10 }}>
                  <h4><i className="fa-solid fa-list"></i> Detalle de Ítems</h4>
                  <button className="doc-btn doc-btn-secondary" onClick={handleAddDetalle} type="button">
                    <i className="fa-solid fa-plus"></i> Agregar Ítem
                  </button>
                </div>
                {docForm.detalles.length === 0 ? (
                  <p style={{ color: "#6b7280", fontSize: "0.9em" }}>Sin ítems. Agregue al menos uno.</p>
                ) : (
                  <table className="doc-items-table">
                    <thead>
                      <tr>
                        <th>Descripción</th>
                        <th style={{ width: 90 }}>Cant.</th>
                        <th style={{ width: 130 }}>Valor Unit.</th>
                        <th style={{ width: 110 }}>IVA Unit.</th>
                        <th style={{ width: 130 }}>Total Línea</th>
                        <th style={{ width: 40 }}></th>
                      </tr>
                    </thead>
                    <tbody>
                      {docForm.detalles.map((d, idx) => (
                        <tr key={idx}>
                          <td><input type="text" value={d.descripcion} onChange={(e) => handleDetalleChange(idx, "descripcion", e.target.value)} className="doc-input" placeholder="Descripción" /></td>
                          <td><input type="number" value={d.cantidad} onChange={(e) => handleDetalleChange(idx, "cantidad", e.target.value)} className="doc-input" min="0" step="0.01" /></td>
                          <td><input type="number" value={d.valor_unitario} onChange={(e) => handleDetalleChange(idx, "valor_unitario", e.target.value)} className="doc-input" min="0" /></td>
                          <td><input type="number" value={d.iva_unitario} onChange={(e) => handleDetalleChange(idx, "iva_unitario", e.target.value)} className="doc-input" min="0" /></td>
                          <td style={{ textAlign: "right", fontWeight: 600 }}>{formatCurrency(d.total)}</td>
                          <td>
                            <button className="btn-icon btn-toggle" onClick={() => handleRemoveDetalle(idx)} title="Eliminar">
                              <i className="fa-solid fa-xmark"></i>
                            </button>
                          </td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                )}
              </div>

              {/* ── Totales ── */}
              <div className="doc-section" style={{ marginTop: 16, background: "#f9fafb", padding: 12, borderRadius: 8 }}>
                <div className="doc-totals">
                  <div className="doc-total-row"><span>Subtotal</span><span>{formatCurrency(docForm.subtotal)}</span></div>
                  <div className="doc-total-row"><span>IVA</span><span>{formatCurrency(docForm.iva)}</span></div>
                  <div className="doc-total-row total-final"><span>Total Documento</span><span>{formatCurrency(docForm.total)}</span></div>
                </div>
              </div>
            </div>
            <div className="modal-footer">
              <button className="doc-btn doc-btn-secondary" onClick={() => setShowCrearDocModal(false)} disabled={creandoDoc}>Cancelar</button>
              <button className="doc-btn doc-btn-primary" onClick={handleCrearDocumentoSubmit} disabled={creandoDoc}>
                <i className="fa-solid fa-floppy-disk"></i> {creandoDoc ? "Guardando..." : "Guardar Documento"}
              </button>
            </div>
          </div>
        </div>
      )}

      {/* ── Modal de Devolución (global, accesible desde cualquier tab) ── */}
      {showDevolverModal && (
        <div className="modal-overlay" onClick={() => setShowDevolverModal(false)}>
          <div className="modal-content" onClick={(e) => e.stopPropagation()}>
            <div className="modal-header">
              <h3><i className="fa-solid fa-reply"></i> Devolver Tarea</h3>
              <button className="modal-close" onClick={() => setShowDevolverModal(false)}><i className="fa-solid fa-xmark"></i></button>
            </div>
            <div className="modal-body">
              <div className="modal-field">
                <label>Devolver a Paso / Usuario <span className="required">*</span></label>
                <select
                  value={devolverForm.tarea_destino_id}
                  onChange={(e) => setDevolverForm(prev => ({ ...prev, tarea_destino_id: e.target.value }))}
                  className="doc-input"
                >
                  <option value="">Seleccione paso de destino...</option>
                  {tareasFlujo
                    .filter(t => {
                      const tareaEnProceso = tareasFlujo.find(x => x.estado?.nombre === "En Proceso");
                      return t.estado?.nombre === "Completada" && t.id < (tareaEnProceso?.id || 999999);
                    })
                    .map(t => (
                      <option key={t.id} value={t.id}>
                        {t.descripcion} — {t.usuario_asignado?.nombre || "Sin usuario"} ({t.usuario_asignado?.email || ""})
                        {" "} [Completada]
                      </option>
                    ))}
                </select>
              </div>

              <div className="modal-field">
                <label>Motivo de la Devolución <span className="required">*</span></label>
                <textarea
                  value={devolverForm.observacion}
                  onChange={(e) => setDevolverForm(prev => ({ ...prev, observacion: e.target.value }))}
                  className="doc-input"
                  rows={3}
                  placeholder="Describa claramente por qué devuelve esta tarea y qué debe ser corregido..."
                  style={{ width: "100%", fontFamily: "inherit", resize: "vertical" }}
                />
              </div>

              <div className="modal-field" style={{ display: "flex", alignItems: "center", gap: 8, marginTop: 12 }}>
                <input
                  type="checkbox"
                  id="retorno_directo"
                  checked={devolverForm.retorno_directo}
                  onChange={(e) => setDevolverForm(prev => ({ ...prev, retorno_directo: e.target.checked }))}
                />
                <label htmlFor="retorno_directo" style={{ fontSize: "0.9em", color: "#374151", cursor: "pointer", margin: 0 }}>
                  Regresar directamente a mi paso una vez corregido (Smart Return)
                </label>
              </div>
            </div>
            <div className="modal-footer">
              <button className="doc-btn doc-btn-secondary" onClick={() => setShowDevolverModal(false)} disabled={devolviendo}>Cancelar</button>
              <button
                className="doc-btn doc-btn-primary"
                onClick={() => {
                  const tareaActiva = tareasFlujo.find(t => t.estado?.nombre === "En Proceso");
                  const radId = selectedRadicadoId || selectedTareaId;
                  if (tareaActiva && radId) handleDevolverTarea(tareaActiva.id, radId);
                }}
                disabled={devolviendo}
                style={{ background: "var(--pardo-red)", borderColor: "var(--pardo-red)" }}
              >
                <i className="fa-solid fa-reply"></i> {devolviendo ? "Devolviendo..." : "Confirmar Devolución"}
              </button>
            </div>
          </div>
        </div>
      )}
    </>
  );
}