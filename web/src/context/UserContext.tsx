import React, {createContext, ReactNode, useEffect, useState} from "react";
import {fetchProfile} from "../services/api";
import {UserProfile} from "../types/User";

interface UserContextProps {
    profile: UserProfile | null;
    isLoading: boolean;
    setProfile: React.Dispatch<React.SetStateAction<UserProfile | null>>;
    refreshProfile: () => Promise<void>;
}

export const UserContext = createContext<UserContextProps>({
    profile: null,
    isLoading: true,
    setProfile: () => {
    },
    refreshProfile: async () => {
    },
});

export const UserProvider: React.FC<{ children: ReactNode }> = ({children}) => {
    const [profile, setProfile] = useState<UserProfile | null>(null);
    const [isLoading, setIsLoading] = useState<boolean>(true);

    const refreshProfile = async () => {
        try {
            const data = await fetchProfile<UserProfile>();
            setProfile(data);
        } catch (error) {
            console.error("プロフィール取得失敗:", error);
            setProfile(null);
        } finally {
            setIsLoading(false);
        }
    };

    useEffect(() => {
        refreshProfile();
    }, []);

    return (
        <UserContext.Provider value={{profile, isLoading, setProfile, refreshProfile}}>
            {children}
        </UserContext.Provider>
    );
};
