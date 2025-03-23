import React from "react";
import {Attendance} from "../types/Attendance";
import {Card, CardContent} from "@/components/ui/card";

interface AttendanceListProps {
    attendances: Attendance[];
}

const AttendanceList: React.FC<AttendanceListProps> = ({attendances}) => {
    if (attendances.length === 0) {
        return (
            <Card className="bg-white shadow-sm text-gray-800">
                <CardContent className="p-4">
                    <p className="text-gray-600">打刻情報はありません。</p>
                </CardContent>
            </Card>
        );
    }

    return (
        <Card className="bg-white shadow-sm text-gray-800">
            <CardContent className="p-4">
                <h3 className="text-lg font-semibold mb-4 text-gray-800">打刻情報</h3>
                <table className="w-full border-collapse text-sm">
                    <thead>
                    <tr className="bg-gray-100 text-left">
                        <th className="border p-2">ID</th>
                        <th className="border p-2">出勤時刻</th>
                        <th className="border p-2">退勤時刻</th>
                    </tr>
                    </thead>
                    <tbody>
                    {attendances.map((att) => (
                        <tr key={att.id} className="text-center hover:bg-gray-50">
                            <td className="border p-2">{att.id}</td>
                            <td className="border p-2">
                                {new Date(att.check_in).toLocaleString()}
                            </td>
                            <td className="border p-2">
                                {new Date(att.check_out).toLocaleString()}
                            </td>
                        </tr>
                    ))}
                    </tbody>
                </table>
            </CardContent>
        </Card>
    );
};

export default AttendanceList;
