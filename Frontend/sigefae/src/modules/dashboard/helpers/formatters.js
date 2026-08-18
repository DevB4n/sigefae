export const formatCurrency = (val) => {
  return new Intl.NumberFormat("es-CO", { style: "currency", currency: "COP", minimumFractionDigits: 0 }).format(val || 0);
};

export const isFinalState = (s) => (s === "Completado" || s === "Devuelto" || s === "Rechazado");