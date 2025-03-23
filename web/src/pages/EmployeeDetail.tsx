import React, {useEffect, useState} from "react";
import {useParams} from "react-router-dom";
import {fetchEmployeeById} from "../services/api";
import {Employee} from "../types/Employee";
import PayrollDisplay from "../components/PayrollDisplay";
import AttendanceList from "../components/AttendanceList";
import AttendanceForm from "../components/AttendanceForm";
import {Card, CardContent, CardHeader, CardTitle} from "@/components/ui/card";
import {BadgeInfo, CalendarCheck2, ClipboardList, Wallet} from "lucide-react";

const EmployeeDetailPage: React.FC = () => {
    const {id} = useParams<{ id: string }>();
    const [employee, setEmployee] = useState<Employee | null>(null);
    const [loading, setLoading] = useState<boolean>(true);
    const [error, setError] = useState<string>("");
    const [refresh, setRefresh] = useState<boolean>(false);

    useEffect(() => {
        const loadEmployee = async (): Promise<void> => {
            try {
                const data: Employee = await fetchEmployeeById<Employee>(id as string);
                setEmployee(data);
            } catch (err) {
                if (err instanceof Error) {
                    setError(err.message);
                }
            } finally {
                setLoading(false);
            }
        };
        void loadEmployee();
    }, [id, refresh]);

    const handleAttendanceAdded = (): void => {
        setRefresh((prev) => !prev);
    };

    if (loading) return <div className="text-center mt-8">Loading...</div>;
    if (error) return <div className="text-center mt-8 text-red-600">{error}</div>;
    if (!employee) return <div className="text-center mt-8">Not Found</div>;

    return (
        <Card className="max-w-4xl mx-auto shadow-md bg-white">
            <CardHeader className="flex flex-row items-center gap-2 border-b pb-4">
                <BadgeInfo className="w-5 h-5 text-blue-600"/>
                <CardTitle className="text-2xl font-bold text-gray-800">
                    {employee.name}さんの詳細
                </CardTitle>
            </CardHeader>

            <CardContent className="space-y-8 pt-4">
                {/* 給与計算 */}
                <section>
                    <div className="flex items-center gap-2 mb-2">
                        <Wallet className="w-5 h-5 text-blue-600"/>
                        <h2 className="text-xl font-semibold text-gray-800">給与計算結果</h2>
                    </div>
                    <PayrollDisplay employeeId={Number(id)} year={2025} month={3}/>
                </section>

                {/* 勤怠一覧 */}
                <section>
                    <div className="flex items-center gap-2 mb-2">
                        <ClipboardList className="w-5 h-5 text-blue-600"/>
                        <h2 className="text-xl font-semibold text-gray-800">勤怠一覧</h2>
                    </div>
                    <AttendanceList attendances={employee.attendances || []}/>
                </section>

                {/* 勤怠登録 */}
                <section>
                    <div className="flex items-center gap-2 mb-2">
                        <CalendarCheck2 className="w-5 h-5 text-blue-600"/>
                        <h2 className="text-xl font-semibold text-gray-800">勤怠登録</h2>
                    </div>
                    <AttendanceForm employeeId={Number(id)} onAttendanceAdded={handleAttendanceAdded}/>
                </section>
            </CardContent>
        </Card>
    );
};

export default EmployeeDetailPage;
