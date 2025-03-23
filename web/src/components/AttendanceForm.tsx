import React, {FormEvent, useState} from "react";
import {createAttendance} from "../services/api";

interface AttendanceFormProps {
    employeeId: number;
    onAttendanceAdded: () => void;
}

const AttendanceForm: React.FC<AttendanceFormProps> = ({employeeId, onAttendanceAdded}) => {
    const [checkIn, setCheckIn] = useState<string>("");
    const [checkOut, setCheckOut] = useState<string>("");
    const [error, setError] = useState<string>("");

    const handleSubmit = async (e: FormEvent<HTMLFormElement>) => {
        e.preventDefault();
        try {
            const checkInISO = new Date(checkIn).toISOString();
            const checkOutISO = new Date(checkOut).toISOString();
            await createAttendance<unknown>(employeeId.toString(), {
                check_in: checkInISO,
                check_out: checkOutISO,
            });
            setCheckIn("");
            setCheckOut("");
            onAttendanceAdded();
        } catch (err) {
            setError((err as Error).message);
        }
    };

    return (
        <div>
            <h3>打刻登録</h3>
            {error && <p style={{color: "red"}}>{error}</p>}
            <form onSubmit={handleSubmit} className="space-y-4">
                <div>
                    <label htmlFor="checkIn" className="block mb-1">出勤時刻</label>
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
                    <label htmlFor="checkOut" className="block mb-1">退勤時刻</label>
                    <input
                        type="datetime-local"
                        id="checkOut"
                        value={checkOut}
                        onChange={(e) => setCheckOut(e.target.value)}
                        className="w-full border p-2 rounded"
                        required
                    />
                </div>
                <button type="submit" className="bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700">
                    登録
                </button>
            </form>
        </div>
    );
};

export default AttendanceForm;
