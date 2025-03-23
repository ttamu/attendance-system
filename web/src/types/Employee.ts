import {Attendance} from "@/types/Attendance.ts";

export interface Employee {
    id: number;
    name: string;
    created_at: string;
    updated_at: string;
    attendances?: Attendance[];
}