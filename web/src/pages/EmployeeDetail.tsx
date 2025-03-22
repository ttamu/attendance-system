import React, {FormEvent, useEffect, useState} from 'react'
import {useParams} from 'react-router-dom'
import {Employee} from '../types/Employee.ts'
import {createAttendance, fetchEmployeeById} from '../services/api'

const EmployeeDetail: React.FC = () => {
    const {id} = useParams<{ id: string }>();
    const [user, setUser] = useState<Employee | null>(null);
    const [loading, setLoading] = useState<boolean>(true);
    const [error, setError] = useState<string>('');
    const [checkIn, setCheckIn] = useState<string>('');
    const [checkOut, setCheckOut] = useState<string>('');
    const [attLoading, setAttLoading] = useState<boolean>(false);

    useEffect(() => {
        const loadUser = async () => {
            try {
                const data = await fetchEmployeeById(id!);
                setUser(data);
            } catch (err: any) {
                setError(err.message);
            } finally {
                setLoading(false);
            }
        };
        void loadUser();
    }, [id]);

    const handleAttendanceSubmit = async (e: FormEvent) => {
        e.preventDefault();
        setAttLoading(true);
        try {
            const checkInISO = new Date(checkIn).toISOString();
            const checkOutISO = new Date(checkOut).toISOString();

            const newAttendance = await createAttendance(id!, {
                check_in: checkInISO,
                check_out: checkOutISO,
            });
            setUser((prev) => prev ? {...prev, attendances: [...(prev.attendances || []), newAttendance]} : prev);
            setCheckIn('');
            setCheckOut('');
        } catch (err: any) {
            alert(err.message);
        } finally {
            setAttLoading(false);
        }
    };

    if (loading) return <div className="text-center mt-8">Loading...</div>;
    if (error) return <div className="text-center mt-8 text-red-600">{error}</div>;
    if (!user) return <div className="text-center mt-8">NotFound</div>;

    return (
        <div className="max-w-2xl mx-auto p-4 text-left">
            <h1 className="text-2xl font-bold mb-4">{user.name}の詳細</h1>

            <section className="mb-8">
                <h2 className="text-xl font-semibold mb-2">勤怠一覧</h2>
                {user.attendances && user.attendances.length > 0 ? (
                    <table className="w-full border-collapse">
                        <thead>
                        <tr>
                            <th className="border p-2">ID</th>
                            <th className="border p-2">出勤時刻</th>
                            <th className="border p-2">退勤時刻</th>
                        </tr>
                        </thead>
                        <tbody>
                        {user.attendances.map((att) => (
                            <tr key={att.id} className="text-center">
                                <td className="border p-2">{att.id}</td>
                                <td className="border p-2">{att.check_in}</td>
                                <td className="border p-2">{att.check_out}</td>
                            </tr>
                        ))}
                        </tbody>
                    </table>
                ) : (
                    <p>勤怠情報はありません。</p>
                )}
            </section>

            <section>
                <h2 className="text-xl font-semibold mb-2">勤怠登録</h2>
                <form onSubmit={handleAttendanceSubmit} className="space-y-4">
                    <div>
                        <label htmlFor="checkIn" className="block mb-1">
                            出勤時刻
                        </label>
                        <input
                            type="datetime-local"
                            id="checkIn"
                            value={checkIn}
                            onChange={(e) => setCheckIn(e.target.value)}
                            className="w-full border p-2 rounded"
                            required
                        />
                    </div>
                    <div>
                        <label htmlFor="checkOut" className="block mb-1">
                            退勤時刻
                        </label>
                        <input
                            type="datetime-local"
                            id="checkOut"
                            value={checkOut}
                            onChange={(e) => setCheckOut(e.target.value)}
                            className="w-full border p-2 rounded"
                            required
                        />
                    </div>
                    <button
                        type="submit"
                        className="bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700"
                        disabled={attLoading}
                    >
                        {attLoading ? '登録中...' : '勤怠登録'}
                    </button>
                </form>
            </section>
        </div>
    );
};

export default EmployeeDetail