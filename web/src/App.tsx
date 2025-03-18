import React from 'react'
import './App.css'
import UserList from './pages/UserList'
import UserDetail from './pages/UserDetail'
import {BrowserRouter, Route, Routes} from "react-router-dom";

const App: React.FC = () => {
    return (
        <BrowserRouter>
            <Routes>
                <Route path="/" element={<UserList/>}/>
                <Route path="/users/:id" element={<UserDetail/>}/>
            </Routes>
        </BrowserRouter>
    );
};

export default App
