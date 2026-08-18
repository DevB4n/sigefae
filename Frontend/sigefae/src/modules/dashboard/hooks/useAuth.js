import { useCallback } from "react";
import { obtenerToken } from "../../auth/token.js";

export const obtenerRol = () => localStorage.getItem("rol") || "";
export const obtenerUserId = () => parseInt(localStorage.getItem("user_id")) || 0;

export function useAuth() {
  const userRol = obtenerRol();
  const userId = obtenerUserId();
  const esAdmin = userRol === "Superadministrador";
  const esUsuario = !esAdmin;
  const showDebug = localStorage.getItem('show_debug') === '1';

  const puedeGestionarRecurso = useCallback((creadoPorId) => {
    if (esAdmin) return true;
    const propietario = Number(creadoPorId || 0);
    return propietario > 0 && propietario === Number(userId);
  }, [esAdmin, userId]);

  return { userRol, userId, esAdmin, esUsuario, showDebug, puedeGestionarRecurso, obtenerToken };
}