import React, { useState, useEffect } from 'react';
import {
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Paper,
  IconButton,
  Button,
  Typography,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  TextField,
  Switch,
  FormControlLabel,
} from '@mui/material';
import { Delete as DeleteIcon, Edit as EditIcon } from '@mui/icons-material';
import axios from 'axios';

const UserList = () => {
  const [users, setUsers] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [editDialogOpen, setEditDialogOpen] = useState(false);
  const [selectedUser, setSelectedUser] = useState(null);

  const fetchUsers = async () => {
    try {
      setLoading(true);
      const response = await axios.get('http://localhost:8081/api/users', {
        headers: {
          Authorization: `Bearer ${localStorage.getItem('token')}`
        }
      });
      setUsers(response.data);
      setError(null);
    } catch (err) {
      setError('Ошибка при загрузке пользователей');
      console.error('Error fetching users:', err);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchUsers();
  }, []);

  const handleEditClick = (user) => {
    setSelectedUser(user);
    setEditDialogOpen(true);
  };

  const handleEditClose = () => {
    setEditDialogOpen(false);
    setSelectedUser(null);
  };

  const handleEditSave = async () => {
    try {
      await axios.put(`http://localhost:8081/api/users/${selectedUser.id}`, selectedUser, {
        headers: {
          Authorization: `Bearer ${localStorage.getItem('token')}`
        }
      });
      fetchUsers();
      handleEditClose();
    } catch (err) {
      console.error('Error updating user:', err);
      setError('Ошибка при обновлении пользователя');
    }
  };

  const handleDeleteUser = async (userId) => {
    if (window.confirm('Вы уверены, что хотите удалить этого пользователя?')) {
      try {
        await axios.delete(`http://localhost:8081/api/users/${userId}`, {
          headers: {
            Authorization: `Bearer ${localStorage.getItem('token')}`
          }
        });
        fetchUsers();
      } catch (err) {
        console.error('Error deleting user:', err);
        setError('Ошибка при удалении пользователя');
      }
    }
  };

  if (loading) return <Typography>Загрузка...</Typography>;
  if (error) return <Typography color="error">{error}</Typography>;

  return (
    <>
      <TableContainer component={Paper}>
        <Table>
          <TableHead>
            <TableRow>
              <TableCell>ID</TableCell>
              <TableCell>Имя пользователя</TableCell>
              <TableCell>Email</TableCell>
              <TableCell>Админ</TableCell>
              <TableCell>Действия</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {users.map((user) => (
              <TableRow key={user.id}>
                <TableCell>{user.id}</TableCell>
                <TableCell>{user.username}</TableCell>
                <TableCell>{user.email}</TableCell>
                <TableCell>{user.is_admin ? 'Да' : 'Нет'}</TableCell>
                <TableCell>
                  <IconButton onClick={() => handleEditClick(user)} color="primary">
                    <EditIcon />
                  </IconButton>
                  <IconButton onClick={() => handleDeleteUser(user.id)} color="error">
                    <DeleteIcon />
                  </IconButton>
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </TableContainer>

      <Dialog open={editDialogOpen} onClose={handleEditClose}>
        <DialogTitle>Редактировать пользователя</DialogTitle>
        <DialogContent>
          {selectedUser && (
            <>
              <TextField
                margin="dense"
                label="Имя пользователя"
                fullWidth
                value={selectedUser.username}
                onChange={(e) => setSelectedUser({ ...selectedUser, username: e.target.value })}
              />
              <TextField
                margin="dense"
                label="Email"
                fullWidth
                value={selectedUser.email}
                onChange={(e) => setSelectedUser({ ...selectedUser, email: e.target.value })}
              />
              <FormControlLabel
                control={
                  <Switch
                    checked={selectedUser.is_admin}
                    onChange={(e) => setSelectedUser({ ...selectedUser, is_admin: e.target.checked })}
                  />
                }
                label="Администратор"
              />
            </>
          )}
        </DialogContent>
        <DialogActions>
          <Button onClick={handleEditClose}>Отмена</Button>
          <Button onClick={handleEditSave} variant="contained" color="primary">
            Сохранить
          </Button>
        </DialogActions>
      </Dialog>
    </>
  );
};

export default UserList;