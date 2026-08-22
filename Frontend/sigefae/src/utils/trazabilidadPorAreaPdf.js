import { jsPDF } from "jspdf";
import autoTable from "jspdf-autotable";

export function generarTrazabilidadPorAreaPDF(areaNombre, fechaDesde, fechaHasta, radicados) {
  const doc = new jsPDF();
  
  // Título principal
  doc.setFontSize(16);
  doc.setTextColor(30, 58, 95);
  doc.text(`Trazabilidad de ${areaNombre}`, 14, 20);
  
  // Subtítulo con fechas
  doc.setFontSize(11);
  doc.setTextColor(80);
  doc.text(`Desde: ${fechaDesde}   Hasta: ${fechaHasta}`, 14, 28);
  
  doc.line(14, 32, 196, 32);
  
  let y = 40;
  
  radicados.forEach((rad, idx) => {
    // Salto de página si es necesario
    if (y > 240) {
      doc.addPage();
      y = 20;
    }
    
    // Encabezado del radicado
    doc.setFontSize(12);
    doc.setTextColor(30, 58, 95);
    doc.text(`${idx + 1}. Radicado: ${rad.numero_radicado}`, 14, y);
    y += 6;
    
    doc.setFontSize(10);
    doc.setTextColor(60);
    const prov = rad.proveedor?.razon_social || "N/A";
    const nit = rad.proveedor?.numero_documento ? ` (NIT: ${rad.proveedor.numero_documento})` : "";
    doc.text(`Proveedor: ${prov}${nit}`, 14, y);
    y += 6;
    
    doc.setFontSize(9);
    doc.setTextColor(100);
    const fechaRad = rad.fecha_radicacion ? new Date(rad.fecha_radicacion).toLocaleDateString("es-CO") : "N/A";
    doc.text(`Estado: ${rad.estado_posesion || "N/A"}  |  Fecha radicación: ${fechaRad}`, 14, y);
    y += 4;
    
    // Tabla de trazabilidad
    if (rad.trazabilidad && rad.trazabilidad.length > 0) {
      const body = rad.trazabilidad.map(t => [
        t.fecha ? new Date(t.fecha).toLocaleString("es-CO") : "",
        t.accion,
        t.descripcion || "",
        t.usuario_nombre || ""
      ]);
      
      autoTable(doc, {
        startY: y,
        head: [["Fecha", "Acción", "Descripción", "Usuario"]],
        body: body,
        theme: "grid",
        headStyles: { fillColor: [41, 128, 185], textColor: 255, fontSize: 9 },
        styles: { fontSize: 8, cellPadding: 2 },
        columnStyles: {
          0: { cellWidth: 32 },
          1: { cellWidth: 30 },
          3: { cellWidth: 28 }
        },
        margin: { left: 14, right: 14 },
      });
      
      y = doc.lastAutoTable.finalY + 10;
    } else {
      doc.setFontSize(9);
      doc.setTextColor(150);
      doc.text("Sin registros de trazabilidad.", 20, y + 6);
      y += 14;
    }
  });
  
  // Pie de página con fecha de generación
  const totalPages = doc.getNumberOfPages();
  for (let i = 1; i <= totalPages; i++) {
    doc.setPage(i);
    doc.setFontSize(8);
    doc.setTextColor(150);
    doc.text(`Generado el ${new Date().toLocaleString("es-CO")} - Página ${i} de ${totalPages}`, 14, 290);
  }
  
  const safeArea = areaNombre.replace(/\s+/g, "_");
  doc.save(`Trazabilidad_${safeArea}_${fechaDesde}_${fechaHasta}.pdf`);
}