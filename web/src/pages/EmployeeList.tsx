import React, {useEffect, useState} from "react";
import {Link} from "react-router-dom";
import {fetchEmployees} from "../services/api";
import {Employee} from "../types/Employee";
import {Table, TableBody, TableCell, TableHead, TableHeader, TableRow,} from "@/components/ui/table";

const EmployeeList: React.FC = () => {
    const [employees, setEmployees] = useState<Employee[]>([]);
    const [loading, setLoading] = useState<boolean>(true);
    const [error, setError] = useState<string>("");

    useEffect(() => {
        const loadEmployees = async (): Promise<void> => {
            try {
                const data = await fetchEmployees<Employee[]>();
                setEmployees(data);
            } catch (err) {
                if (err instanceof Error) {
                    setError(err.message);
                }
            } finally {
                setLoading(false);
            }
        };
        void loadEmployees();
    }, []);

    if (loading) return <p>Loading...</p>;
    if (error) return <p className="text-red-600">{error}</p>;

    return (
        <div className="max-w-2xl mx-auto p-4">
            <h1 className="text-xl font-semibold mb-4">従業員一覧</h1>
            <Table className="w-full">
                <TableHeader>
                    <TableRow>
                        <TableHead>ID</TableHead>
                        <TableHead>Name</TableHead>
                        <TableHead>Created At</TableHead>
                    </TableRow>
                </TableHeader>
                <TableBody>
                    {employees.map((emp) => (
                        <TableRow key={emp.id}>
                            <TableCell>{emp.id}</TableCell>
                            <TableCell>
                                <Link to={`/employees/${emp.id}`} className="text-blue-600 hover:underline">
                                    {emp.name}
                                </Link>
                            </TableCell>
                            <TableCell>{emp.created_at}</TableCell>
                        </TableRow>
                    ))}
                </TableBody>
            </Table>
        </div>
    );
};

export default EmployeeList;
