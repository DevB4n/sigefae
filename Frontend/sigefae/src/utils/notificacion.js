
let audioCtx = null;
// ============================================
// SONIDO de notificación (archivo personalizado)
// ============================================
let notifAudio = null;

export function inicializarAudio() {
  // Pre-cargamos el audio para que esté listo al primer uso
  if (!notifAudio) {
    notifAudio = new Audio("/notification.mp3");
    notifAudio.volume = 0.6; // ajusta el volumen 0.0 - 1.0
  }

  // Desbloqueamos autoplay con la primera interacción del usuario
  const desbloquear = () => {
    notifAudio.play().then(() => {
      notifAudio.pause();
      notifAudio.currentTime = 0;
    }).catch(() => {});
    window.removeEventListener("click", desbloquear);
    window.removeEventListener("keydown", desbloquear);
  };
  window.addEventListener("click", desbloquear);
  window.addEventListener("keydown", desbloquear);
}

export function playNotificacionSound() {
  try {
    if (!notifAudio) {
      notifAudio = new Audio("/notification.mp3");
      notifAudio.volume = 0.6;
    }
    notifAudio.currentTime = 0; // reinicia si ya estaba sonando
    notifAudio.play().catch((e) => console.log("No se pudo reproducir:", e));
  } catch (e) {
    console.log("Audio no disponible", e);
  }
}

// ============================================
// PERMISO de notificaciones nativas
// ============================================
export async function pedirPermisoNotificaciones() {
  if (!("Notification" in window)) return false;
  const permiso = await Notification.requestPermission();
  return permiso === "granted";
}

// ============================================
// NOTIFICACIÓN NATIVA del navegador
// ============================================
export function mostrarNotificacionNativa(titulo, cuerpo, onClick) {
  if (!("Notification" in window)) return;
  if (Notification.permission !== "granted") return;

  const notif = new Notification(titulo, {
    body: cuerpo,
    icon: "/logo.png",
    tag: "sigefae-" + Date.now(),
  });

  notif.onclick = () => {
    window.focus();
    if (onClick) onClick();
    notif.close();
  };
}