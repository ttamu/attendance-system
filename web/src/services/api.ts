import {UserProfile} from "../types/User";

const API_BASE_URL = import.meta.env.VITE_API_URL as string;

export async function fetchAPI<T>(endpoint: string, options: RequestInit = {}): Promise<T> {
    const response = await fetch(`${API_BASE_URL}${endpoint}`, {
        ...options,
        credentials: "include",
    });
    if (!response.ok) {
        throw new Error(`HTTP error! Status: ${response.status}`);
    }
    return await response.json() as Promise<T>;
}

export async function fetchEmployees<T>(): Promise<T> {
    return fetchAPI<T>("/employees");
}

export async function fetchEmployeeById<T>(id: string): Promise<T> {
    return fetchAPI<T>(`/employees/${id}`);
}

export async function createAttendance<T>(
    employeeId: string,
    data: { check_in: string; check_out: string }
): Promise<T> {
    return fetchAPI<T>(`/employees/${employeeId}/attendances`, {
        method: 'POST',
        headers: {"Content-Type": "application/json"},
        body: JSON.stringify(data),
    });
}

export async function fetchPayroll<T>(employeeId: string, year: number, month: number): Promise<T> {
    return fetchAPI<T>(`/employees/${employeeId}/payroll?year=${year}&month=${month}`);
}

export async function login<T>(credential: { email: string, password: string }): Promise<T> {
    return fetchAPI<T>("/login", {
        method: 'POST',
        headers: {"Content-Type": "application/json"},
        body: JSON.stringify(credential),
    });
}

export async function fetchProfile<T = UserProfile>(): Promise<T> {
    return fetchAPI<T>("/current_account");
}
