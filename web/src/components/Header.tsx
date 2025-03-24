import React, {useContext} from "react"
import {UserContext} from "../context/UserContext"
import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuItem,
    DropdownMenuLabel,
    DropdownMenuSeparator,
    DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import {Avatar, AvatarFallback, AvatarImage} from "@/components/ui/avatar"
import {Building2, LogOut, UserIcon} from "lucide-react"

const Header: React.FC = () => {
    const {profile} = useContext(UserContext)

    const handleLogout = () => {
        console.log("ログアウト処理の実装")
    }

    return (
        <header className="bg-white shadow px-6 py-3 flex justify-end items-center gap-4">
            <div className="flex items-center gap-2 text-gray-800">
                <Building2 className="w-5 h-5 text-blue-600"/>
                <span className="text-sm font-medium">
          {profile?.company?.name || ""}
        </span>
            </div>

            <DropdownMenu>
                <DropdownMenuTrigger asChild>
                    <button className="flex items-center focus:outline-none">
                        <Avatar className="w-6 h-6 text-blue-600">
                            <AvatarImage src="/path/to/default/avatar.png" alt="User Avatar"/>
                            <AvatarFallback className="bg-transparent"><UserIcon/></AvatarFallback>
                        </Avatar>
                    </button>
                </DropdownMenuTrigger>
                <DropdownMenuContent align="end" sideOffset={8} alignOffset={-8} className="w-56">
                    <DropdownMenuLabel
                        className="text-gray-500">
                        ログイン中のユーザー
                    </DropdownMenuLabel>
                    <div className="px-4 py-2 text-sm text-gray-700">{profile?.email || ""}</div>
                    <DropdownMenuSeparator/>
                    <DropdownMenuItem onClick={handleLogout} className="cursor-pointer">
                        <LogOut className="mr-2 w-4 h-4"/>
                        ログアウト
                    </DropdownMenuItem>
                </DropdownMenuContent>
            </DropdownMenu>
        </header>
    )
}

export default Header
