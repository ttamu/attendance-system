import React, {useState} from "react";
import {createTimeClock} from "../services/api";
import {Card, CardContent} from "@/components/ui/card";
import {Label} from "@/components/ui/label";
import {Input} from "@/components/ui/input";
import {Button} from "@/components/ui/button";
import {Select, SelectContent, SelectItem, SelectTrigger, SelectValue} from "@/components/ui/select";
import {Tooltip, TooltipContent, TooltipTrigger} from "@/components/ui/tooltip";
import {Bell, BellOff, Clock, Coffee, LogIn, LogOut} from "lucide-react";

interface TimeClockFormProps {
    employeeId: number;
    onTimeClockAdded: () => void;
    isLineLinked: boolean;
}

const hours = Array.from({length: 24}, (_, i) => i.toString().padStart(2, "0"));
const minutes = Array.from({length: 60}, (_, i) => i.toString().padStart(2, "0"));

const TimeClockForm: React.FC<TimeClockFormProps> = ({employeeId, onTimeClockAdded, isLineLinked}) => {
    const [error, setError] = useState<string>("");
    const [success, setSuccess] = useState<boolean>(false);
    const [notify, setNotify] = useState<boolean>(false);
    const [notifyHour, setNotifyHour] = useState<string>("17");
    const [notifyMinute, setNotifyMinute] = useState<string>("00");
    const [useTime, setUseTime] = useState<boolean>(false);
    const [customTime, setCustomTime] = useState<string>("");

    const handleTimeClockAction = async (clockType: string) => {
        setError("");
        setSuccess(false);
        try {
            const timestampISO = useTime && customTime ? new Date(customTime).toISOString() : new Date().toISOString();
            const payload: any = {
                type: clockType,
                timestamp: timestampISO,
                notify,
            };
            if (notify) {
                payload.notify_at = `${notifyHour.padStart(2, "0")}:${notifyMinute.padStart(2, "0")}`;
            }
            await createTimeClock<unknown>(employeeId.toString(), payload);
            setSuccess(true);
            setCustomTime("");
            setNotify(false);
            setNotifyHour("17");
            setNotifyMinute("00");
            onTimeClockAdded();
        } catch (err) {
            setError((err as Error).message);
        }
    };

    return (
        <Card className="bg-white shadow-sm text-gray-800">
            <CardContent className="space-y-2">
                {success && (
                    <p className="text-green-600 bg-green-50 border border-green-200 p-2 rounded">登録に成功しました！</p>
                )}
                {error && (
                    <p className="text-red-600 bg-red-50 border border-red-200 p-2 rounded">エラー: {error}</p>
                )}

                <div className="flex flex-wrap gap-4">
                    <Tooltip>
                        <TooltipTrigger asChild>
                            <Button onClick={() => setUseTime((p) => !p)}
                                    className={`w-10 h-10 rounded-full p-2 ${useTime ? "bg-blue-500" : "bg-gray-300"}`}>
                                <Clock className="w-5 h-5 text-white"/>
                            </Button>
                        </TooltipTrigger>
                        <TooltipContent>任意の日時を指定</TooltipContent>
                    </Tooltip>
                    <Tooltip>
                        <TooltipTrigger asChild>
                            {isLineLinked ? (
                                <Button onClick={() => setNotify((p) => !p)} disabled={!isLineLinked}
                                        className={`w-10 h-10 rounded-full p-2 ${notify ? "bg-green-600" : "bg-gray-300"} ${!isLineLinked ? "opacity-50 cursor-not-allowed" : ""}`}>
                                    {notify ? <Bell className="w-5 h-5 text-white"/> :
                                        <BellOff className="w-5 h-5 text-white"/>}
                                </Button>
                            ) : (
                                <span className="inline-block">
                  <Button disabled className="w-10 h-10 rounded-full p-2 bg-gray-300 opacity-50 cursor-not-allowed">
                    <BellOff className="w-5 h-5 text-white"/>
                  </Button>
                </span>
                            )}
                        </TooltipTrigger>
                        <TooltipContent>{isLineLinked ? "LINEのリマインダー設定" : "LINE連携してください"}</TooltipContent>
                    </Tooltip>
                </div>

                <div
                    className={`transition-all duration-300 overflow-hidden ${useTime ? "max-h-24 opacity-100 mt-2" : "max-h-0 opacity-0"}`}>
                    <Label htmlFor="customTime" className="block mb-1">指定日時</Label>
                    <Input type="datetime-local" id="customTime" value={customTime}
                           onChange={(e) => setCustomTime(e.target.value)} className="w-full"/>
                </div>

                {isLineLinked && notify && (
                    <div className="mt-2">
                        <Label className="block mb-1">通知時刻</Label>
                        <div className="flex items-center space-x-2">
                            <Select onValueChange={setNotifyHour} value={notifyHour}>
                                <SelectTrigger className="w-24 px-2 py-1 text-center"> <SelectValue
                                    placeholder="HH"/></SelectTrigger>
                                <SelectContent className="max-h-60 overflow-auto">
                                    {hours.map((h) => (<SelectItem key={h} value={h}>{h}</SelectItem>))}
                                </SelectContent>
                            </Select>
                            <span>:</span>
                            <Select onValueChange={setNotifyMinute} value={notifyMinute}>
                                <SelectTrigger className="w-24 px-2 py-1 text-center"> <SelectValue
                                    placeholder="MM"/></SelectTrigger>
                                <SelectContent className="max-h-60 overflow-auto">
                                    {minutes.map((m) => (<SelectItem key={m} value={m}>{m}</SelectItem>))}
                                </SelectContent>
                            </Select>
                        </div>
                    </div>
                )}

                <div className="grid grid-cols-2 gap-4 mt-4">
                    <button onClick={() => handleTimeClockAction("clock_in")}
                            className="flex flex-col items-center justify-center h-28 p-2 bg-white border border-gray-200 shadow-sm rounded-lg hover:bg-gray-50">
                        <LogIn className="text-3xl text-blue-600"/>
                        <span className="mt-2 text-base font-semibold text-gray-800">出勤</span>
                    </button>
                    <button onClick={() => handleTimeClockAction("clock_out")}
                            className="flex flex-col items-center justify-center h-28 p-2 bg-white border border-gray-200 shadow-sm rounded-lg hover:bg-gray-50">
                        <LogOut className="text-3xl text-red-600"/>
                        <span className="mt-2 text-base font-semibold text-gray-800">退勤</span>
                    </button>
                    <button onClick={() => handleTimeClockAction("break_begin")}
                            className="flex flex-col items-center justify-center h-28 p-2 bg-white border border-gray-200 shadow-sm rounded-lg hover:bg-gray-50">
                        <Coffee className="text-3xl text-amber-600"/>
                        <span className="mt-2 text-base font-semibold text-gray-800">休憩開始</span>
                    </button>
                    <button onClick={() => handleTimeClockAction("break_end")}
                            className="flex flex-col items-center justify-center h-28 p-2 bg-white border border-gray-200 shadow-sm rounded-lg hover:bg-gray-50 relative">
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
