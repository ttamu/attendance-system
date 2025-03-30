export interface TimeClock {
    id: number;
    employee_id: number;
    type: "clock_in" | "clock_out" | "break_begin" | "break_end";
    timestamp: string;
    created_at: string;
}