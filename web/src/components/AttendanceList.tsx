import React from "react";
import {Attendance} from "../types/Employee";

interface AttendanceListProps {
    attendances: Attendance[];
}

const AttendanceList: React.FC<AttendanceListProps> = ({attendances}) => {
    if (attendances.length === 0) {
        return <p>打刻情報はありません。</p>;
    }

    return (
        <div>
            <h3 className="text-lg font-semibold mb-2">打刻情報</h3>
            <table className="w-full border-collapse">
                <thead>
                <tr>
                    <th className="border p-2">ID</th>
                    <th className="border p-2">出勤時刻</th>
                    <th className="border p-2">退勤時刻</th>
                </tr>
                </thead>
                <tbody>
                {attendances.map((att) => (
                    <tr key={att.id} className="text-center">
                        <td className="border p-2">{att.id}</td>
                        <td className="border p-2">{new Date(att.check_in).toLocaleString()}</td>
                        <td className="border p-2">{new Date(att.check_out).toLocaleString()}</td>
                    </tr>
                ))}
                </tbody>
            </table>
        </div>
    );
};

export default AttendanceList;
