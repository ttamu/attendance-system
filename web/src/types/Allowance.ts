export interface AllowanceType {
    id: number;
    company_id: number;
    name: string;
    type: string;
    description: string;
    commission_rate?: number;
    created_at: string;
    updated_at: string;
}

export interface EmployeeAllowance {
    id?: number;
    employee_id: number;
    allowance_type_id: number;
    amount: number;
    commission_rate?: number;
    year: number;
    month: number;
    created_at?: string;
    updated_at?: string;
    employee_name?: string;
    allowance_type_name?: string;
    allowance_type?: string;
}
