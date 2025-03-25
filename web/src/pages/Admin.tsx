import React, {useContext} from 'react';
import {Card, CardContent, CardHeader, CardTitle} from '@/components/ui/card';
import {UserContext} from '../context/UserContext';
import {Link} from 'react-router-dom';
import {ClipboardList, PlusCircle} from 'lucide-react';

const AdminPage: React.FC = () => {
    const {profile} = useContext(UserContext);

    if (!profile || !profile.is_admin) {
        return (
            <div className="flex items-center justify-center min-h-screen bg-gray-50">
                <p className="text-lg text-red-600">アクセス権限がありません</p>
            </div>
        );
    }

    return (
        <div className="p-8 bg-gray-50 min-h-screen">
            <Card className="max-w-3xl mx-auto">
                <CardHeader>
                    <CardTitle className="text-2xl font-bold">管理者ページ</CardTitle>
                </CardHeader>
                <CardContent>
                    <p className="mb-6 text-gray-700">
                        ここでは手当タイプの管理や従業員への手当割り当てが行えます。
                    </p>
                    <div className="flex flex-col sm:flex-row gap-4">
                        <Link
                            to="/admin/allowance-types"
                            className="flex items-center gap-2 text-blue-600 hover:underline"
                        >
                            <PlusCircle className="w-5 h-5"/>
                            手当タイプの追加・管理
                        </Link>
                        <Link
                            to="/admin/assign-allowance"
                            className="flex items-center gap-2 text-blue-600 hover:underline"
                        >
                            <ClipboardList className="w-5 h-5"/>
                            従業員への手当割り当て
                        </Link>
                    </div>
                </CardContent>
            </Card>
        </div>
    );
};

export default AdminPage;
