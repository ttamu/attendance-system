import React, {useEffect, useState} from 'react';
import {Link} from 'react-router-dom'
import {User} from '../types/User';
import {fetchUsers} from '../services/api';

const UserList: React.FC = () => {
    const [users, setUsers] = useState<User[]>([]);
    const [loading, setLoading] = useState<boolean>(true);

    useEffect(() => {
        const loadUsers = async () => {
            try {
                const usersData = await fetchUsers();
                setUsers(usersData);
            } catch (error) {
                console.error('Error fetching users:', error);
            } finally {
                setLoading(false);
            }
        };

        loadUsers();
    }, []);

    if (loading) {
        return <div>Loading...</div>;
    }

    return (
        <div className="max-w-2xl mx-auto p-4">
            <h1 className="text-xl font-semibold mb-4">User List</h1>
            <table className="w-full border-collapse">
                <thead>
                <tr>
                    <th className="border p-2 text-left">ID</th>
                    <th className="border p-2 text-left">Name</th>
                    <th className="border p-2 text-left">Created At</th>
                </tr>
                </thead>
                <tbody>
                {users.map(user => (
                    <tr key={user.id}>
                        <td className="border p-2 text-left">{user.id}</td>
                        <td className="border p-2 text-left">
                            <Link to={`/users/${user.id}`} className="text-blue-600 hover:underline">
                                {user.name}
                            </Link>
                        </td>
                        <td className="border p-2 text-left">{user.created_at}</td>
                    </tr>
                ))}
                </tbody>
            </table>
        </div>
    );
};

export default UserList;
