import { useState, useRef, useEffect, useCallback } from "react";
import { PDFDocument, rgb, StandardFonts } from "pdf-lib";
import * as pdfjsLib from "pdfjs-dist";
import SignatureCanvas from "react-signature-canvas";
import { obtenerToken } from "../modules/auth/token.js";
import "./PdfEditor.css";

// Worker de pdf.js — esta forma funciona tanto con Vite como con Webpack 5
pdfjsLib.GlobalWorkerOptions.workerSrc = new URL(
  "pdfjs-dist/build/pdf.worker.min.mjs",
  import.meta.url
).toString();

const API = "http://localhost:8080/api";

const hexToRgb = (hex) => {
  const r = parseInt(hex.slice(1, 3), 16) / 255;
  const g = parseInt(hex.slice(3, 5), 16) / 255;
  const b = parseInt(hex.slice(5, 7), 16) / 255;
  return rgb(r, g, b);
};

const STAMP_STORAGE_KEY = "pdfEditor:customStamp";
const SIGNATURE_STORAGE_KEY = "pdfEditor:customSignature"; 

const persistStamp = (stamp) => {
  try {
    localStorage.setItem(
      STAMP_STORAGE_KEY,
      JSON.stringify({ dataUrl: stamp.dataUrl, mimeType: stamp.mimeType })
    );
  } catch (e) {
    console.warn("No se pudo guardar el sello en localStorage:", e);
  }
};

const loadStampFromStorage = () => {
  try {
    const raw = localStorage.getItem(STAMP_STORAGE_KEY);
    if (!raw) return null;
    const { dataUrl, mimeType } = JSON.parse(raw);
    const base64 = dataUrl.split(",")[1];
    const bytes = Uint8Array.from(atob(base64), (c) => c.charCodeAt(0));
    return { dataUrl, bytes, mimeType };
  } catch {
    return null;
  }
};

const persistSignature = (signature) => {
  try {
    localStorage.setItem(
      SIGNATURE_STORAGE_KEY,
      JSON.stringify({ dataUrl: signature.dataUrl, mimeType: signature.mimeType })
    );
  } catch (e) {
    console.warn("No se pudo guardar la firma en localStorage:", e);
  }
};

const loadSignatureFromStorage = () => {
  try {
    const raw = localStorage.getItem(SIGNATURE_STORAGE_KEY);
    if (!raw) return null;
    const { dataUrl, mimeType } = JSON.parse(raw);
    const base64 = dataUrl.split(",")[1];
    const bytes = Uint8Array.from(atob(base64), (c) => c.charCodeAt(0));
    return { dataUrl, bytes, mimeType };
  } catch {
    return null;
  }
};

