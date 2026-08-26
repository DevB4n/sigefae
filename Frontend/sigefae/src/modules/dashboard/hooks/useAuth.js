import { useCallback } from "react";
import { obtenerToken } from "../../auth/token.js";

export const obtenerRol = () => localStorage.getItem("rol") || "";
export const obtenerUserId = () => parseInt(localStorage.getItem("user_id")) || 0;
export const obtenerNombre = () => localStorage.getItem("user_name") || "";

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

  const userName = obtenerNombre();

  return { userRol, userId, userName, esAdmin, esUsuario, showDebug, puedeGestionarRecurso, obtenerToken };
}