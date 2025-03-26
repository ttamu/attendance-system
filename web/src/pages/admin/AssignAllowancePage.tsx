import React, {useEffect, useState} from 'react'
import {createEmployeeAllowance, fetchAllowanceTypes, fetchEmployees} from '@/services/api'
import {Employee} from '@/types/Employee'
import {AllowanceType, EmployeeAllowance} from '@/types/Allowance'
import {Card, CardContent, CardHeader, CardTitle} from '@/components/ui/card'
import {Input} from '@/components/ui/input'
import {Button} from '@/components/ui/button'
import {Label} from '@/components/ui/label'

const AssignAllowancePage: React.FC = () => {
    const [employees, setEmployees] = useState<Employee[]>([])
    const [allowanceTypes, setAllowanceTypes] = useState<AllowanceType[]>([])
    const [assignment, setAssignment] = useState<EmployeeAllowance>({
        employee_id: 0,
        allowance_type_id: 0,
        amount: 0,
        commission_rate: 0,
        year: new Date().getFullYear(),
        month: new Date().getMonth() + 1,
    })

    const loadData = async () => {
        try {
            const [emps, ats] = await Promise.all([
                fetchEmployees<Employee[]>(),
                fetchAllowanceTypes<AllowanceType[]>()
            ])
            setEmployees(emps)
            setAllowanceTypes(ats)
        } catch (error) {
            console.error('データの取得に失敗:', error)
        }
    }

    useEffect(() => {
        loadData()
    }, [])

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault()
        try {
            await createEmployeeAllowance<EmployeeAllowance>(assignment)
            setAssignment({
                employee_id: 0,
                allowance_type_id: 0,
                amount: 0,
                commission_rate: 0,
                year: new Date().getFullYear(),
                month: new Date().getMonth() + 1,
            })
            alert('手当割り当てが登録されました')
        } catch (error) {
            console.error('手当割り当ての登録に失敗:', error)
        }
    }

    return (
        <div className="p-8">
            <Card className="w-full max-w-xl mx-auto shadow-lg">
                <CardHeader className="border-b">
                    <CardTitle className="text-xl font-bold text-gray-900">従業員への手当割り当て</CardTitle>
                </CardHeader>
                <CardContent>
                    <form onSubmit={handleSubmit} className="mb-6 space-y-4">
                        <div>
                            <Label className="mb-1 text-gray-800">従業員</Label>
                            <select
                                value={assignment.employee_id}
                                onChange={(e) => setAssignment({...assignment, employee_id: Number(e.target.value)})}
                                className="border p-2 w-full rounded"
                                required
                            >
                                <option value="">選択してください</option>
                                {employees.map(emp => (
                                    <option key={emp.id} value={emp.id}>
                                        {emp.name}
                                    </option>
                                ))}
                            </select>
                        </div>
                        <div>
                            <Label className="mb-1 text-gray-800">手当タイプ</Label>
                            <select
                                value={assignment.allowance_type_id}
                                onChange={(e) => setAssignment({
                                    ...assignment,
                                    allowance_type_id: Number(e.target.value)
                                })}
                                className="border p-2 w-full rounded"
                                required
                            >
                                <option value="">選択してください</option>
                                {allowanceTypes.map(at => (
                                    <option key={at.id} value={at.id}>
                                        {at.name}
                                    </option>
                                ))}
                            </select>
                        </div>
                        <div>
                            <Label className="mb-1 text-gray-800">金額</Label>
                            <Input
                                type="number"
                                value={assignment.amount}
                                onChange={(e) => setAssignment({...assignment, amount: Number(e.target.value)})}
                                required
                            />
                        </div>
                        <div>
                            <Label className="mb-1 text-gray-800">歩合率</Label>
                            <Input
                                type="number"
                                value={assignment.commission_rate}
                                onChange={(e) => setAssignment({
                                    ...assignment,
                                    commission_rate: Number(e.target.value)
                                })}
                            />
                        </div>
                        <div className="flex space-x-2">
                            <div className="flex-1">
                                <Label className="mb-1 text-gray-800">年</Label>
                                <Input
                                    type="number"
                                    value={assignment.year}
                                    onChange={(e) => setAssignment({...assignment, year: Number(e.target.value)})}
                                    required
                                />
                            </div>
                            <div className="flex-1">
                                <Label className="mb-1 text-gray-800">月</Label>
                                <Input
                                    type="number"
                                    value={assignment.month}
                                    onChange={(e) => setAssignment({...assignment, month: Number(e.target.value)})}
                                    required
                                />
                            </div>
                        </div>
                        <Button type="submit" className="bg-black text-white w-full">
                            割り当て登録
                        </Button>
                    </form>
                </CardContent>
            </Card>
        </div>
    )
}

export default AssignAllowancePage
