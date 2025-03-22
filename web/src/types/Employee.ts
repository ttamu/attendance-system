export interface Attendance {
    id: number;
    check_in: string;
    check_out: string;
    created_at: string;
    updated_at: string;
}

export interface Employee {
    id: number;
    name: string;
    created_at: string;
    updated_at: string;
    attendances?: Attendance[];
}