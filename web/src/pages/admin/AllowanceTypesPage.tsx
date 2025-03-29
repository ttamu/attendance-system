import React, {useEffect, useState} from 'react'
import {createAllowanceType, deleteAllowanceType, fetchAllowanceTypes, updateAllowanceType} from '@/services/api'
import {AllowanceType} from '@/types/Allowance'
import {Card, CardContent, CardHeader, CardTitle} from '@/components/ui/card'
import {Input} from '@/components/ui/input'
import {Button} from '@/components/ui/button'
import {Label} from '@/components/ui/label'
import {Pencil, Trash2} from 'lucide-react'

export default function AllowanceTypesPage() {
    const [list, setList] = useState<AllowanceType[]>([])
    const [form, setForm] = useState<Partial<AllowanceType>>({name: '', type: '', description: '', commission_rate: 0})
    const [edit, setEdit] = useState<AllowanceType | null>(null)

    useEffect(() => {
        fetchAllowanceTypes<AllowanceType[]>().then(setList)
    }, [])

    const reload = () => fetchAllowanceTypes<AllowanceType[]>().then(setList)

    const submit = (e: React.FormEvent) => {
        e.preventDefault()
        if (edit) {
            updateAllowanceType<AllowanceType>(edit.id, edit).then(() => {
                setEdit(null)
                reload()
            })
        } else {
            createAllowanceType<AllowanceType>(form).then(() => {
                setForm({name: '', type: '', description: '', commission_rate: 0})
                reload()
            })
        }
    }

    const del = (id: number) => deleteAllowanceType(id).then(reload)

    return (
        <div className="p-8 flex flex-col gap-6 items-center">
            <Card className="w-full max-w-xl shadow-lg">
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
                                onChange={e => edit
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
                                        : setForm({...form, type: e.target.value})
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
                                onChange={e => edit
                                    ? setEdit({...edit, description: e.target.value})
                                    : setForm({...form, description: e.target.value})
                                }
                            />
                        </div>
                        <div>
                            <Label className="mb-1 text-gray-800">歩合率</Label>
                            <Input
                                type="number"
                                value={edit ? edit.commission_rate : form.commission_rate}
                                onChange={e => {
                                    const val = Number(e.target.value)
                                    edit
                                        ? setEdit({...edit, commission_rate: val})
                                        : setForm({...form, commission_rate: val})
                                }}
                            />
                        </div>
                        <div className="flex items-center justify-end space-x-2">
                            <Button type="submit" className="bg-black text-white flex items-center gap-1 px-4">
                                {edit ? <><Pencil className="w-4 h-4"/>更新</> : <>追加</>}
                            </Button>
                            {edit &&
                                <Button variant="outline" onClick={() => setEdit(null)} className="text-gray-700 px-4">
                                    キャンセル
                                </Button>
                            }
                        </div>
                    </form>
                </CardContent>
            </Card>

            <Card className="w-full max-w-xl shadow-lg">
                <CardHeader className="border-b">
                    <CardTitle className="text-xl font-bold text-gray-900">手当タイプ一覧</CardTitle>
                </CardHeader>
                <CardContent>
                    <ul className="space-y-4">
                        {list.map(at =>
                            <li key={at.id} className="border p-4 rounded flex justify-between items-center">
                                <div>
                                    <strong className="text-gray-900">{at.name}</strong> ({at.type})<br/>
                                    <span className="text-gray-600">{at.description}</span> -
                                    コミッション率: {at.commission_rate}
                                </div>
                                <div className="flex space-x-2">
                                    <Button variant="link" onClick={() => setEdit(at)}
                                            className="flex items-center gap-1 text-black">
                                        <Pencil className="w-4 h-4"/> 編集
                                    </Button>
                                    <Button variant="link" onClick={() => del(at.id)}
                                            className="flex items-center gap-1 text-black">
                                        <Trash2 className="w-4 h-4"/> 削除
                                    </Button>
                                </div>
                            </li>
                        )}
                    </ul>
                </CardContent>
            </Card>
        </div>
    )
}
