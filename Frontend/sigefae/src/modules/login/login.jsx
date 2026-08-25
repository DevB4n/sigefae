import { useState } from "react";
import { useNavigate } from "react-router-dom";

import { guardarToken } from "../auth/token";
import { login } from "./login";
import "./login.css";

function Login() {

    const navigate = useNavigate();

    const [email, setEmail] = useState("");
    const [password, setPassword] = useState("");

    async function iniciarSesion() {

        try {

            const data = await login(email, password);

            guardarToken(data.token);

            // ── Guardar datos del usuario para el rol-based dashboard ──
            localStorage.setItem("rol", data.usuario.rol);
            localStorage.setItem("user_id", data.usuario.id);

            navigate("/dashboard");

        } catch (error) {

            alert(error.message);

        }

    }

    return (
        <div className="login-container">

            <div className="logo-container">
                <img
                    src={`${import.meta.env.BASE_URL}logo.png`}
                    alt="Logo"
                />
            </div>

            <div className="form-container">

                <input
                    type="email"
                    placeholder="Correo"
                    value={email}
                    onChange={(e) => setEmail(e.target.value)}
                />

                <input
                    type="password"
                    placeholder="Contraseña"
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                />

                <button
                    onClick={iniciarSesion}
                >
                    Iniciar Sesion
                </button>

            </div>

        </div>
    );

}

export default Login;