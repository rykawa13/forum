import React from 'react';
import { Switch, message } from 'antd';
import { api } from '../../api';

interface UserRoleManagerProps {
  userId: number;
  isAdmin: boolean;
  onRoleUpdate?: () => void;
}

const UserRoleManager: React.FC<UserRoleManagerProps> = ({ userId, isAdmin, onRoleUpdate }) => {
  const handleRoleChange = async (checked: boolean) => {
    try {
      await api.put(`/api/admin/users/${userId}/role`, { is_admin: checked });
      message.success('Роль пользователя успешно обновлена');
      if (onRoleUpdate) {
        onRoleUpdate();
      }
    } catch (error) {
      message.error('Ошибка при обновлении роли пользователя');
      console.error('Error updating user role:', error);
    }
  };

  return (
    <Switch
      checked={isAdmin}
      onChange={handleRoleChange}
      checkedChildren="Админ"
      unCheckedChildren="Пользователь"
    />
  );
};

export default UserRoleManager; 