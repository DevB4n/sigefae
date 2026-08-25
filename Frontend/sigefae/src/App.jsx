import { Routes, Route } from "react-router-dom";

import AuthGuard from "./modules/authGuard/authGuard.jsx";
import Login from "./modules/login/login.jsx";
import Dashboard from "./modules/dashboard/components/dashboard.jsx";
import VerificarRadicado from "./pages/VerificarRadicado";

function App() {
    return (
        <Routes>
            <Route
                path="/"
                element={<Login />}
            />

            <Route 
                path="/radicado/:numero" 
                element={<VerificarRadicado />} 
            />

            <Route
                path="/dashboard"
                element={
                    <AuthGuard>
                        <Dashboard />
                    </AuthGuard>
                }
            />
        </Routes>
    );
}

export default App;