import React, {useEffect, useState} from "react";
import {fetchPayroll} from "../services/api";
import {Card, CardContent} from "@/components/ui/card";

interface PayrollResponse {
    employee_name: string;
    gross_salary: number;
    total_allowance: number;
    health_insurance: number;
    pension: number;
    total_deductions: number;
    net_salary: number;
}

interface PayrollDisplayProps {
    employeeId: number;
    year: number;
    month: number;
}

const PayrollDisplay: React.FC<PayrollDisplayProps> = ({employeeId, year, month}) => {
    const [payroll, setPayroll] = useState<PayrollResponse | null>(null);
    const [loading, setLoading] = useState<boolean>(true);
    const [error, setError] = useState<string>("");

    useEffect(() => {
        const loadPayroll = async () => {
            try {
                const data = await fetchPayroll<PayrollResponse>(employeeId.toString(), year, month);
                setPayroll(data);
            } catch (err) {
                setError((err as Error).message);
            } finally {
                setLoading(false);
            }
        };
        loadPayroll();
    }, [employeeId, year, month]);

    if (loading) return <p className="text-gray-600">給与情報を読み込み中...</p>;
    if (error) return <p className="text-red-600">{error}</p>;
    if (!payroll) return null;

    return (
        <Card className="bg-white shadow-sm border text-gray-800">
            <CardContent className="space-y-2 p-4">
                <div className="flex justify-between">
                    <span>従業員名</span>
                    <span>{payroll.employee_name}</span>
                </div>
                <div className="flex justify-between">
                    <span>支給額</span>
                    <span>{payroll.gross_salary.toLocaleString()} 円</span>
                </div>
                <div className="flex justify-between">
                    <span>手当合計</span>
                    <span>{payroll.total_allowance.toLocaleString()} 円</span>
                </div>
                <div className="flex justify-between">
                    <span>健康保険（従業員負担）</span>
                    <span>{payroll.health_insurance.toLocaleString()} 円</span>
                </div>
                <div className="flex justify-between">
                    <span>厚生年金（従業員負担）</span>
                    <span>{payroll.pension.toLocaleString()} 円</span>
                </div>
                <div className="flex justify-between">
                    <span>控除合計</span>
                    <span>{payroll.total_deductions.toLocaleString()} 円</span>
                </div>
                <div className="flex justify-between border-t pt-2 mt-2">
                    <span>手取り給与</span>
                    <span>{payroll.net_salary.toLocaleString()} 円</span>
                </div>
            </CardContent>
        </Card>
    );
};

export default PayrollDisplay;
