import React, {useEffect, useState} from "react";
import {fetchPayroll} from "../services/api";

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

    if (loading) return <p>給与情報を読み込み中...</p>;
    if (error) return <p style={{color: "red"}}>{error}</p>;
    if (!payroll) return null;

    return (
        <div>
            <p>従業員名: {payroll.employee_name}</p>
            <p>支給額: {payroll.gross_salary} 円</p>
            <p>手当合計: {payroll.total_allowance} 円</p>
            <p>健康保険（従業員負担分）: {payroll.health_insurance} 円</p>
            <p>厚生年金（従業員負担分）: {payroll.pension} 円</p>
            <p>控除合計: {payroll.total_deductions} 円</p>
            <p>手取り給与: {payroll.net_salary} 円</p>
        </div>
    );
};

export default PayrollDisplay;
