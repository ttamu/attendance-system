import React, {useEffect, useState} from "react";
import {Card, CardContent, CardHeader, CardTitle} from "@/components/ui/card";
import {Table, TableBody, TableCell, TableHead, TableHeader, TableRow} from "@/components/ui/table";
import {approveClockReq, fetchClockRequests, rejectClockReq} from "@/services/api";
import {ClockRequest} from "@/types/ClockRequest";
import {Button} from "@/components/ui/button";
import {CheckCircle, Clock, XCircle} from 'lucide-react';

const typeLabels: Record<ClockRequest["type"], string> = {
    clock_in: "出勤",
    clock_out: "退勤",
    break_begin: "休憩開始",
    break_end: "休憩終了",
};

const statusIcons: Record<ClockRequest["status"], React.ReactNode> = {
    pending: <Clock className="w-4 h-4 text-gray-500"/>,
    approved: <CheckCircle className="w-4 h-4 text-green-500"/>,
    rejected: <XCircle className="w-4 h-4 text-red-500"/>,
};

const formatDate = (iso: string) =>
    new Date(iso).toLocaleString("ja-JP", {
        year: "numeric",
        month: "short",
        day: "numeric",
        hour: "2-digit",
        minute: "2-digit",
    });

const ClockRequestList: React.FC = () => {
    const [requests, setRequests] = useState<ClockRequest[]>([]);
    const [loading, setLoading] = useState<boolean>(true);
    const [error, setError] = useState<string>("");

    const loadRequests = async () => {
        setLoading(true);
        try {
            const data = await fetchClockRequests<ClockRequest[]>();
            setRequests(data);
        } catch (err) {
            if (err instanceof Error) {
                setError(err.message);
            }
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        loadRequests();
    }, []);

    const handleApprove = async (id: number) => {
        await approveClockReq(id);
        loadRequests();
    };

    const handleReject = async (id: number) => {
        await rejectClockReq(id);
        loadRequests();
    };

    if (loading) return <p className="text-center py-4">Loading...</p>;
    if (error) return <p className="text-center text-red-600 py-4">{error}</p>;

    return (
        <Card className="mt-8">
            <CardHeader>
                <CardTitle>打刻修正申請一覧</CardTitle>
            </CardHeader>
            <CardContent>
                <Table className="w-full">
                    <TableHeader>
                        <TableRow>
                            <TableHead className="px-4 py-2">従業員名</TableHead>
                            <TableHead className="px-4 py-2">申請日</TableHead>
                            <TableHead className="px-4 py-2">種別</TableHead>
                            <TableHead className="px-4 py-2">修正時刻</TableHead>
                            <TableHead className="px-4 py-2">理由</TableHead>
                            <TableHead className="px-4 py-2 text-center">ステータス</TableHead>
                            <TableHead className="px-4 py-2 text-center">操作</TableHead>
                        </TableRow>
                    </TableHeader>
                    <TableBody>
                        {requests.map((req) => (
                            <TableRow key={req.id}>
                                <TableCell className="px-4 py-2">{req.employee_name}</TableCell>
                                <TableCell className="px-4 py-2">{formatDate(req.created_at)}</TableCell>
                                <TableCell className="px-4 py-2">{typeLabels[req.type]}</TableCell>
                                <TableCell className="px-4 py-2">{formatDate(req.time)}</TableCell>
                                <TableCell className="px-4 py-2">{req.reason || "―"}</TableCell>
                                <TableCell className="px-4 py-2 text-center align-middle">
                                    {statusIcons[req.status]}
                                </TableCell>
                                <TableCell className="px-4 py-2 text-center align-middle">
                                    {req.status === "pending" && (
                                        <>
                                            <Button size="icon"
                                                    className="mr-2 bg-green-100 hover:bg-green-200 text-green-600"
                                                    onClick={() => handleApprove(req.id)}>
                                                <CheckCircle className="w-4 h-4"/>
                                            </Button>
                                            <Button size="icon" className="bg-red-100 hover:bg-red-200 text-red-600"
                                                    onClick={() => handleReject(req.id)}>
                                                <XCircle className="w-4 h-4"/>
                                            </Button>
                                        </>
                                    )}
                                </TableCell>
                            </TableRow>
                        ))}
                    </TableBody>
                </Table>
            </CardContent>
        </Card>
    );
};

export default ClockRequestList;