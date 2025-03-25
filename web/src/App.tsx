import React from 'react'
import EmployeeList from './pages/EmployeeList.tsx'
import EmployeeDetail from './pages/EmployeeDetail.tsx'
import Login from './pages/Login.tsx'
import Header from './components/Header'
import Admin from './pages/Admin'
import ProtectedRoute from "@/components/ProtectedRoute.tsx";
import {UserProvider} from "./context/UserContext"
import {BrowserRouter, Route, Routes} from "react-router-dom";

const App: React.FC = () => {
    return (
        <UserProvider>
            <BrowserRouter>
                <Header/>
                <main className="pt-8">
                    <Routes>
                        <Route path="/login" element={<Login/>}/>

                        <Route element={<ProtectedRoute/>}>
                            <Route path="/" element={<EmployeeList/>}/>
                            <Route path="/employees/:id" element={<EmployeeDetail/>}/>
                        </Route>

                        <Route element={<ProtectedRoute requireAdmin={true}/>}>
                            <Route path="admin" element={<Admin/>}/>
                        </Route>
                    </Routes>
                </main>
            </BrowserRouter>
        </UserProvider>
    );
};

export default App
