// src/constants/api.js

// Detecta la ruta base de Vite o ajusta a la subcarpeta actual
const BASE_URL = import.meta.env.BASE_URL || '/';

// Elimina barras duplicadas para evitar //api
export const API = `${BASE_URL.replace(/\/$/, '')}/api`;