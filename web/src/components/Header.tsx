import React, {useContext} from "react";
import {Link, useNavigate} from "react-router-dom";
import {UserContext} from "../context/UserContext";
import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuItem,
    DropdownMenuLabel,
    DropdownMenuSeparator,
    DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import {Avatar, AvatarFallback, AvatarImage} from "@/components/ui/avatar";
import {Building2, Home, LogOut, Settings, UserIcon} from "lucide-react";
import {logout} from "../services/api";

const Header: React.FC = () => {
    const {profile, setProfile} = useContext(UserContext);
    const navigate = useNavigate();

    const handleLogout = async () => {
        try {
            await logout<{ message: string }>();
            setProfile(null);
            navigate("/login");
        } catch (error) {
            console.error("ログアウトに失敗しました", error);
        }
    };

    const handleAdminPage = () => {
        navigate("/admin");
    };

    return (
        <header className="bg-white shadow px-6 py-3 flex justify-between items-center">
            <div className="flex items-center gap-2">
                <Link to="/"
                      className="flex items-center text-gray-600 px-2 py-1 rounded transition-colors hover:bg-gray-100 hover:text-blue-600">
                    <Home className="w-5 h-5 mr-1 text-blue-600"/>
                    <span className="text-sm font-medium">ホーム</span>
                </Link>
            </div>

            <div className="flex items-center gap-4">
                <div className="flex items-center gap-2 text-gray-800">
                    <Building2 className="w-5 h-5 text-blue-600"/>
                    <span className="text-sm font-medium">{profile?.company?.name}</span>
                </div>

                <DropdownMenu>
                    <DropdownMenuTrigger asChild>
                        <button className="flex items-center focus:outline-none">
                            <Avatar className="w-6 h-6 text-blue-600">
                                <AvatarImage src="/path/to/default/avatar.png" alt="User Avatar"/>
                                <AvatarFallback className="bg-transparent">
                                    <UserIcon/>
                                </AvatarFallback>
                            </Avatar>
                        </button>
                    </DropdownMenuTrigger>
                    <DropdownMenuContent align="end" sideOffset={8} alignOffset={-8} className="w-56">
                        <DropdownMenuLabel className="text-gray-500">ログイン中のユーザー</DropdownMenuLabel>
                        <div className="px-4 py-2 text-sm text-gray-700">{profile?.email}</div>
                        {profile?.is_admin && (
                            <>
                                <DropdownMenuSeparator/>
                                <DropdownMenuItem onClick={handleAdminPage} className="cursor-pointer">
                                    <Settings className="mr-2 w-4 h-4"/>
                                    管理者ページ
                                </DropdownMenuItem>
                            </>
                        )}
                        <DropdownMenuSeparator/>
                        <DropdownMenuItem onClick={handleLogout} className="cursor-pointer">
                            <LogOut className="mr-2 w-4 h-4"/>
                            ログアウト
                        </DropdownMenuItem>
                    </DropdownMenuContent>
                </DropdownMenu>
            </div>
        </header>
    );
};

export default Header;
