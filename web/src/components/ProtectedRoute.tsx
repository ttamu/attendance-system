import React, {useContext} from "react";
import {Navigate, Outlet} from "react-router-dom";
import {UserContext} from "../context/UserContext";

interface ProtectedRouteProps {
    requireAdmin?: boolean;
}

const ProtectedRoute: React.FC<ProtectedRouteProps> = ({requireAdmin = false}) => {
    const {profile, isLoading} = useContext(UserContext);

    if (isLoading) {
        return <div className="flex items-center justify-center min-h-screen">Loading...</div>;
    }

    if (!profile) {
        return <Navigate to="/login" replace/>;
    }

    if (requireAdmin && !profile.is_admin) {
        return <Navigate to="/" replace/>;
    }

    return <Outlet/>;
};

export default ProtectedRoute;