export default function PdfEditor({ archivoId, radicadoId, onClose, onSaved }) {
  const [pdfBytes, setPdfBytes] = useState(null);
  const [tool, setTool] = useState(null);
  const [annotations, setAnnotations] = useState([]);
  const [showTextModal, setShowTextModal] = useState(false);
  const [showSigModal, setShowSigModal] = useState(false);
  const [textValue, setTextValue] = useState("");
  const [textColor, setTextColor] = useState("#000000");
  const [tempPos, setTempPos] = useState(null);

  const [customStamp, setCustomStamp] = useState(loadStampFromStorage); // { dataUrl, bytes, mimeType }
  const [customSignature, setCustomSignature] = useState(loadSignatureFromStorage); // { dataUrl, bytes, mimeType }
  const stampInputRef = useRef(null);

  const [dragging, setDragging] = useState(null);
  const [resizing, setResizing] = useState(null);

  // Tamaño real de la página del PDF (en puntos) y escala de renderizado en pantalla
  const [pageSize, setPageSize] = useState({ w: 0, h: 0 });
  const [renderScale, setRenderScale] = useState(1);

  const wrapperRef = useRef(null);
  const viewerWrapRef = useRef(null);
  const overlayRef = useRef(null);
  const canvasRef = useRef(null);
  const sigCanvas = useRef(null);
  const pdfDocRef = useRef(null); // documento pdf.js cacheado, para re-renderizar en resize

  // ── Cargar y renderizar el PDF en un <canvas> (reemplaza al iframe) ──
  useEffect(() => {
    let cancelled = false;

    const load = async () => {
      const res = await fetch(`${API}/archivo/${archivoId}/download?t=${Date.now()}`, {
        headers: { Authorization: `Bearer ${obtenerToken()}` },
      });
      const buf = await res.arrayBuffer();
      if (cancelled) return;

      setPdfBytes(new Uint8Array(buf)); // Uint8Array nunca se transfiere/consume

      // pdf.js necesita su propia copia del buffer (lo consume internamente)
      const pdf = await pdfjsLib.getDocument({ data: buf.slice(0) }).promise;
      if (cancelled) return;
      pdfDocRef.current = pdf;

      await renderPage();
    };

    load();
    return () => {
      cancelled = true;
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [archivoId]);

  const renderPage = useCallback(async () => {
    const pdf = pdfDocRef.current;
    if (!pdf) return;

    const page = await pdf.getPage(1);
    const viewport1 = page.getViewport({ scale: 1 });

    const targetWidth = viewerWrapRef.current?.clientWidth || viewport1.width;
    const scale = targetWidth / viewport1.width;

    setPageSize({ w: viewport1.width, h: viewport1.height });
    setRenderScale(scale);

    const viewport = page.getViewport({ scale });
    const canvas = canvasRef.current;
    canvas.width = viewport.width;
    canvas.height = viewport.height;
    const ctx = canvas.getContext("2d");
    await page.render({ canvasContext: ctx, viewport }).promise;
  }, []);

  // Re-renderizar si cambia el ancho disponible (p. ej. resize de ventana)
  useEffect(() => {
    const onResize = () => renderPage();
    window.addEventListener("resize", onResize);
    return () => window.removeEventListener("resize", onResize);
  }, [renderPage]);

  // ── Drag & Resize globales ──
  useEffect(() => {
    const onMove = (e) => {
      if (dragging) {
        const dx = e.clientX - dragging.startX;
        const dy = e.clientY - dragging.startY;
        setAnnotations((prev) =>
          prev.map((a) =>
            a.id === dragging.id
              ? { ...a, x: dragging.startAx + dx, y: dragging.startAy + dy }
              : a
          )
        );
      }
      if (resizing) {
        const dx = e.clientX - resizing.startX;
        const dy = e.clientY - resizing.startY;
        setAnnotations((prev) =>
          prev.map((a) =>
            a.id === resizing.id
              ? {
                  ...a,
                  width: Math.max(60, resizing.startW + dx),
                  height: Math.max(50, resizing.startH + dy),
                }
              : a
          )
        );
      }
    };
    const onUp = () => {
      setDragging(null);
      setResizing(null);
    };
    window.addEventListener("mousemove", onMove);
    window.addEventListener("mouseup", onUp);
    return () => {
      window.removeEventListener("mousemove", onMove);
      window.removeEventListener("mouseup", onUp);
    };
  }, [dragging, resizing]);

  // ── Sello PNG válido ──
  const generateStamp = useCallback(() => {
    const canvas = document.createElement("canvas");
    canvas.width = 220;
    canvas.height = 220;
    const ctx = canvas.getContext("2d");

    ctx.strokeStyle = "#dc2626";
    ctx.lineWidth = 5;
    ctx.beginPath();
    ctx.arc(110, 110, 100, 0, Math.PI * 2);
    ctx.stroke();

    ctx.fillStyle = "#dc2626";
    ctx.font = "bold 26px sans-serif";
    ctx.textAlign = "center";
    ctx.textBaseline = "middle";
    ctx.fillText("APROBADO", 110, 110);

    const dataUrl = canvas.toDataURL("image/png");
    const base64 = dataUrl.split(",")[1];
    const bytes = Uint8Array.from(atob(base64), (c) => c.charCodeAt(0));
    return { dataUrl, bytes, mimeType: "image/png" };
  }, []);

  const addSignatureWithDate = (signatureData, pos) => {
  const sigId = Date.now();
  const sigWidth = 200;
  const sigHeight = 80;

  const fecha = new Date().toLocaleDateString("es-CO"); // "30/07/2026"

  setAnnotations((prev) => [
    ...prev,
    {
      id: sigId,
      type: "signature",
      imageBytes: signatureData.bytes,
      imageDataUrl: signatureData.dataUrl,
      imageMime: signatureData.mimeType || "image/png",
      x: pos.x,
      y: pos.y,
      width: sigWidth,
      height: sigHeight,
    },
    {
      id: sigId + 1,
      type: "text",
      text: fecha,
      color: "#374151",
      x: pos.x,
      y: pos.y + sigHeight + 4, // justo debajo de la firma
      width: sigWidth,
      height: 20,
    },
  ]);
  };

  // ── Click en overlay: coordenadas EXACTAS respecto al canvas ──
  const handleOverlayClick = (e) => {
    if (!tool || !overlayRef.current) return;
    if (e.target.closest(".pdf-annotation")) return;

    // offsetX/Y son exactas respecto al elemento overlay, que ahora tiene
    // EXACTAMENTE el mismo tamaño en píxeles que el canvas del PDF.
    const x = e.nativeEvent.offsetX;
    const y = e.nativeEvent.offsetY;

    const base = {
      id: Date.now(),
      x,
      y,
      width: tool === "text" ? 200 : 120,
      height: tool === "text" ? 36 : 120,
    };

    setTempPos({ x, y });

    if (tool === "text") {
      setTextValue("");
      setTextColor("#000000");
      setShowTextModal(true);
    } else if (tool === "signature") {
      setTempPos({ x, y });
      setShowSigModal(true);
    } else if (tool === "stamp") {
      let stampData = customStamp;
      if (!stampData) stampData = generateStamp(); // { dataUrl, bytes, mimeType: "image/png" }
      setAnnotations((prev) => [
        ...prev,
        {
          ...base,
          type: "stamp",
          imageBytes: stampData.bytes,
          imageDataUrl: stampData.dataUrl,
          imageMime: stampData.mimeType || "image/png",
        },
      ]);
      setTool(null);
    }
  };

  const confirmText = () => {
    if (!textValue.trim()) return;
    setAnnotations((prev) => [
      ...prev,
      {
        id: Date.now(),
        type: "text",
        text: textValue,
        color: textColor,
        x: tempPos.x,
        y: tempPos.y,
        width: Math.max(200, textValue.length * 9),
        height: 36,
      },
    ]);
    setShowTextModal(false);
    setTool(null);
    setTextValue("");
  };

  const confirmSignature = () => {
    let signatureData = customSignature;

    if (sigCanvas.current && !sigCanvas.current.isEmpty()) {
      const dataUrl = sigCanvas.current.toDataURL("image/png");
      const base64 = dataUrl.split(",")[1];
      const bytes = Uint8Array.from(atob(base64), (c) => c.charCodeAt(0));
      signatureData = { dataUrl, bytes, mimeType: "image/png" };
    }

    if (!signatureData) {
      alert("Dibuja algo primero o guarda una firma antes de continuar");
      return;
    }

    setCustomSignature(signatureData);
    persistSignature(signatureData);

    addSignatureWithDate(signatureData, tempPos);

    setShowSigModal(false);
    setTool(null);
  };

  const removeAnnotation = (id, e) => {
    e.stopPropagation();
    setAnnotations((prev) => prev.filter((a) => a.id !== id));
  };

  const startDrag = (e, id) => {
    e.stopPropagation();
    const ann = annotations.find((a) => a.id === id);
    setDragging({
      id,
      startX: e.clientX,
      startY: e.clientY,
      startAx: ann.x,
      startAy: ann.y,
    });
  };

  const startResize = (e, id) => {
    e.stopPropagation();
    const ann = annotations.find((a) => a.id === id);
    setResizing({
      id,
      startX: e.clientX,
      startY: e.clientY,
      startW: ann.width,
      startH: ann.height,
    });
  };

  const applyChanges = async () => {
    if (!pdfBytes || annotations.length === 0) {
      alert("No hay cambios para aplicar");
      return;
    }
    const doc = await PDFDocument.load(pdfBytes);
    const page = doc.getPages()[0];
    const { height: pdfH } = page.getSize();
    const font = await doc.embedFont(StandardFonts.Helvetica);
    const scale = 1 / renderScale;

    try {
      for (const ann of annotations) {
        const pdfX = ann.x * scale;
        const pdfY = pdfH - (ann.y + ann.height) * scale;

    if (ann.type === "text") {
      const fontSize = Math.round(ann.height * scale * 0.8);
      page.drawText(ann.text, {
        x: pdfX, y: pdfY,
        size: Math.max(8, fontSize),
        font, color: hexToRgb(ann.color),
      });
    } else if ((ann.type === "signature" || ann.type === "stamp") && ann.imageBytes) {
      const mime = ann.imageMime || "image/png";
      const img = mime.includes("jpeg") || mime.includes("jpg")
        ? await doc.embedJpg(ann.imageBytes)
        : await doc.embedPng(ann.imageBytes);

      page.drawImage(img, {
        x: pdfX, y: pdfY,
        width: ann.width * scale,
        height: ann.height * scale,
      });
    }
  }

      const newBytes = await doc.save();

      // Actualizar pdfBytes con el PDF modificado, para que ediciones
      // sucesivas no sobreescriban con el original.
      setPdfBytes(new Uint8Array(newBytes));
      setAnnotations([]);

      const blob = new Blob([newBytes], { type: "application/pdf" });
      const formData = new FormData();
      formData.append("file", blob, "documento_editado.pdf");

      const res = await fetch(`${API}/archivo/${archivoId}/reemplazar`, {
      method: "PATCH",
      headers: { Authorization: `Bearer ${obtenerToken()}` },
      body: formData,
  });
    if (!res.ok) {
      const errData = await res.json();
      throw new Error(errData.error || "Error reemplazando archivo");
  }
    alert("PDF actualizado correctamente");
    onSaved && onSaved();
} catch (err) {
  console.error(err);
  alert("Error al generar el PDF: " + err.message);
}
  };
    const handleStampUpload = (e) => {
    const file = e.target.files[0];
    if (!file) return;

    const reader = new FileReader();
    reader.onload = (ev) => {
      const dataUrl = ev.target.result;
      const base64 = dataUrl.split(",")[1];
      const bytes = Uint8Array.from(atob(base64), (c) => c.charCodeAt(0));
      const mimeType = file.type; // "image/png" | "image/jpeg"
      const stamp = { dataUrl, bytes, mimeType };
      setCustomStamp(stamp);
      persistStamp(stamp); // ver punto 2
      setTool("stamp");
    };
    reader.readAsDataURL(file);
    e.target.value = "";
  };

  const handleRemoveStamp = () => {
  setCustomStamp(null);
  localStorage.removeItem(STAMP_STORAGE_KEY);
  if (tool === "stamp") setTool(null);
};

  const handleRemoveSignature = () => {
    setCustomSignature(null);
    localStorage.removeItem(SIGNATURE_STORAGE_KEY);
    if (tool === "signature") setTool(null);
  };

  useEffect(() => {
    if (!showSigModal || !sigCanvas.current) return;

    try {
      sigCanvas.current.clear();
      if (customSignature?.dataUrl) {
        sigCanvas.current.fromDataURL(customSignature.dataUrl);
      }
    } catch {
      // Si la imagen guardada no es válida, simplemente se deja el canvas vacío.
    }
  }, [showSigModal, customSignature]);

  return (
    <div className="pdf-editor-overlay">
      <div className="pdf-editor-box">
        {/* Toolbar */}
        <div className="pdf-toolbar">
          <button
            className={tool === "text" ? "active" : ""}
            onClick={() => setTool(tool === "text" ? null : "text")}
          >
            <i className="fa-solid fa-font"></i> Texto
          </button>
          <button
            className={tool === "signature" ? "active" : ""}
            onClick={() => setTool(tool === "signature" ? null : "signature")}
          >
            <i className="fa-solid fa-signature"></i> {customSignature ? "Firma" : "Crear Firma"}
          </button>
          {customSignature && (
            <button onClick={handleRemoveSignature} title="Quitar firma guardada">
              <i className="fa-solid fa-trash"></i>
            </button>
          )}
                    {/* Input oculto para subir sello personalizado */}
          <input
            type="file"
            accept="image/png,image/jpeg,image/jpg"
            style={{ display: "none" }}
            ref={stampInputRef}
            onChange={handleStampUpload}
          />

          <button
            className={tool === "stamp" ? "active" : ""}
            onClick={() => {
              if (!customStamp) {
                stampInputRef.current?.click();
              } else {
                setTool(tool === "stamp" ? null : "stamp");
              }
            }}
          >
            <i className="fa-solid fa-stamp"></i>{" "}
            {customStamp ? "Sello" : "Cargar Sello"}
          </button>
          {customStamp && (
          <button onClick={handleRemoveStamp} title="Quitar sello guardado">
            <i className="fa-solid fa-trash"></i>
          </button>
            )}

          <div style={{ flex: 1 }} />

          <button className="btn-save-pdf" onClick={applyChanges}>
            <i className="fa-solid fa-floppy-disk"></i> Guardar PDF
          </button>
          <button className="btn-close-pdf" onClick={onClose}>
            <i className="fa-solid fa-xmark"></i> Cerrar
          </button>
        </div>

        {/* Visor scrolleable */}
        <div className="pdf-viewer-wrap" ref={viewerWrapRef}>
          <div
            className="pdf-page-wrapper"
            ref={wrapperRef}
            style={{
              width: pageSize.w * renderScale || "100%",
              height: pageSize.h * renderScale || "auto",
            }}
          >
            <canvas ref={canvasRef} className="pdf-canvas" />

            <div
              ref={overlayRef}
              className={`pdf-overlay ${tool ? "editing" : ""}`}
              onClick={handleOverlayClick}
            />

            {annotations.map((ann) => (
              <div
                key={ann.id}
                className="pdf-annotation"
                style={{
                  left: ann.x,
                  top: ann.y,
                  width: ann.width,
                  height: ann.height,
                  color: ann.type === "text" ? ann.color : "inherit",
                }}
                onMouseDown={(e) => startDrag(e, ann.id)}
              >
                {ann.type === "text" && (
                  <span className="ann-text-content">{ann.text}</span>
                )}
                {(ann.type === "signature" || ann.type === "stamp") && (
                  <img
                    src={ann.imageDataUrl}
                    alt={ann.type}
                    draggable={false}
                  />
                )}

                <div
                  className="resize-handle"
                  onMouseDown={(e) => startResize(e, ann.id)}
                />
                <button
                  className="ann-delete"
                  onClick={(e) => removeAnnotation(ann.id, e)}
                  title="Eliminar"
                >
                  <i className="fa-solid fa-xmark"></i>
                </button>
              </div>
            ))}
          </div>
        </div>
      </div>

      {/* Modal Texto */}
      {showTextModal && (
        <div className="pdf-modal">
          <div className="pdf-modal-content">
            <h4>Insertar texto</h4>
            <input
              autoFocus
              value={textValue}
              onChange={(e) => setTextValue(e.target.value)}
              placeholder="Escribe aquí..."
            />
            <div className="color-row">
              <label>Color:</label>
              <input
                type="color"
                value={textColor}
                onChange={(e) => setTextColor(e.target.value)}
              />
            </div>
            <div className="pdf-modal-actions">
              <button onClick={() => { setShowTextModal(false); setTool(null); }}>
                Cancelar
              </button>
              <button className="btn-save-pdf" onClick={confirmText}>
                Añadir
              </button>
            </div>
          </div>
        </div>
      )}

      {/* Modal Firma */}
      {showSigModal && (
        <div className="pdf-modal">
          <div className="pdf-modal-content">
            <h4>Dibuja tu firma</h4>
            <div className="sig-canvas-wrap">
              <SignatureCanvas
                ref={sigCanvas}
                canvasProps={{ width: 400, height: 150 }}
              />
            </div>
            <div className="pdf-modal-actions">
              <button onClick={() => sigCanvas.current.clear()}>Limpiar</button>
              <button onClick={() => { setShowSigModal(false); setTool(null); }}>
                Cancelar
              </button>
              <button className="btn-save-pdf" onClick={confirmSignature}>
                Confirmar firma
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}