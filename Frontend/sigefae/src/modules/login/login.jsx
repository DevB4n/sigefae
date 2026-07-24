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

            navigate("/dashboard");

        } catch (error) {

            alert(error.message);

        }

    }

    return (
        <div className="login-container">

            <div className="logo-container">
                <img
                    src="./src/assets/login/logo.png"
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