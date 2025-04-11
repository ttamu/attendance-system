import React, {useEffect, useState} from "react";
import {useParams} from "react-router-dom";
import {fetchEmployeeById, fetchTimeClocks} from "../services/api";
import {Employee} from "../types/Employee";
import {TimeClock} from "../types/TimeClock";
import PayrollDisplay from "../components/PayrollDisplay";
import TimeClockList from "../components/TimeClockList";
import TimeClockForm from "../components/TimeClockForm";
import DateSelector from "../components/DateSelector";
import {Card, CardContent, CardHeader, CardTitle} from "@/components/ui/card";
import {BadgeInfo, CalendarCheck2, ClipboardList, Wallet} from "lucide-react";

const EmployeeDetailPage: React.FC = () => {
    const {id} = useParams<{ id: string }>();
    const [employee, setEmployee] = useState<Employee | null>(null);
    const [timeClocks, setTimeClocks] = useState<TimeClock[]>([]);
    const [loading, setLoading] = useState<boolean>(true);
    const [error, setError] = useState<string>("");
    const [refresh, setRefresh] = useState<boolean>(false);
    const [year, setYear] = useState<number>(new Date().getFullYear());
    const [month, setMonth] = useState<number>(new Date().getMonth() + 1);

    useEffect(() => {
        const loadEmployee = async (): Promise<void> => {
            try {
                const data = await fetchEmployeeById<Employee>(id as string);
                setEmployee(data);
            } catch (err) {
                if (err instanceof Error) setError(err.message);
            } finally {
                setLoading(false);
            }
        };
        void loadEmployee();
    }, [id]);

    useEffect(() => {
        const loadTimeClocks = async (): Promise<void> => {
            try {
                const data = await fetchTimeClocks<TimeClock[]>(Number(id), year, month);
                setTimeClocks(data ?? []);
            } catch (err) {
                if (err instanceof Error) setError(err.message);
            }
        };
        void loadTimeClocks();
    }, [id, year, month, refresh]);

    const handleTimeClockAdded = (): void => {
        setRefresh((prev) => !prev);
    };

    const handleDateChange = (newYear: number, newMonth: number) => {
        setYear(newYear);
        setMonth(newMonth);
    };

    if (loading) return <div className="text-center mt-8">Loading...</div>;
    if (error) return <div className="text-center mt-8 text-red-600">{error}</div>;
    if (!employee) return <div className="text-center mt-8">Not Found</div>;

    return (
        <div className="container mx-auto px-4 py-4">
            <Card className="w-full shadow-md bg-white">
                <CardHeader className="flex flex-row items-center gap-2 border-b pb-4">
                    <BadgeInfo className="w-5 h-5 text-blue-600"/>
                    <CardTitle className="text-2xl font-bold text-gray-800">
                        {employee.name}さんの詳細
                    </CardTitle>
                </CardHeader>

                <CardContent className="space-y-8 pt-4">
                    {/* 年・月の選択 */}
                    <section>
                        <DateSelector year={year} month={month} onChange={handleDateChange}/>
                    </section>

                    {/* 給与計算 */}
                    <section>
                        <div className="flex items-center gap-2 mb-2">
                            <Wallet className="w-5 h-5 text-blue-600"/>
                            <h2 className="text-xl font-semibold text-gray-800">給与計算結果</h2>
                        </div>
                        <PayrollDisplay employeeId={Number(id)} year={year} month={month}/>
                    </section>

                    {/* 打刻一覧 */}
                    <section>
                        <div className="flex items-center gap-2 mb-2">
                            <ClipboardList className="w-5 h-5 text-blue-600"/>
                            <h2 className="text-xl font-semibold text-gray-800">打刻一覧</h2>
                        </div>
                        <TimeClockList timeClocks={timeClocks}/>
                    </section>

                    {/* 打刻登録 */}
                    <section>
                        <div className="flex items-center gap-2 mb-2">
                            <CalendarCheck2 className="w-5 h-5 text-blue-600"/>
                            <h2 className="text-xl font-semibold text-gray-800">打刻登録</h2>
                        </div>
                        <TimeClockForm
                            employeeId={Number(id)}
                            onTimeClockAdded={handleTimeClockAdded}
                            isLineLinked={employee.line_linked}
                        />
                    </section>
                </CardContent>
            </Card>
        </div>
    );
};

export default EmployeeDetailPage;
