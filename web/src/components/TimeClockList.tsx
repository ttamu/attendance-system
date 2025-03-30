import React, {useState} from "react";
import {TimeClock} from "@/types/TimeClock";
import {Card, CardContent, CardHeader, CardTitle} from "@/components/ui/card";
import {Table, TableBody, TableCell, TableHead, TableHeader, TableRow} from "@/components/ui/table";
import {Button} from "@/components/ui/button";
import ClockRequestForm from "@/components/ClockRequestForm";
import {createClockRequest} from "../services/api";

const clockTypeLabels: Record<TimeClock["type"], string> = {
    clock_in: "出勤",
    clock_out: "退勤",
    break_begin: "休憩開始",
    break_end: "休憩終了",
};

interface TimeClockListProps {
    timeClocks: TimeClock[];
    onRefresh?: () => void;
}

const TimeClockList: React.FC<TimeClockListProps> = ({timeClocks, onRefresh}) => {
    const [editingClockIndex, setEditingClockIndex] = useState<number | null>(null);

    const handleStartEditing = (index: number) => {
        setEditingClockIndex(index);
    };

    const handleCancelEditing = () => {
        setEditingClockIndex(null);
    };

    const handleRequestSubmit = async (
        clock: TimeClock,
        data: { type: string; time: string; reason: string }
    ) => {
        try {
            await createClockRequest(clock.id, {
                employee_id: clock.employee_id,
                type: data.type,
                time: data.time,
                reason: data.reason,
            });
            setEditingClockIndex(null);
            if (onRefresh) {
                onRefresh();
            }
        } catch (error) {
            console.error("打刻修正申請の送信に失敗しました:", error);
        }
    };

    return (
        <Card className="mt-4">
            <CardHeader>
                <CardTitle>打刻一覧</CardTitle>
            </CardHeader>
            <CardContent>
                <Table className="rounded-xl overflow-hidden">
                    <TableHeader>
                        <TableRow className="bg-gray-100">
                            <TableHead className="px-4 py-2 text-gray-700">打刻種別</TableHead>
                            <TableHead className="px-4 py-2 text-gray-700">打刻日時</TableHead>
                            <TableHead className="px-4 py-2 text-gray-700">操作</TableHead>
                        </TableRow>
                    </TableHeader>
                    <TableBody>
                        {timeClocks.map((clock, idx) => (
                            <React.Fragment key={idx}>
                                <TableRow className="hover:bg-gray-50">
                                    <TableCell className="px-4 py-2">
                                        {clockTypeLabels[clock.type] ?? clock.type}
                                    </TableCell>
                                    <TableCell className="px-4 py-2">{clock.timestamp}</TableCell>
                                    <TableCell className="px-4 py-2">
                                        <Button size="sm" onClick={() => handleStartEditing(idx)}>
                                            修正申請
                                        </Button>
                                    </TableCell>
                                </TableRow>
                                {editingClockIndex === idx && (
                                    <ClockRequestForm
                                        open={true}
                                        onOpenChange={(open) => {
                                            if (!open) handleCancelEditing();
                                        }}
                                        clock={clock}
                                        onSubmit={(data) => handleRequestSubmit(clock, data)}
                                        onCancel={handleCancelEditing}
                                    />
                                )}
                            </React.Fragment>
                        ))}
                    </TableBody>
                </Table>
            </CardContent>
        </Card>
    );
};

export default TimeClockList;
