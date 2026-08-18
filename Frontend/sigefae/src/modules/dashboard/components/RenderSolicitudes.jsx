export default function RenderSolicitudes({ solicitudes, loadingSolicitudes, esAdmin, handleDecidir }) {
  return (
    <div className="content-body">
      <h3>Solicitudes de Rechazo</h3>
      {loadingSolicitudes ? <p>Cargando solicitudes...</p> : solicitudes.length === 0 ? <p>No hay solicitudes de rechazo.</p> : (
        <div className="solicitudes-list">
          {solicitudes.map((s) => (
            <div key={s.id} className="solicitud-card">
              <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                <div>
                  <strong>#{s.id}</strong> — <span style={{ color: '#374151' }}>{s.mensaje || 'Sin mensaje'}</span>
                  <div style={{ fontSize: '0.85em', color: '#6b7280' }}>{s.documento_radicado ? `Radicado #${s.documento_radicado.numero_radicado}` : ''}</div>
                </div>
                <div style={{ display: 'flex', gap: 8, alignItems: 'center' }}>
                  <span className={`status-badge ${s.estado === 'Pendiente' ? 'doc-pendiente' : s.estado === 'Aceptada' ? 'radicado' : 'doc-rechazado'}`}>{s.estado}</span>
                  {esAdmin && s.estado === 'Pendiente' && (
                    <>
                      <button className="doc-btn doc-btn-primary" onClick={() => handleDecidir(s.id, true)}>Aceptar</button>
                      <button className="doc-btn doc-btn-danger" onClick={() => handleDecidir(s.id, false)}>Rechazar</button>
                    </>
                  )}
                </div>
              </div>
              <div style={{ marginTop: 8, fontSize: '0.85em', color: '#6b7280' }}>Solicitante: {s.usuario?.nombre || s.usuario_id} — {new Date(s.fecha_creacion).toLocaleString()}</div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}