import React, { FormEvent, useState } from "react";
import { createTimeClock } from "../services/api";
import { Card, CardContent } from "@/components/ui/card";
import { Label } from "@/components/ui/label";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";

interface TimeClockFormProps {
    employeeId: number;
    onTimeClockAdded: () => void;
}

const TimeClockForm: React.FC<TimeClockFormProps> = ({ employeeId, onTimeClockAdded }) => {
    const [type, setType] = useState<string>("clock_in");
    const [timestamp, setTimestamp] = useState<string>("");
    const [error, setError] = useState<string>("");
    const [success, setSuccess] = useState<boolean>(false);

    const handleSubmit = async (e: FormEvent<HTMLFormElement>) => {
        e.preventDefault();
        setError("");
        setSuccess(false);

        try {
            const timestampISO = new Date(timestamp).toISOString();
            await createTimeClock<unknown>(employeeId.toString(), {
                type,
                timestamp: timestampISO,
            });

            setTimestamp("");
            setSuccess(true);
            onTimeClockAdded();
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
                        <Label htmlFor="type">打刻種別</Label>
                        <select
                            id="type"
                            value={type}
                            onChange={(e) => setType(e.target.value)}
                            className="border p-2 rounded w-full"
                        >
                            <option value="clock_in">出勤</option>
                            <option value="clock_out">退勤</option>
                            <option value="break_begin">休憩開始</option>
                            <option value="break_end">休憩終了</option>
                        </select>
                    </div>
                    <div>
                        <Label htmlFor="timestamp">打刻日時</Label>
                        <Input
                            type="datetime-local"
                            id="timestamp"
                            value={timestamp}
                            onChange={(e) => setTimestamp(e.target.value)}
                            required
                        />
                    </div>
                    <Button type="submit">登録</Button>
                </form>
            </CardContent>
        </Card>
    );
};

export default TimeClockForm;
