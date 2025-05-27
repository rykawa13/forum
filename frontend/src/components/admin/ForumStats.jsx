import React, { useState, useEffect } from 'react';
import {
  Card,
  CardContent,
  Typography,
  Grid,
  CircularProgress,
} from '@mui/material';
import axios from 'axios';

const ForumStats = () => {
  const [stats, setStats] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    const fetchStats = async () => {
      try {
        setLoading(true);
        const response = await axios.get('http://localhost:8081/api/stats', {
          headers: {
            Authorization: `Bearer ${localStorage.getItem('token')}`
          }
        });
        setStats(response.data);
        setError(null);
      } catch (err) {
        setError('Ошибка при загрузке статистики');
        console.error('Error fetching stats:', err);
      } finally {
        setLoading(false);
      }
    };

    fetchStats();
  }, []);

  if (loading) return <CircularProgress />;
  if (error) return <Typography color="error">{error}</Typography>;
  if (!stats) return null;

  return (
    <Grid container spacing={2}>
      <Grid item xs={12} md={4}>
        <Card>
          <CardContent>
            <Typography variant="h6" gutterBottom>
              Пользователи
            </Typography>
            <Typography variant="h4">
              {stats.totalUsers}
            </Typography>
            <Typography variant="body2" color="text.secondary">
              Всего зарегистрировано
            </Typography>
          </CardContent>
        </Card>
      </Grid>
      <Grid item xs={12} md={4}>
        <Card>
          <CardContent>
            <Typography variant="h6" gutterBottom>
              Темы
            </Typography>
            <Typography variant="h4">
              {stats.totalTopics}
            </Typography>
            <Typography variant="body2" color="text.secondary">
              Создано на форуме
            </Typography>
          </CardContent>
        </Card>
      </Grid>
      <Grid item xs={12} md={4}>
        <Card>
          <CardContent>
            <Typography variant="h6" gutterBottom>
              Сообщения
            </Typography>
            <Typography variant="h4">
              {stats.totalPosts}
            </Typography>
            <Typography variant="body2" color="text.secondary">
              Опубликовано всего
            </Typography>
          </CardContent>
        </Card>
      </Grid>
    </Grid>
  );
};

export default ForumStats; 