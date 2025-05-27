import React, { useState, useEffect } from 'react';
import {
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Paper,
  Button,
  IconButton,
  Typography,
  Box,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Switch,
  Pagination,
} from '@mui/material';
import { Block, LockOpen, AdminPanelSettings, RemoveCircle, Visibility } from '@mui/icons-material';
import { useAuth } from '../../hooks/useAuth';

interface User {
  id: number;
  username: string;
  email: string;
  is_admin: boolean;
  is_blocked: boolean;
  created_at: string;
}

interface Session {
  id: number;
  user_agent: string;
  ip: string;
  created_at: string;
  is_active: boolean;
}

const UserList: React.FC = () => {
  const { token } = useAuth();
  const [users, setUsers] = useState<User[]>([]);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [selectedUser, setSelectedUser] = useState<number | null>(null);
  const [sessions, setSessions] = useState<Session[]>([]);
  const [sessionsDialogOpen, setSessionsDialogOpen] = useState(false);
  const limit = 10;

  const fetchUsers = async () => {
    try {
      const response = await fetch(
        `http://localhost:8081/api/admin/users?offset=${(page - 1) * limit}&limit=${limit}`,
        {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        }
      );
      const data = await response.json();
      setUsers(data.users);
      setTotal(Math.ceil(data.total / limit));
    } catch (error) {
      console.error('Error fetching users:', error);
    }
  };

  const fetchSessions = async (userId: number) => {
    try {
      const response = await fetch(
        `http://localhost:8081/api/admin/users/${userId}/sessions`,
        {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        }
      );
      const data = await response.json();
      setSessions(data.sessions);
      setSelectedUser(userId);
      setSessionsDialogOpen(true);
    } catch (error) {
      console.error('Error fetching sessions:', error);
    }
  };

  const handleRoleChange = async (userId: number, isAdmin: boolean) => {
    try {
      await fetch(`http://localhost:8081/api/admin/users/${userId}/role`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({ is_admin: isAdmin }),
      });
      fetchUsers();
    } catch (error) {
      console.error('Error updating user role:', error);
    }
  };

  const handleBlockUser = async (userId: number, isBlocked: boolean) => {
    try {
      await fetch(`http://localhost:8081/api/admin/users/${userId}/status`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({ is_blocked: isBlocked }),
      });
      fetchUsers();
    } catch (error) {
      console.error('Error updating user status:', error);
    }
  };

  const handleTerminateSessions = async (userId: number) => {
    try {
      await fetch(`http://localhost:8081/api/admin/users/${userId}/sessions`, {
        method: 'DELETE',
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });
      setSessionsDialogOpen(false);
      fetchUsers();
    } catch (error) {
      console.error('Error terminating sessions:', error);
    }
  };

  useEffect(() => {
    fetchUsers();
  }, [page]);

  return (
    <Box sx={{ p: 3 }}>
      <Typography variant="h4" gutterBottom>
        Управление пользователями
      </Typography>
      <TableContainer component={Paper}>
        <Table>
          <TableHead>
            <TableRow>
              <TableCell>ID</TableCell>
              <TableCell>Имя пользователя</TableCell>
              <TableCell>Email</TableCell>
              <TableCell>Дата регистрации</TableCell>
              <TableCell>Администратор</TableCell>
              <TableCell>Статус</TableCell>
              <TableCell>Действия</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {users.map((user) => (
              <TableRow key={user.id}>
                <TableCell>{user.id}</TableCell>
                <TableCell>{user.username}</TableCell>
                <TableCell>{user.email}</TableCell>
                <TableCell>
                  {new Date(user.created_at).toLocaleDateString()}
                </TableCell>
                <TableCell>
                  <Switch
                    checked={user.is_admin}
                    onChange={(e) => handleRoleChange(user.id, e.target.checked)}
                    color="primary"
                  />
                </TableCell>
                <TableCell>
                  {user.is_blocked ? (
                    <Typography color="error">Заблокирован</Typography>
                  ) : (
                    <Typography color="success">Активен</Typography>
                  )}
                </TableCell>
                <TableCell>
                  <IconButton
                    onClick={() => handleBlockUser(user.id, !user.is_blocked)}
                    color={user.is_blocked ? 'success' : 'error'}
                  >
                    {user.is_blocked ? <LockOpen /> : <Block />}
                  </IconButton>
                  <IconButton
                    onClick={() => fetchSessions(user.id)}
                    color="primary"
                  >
                    <Visibility />
                  </IconButton>
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </TableContainer>
      <Box sx={{ mt: 2, display: 'flex', justifyContent: 'center' }}>
        <Pagination
          count={total}
          page={page}
          onChange={(_, value) => setPage(value)}
        />
      </Box>

      <Dialog
        open={sessionsDialogOpen}
        onClose={() => setSessionsDialogOpen(false)}
        maxWidth="md"
        fullWidth
      >
        <DialogTitle>Активные сессии пользователя</DialogTitle>
        <DialogContent>
          <TableContainer>
            <Table>
              <TableHead>
                <TableRow>
                  <TableCell>ID</TableCell>
                  <TableCell>User Agent</TableCell>
                  <TableCell>IP</TableCell>
                  <TableCell>Дата создания</TableCell>
                  <TableCell>Статус</TableCell>
                </TableRow>
              </TableHead>
              <TableBody>
                {sessions.map((session) => (
                  <TableRow key={session.id}>
                    <TableCell>{session.id}</TableCell>
                    <TableCell>{session.user_agent}</TableCell>
                    <TableCell>{session.ip}</TableCell>
                    <TableCell>
                      {new Date(session.created_at).toLocaleDateString()}
                    </TableCell>
                    <TableCell>
                      {session.is_active ? 'Активна' : 'Завершена'}
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </TableContainer>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setSessionsDialogOpen(false)}>Закрыть</Button>
          {selectedUser && (
            <Button
              onClick={() => handleTerminateSessions(selectedUser)}
              color="error"
              startIcon={<RemoveCircle />}
            >
              Завершить все сессии
            </Button>
          )}
        </DialogActions>
      </Dialog>
    </Box>
  );
};

export default UserList; 