import React from 'react';
import { Grid, Paper, Typography, Box } from '@mui/material';
import { 
  People as PeopleIcon,
  Chat as ChatIcon,
  Forum as ForumIcon,
  AccessTime as AccessTimeIcon 
} from '@mui/icons-material';

const StatsPanel = ({ stats }) => {
  const statItems = [
    { 
      title: 'Пользователи', 
      value: stats?.usersCount || 0,
      icon: <PeopleIcon color="primary" sx={{ fontSize: 40 }} />
    },
    { 
      title: 'Сообщения чата', 
      value: stats?.messagesCount || 0,
      icon: <ChatIcon color="primary" sx={{ fontSize: 40 }} />
    },
    { 
      title: 'Посты форума', 
      value: stats?.postsCount || 0,
      icon: <ForumIcon color="primary" sx={{ fontSize: 40 }} />
    },
    { 
      title: 'Активных сегодня', 
      value: stats?.activeToday || 0,
      icon: <AccessTimeIcon color="primary" sx={{ fontSize: 40 }} />
    }
  ];

  return (
    <Grid container spacing={3} sx={{ mb: 4 }}>
      {statItems.map((item, index) => (
        <Grid item xs={12} sm={6} md={3} key={index}>
          <Paper sx={{ p: 3, height: '100%' }}>
            <Box sx={{ display: 'flex', alignItems: 'center', mb: 2 }}>
              {item.icon}
              <Typography variant="h6" sx={{ ml: 2 }}>
                {item.title}
              </Typography>
            </Box>
            <Typography variant="h4" component="div" sx={{ fontWeight: 'bold' }}>
              {item.value}
            </Typography>
          </Paper>
        </Grid>
      ))}
    </Grid>
  );
};

export default StatsPanel;