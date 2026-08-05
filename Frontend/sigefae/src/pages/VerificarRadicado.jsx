import { useState, useEffect } from "react";
import { useParams } from "react-router-dom";
import styles from "./VerificarRadicado.module.css";

const API = "http://localhost:8080/api";

export default function VerificarRadicado() {
  const { numero } = useParams();
  const [data, setData] = useState(null);
  const [error, setError] = useState(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetch(`${API}/documentoradicado/verificar/${numero}`)
      .then((r) => r.json())
      .then((res) => {
        if (res.valido) {
          setData(res);
        } else {
          setError(res.error || "Documento no válido");
        }
        setLoading(false);
      })
      .catch(() => {
        setError("No se pudo conectar con el servidor de verificación");
        setLoading(false);
      });
  }, [numero]);

  if (loading) {
    return (
      <div className={styles.loading}>
        <div className={styles.spinner}></div>
        <p>Verificando documento en SIGEFAE...</p>
      </div>
    );
  }

  if (error) {
    return (
      <div className={styles.page}>
        <div className={styles.container}>
          <div className={styles.errorBox}>
            <div className={styles.errorIcon}>
              <i className="fa-solid fa-triangle-exclamation"></i>
            </div>
            <h2 className={styles.errorTitle}>Documento No Encontrado</h2>
            <p className={styles.errorText}>{error}</p>
            <span className={styles.errorNumero}>{numero}</span>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className={styles.page}>
      <div className={styles.container}>

        {/* Header */}
        <div className={styles.header}>
          <div className={styles.brand}>
            SIGE<span className={styles.brandAccent}>FAE</span>
          </div>
          <div className={styles.subtitle}>
            Sistema de Gestión de Facturas Electrónicas
          </div>
          <div className={styles.badge}>
            <div className={styles.badgeIcon}>
              <i className="fa-solid fa-check"></i>
            </div>
            Documento Verificado — Existe y es auténtico
          </div>
        </div>

        {/* Radicado */}
        <div className={styles.card}>
          <div className={styles.cardHeader}>
            <div className={styles.cardIcon}>
              <i className="fa-solid fa-stamp"></i>
            </div>
            <h3 className={styles.cardTitle}>Información del Radicado</h3>
          </div>
          <div className={styles.grid}>
            <div className={styles.field}>
              <span className={styles.fieldLabel}>Número de Radicado</span>
              <span className={`${styles.fieldValue} ${styles.fieldValueHighlight}`}>
                {data.numero_radicado}
              </span>
            </div>
            <div className={styles.field}>
              <span className={styles.fieldLabel}>Fecha de Radicación</span>
              <span className={styles.fieldValue}>
                {new Date(data.fecha_radicacion).toLocaleString()}
              </span>
            </div>
            <div className={styles.field}>
              <span className={styles.fieldLabel}>Tipo de Radicación</span>
              <span className={styles.fieldValue}>{data.tipo_radicacion}</span>
            </div>
            <div className={styles.field}>
              <span className={styles.fieldLabel}>Estado del Proceso</span>
              <span className={styles.fieldValue}>{data.estado_posesion}</span>
            </div>
          </div>
        </div>

        {/* Documento Comercial */}
        <div className={styles.card}>
          <div className={styles.cardHeader}>
            <div className={styles.cardIcon}>
              <i className="fa-solid fa-file-invoice"></i>
            </div>
            <h3 className={styles.cardTitle}>Documento Comercial</h3>
          </div>
          <div className={styles.grid}>
            <div className={styles.field}>
              <span className={styles.fieldLabel}>Tipo</span>
              <span className={styles.fieldValue}>{data.documento.tipo}</span>
            </div>
            <div className={styles.field}>
              <span className={styles.fieldLabel}>Número</span>
              <span className={`${styles.fieldValue} ${styles.fieldValueHighlight}`}>
                {data.documento.numero}
              </span>
            </div>
            <div className={styles.field}>
              <span className={styles.fieldLabel}>Fecha de Emisión</span>
              <span className={styles.fieldValue}>
                {new Date(data.documento.fecha_emision).toLocaleDateString()}
              </span>
            </div>
            <div className={styles.field}>
              <span className={styles.fieldLabel}>Valor Total</span>
              <span className={`${styles.fieldValue} ${styles.fieldValueCurrency}`}>
                {new Intl.NumberFormat("es-CO", { style: "currency", currency: "COP" }).format(data.documento.total)} {data.documento.moneda}
              </span>
            </div>
          </div>
        </div>

        {/* Partes */}
        <div className={styles.card}>
          <div className={styles.cardHeader}>
            <div className={styles.cardIcon}>
              <i className="fa-solid fa-building"></i>
            </div>
            <h3 className={styles.cardTitle}>Partes del Documento</h3>
          </div>
          <div className={styles.grid}>
            <div className={styles.field}>
              <span className={styles.fieldLabel}>Proveedor</span>
              <span className={styles.fieldValue}>{data.proveedor.razon_social || "No disponible"}</span>
            </div>
            <div className={styles.field}>
              <span className={styles.fieldLabel}>NIT Proveedor</span>
              <span className={styles.fieldValue}>{data.proveedor.nit || "No disponible"}</span>
            </div>
            <div className={styles.field}>
              <span className={styles.fieldLabel}>Receptor</span>
              <span className={styles.fieldValue}>{data.receptor.nombre || "No disponible"}</span>
            </div>
            <div className={styles.field}>
              <span className={styles.fieldLabel}>NIT Receptor</span>
              <span className={styles.fieldValue}>{data.receptor.nit || "No disponible"}</span>
            </div>
          </div>
        </div>

        {/* CUFE */}
        {data.documento.cufe && (
          <div className={styles.card}>
            <div className={styles.cardHeader}>
              <div className={styles.cardIcon} style={{ background: 'linear-gradient(135deg, #f59e0b, #d97706)' }}>
                <i className="fa-solid fa-shield-halved"></i>
              </div>
              <h3 className={styles.cardTitle}>Verificación de Integridad (DIAN)</h3>
            </div>
            <div className={styles.cufeBox}>
              <div className={styles.cufeHeader}>
                <i className="fa-solid fa-lock"></i>
                Código Único de Factura Electrónica (CUFE)
              </div>
              <code className={styles.cufeCode}>{data.documento.cufe}</code>
              <p className={styles.cufeNote}>
                Este código único permite verificar la autenticidad del documento directamente en los sistemas de la DIAN. Cualquier alteración al documento original invalidaría este hash.
              </p>
            </div>
          </div>
        )}

        {/* Footer */}
        <div className={styles.footer}>
          Verificación generada por <strong>SIGEFAE</strong> • {new Date().toLocaleDateString()}
        </div>

      </div>
    </div>
  );
}