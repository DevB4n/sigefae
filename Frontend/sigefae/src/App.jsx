import { BrowserRouter, Routes, Route } from "react-router-dom";

import AuthGuard from "./modules/authGuard/authGuard.jsx";
import Login from "./modules/login/login.jsx";
import Dashboard from "./modules/dashboard/dashboard.jsx";

function App() {

    return (

        <BrowserRouter>

            <Routes>

                <Route
                    path="/"
                    element={<Login />}
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

        </BrowserRouter>

    );

}

export default App;