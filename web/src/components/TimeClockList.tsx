import React from "react"
import {TimeClock} from "@/types/TimeClock"
import {Card, CardContent, CardHeader, CardTitle} from "@/components/ui/card"
import {Table, TableBody, TableCell, TableHead, TableHeader, TableRow,} from "@/components/ui/table"

const clockTypeLabels: Record<TimeClock["type"], string> = {
    clock_in: "出勤",
    clock_out: "退勤",
    break_begin: "休憩開始",
    break_end: "休憩終了",
}

interface TimeClockListProps {
    timeClocks: TimeClock[]
}

const TimeClockList: React.FC<TimeClockListProps> = ({timeClocks}) => {
    return (
        <Card className="mt-4 ">
            <CardHeader>
                <CardTitle>打刻一覧</CardTitle>
            </CardHeader>
            <CardContent>
                <Table className="rounded-xl overflow-hidden">
                    <TableHeader>
                        <TableRow className="bg-gray-100">
                            <TableHead className="px-4 py-2 text-gray-700">
                                打刻種別
                            </TableHead>
                            <TableHead className="px-4 py-2 text-gray-700">
                                打刻日時
                            </TableHead>
                        </TableRow>
                    </TableHeader>
                    <TableBody>
                        {timeClocks.map((clock, index) => (
                            <TableRow key={index} className="hover:bg-gray-50">
                                <TableCell className="px-4 py-2">
                                    {clockTypeLabels[clock.type] ?? clock.type}
                                </TableCell>
                                <TableCell className="px-4 py-2">
                                    {clock.timestamp}
                                </TableCell>
                            </TableRow>
                        ))}
                    </TableBody>
                </Table>
            </CardContent>
        </Card>
    )
}

export default TimeClockList
