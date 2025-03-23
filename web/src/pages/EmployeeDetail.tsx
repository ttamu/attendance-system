import React, {useEffect, useState} from "react";
import {useParams} from "react-router-dom";
import {fetchEmployeeById} from "../services/api";
import {Employee} from "../types/Employee";
import AttendanceList from "../components/AttendanceList";
import AttendanceForm from "../components/AttendanceForm";

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
        <div className="max-w-2xl mx-auto p-4">
            <h1 className="text-2xl font-bold mb-4">{employee.name}の詳細</h1>

            {/* 勤怠一覧表示 */}
            <section className="mb-8">
                <h2 className="text-xl font-semibold mb-2">勤怠一覧</h2>
                <AttendanceList attendances={employee.attendances || []}/>
            </section>

            {/* 勤怠登録フォーム */}
            <section>
                <h2 className="text-xl font-semibold mb-2">勤怠登録</h2>
                <AttendanceForm employeeId={Number(id)} onAttendanceAdded={handleAttendanceAdded}/>
            </section>
        </div>
    );
};

export default EmployeeDetailPage;
