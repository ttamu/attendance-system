export interface Company {
    id: number;
    name: string;
}

export interface UserProfile {
    account_id: number;
    company_id: number;
    is_admin: boolean;
    email: string;
    company: Company;
}
