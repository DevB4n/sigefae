import { jsPDF } from "jspdf";
import autoTable from "jspdf-autotable";
import { PDFDocument } from "pdf-lib";

/**
 * Genera y descarga un Expediente PDF completo para un documento radicado.
 * 
 * @param {Object} radicado - Información del radicado
 * @param {Array} flujo - Lista de tareas del flujo de aprobación
 * @param {Array} trazabilidad - Historial de eventos
 * @param {Array} anexosUrls - Lista de URLs de los anexos PDF para fusionar
 */
export async function generarExpedientePDF(radicado, flujo, trazabilidad, anexosUrls) {
  // 1. Generar Portada
  const docPortada = new jsPDF();
  docPortada.setFontSize(18);
  docPortada.text(`Expediente: ${radicado.numero_radicado}`, 14, 20);
  
  const formatCurrency = (val) => val ? new Intl.NumberFormat('es-CO', { style: 'currency', currency: 'COP', minimumFractionDigits: 0 }).format(val) : "$ 0";

  // Tabla de Información Básica
  const infoData = [
    ["Número Radicado", radicado.numero_radicado || ""],
    ["Fecha Radicación", radicado.fecha_radicacion ? new Date(radicado.fecha_radicacion).toLocaleString() : ""],
    ["Tipo Radicación", radicado.tipo_radicacion?.nombre || ""],
    ["Ruta", radicado.ruta?.nombre || ""],
    ["Método de Pago", radicado.metodo_pago?.nombre || ""],
    ["Estado Actual", radicado.estado_posesion || ""],
    ["Paso Actual", radicado.paso_actual?.nombre || "Inicio"],
    ["Responsable Actual", radicado.usuario_actual?.nombre || ""],
  ];
  if (radicado.documento_comercial) {
    infoData.push(["Doc. Comercial", `${radicado.documento_comercial.tipo || ''} - ${radicado.documento_comercial.numero_documento || ''}`]);
    infoData.push(["Proveedor", radicado.documento_comercial.proveedor?.razon_social || ""]);
    infoData.push(["Receptor", radicado.documento_comercial.receptor?.nombre || "HARINERA PARDO S.A"]);
    infoData.push(["Subtotal", formatCurrency(radicado.documento_comercial.subtotal)]);
    infoData.push(["IVA", formatCurrency(radicado.documento_comercial.iva)]);
    infoData.push(["Total", formatCurrency(radicado.documento_comercial.total)]);
  }

  autoTable(docPortada, {
    startY: 30,
    head: [["Atributo", "Valor"]],
    body: infoData,
    theme: 'grid',
    headStyles: { fillColor: [41, 128, 185] } // Azul
  });

  const portadaBytes = docPortada.output('arraybuffer');

  // 2. Inicializar pdf-lib
  const mergedPdf = await PDFDocument.create();
  
  // Añadir portada
  const portadaDoc = await PDFDocument.load(portadaBytes);
  const portadaPages = await mergedPdf.copyPages(portadaDoc, portadaDoc.getPageIndices());
  portadaPages.forEach((page) => mergedPdf.addPage(page));

  // 3. Añadir anexos (solo PDFs)
  for (const url of anexosUrls) {
    // Si la URL no termina explícitamente en .pdf, asumimos que puede serlo si viene del servidor
    // Para mayor seguridad en la UI filtraremos antes.
    try {
      const fetchRes = await fetch(url, {
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('token') || ''}`
        }
      });
      if (fetchRes.ok) {
        const anexoBytes = await fetchRes.arrayBuffer();
        const anexoDoc = await PDFDocument.load(anexoBytes);
        const anexoPages = await mergedPdf.copyPages(anexoDoc, anexoDoc.getPageIndices());
        anexoPages.forEach((page) => mergedPdf.addPage(page));
      }
    } catch (err) {
      console.error("Error cargando anexo para expediente:", url, err);
    }
  }

  // 4. Generar Flujo y Trazabilidad (Penúltimo y Último)
  const docFlujoTraza = new jsPDF();
  let finalY = 15;

  // Flujo
  if (flujo && flujo.length > 0) {
    docFlujoTraza.setFontSize(14);
    docFlujoTraza.text("Flujo de Aprobación", 14, finalY);

    const flujoData = flujo.map(t => [
      t.usuario_asignado?.nombre || "",
      t.descripcion || "",
      t.estado?.nombre || t.estado_posesion || "Pendiente",
      t.fecha_finalizacion ? new Date(t.fecha_finalizacion).toLocaleString() : "Pendiente"
    ]);

    autoTable(docFlujoTraza, {
      startY: finalY + 5,
      head: [["Usuario", "Acción", "Estado", "Fecha Fin"]],
      body: flujoData,
      theme: 'grid',
      headStyles: { fillColor: [90, 150, 110] } // Verde suave
    });
    finalY = docFlujoTraza.lastAutoTable.finalY;
  }

  // Trazabilidad
  if (trazabilidad && trazabilidad.length > 0) {
    // Si no cabe en la página, forzar salto
    if (finalY > 230) {
      docFlujoTraza.addPage();
      finalY = 15;
    } else {
      finalY += 15;
    }

    docFlujoTraza.setFontSize(14);
    docFlujoTraza.text("Historial de Trazabilidad", 14, finalY);

    const trazaData = trazabilidad.map(t => [
      new Date(t.fecha).toLocaleString(),
      t.usuario_nombre || "Sistema",
      t.accion,
      t.descripcion
    ]);

    autoTable(docFlujoTraza, {
      startY: finalY + 5,
      head: [["Fecha", "Usuario", "Acción", "Detalle"]],
      body: trazaData,
      theme: 'grid',
      headStyles: { fillColor: [180, 90, 90] }, // Rojo suave
      styles: { cellPadding: 2, fontSize: 8 },
      columnStyles: { 
        0: { cellWidth: 35 },
        1: { cellWidth: 35 },
        2: { cellWidth: 40 },
        3: { cellWidth: 'auto' } 
      }
    });
  }

  const trazaBytes = docFlujoTraza.output('arraybuffer');
  const trazaDoc = await PDFDocument.load(trazaBytes);
  const trazaPages = await mergedPdf.copyPages(trazaDoc, trazaDoc.getPageIndices());
  trazaPages.forEach((page) => mergedPdf.addPage(page));

  // 5. Guardar y forzar descarga
  const finalBytes = await mergedPdf.save();
  const blob = new Blob([finalBytes], { type: 'application/pdf' });
  const downloadUrl = URL.createObjectURL(blob);
  
  const a = document.createElement('a');
  a.href = downloadUrl;
  a.download = `Expediente_${radicado.numero_radicado}.pdf`;
  document.body.appendChild(a);
  a.click();
  document.body.removeChild(a);
  URL.revokeObjectURL(downloadUrl);
}
