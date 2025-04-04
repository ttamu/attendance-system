import {TimeClock} from "@/types/TimeClock.ts";

export interface Employee {
    id: number;
    line_linked: boolean;
    name: string;
    created_at: string;
    updated_at: string;
    time_clocks?: TimeClock[];
}