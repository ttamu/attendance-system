import React from 'react'
import EmployeeList from './pages/EmployeeList.tsx'
import EmployeeDetail from './pages/EmployeeDetail.tsx'
import Login from './pages/Login.tsx'
import {UserProvider} from "./context/UserContext"
import {BrowserRouter, Route, Routes} from "react-router-dom";

const App: React.FC = () => {
    return (
        <UserProvider>
            <BrowserRouter>
                <main className="pt-8">
                    <Routes>
                        <Route path="/" element={<EmployeeList/>}/>
                        <Route path="/login" element={<Login/>}/>
                        <Route path="/employees/:id" element={<EmployeeDetail/>}/>
                    </Routes>
                </main>
            </BrowserRouter>
        </UserProvider>
    );
};

export default App
