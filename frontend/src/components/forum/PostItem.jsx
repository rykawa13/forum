import React from 'react';
import { Avatar, Typography, Box } from '@mui/material';
import { format } from 'date-fns';
import { ru } from 'date-fns/locale';

const PostItem = ({ post }) => {
  return (
    <Box sx={{ width: '100%' }}>
      <Box sx={{ display: 'flex', alignItems: 'center', mb: 1 }}>
        <Avatar sx={{ mr: 1 }}>{post.author.username.charAt(0)}</Avatar>
        <Box>
          <Typography variant="subtitle1" fontWeight="bold">
            {post.author.username}
          </Typography>
          <Typography variant="caption" color="text.secondary">
            {format(new Date(post.createdAt), 'dd MMMM yyyy HH:mm', { locale: ru })}
          </Typography>
        </Box>
      </Box>
      <Typography variant="body1" sx={{ mb: 1 }}>
        {post.title}
      </Typography>
      <Typography variant="body2" color="text.secondary">
        {post.content}
      </Typography>
    </Box>
  );
};

export default PostItem;