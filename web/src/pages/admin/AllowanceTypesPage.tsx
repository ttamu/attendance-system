import React, {useEffect, useState} from 'react';
import {createAllowanceType, deleteAllowanceType, fetchAllowanceTypes, updateAllowanceType,} from '@/services/api';
import {AllowanceType} from '@/types/Allowance';
import {Card, CardContent, CardHeader, CardTitle} from '@/components/ui/card';
import {Input} from '@/components/ui/input';
import {Button} from '@/components/ui/button';
import {Label} from '@/components/ui/label';
import {Pencil, Trash2} from 'lucide-react';
import {Table, TableBody, TableCell, TableHead, TableHeader, TableRow,} from '@/components/ui/table';

const AllowanceTypesPage: React.FC = () => {
    // 一覧表示用 state
    const [list, setList] = useState<AllowanceType[]>([]);
    // 新規登録用のフォーム状態
    const [form, setForm] = useState<Partial<AllowanceType>>({
        name: '',
        type: '',
        description: '',
        commission_rate: undefined,
    });
    // 編集対象があれば edit に既存データをセットする
    const [edit, setEdit] = useState<AllowanceType | null>(null);

    useEffect(() => {
        reload();
    }, []);

    const reload = () => {
        fetchAllowanceTypes<AllowanceType[]>().then(setList);
    };

    // 選択中のタイプ: 編集中なら edit.type、なければ form.type を参照
    const selType = edit ? edit.type : form.type;

    const submit = (e: React.FormEvent) => {
        e.preventDefault();
        if (edit) {
            // 編集時は既存データをそのまま update する
            updateAllowanceType<AllowanceType>(edit.id, edit).then(() => {
                setEdit(null);
                reload();
            });
        } else {
            // 新規作成時：
            // 歩合制の場合は、ユーザーには「%」単位で入力してもらい、onChange で既に 100 で割った値が state に入っているのでそのまま送信
            const data: Partial<AllowanceType> = {
                name: form.name,
                type: form.type,
                description: form.description,
                commission_rate: selType === 'commission' ? form.commission_rate ?? 0 : 0,
            };
            createAllowanceType<AllowanceType>(data).then(() => {
                setForm({name: '', type: '', description: '', commission_rate: undefined});
                reload();
            });
        }
    };

    const del = (id: number) => deleteAllowanceType(id).then(reload);

    return (
        <div className="container mx-auto px-4 py-4 space-y-6">
            {/* 登録／編集フォーム */}
            <Card className="w-full shadow-lg">
                <CardHeader className="border-b">
                    <CardTitle className="text-xl font-bold text-gray-900">
                        {edit ? '手当タイプの編集' : '新規手当タイプ追加'}
                    </CardTitle>
                </CardHeader>
                <CardContent>
                    <form onSubmit={submit} className="space-y-4">
                        <div>
                            <Label className="mb-1 text-gray-800">名前</Label>
                            <Input
                                required
                                value={edit ? edit.name : form.name}
                                onChange={(e) =>
                                    edit
                                        ? setEdit({...edit, name: e.target.value})
                                        : setForm({...form, name: e.target.value})
                                }
                            />
                        </div>
                        <div>
                            <Label className="mb-1 text-gray-800">タイプ</Label>
                            <select
                                required
                                className="w-full border border-gray-300 rounded-md px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-black focus:border-black transition"
                                value={edit ? edit.type : form.type}
                                onChange={(e) =>
                                    edit
                                        ? setEdit({...edit, type: e.target.value})
                                        : setForm({...form, type: e.target.value, commission_rate: undefined})
                                }
                            >
                                <option value="">選択してください</option>
                                <option value="commission">歩合制</option>
                                <option value="fixed">固定額</option>
                            </select>
                        </div>
                        <div>
                            <Label className="mb-1 text-gray-800">説明</Label>
                            <Input
                                value={edit ? edit.description : form.description}
                                onChange={(e) =>
                                    edit
                                        ? setEdit({...edit, description: e.target.value})
                                        : setForm({...form, description: e.target.value})
                                }
                            />
                        </div>
                        {/* 歩合率入力欄：歩合制の場合のみ表示、CSS transition で簡易アニメーション */}
                        <div
                            className={`overflow-hidden transition-all duration-300 ${
                                selType === 'commission' ? 'max-h-40 opacity-100 mt-2' : 'max-h-0 opacity-0'
                            }`}
                        >
                            <Label className="mb-1 text-gray-800">歩合率（%）</Label>
                            <Input
                                type="number"
                                value={
                                    edit
                                        ? edit.commission_rate !== undefined
                                            ? edit.commission_rate * 100
                                            : ''
                                        : form.commission_rate !== undefined
                                            ? form.commission_rate * 100
                                            : ''
                                }
                                onChange={(e) => {
                                    const val = e.target.value;
                                    edit
                                        ? setEdit({...edit, commission_rate: val ? Number(val) / 100 : 0})
                                        : setForm({...form, commission_rate: val ? Number(val) / 100 : undefined});
                                }}
                            />
                        </div>
                        <div className="flex items-center justify-end space-x-2">
                            <Button type="submit" className="bg-black text-white flex items-center gap-1 px-4">
                                {edit ? (
                                    <>
                                        <Pencil className="w-4 h-4"/> 更新
                                    </>
                                ) : (
                                    <>追加</>
                                )}
                            </Button>
                            {edit && (
                                <Button variant="outline" onClick={() => setEdit(null)} className="text-gray-700 px-4">
                                    キャンセル
                                </Button>
                            )}
                        </div>
                    </form>
                </CardContent>
            </Card>

            {/* 一覧表示テーブル */}
            <Card className="w-full shadow-lg">
                <CardHeader className="border-b">
                    <CardTitle className="text-xl font-bold text-gray-900">手当タイプ一覧</CardTitle>
                </CardHeader>
                <CardContent>
                    <Table className="w-full">
                        <TableHeader>
                            <TableRow>
                                <TableHead className="px-4 py-2">名前</TableHead>
                                <TableHead className="px-4 py-2">タイプ</TableHead>
                                <TableHead className="px-4 py-2">説明</TableHead>
                                <TableHead className="px-4 py-2">歩合率</TableHead>
                                <TableHead className="px-4 py-2 text-center">操作</TableHead>
                            </TableRow>
                        </TableHeader>
                        <TableBody>
                            {list.length === 0 ? (
                                <TableRow>
                                    <TableCell colSpan={5} className="text-center py-6 text-gray-500">
                                        手当タイプが登録されていません
                                    </TableCell>
                                </TableRow>
                            ) : (
                                list.map((at) => (
                                    <TableRow key={at.id}>
                                        <TableCell className="px-4 py-2">{at.name}</TableCell>
                                        <TableCell className="px-4 py-2">
                                            {at.type === 'commission' ? '歩合制' : '固定額'}
                                        </TableCell>
                                        <TableCell className="px-4 py-2">{at.description}</TableCell>
                                        <TableCell className="px-4 py-2">
                                            {at.type === 'commission'
                                                ? `${(at.commission_rate ?? 0) * 100}%`
                                                : '―'}
                                        </TableCell>
                                        <TableCell className="px-4 py-2 text-center space-x-2">
                                            <Button
                                                variant="link"
                                                onClick={() => setEdit(at)}
                                                className="flex items-center gap-1 text-black"
                                            >
                                                <Pencil className="w-4 h-4"/> 編集
                                            </Button>
                                            <Button
                                                variant="link"
                                                onClick={() => del(at.id)}
                                                className="flex items-center gap-1 text-black"
                                            >
                                                <Trash2 className="w-4 h-4"/> 削除
                                            </Button>
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
}

export default AllowanceTypesPage;
