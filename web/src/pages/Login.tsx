import React, {useContext, useState} from "react";
import {useNavigate} from "react-router-dom";
import {Button} from "@/components/ui/button";
import {Input} from "@/components/ui/input";
import {Card, CardContent, CardHeader, CardTitle} from "@/components/ui/card";
import {Lock, Mail} from "lucide-react";
import {login} from "../services/api";
import {UserContext} from "../context/UserContext";

const Login: React.FC = () => {
    const navigate = useNavigate();
    const [email, setEmail] = useState("");
    const [password, setPassword] = useState("");
    const [error, setError] = useState("");
    const {refreshProfile} = useContext(UserContext);

    const handleLogin = async (e: React.FormEvent) => {
        e.preventDefault();
        try {
            await login<{ message: string }>({email, password});
            await refreshProfile();
            navigate("/");
        } catch (err) {
            setError((err as Error).message);
        }
    };

    return (
        <div className="flex min-h-screen items-center justify-center bg-gray-50">
            <Card className="w-full max-w-sm shadow-lg">
                <CardHeader className="text-center">
                    <CardTitle className="text-2xl font-bold">ログイン</CardTitle>
                </CardHeader>
                <CardContent>
                    {error && <p className="mb-4 text-center text-red-500">{error}</p>}
                    <form onSubmit={handleLogin} className="space-y-4">
                        <div className="flex items-center gap-2">
                            <Mail className="w-5 h-5 text-gray-500"/>
                            <Input
                                type="email"
                                placeholder="メールアドレス"
                                value={email}
                                onChange={(e) => setEmail(e.target.value)}
                                className="w-full"
                                required
                            />
                        </div>
                        <div className="flex items-center gap-2">
                            <Lock className="w-5 h-5 text-gray-500"/>
                            <Input
                                type="password"
                                placeholder="パスワード"
                                value={password}
                                onChange={(e) => setPassword(e.target.value)}
                                className="w-full"
                                required
                            />
                        </div>
                        <Button type="submit" className="w-full">ログイン</Button>
                    </form>
                </CardContent>
            </Card>
        </div>
    );
};

export default Login;
