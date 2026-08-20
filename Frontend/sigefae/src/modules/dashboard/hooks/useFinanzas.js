import { useState } from "react";
import { API } from "../constants/api.js";

export function useFinanzas(obtenerToken, radicadosHook) {
  const [searchFinanzas, setSearchFinanzas] = useState("");
  const [sortFinanzas, setSortFinanzas] = useState("fecha_fact_desc");

  const handleCausar = async (id, causado, numero_egreso) => {
    try {
      const token = obtenerToken();
      const res = await fetch(`${API}/documentoradicado/${id}/causar`, {
        method: "PATCH",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`
        },
        body: JSON.stringify({ causado, numero_egreso })
      });

      if (!res.ok) {
        const err = await res.json();
        throw new Error(err.error || "Error al actualizar la causación");
      }

      // Actualizar localmente el radicado en el hook de radicados
      radicadosHook.setRadicados((prev) => 
        prev.map((rad) => 
          rad.id === id 
            ? { ...rad, causado, numero_egreso, fecha_causacion: causado && !rad.causado ? new Date().toISOString() : (!causado ? null : rad.fecha_causacion) } 
            : rad
        )
      );

      // Si está viéndolo en detalle, también actualizar
      if (radicadosHook.radicadoDetail && radicadosHook.radicadoDetail.id === id) {
        radicadosHook.setRadicadoDetail({
          ...radicadosHook.radicadoDetail,
          causado,
          numero_egreso,
          fecha_causacion: causado && !radicadosHook.radicadoDetail.causado ? new Date().toISOString() : (!causado ? null : radicadosHook.radicadoDetail.fecha_causacion)
        });
      }

      const AdminToast = document.querySelector(".admin-toast");
      if (AdminToast) {
        const event = new CustomEvent("showToast", { detail: { message: "Causación actualizada", type: "success" } });
        document.dispatchEvent(event);
      }
    } catch (err) {
      alert("Error: " + err.message);
    }
  };

  const handlePagar = async (id, pagado) => {
    try {
      const token = obtenerToken();
      const res = await fetch(`${API}/documentoradicado/${id}/pagar`, {
        method: "PATCH",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`
        },
        body: JSON.stringify({ pagado })
      });

      if (!res.ok) {
        const err = await res.json();
        throw new Error(err.error || "Error al actualizar el pago");
      }

      // Actualizar localmente el radicado en el hook de radicados
      radicadosHook.setRadicados((prev) => 
        prev.map((rad) => 
          rad.id === id 
            ? { ...rad, pagado, fecha_pago: pagado && !rad.pagado ? new Date().toISOString() : (!pagado ? null : rad.fecha_pago) } 
            : rad
        )
      );

      // Si está viéndolo en detalle, también actualizar
      if (radicadosHook.radicadoDetail && radicadosHook.radicadoDetail.id === id) {
        radicadosHook.setRadicadoDetail({
          ...radicadosHook.radicadoDetail,
          pagado,
          fecha_pago: pagado && !radicadosHook.radicadoDetail.pagado ? new Date().toISOString() : (!pagado ? null : radicadosHook.radicadoDetail.fecha_pago)
        });
      }

      const AdminToast = document.querySelector(".admin-toast");
      if (AdminToast) {
        const event = new CustomEvent("showToast", { detail: { message: "Pago actualizado", type: "success" } });
        document.dispatchEvent(event);
      }
    } catch (err) {
      alert("Error: " + err.message);
    }
  };

  const handleComprobantesSubidosLocal = (id) => {
    radicadosHook.setRadicados((prev) => 
      prev.map((rad) => 
        rad.id === id 
          ? { ...rad, comprobantes_subidos: true } 
          : rad
      )
    );
    if (radicadosHook.radicadoDetail && radicadosHook.radicadoDetail.id === id) {
      radicadosHook.setRadicadoDetail({
        ...radicadosHook.radicadoDetail,
        comprobantes_subidos: true
      });
    }
  };

  return {
    searchFinanzas,
    setSearchFinanzas,
    sortFinanzas,
    setSortFinanzas,
    handleCausar,
    handlePagar,
    handleComprobantesSubidosLocal
  };
}
