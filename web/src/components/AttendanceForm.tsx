import React, {FormEvent, useState} from "react";
import {createAttendance} from "../services/api";
import {Card, CardContent} from "@/components/ui/card";
import {Label} from "@/components/ui/label";
import {Input} from "@/components/ui/input";
import {Button} from "@/components/ui/button";

interface AttendanceFormProps {
    employeeId: number;
    onAttendanceAdded: () => void;
}

const AttendanceForm: React.FC<AttendanceFormProps> = ({employeeId, onAttendanceAdded}) => {
    const [checkIn, setCheckIn] = useState<string>("");
    const [checkOut, setCheckOut] = useState<string>("");
    const [error, setError] = useState<string>("");
    const [success, setSuccess] = useState<boolean>(false);

    const handleSubmit = async (e: FormEvent<HTMLFormElement>) => {
        e.preventDefault();
        setError("");
        setSuccess(false);

        try {
            const checkInISO = new Date(checkIn).toISOString();
            const checkOutISO = new Date(checkOut).toISOString();
            await createAttendance<unknown>(employeeId.toString(), {
                check_in: checkInISO,
                check_out: checkOutISO,
            });

            setCheckIn("");
            setCheckOut("");
            setSuccess(true);
            onAttendanceAdded();
        } catch (err) {
            setError((err as Error).message);
        }
    };

    return (
        <Card className="bg-white shadow-sm text-gray-800">
            <CardContent className="p-4 space-y-4">
                <h3 className="text-lg font-semibold text-gray-800">打刻登録</h3>

                {success && (
                    <p className="text-green-600 bg-green-50 border border-green-200 p-2 rounded">
                        登録に成功しました！
                    </p>
                )}
                {error && (
                    <p className="text-red-600 bg-red-50 border border-red-200 p-2 rounded">
                        エラー: {error}
                    </p>
                )}

                <form onSubmit={handleSubmit} className="space-y-4">
                    <div>
                        <Label htmlFor="checkIn">出勤時刻</Label>
                        <Input
                            type="datetime-local"
                            id="checkIn"
                            value={checkIn}
                            onChange={(e) => setCheckIn(e.target.value)}
                            required
                        />
                    </div>
                    <div>
                        <Label htmlFor="checkOut">退勤時刻</Label>
                        <Input
                            type="datetime-local"
                            id="checkOut"
                            value={checkOut}
                            onChange={(e) => setCheckOut(e.target.value)}
                            required
                        />
                    </div>
                    <Button type="submit">登録</Button>
                </form>
            </CardContent>
        </Card>
    );
};

export default AttendanceForm;
