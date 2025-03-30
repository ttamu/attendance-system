import {TimeClock} from "@/types/TimeClock.ts";

export interface Employee {
    id: number;
    name: string;
    created_at: string;
    updated_at: string;
    time_clocks?: TimeClock[];
}