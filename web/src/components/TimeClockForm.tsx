import React, {useState} from "react";
import {createTimeClock} from "../services/api";
import {Card, CardContent} from "@/components/ui/card";
import {Label} from "@/components/ui/label";
import {Input} from "@/components/ui/input";
import {Button} from "@/components/ui/button";
import {Bell, BellOff, Clock, Coffee, LogIn, LogOut,} from "lucide-react";

interface TimeClockFormProps {
    employeeId: number;
    onTimeClockAdded: () => void;
    isLineLinked: boolean;
}

const TimeClockForm: React.FC<TimeClockFormProps> = ({employeeId, onTimeClockAdded, isLineLinked,}) => {
    const [error, setError] = useState<string>("");
    const [success, setSuccess] = useState<boolean>(false);
    const [notify, setNotify] = useState<boolean>(false);
    const [delayH, setDelayH] = useState<number>(0);
    const [delayM, setDelayM] = useState<number>(0);
    const [useTime, setUseTime] = useState<boolean>(false);
    const [customTime, setCustomTime] = useState<string>("");

    const handleTimeClockAction = async (clockType: string) => {
        setError("");
        setSuccess(false);
        try {
            const timestampISO = useTime && customTime ? new Date(customTime).toISOString() : new Date().toISOString();
            await createTimeClock<unknown>(employeeId.toString(), {
                type: clockType,
                timestamp: timestampISO,
                notify,
                delay_h: delayH,
                delay_m: delayM,
            });
            setSuccess(true);
            setCustomTime("");
            onTimeClockAdded();
        } catch (err) {
            setError((err as Error).message);
        }
    };

    return (
        <Card className="bg-white shadow-sm text-gray-800">
            <CardContent className="p-4 space-y-4">
                <h3 className="text-lg font-semibold">打刻登録</h3>

                {success && (
                    <p className="text-green-600 bg-green-50 border border-green-200 p-2 rounded">登録に成功しました！</p>
                )}
                {error && (
                    <p className="text-red-600 bg-red-50 border border-red-200 p-2 rounded">エラー: {error}</p>
                )}

                <div className="flex flex-wrap gap-4">
                    {/* 日時指定トグル */}
                    <Button
                        onClick={() => setUseTime((prev) => !prev)}
                        className={`w-10 h-10 rounded-full p-2 ${
                            useTime ? "bg-blue-500" : "bg-gray-300"
                        }`}
                    >
                        <Clock className="w-5 h-5 text-white"/>
                    </Button>

                    {/* 通知トグル */}
                    <Button
                        onClick={() => setNotify((prev) => !prev)}
                        disabled={!isLineLinked}
                        className={`w-10 h-10 rounded-full p-2 ${
                            notify ? "bg-green-600" : "bg-gray-300"
                        } ${!isLineLinked ? "opacity-50 cursor-not-allowed" : ""}`}
                    >
                        {notify ? (<Bell className="w-5 h-5 text-white"/>) : (
                            <BellOff className="w-5 h-5 text-white"/>)}
                    </Button>
                </div>

                {/* 日時指定の入力 */}
                <div
                    className={`transition-all duration-300 overflow-hidden ${
                        useTime ? "max-h-24 opacity-100 mt-2" : "max-h-0 opacity-0"
                    }`}
                >
                    <Label htmlFor="customTime" className="block mb-1">
                        指定日時
                    </Label>
                    <Input
                        type="datetime-local"
                        id="customTime"
                        value={customTime}
                        onChange={(e) => setCustomTime(e.target.value)}
                        className="w-full"
                    />
                </div>

                {/* 通知設定 */}
                {isLineLinked && (
                    <div
                        className={`transition-all duration-300 overflow-hidden ${
                            notify ? "max-h-40 opacity-100 mt-2" : "max-h-0 opacity-0"
                        }`}
                    >
                        <div className="grid grid-cols-2 gap-4">
                            <div>
                                <Label htmlFor="delayH" className="block mb-1">通知までの時間(時)</Label>
                                <select
                                    id="delayH"
                                    value={delayH}
                                    onChange={(e) => setDelayH(Number(e.target.value))}
                                    className="border p-2 rounded w-full"
                                >
                                    {[...Array(13).keys()].map((h) => (<option key={h} value={h}>{h}</option>))}
                                </select>
                            </div>
                            <div>
                                <Label htmlFor="delayM" className="block mb-1">通知までの時間(分)</Label>
                                <select
                                    id="delayM"
                                    value={delayM}
                                    onChange={(e) => setDelayM(Number(e.target.value))}
                                    className="border p-2 rounded w-full"
                                >
                                    {[0, 1, 5, 10, 15, 20, 30, 45, 50, 55].map((m) => (
                                        <option key={m} value={m}>{m}</option>))}
                                </select>
                            </div>
                        </div>
                    </div>
                )}

                {/* 打刻用のボタン */}
                <div className="grid grid-cols-2 gap-4 mt-4">
                    <button
                        onClick={() => handleTimeClockAction("clock_in")}
                        className="flex flex-col items-center justify-center h-28 p-2 bg-white border border-gray-200 shadow-sm rounded-lg hover:bg-gray-50"
                    >
                        <LogIn className="text-3xl text-blue-600"/>
                        <span className="mt-2 text-base font-semibold text-gray-800">出勤</span>
                    </button>

                    <button
                        onClick={() => handleTimeClockAction("clock_out")}
                        className="flex flex-col items-center justify-center h-28 p-2 bg-white border border-gray-200 shadow-sm rounded-lg hover:bg-gray-50"
                    >
                        <LogOut className="text-3xl text-red-600"/>
                        <span className="mt-2 text-base font-semibold text-gray-800">退勤</span>
                    </button>
                    <button
                        onClick={() => handleTimeClockAction("break_begin")}
                        className="flex flex-col items-center justify-center h-28 p-2 bg-white border border-gray-200 shadow-sm rounded-lg hover:bg-gray-50"
                    >
                        <Coffee className="text-3xl text-amber-600"/>
                        <span className="mt-2 text-base font-semibold text-gray-800">休憩開始</span>
                    </button>
                    <button
                        onClick={() => handleTimeClockAction("break_end")}
                        className="flex flex-col items-center justify-center h-28 p-2 bg-white border border-gray-200 shadow-sm rounded-lg hover:bg-gray-50 relative"
                    >
                        <div className="relative">
                            <Coffee className="text-3xl text-amber-600"/>
                            <div className="absolute inset-0 flex items-center justify-center pointer-events-none">
                                <div className="w-full h-0.5 bg-gray-600 rotate-45"/>
                            </div>
                        </div>
                        <span className="mt-2 text-base font-semibold text-gray-800">休憩終了</span>
                    </button>
                </div>
            </CardContent>
        </Card>
    );
};

export default TimeClockForm;
