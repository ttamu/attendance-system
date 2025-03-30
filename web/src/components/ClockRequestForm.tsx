import React, {useEffect, useState} from "react";
import {TimeClock} from "@/types/TimeClock";
import {Button} from "@/components/ui/button";
import {Input} from "@/components/ui/input";
import {Dialog, DialogClose, DialogContent, DialogFooter, DialogHeader, DialogTitle,} from "@/components/ui/dialog";

interface ClockRequestFormProps {
    open: boolean;
    onOpenChange: (open: boolean) => void;
    clock: TimeClock;
    onSubmit: (
        data: {
            type: "clock_in" | "clock_out" | "break_begin" | "break_end";
            time: string;
            reason: string;
        }
    ) => void;
    onCancel?: () => void;
}

const ClockRequestForm: React.FC<ClockRequestFormProps> = ({open, onOpenChange, clock, onSubmit, onCancel,}) => {
    const initTime = new Date(clock.timestamp).toISOString().slice(0, 16);
    const [type, setType] = useState<"clock_in" | "clock_out" | "break_begin" | "break_end">(clock.type);
    const [time, setTime] = useState(initTime);
    const [reason, setReason] = useState("");

    useEffect(() => {
        const newInitTime = new Date(clock.timestamp).toISOString().slice(0, 16);
        setType(clock.type);
        setTime(newInitTime);
        setReason("");
    }, [clock]);

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault();
        const fTime = `${time}:00+09:00`;
        onSubmit({type, time: fTime, reason});
    };

    const handleCancel = () => {
        if (onCancel) onCancel();
        onOpenChange(false);
    };

    return (
        <Dialog open={open} onOpenChange={onOpenChange}>
            <DialogContent>
                <form onSubmit={handleSubmit}>
                    <DialogHeader>
                        <DialogTitle>打刻修正申請</DialogTitle>
                        <DialogClose/>
                    </DialogHeader>
                    <div className="mb-4">
                        <label className="block text-sm font-medium text-gray-700">打刻種別</label>
                        <select
                            value={type}
                            onChange={(e) =>
                                setType(e.target.value as "clock_in" | "clock_out" | "break_begin" | "break_end")
                            }
                            className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm"
                        >
                            <option value="clock_in">出勤</option>
                            <option value="clock_out">退勤</option>
                            <option value="break_begin">休憩開始</option>
                            <option value="break_end">休憩終了</option>
                        </select>
                    </div>
                    <div className="mb-4">
                        <label className="block text-sm font-medium text-gray-700">打刻日時</label>
                        <Input
                            type="datetime-local"
                            value={time}
                            onChange={(e) => setTime(e.target.value)}
                        />
                    </div>
                    <div className="mb-4">
                        <label className="block text-sm font-medium text-gray-700">申請理由</label>
                        <Input
                            type="text"
                            value={reason}
                            onChange={(e) => setReason(e.target.value)}
                        />
                    </div>
                    <DialogFooter className="flex justify-end gap-2">
                        <Button variant="outline" type="button" onClick={handleCancel}>
                            キャンセル
                        </Button>
                        <Button type="submit">送信</Button>
                    </DialogFooter>
                </form>
            </DialogContent>
        </Dialog>
    );
};

export default ClockRequestForm;
