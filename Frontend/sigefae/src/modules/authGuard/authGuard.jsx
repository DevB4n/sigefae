import { Navigate } from "react-router-dom";
import { obtenerToken } from "../auth/token.js";

function AuthGuard({ children }) {

    const token = obtenerToken();

    if (!token) {
        return <Navigate to="/" replace />;
    }

    return children;
}

export default AuthGuard;