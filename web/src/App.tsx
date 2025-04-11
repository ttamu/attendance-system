import React from 'react'
import EmployeeList from './pages/EmployeeList.tsx'
import EmployeeDetail from './pages/EmployeeDetail.tsx'
import Login from './pages/Login.tsx'
import Header from './components/Header'
import Admin from './pages/Admin'
import AllowanceTypesPage from './pages/admin/AllowanceTypesPage'
import AssignAllowancePage from './pages/admin/AssignAllowancePage'
import ClockRequestList from './pages/admin/ClockRequestList'
import ProtectedRoute from "@/components/ProtectedRoute.tsx";
import {UserProvider} from "./context/UserContext"
import {BrowserRouter, Route, Routes} from "react-router-dom";

const App: React.FC = () => {
    return (
        <UserProvider>
            <BrowserRouter>
                <Header/>
                <main className="p-4 bg-gray-50 min-h-screen">
                    <Routes>
                        <Route path="/login" element={<Login/>}/>

                        <Route element={<ProtectedRoute/>}>
                            <Route path="/" element={<EmployeeList/>}/>
                            <Route path="/employees/:id" element={<EmployeeDetail/>}/>
                        </Route>

                        <Route element={<ProtectedRoute requireAdmin={true}/>}>
                            <Route path="admin" element={<Admin/>}>
                                <Route path="allowance-types" element={<AllowanceTypesPage/>}/>
                                <Route path="assign-allowance" element={<AssignAllowancePage/>}/>
                                <Route path="clock-requests" element={<ClockRequestList/>}/>
                            </Route>
                        </Route>
                    </Routes>
                </main>
            </BrowserRouter>
        </UserProvider>
    );
};

export default App;
