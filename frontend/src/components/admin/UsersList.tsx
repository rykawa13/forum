import React, { useEffect, useState } from 'react';
import { Table, Space, Button, Popconfirm, message } from 'antd';
import { DeleteOutlined } from '@ant-design/icons';
import { api } from '../../api';
import UserRoleManager from './UserRoleManager';

interface User {
  id: number;
  username: string;
  email: string;
  is_admin: boolean;
  created_at: string;
}

const UsersList: React.FC = () => {
  const [users, setUsers] = useState<User[]>([]);
  const [loading, setLoading] = useState(false);

  const fetchUsers = async () => {
    setLoading(true);
    try {
      const response = await api.get('/api/admin/users');
      setUsers(response.data);
    } catch (error) {
      message.error('Ошибка при загрузке пользователей');
      console.error('Error fetching users:', error);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchUsers();
  }, []);

  const handleDelete = async (userId: number) => {
    try {
      await api.delete(`/api/admin/users/${userId}`);
      message.success('Пользователь успешно удален');
      fetchUsers();
    } catch (error) {
      message.error('Ошибка при удалении пользователя');
      console.error('Error deleting user:', error);
    }
  };

  const columns = [
    {
      title: 'ID',
      dataIndex: 'id',
      key: 'id',
    },
    {
      title: 'Имя пользователя',
      dataIndex: 'username',
      key: 'username',
    },
    {
      title: 'Email',
      dataIndex: 'email',
      key: 'email',
    },
    {
      title: 'Роль',
      key: 'role',
      render: (text: string, record: User) => (
        <UserRoleManager
          userId={record.id}
          isAdmin={record.is_admin}
          onRoleUpdate={fetchUsers}
        />
      ),
    },
    {
      title: 'Дата регистрации',
      dataIndex: 'created_at',
      key: 'created_at',
      render: (date: string) => new Date(date).toLocaleDateString(),
    },
    {
      title: 'Действия',
      key: 'actions',
      render: (text: string, record: User) => (
        <Space>
          <Popconfirm
            title="Вы уверены, что хотите удалить этого пользователя?"
            onConfirm={() => handleDelete(record.id)}
            okText="Да"
            cancelText="Нет"
          >
            <Button
              type="link"
              danger
              icon={<DeleteOutlined />}
              disabled={record.is_admin}
            />
          </Popconfirm>
        </Space>
      ),
    },
  ];

  return (
    <Table
      columns={columns}
      dataSource={users}
      rowKey="id"
      loading={loading}
      pagination={{ pageSize: 10 }}
    />
  );
};

export default UsersList; 