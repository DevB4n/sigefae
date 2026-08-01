// ============================================
// SONIDO de notificación (Web Audio API)
// ============================================
export function playNotificacionSound() {
  try {
    const AudioContext = window.AudioContext || window.webkitAudioContext;
    if (!AudioContext) return;

    const ctx = new AudioContext();
    const osc = ctx.createOscillator();
    const gain = ctx.createGain();

    osc.connect(gain);
    gain.connect(ctx.destination);

    osc.type = "sine";
    osc.frequency.setValueAtTime(880, ctx.currentTime);
    osc.frequency.exponentialRampToValueAtTime(440, ctx.currentTime + 0.15);

    gain.gain.setValueAtTime(0.25, ctx.currentTime);
    gain.gain.exponentialRampToValueAtTime(0.001, ctx.currentTime + 0.2);

    osc.start();
    osc.stop(ctx.currentTime + 0.2);
  } catch (e) {
    console.log("Audio no disponible", e);
  }
}

// ============================================
// NOTIFICACIÓN NATIVA del navegador
// ============================================
export async function pedirPermisoNotificaciones() {
  if (!("Notification" in window)) return false;
  const permiso = await Notification.requestPermission();
  return permiso === "granted";
}

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