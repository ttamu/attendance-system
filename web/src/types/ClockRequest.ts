export interface ClockRequest {
    id: number;
    employee_id: number;
    employee_name: string;
    clock_id: number;
    type: "clock_in" | "clock_out" | "break_begin" | "break_end";
    time: string;
    status: "pending" | "approved" | "rejected";
    reason: string;
    reviewed_by?: number;
    reviewed_at?: string;
    created_at: string;
    updated_at: string;
}
