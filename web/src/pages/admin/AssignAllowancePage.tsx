import React, {useEffect, useState} from 'react';
import {
    createEmployeeAllowance,
    deleteEmployeeAllowance,
    fetchAllowanceTypes,
    fetchEmployeeAllowances,
    fetchEmployees,
} from '@/services/api';
import {Employee} from '@/types/Employee';
import {AllowanceType, EmployeeAllowance} from '@/types/Allowance';
import {Card, CardContent, CardHeader, CardTitle} from '@/components/ui/card';
import {Input} from '@/components/ui/input';
import {Button} from '@/components/ui/button';
import {Label} from '@/components/ui/label';
import {Trash2} from 'lucide-react';
import {Table, TableBody, TableCell, TableHead, TableHeader, TableRow} from '@/components/ui/table';

const AssignAllowancePage: React.FC = () => {
    const [emps, setEmps] = useState<Employee[]>([]);
    const [types, setTypes] = useState<AllowanceType[]>([]);
    const [asgs, setAsgs] = useState<EmployeeAllowance[]>([]);
    const [form, setForm] = useState<Partial<EmployeeAllowance>>({
        employee_id: undefined,
        allowance_type_id: undefined,
        amount: undefined,
        commission_rate: undefined,
        year: new Date().getFullYear(),
        month: new Date().getMonth() + 1,
    });

    const selType = types.find(t => t.id === form.allowance_type_id);

    const loadData = async () => {
        try {
            const [e, t] = await Promise.all([
                fetchEmployees<Employee[]>(),
                fetchAllowanceTypes<AllowanceType[]>(),
            ]);
            setEmps(e || []);
            setTypes(t || []);
        } catch (err) {
            console.error('データ取得エラー:', err);
        }
    };

    const loadAsgs = async () => {
        try {
            const data = await fetchEmployeeAllowances<EmployeeAllowance[]>();
            setAsgs(data || []);
        } catch (err) {
            console.error('割り当て取得エラー:', err);
            setAsgs([]);
        }
    };

    useEffect(() => {
        loadData();
        loadAsgs();
    }, []);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        if (
            form.employee_id === undefined ||
            form.allowance_type_id === undefined ||
            form.amount === undefined ||
            form.year === undefined ||
            form.month === undefined
        ) {
            alert('必須項目をすべて入力してください');
            return;
        }
        try {
            const sendData: EmployeeAllowance = {
                employee_id: Number(form.employee_id),
                allowance_type_id: Number(form.allowance_type_id),
                amount: Number(form.amount),
                commission_rate:
                    selType?.type === 'commission'
                        ? Number(form.commission_rate) / 100
                        : 0,
                year: Number(form.year),
                month: Number(form.month),
            };
            await createEmployeeAllowance<EmployeeAllowance>(sendData);
            alert('登録しました');
            setForm({
                employee_id: undefined,
                allowance_type_id: undefined,
                amount: undefined,
                commission_rate: undefined,
                year: new Date().getFullYear(),
                month: new Date().getMonth() + 1,
            });
            loadAsgs();
        } catch (err) {
            console.error('登録失敗:', err);
        }
    };

    const handleDelete = async (id: number) => {
        try {
            await deleteEmployeeAllowance(id);
            loadAsgs();
        } catch (err) {
            console.error('削除失敗:', err);
        }
    };

    return (
        <div className="container mx-auto px-4 py-4 space-y-6">
            {/* 登録フォーム */}
            <Card className="w-full shadow-lg">
                <CardHeader className="border-b">
                    <CardTitle className="text-xl font-bold text-gray-900">
                        従業員への手当割り当て
                    </CardTitle>
                </CardHeader>
                <CardContent>
                    <form onSubmit={handleSubmit} className="mb-6 space-y-4">
                        <div>
                            <Label className="mb-1 text-gray-800">従業員</Label>
                            <select
                                value={form.employee_id ?? ""}
                                onChange={e =>
                                    setForm({...form, employee_id: e.target.value ? Number(e.target.value) : undefined})
                                }
                                className="border p-2 w-full rounded"
                                required
                            >
                                <option value="">選択してください</option>
                                {emps.map(emp => (
                                    <option key={emp.id} value={emp.id}>{emp.name}</option>
                                ))}
                            </select>
                        </div>
                        <div>
                            <Label className="mb-1 text-gray-800">手当タイプ</Label>
                            <select
                                value={form.allowance_type_id ?? ""}
                                onChange={e =>
                                    setForm({
                                        ...form,
                                        allowance_type_id: e.target.value ? Number(e.target.value) : undefined
                                    })
                                }
                                className="border p-2 w-full rounded"
                                required
                            >
                                <option value="">選択してください</option>
                                {types.map(t => (
                                    <option key={t.id} value={t.id}>
                                        {t.name} ({t.type === 'commission' ? '歩合制' : '固定'})
                                    </option>
                                ))}
                            </select>
                        </div>
                        <div>
                            <Label className="mb-1 text-gray-800">金額</Label>
                            <Input
                                type="number"
                                value={form.amount ?? ""}
                                onChange={e =>
                                    setForm({...form, amount: e.target.value ? Number(e.target.value) : undefined})
                                }
                                required
                            />
                        </div>
                        <div
                            className={`overflow-hidden transition-all duration-300 ${
                                selType?.type === 'commission' ? 'max-h-40 opacity-100 mt-2' : 'max-h-0 opacity-0'
                            }`}
                        >
                            <Label className="mb-1 text-gray-800">歩合率（%）</Label>
                            <Input
                                type="number"
                                value={form.commission_rate ?? ""}
                                onChange={e =>
                                    setForm({
                                        ...form,
                                        commission_rate: e.target.value ? Number(e.target.value) : undefined
                                    })
                                }
                            />
                        </div>
                        <div className="flex space-x-2">
                            <div className="flex-1">
                                <Label className="mb-1 text-gray-800">年</Label>
                                <Input
                                    type="number"
                                    value={form.year ?? ""}
                                    onChange={e =>
                                        setForm({...form, year: e.target.value ? Number(e.target.value) : undefined})
                                    }
                                    required
                                />
                            </div>
                            <div className="flex-1">
                                <Label className="mb-1 text-gray-800">月</Label>
                                <Input
                                    type="number"
                                    value={form.month ?? ""}
                                    onChange={e =>
                                        setForm({...form, month: e.target.value ? Number(e.target.value) : undefined})
                                    }
                                    required
                                />
                            </div>
                        </div>
                        <Button type="submit" className="bg-black text-white w-full">
                            登録
                        </Button>
                    </form>
                </CardContent>
            </Card>

            {/* 一覧表示テーブル */}
            <Card className="w-full shadow-lg">
                <CardHeader className="border-b">
                    <CardTitle className="text-xl font-bold text-gray-900">
                        付与済み手当一覧
                    </CardTitle>
                </CardHeader>
                <CardContent>
                    <Table className="w-full">
                        <TableHeader>
                            <TableRow>
                                <TableHead className="px-4 py-2">従業員</TableHead>
                                <TableHead className="px-4 py-2">手当タイプ</TableHead>
                                <TableHead className="px-4 py-2">金額</TableHead>
                                <TableHead className="px-4 py-2">歩合率</TableHead>
                                <TableHead className="px-4 py-2">年月</TableHead>
                                <TableHead className="px-4 py-2 text-center">操作</TableHead>
                            </TableRow>
                        </TableHeader>
                        <TableBody>
                            {asgs.length === 0 ? (
                                <TableRow>
                                    <TableCell colSpan={6} className="text-center py-6 text-gray-500">
                                        まだ割り当てがありません
                                    </TableCell>
                                </TableRow>
                            ) : (
                                asgs.map(a => (
                                    <TableRow key={a.id}>
                                        <TableCell className="px-4 py-2">{a.employee_name || a.employee_id}</TableCell>
                                        <TableCell className="px-4 py-2">{a.allowance_type_name}</TableCell>
                                        <TableCell className="px-4 py-2">{a.amount} 円</TableCell>
                                        <TableCell className="px-4 py-2">
                                            {a.allowance_type === 'commission'
                                                ? a.commission_rate !== undefined
                                                    ? `${a.commission_rate * 100}%`
                                                    : ""
                                                : "―"}
                                        </TableCell>
                                        <TableCell className="px-4 py-2">{a.year}年{a.month}月</TableCell>
                                        <TableCell className="px-4 py-2 text-center">
                                            <div className="flex justify-center items-center space-x-2">
                                                <Button
                                                    variant="outline"
                                                    title="削除"
                                                    onClick={() => handleDelete(a.id!)}
                                                    className="p-2 hover:bg-gray-100"
                                                ><Trash2 className="w-5 h-5"/>
                                                </Button>
                                            </div>
                                        </TableCell>
                                    </TableRow>
                                ))
                            )}
                        </TableBody>
                    </Table>
                </CardContent>
            </Card>
        </div>
    );
};

export default AssignAllowancePage;
