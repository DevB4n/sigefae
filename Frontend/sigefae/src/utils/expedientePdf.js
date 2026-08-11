import { jsPDF } from "jspdf";
import autoTable from "jspdf-autotable";
import { PDFDocument } from "pdf-lib";
import QRCode from "qrcode";

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
  // ═══════════════════════════════════════════════════════════════
  // ═══ QR DE VERIFICACIÓN — SE AGREGA DESPUÉS DE LA TABLA ======
  // ═══════════════════════════════════════════════════════════════
  if (radicado.qr?.url) {
    try {
      // Generar QR como imagen base64 localmente (sin depender de internet)
      const qrDataUrl = await QRCode.toDataURL(radicado.qr.url, {
        width: 200,
        margin: 2,
        color: { dark: '#1e3a5f', light: '#ffffff' }
      });

      let qrY = docPortada.lastAutoTable.finalY + 12;

      // Si no cabe en la página actual, saltar a nueva
      if (qrY > 220) {
        docPortada.addPage();
        qrY = 20;
      }

      // Título del QR
      docPortada.setFontSize(12);
      docPortada.setTextColor(30, 58, 95);
      docPortada.text("Código QR de Verificación del Expediente", 14, qrY);

      // Imagen del QR (50x50 mm)
      docPortada.addImage(qrDataUrl, 'PNG', 14, qrY + 6, 50, 50);

      // URL debajo del QR
      docPortada.setFontSize(8);
      docPortada.setTextColor(100);
      docPortada.text(radicado.qr.url, 14, qrY + 60, { maxWidth: 180 });

    } catch (err) {
      console.error("Error generando QR para expediente:", err);
    }
  }

    // ── NORMAS DE REPARTO ──
  if (radicado.normas_reparto && radicado.normas_reparto.length > 0) {
    const normasData = radicado.normas_reparto.map(n => [
      n.norma_reparto?.codigo || n.codigo || "",
      n.norma_reparto?.nombre || n.nombre || "",
      n.norma_reparto?.sucursal || n.sucursal || "",
      n.norma_reparto?.departamento || n.departamento || "",
      `${parseFloat(n.porcentaje).toFixed(2)}%`
    ]);

    docPortada.addPage();
    docPortada.setFontSize(14);
    docPortada.text("Normas de Reparto", 14, 20);
    
    autoTable(docPortada, {
      startY: 30,
      head: [["Código", "Nombre", "Sede", "Área", "%"]],
      body: normasData,
      theme: 'grid',
      headStyles: { fillColor: [200, 160, 30] },
      styles: { fontSize: 9 },
      columnStyles: { 4: { halign: 'right', fontStyle: 'bold' } }
    });
  }

  const portadaBytes = docPortada.output('arraybuffer');

  // 2. Inicializar pdf-lib
  const mergedPdf = await PDFDocument.create();
  
  // Añadir portada (ahora incluye el QR)
  const portadaDoc = await PDFDocument.load(portadaBytes);
  const portadaPages = await mergedPdf.copyPages(portadaDoc, portadaDoc.getPageIndices());
  portadaPages.forEach((page) => mergedPdf.addPage(page));

  // 3. Añadir anexos (solo PDFs)
  for (const url of anexosUrls) {
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