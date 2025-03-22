const API_BASE_URL = import.meta.env.VITE_API_URL;

export const fetchEmployees = async () => {
    const response = await fetch(`${API_BASE_URL}/employees`);
    if (!response.ok) {
        throw new Error(`HTTP error! Status: ${response.status}`);
    }
    return response.json();
}

export const fetchEmployeeById = async (id: string) => {
    const response = await fetch(`${API_BASE_URL}/employees/${id}`);
    if (!response.ok) {
        throw new Error(`HTTP error! Status: ${response.status}`);
    }
    return response.json();
}

export const createAttendance = async (userId: string, data: { check_in: string; check_out: string }) => {
    const response = await fetch(`${API_BASE_URL}/employees/${userId}/attendances`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(data),
    });
    if (!response.ok) {
        throw new Error(`HTTP error! Status: ${response.status}`);
    }
    return response.json();
};