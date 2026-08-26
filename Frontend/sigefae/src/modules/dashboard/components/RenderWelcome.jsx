export default function RenderWelcome() {
  const userName = localStorage.getItem("user_name") || "Usuario";

  return (
    <div className="content-body">
      <div className="welcome-wrap">
        <div className="welcome-icon"><i className="fa-solid fa-user"></i></div>
        <h2>Bienvenido, {userName}</h2>
        <p>En la barra a tu izquierda encontrarás todos los procesos de SIGEFAE.</p>
      </div>
    </div>
  );
}