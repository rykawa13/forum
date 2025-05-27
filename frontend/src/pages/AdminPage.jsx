import React, { useEffect } from 'react';
import { useSelector, useDispatch } from 'react-redux';
import { fetchUsers, deleteUser } from '../store/adminSlice';
import UserList from '../components/admin/UserList';
import { Box, Typography } from '@mui/material';

const AdminPage = () => {
  const dispatch = useDispatch();
  const { users, loading, error } = useSelector(state => state.admin);
  const { user } = useSelector(state => state.auth);

  useEffect(() => {
    if (user?.isAdmin) {
      dispatch(fetchUsers());
    }
  }, [dispatch, user]);

  const handleDeleteUser = (userId) => {
    if (window.confirm('Вы уверены, что хотите удалить этого пользователя?')) {
      dispatch(deleteUser(userId));
    }
  };

  if (!user?.isAdmin) {
    return (
      <Box sx={{ p: 3 }}>
        <Typography color="error">У вас нет доступа к этой странице</Typography>
      </Box>
    );
  }

  if (loading) return <Typography>Загрузка пользователей...</Typography>;
  if (error) return <Typography color="error">Ошибка: {error}</Typography>;

  return (
    <Box sx={{ p: 3 }}>
      <Typography variant="h4" gutterBottom>
        Панель администратора
      </Typography>
      <UserList users={users} onDeleteUser={handleDeleteUser} />
    </Box>
  );
};

export default AdminPage;