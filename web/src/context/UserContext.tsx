import React, {createContext, ReactNode, useEffect, useState} from "react"
import {fetchProfile} from "../services/api"
import {UserProfile} from "../types/User"

interface UserContextProps {
    profile: UserProfile | null
    setProfile: React.Dispatch<React.SetStateAction<UserProfile | null>>
}

export const UserContext = createContext<UserContextProps>({
    profile: null,
    setProfile: () => {
    },
})

export const UserProvider: React.FC<{ children: ReactNode }> = ({children}) => {
    const [profile, setProfile] = useState<UserProfile | null>(null)

    useEffect(() => {
        async function loadProfile() {
            try {
                const data = await fetchProfile<UserProfile>()
                setProfile(data)
            } catch (error) {
                console.error("プロフィール取得失敗:", error)
            }
        }

        loadProfile()
    }, [])

    return (
        <UserContext.Provider value={{profile, setProfile}}>
            {children}
        </UserContext.Provider>
    )
}
