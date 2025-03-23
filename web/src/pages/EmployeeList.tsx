import React, {useEffect, useState} from "react";
import {Link} from "react-router-dom";
import {fetchEmployees} from "../services/api";
import {Employee} from "../types/Employee";
import {Table, TableBody, TableCell, TableHead, TableHeader, TableRow,} from "@/components/ui/table";
import {Card, CardContent, CardHeader, CardTitle} from "@/components/ui/card";
import {Users} from "lucide-react";

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

    if (loading) return <p className="text-center py-4">Loading...</p>;
    if (error) return <p className="text-center text-red-600 py-4">{error}</p>;

    return (
        <Card className="max-w-4xl mx-auto shadow-md rounded-xl bg-white">
            <CardHeader className="flex flex-row items-center gap-2 border-b pb-4">
                <Users className="w-6 h-6 text-blue-600"/>
                <CardTitle className="text-2xl font-bold text-left text-gray-800">従業員一覧</CardTitle>
            </CardHeader>
            <CardContent className="p-4">
                <Table className="w-full">
                    <TableHeader>
                        <TableRow>
                            <TableHead className="text-left px-4 py-2">ID</TableHead>
                            <TableHead className="text-left px-4 py-2">名前</TableHead>
                            <TableHead className="text-left px-4 py-2">作成日時</TableHead>
                        </TableRow>
                    </TableHeader>
                    <TableBody>
                        {employees.map((emp) => (
                            <TableRow key={emp.id} className="hover:bg-blue-50 transition-colors duration-200">
                                <TableCell className="text-left px-4 py-2 font-medium">{emp.id}</TableCell>
                                <TableCell className="text-left px-4 py-2">
                                    <Link to={`/employees/${emp.id}`} className="text-blue-600 hover:underline">
                                        {emp.name}
                                    </Link>
                                </TableCell>
                                <TableCell className="text-left px-4 py-2">
                                    {new Date(emp.created_at).toLocaleString()}
                                </TableCell>
                            </TableRow>
                        ))}
                    </TableBody>
                </Table>
            </CardContent>
        </Card>
    );
};

export default EmployeeList;
