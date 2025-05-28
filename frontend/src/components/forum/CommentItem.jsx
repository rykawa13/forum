import React from 'react';
import { Avatar, Typography, Box } from '@mui/material';
import { format, parseISO } from 'date-fns';
import { ru } from 'date-fns/locale';

const CommentItem = ({ comment }) => {
  const authorName = comment.author?.username || `User ${comment.author_id}`;
  const avatarLetter = authorName.charAt(0).toUpperCase();

  const formatDate = (dateString) => {
    try {
      if (!dateString) return '';
      const date = parseISO(dateString);
      return format(date, 'dd MMMM yyyy HH:mm', { locale: ru });
    } catch (error) {
      console.error('Error formatting date:', error);
      return '';
    }
  };

  return (
    <Box sx={{ display: 'flex', alignItems: 'flex-start', mb: 2, pl: 2 }}>
      <Avatar sx={{ width: 32, height: 32, mr: 1, fontSize: '0.875rem' }}>
        {avatarLetter}
      </Avatar>
      <Box>
        <Box sx={{ display: 'flex', alignItems: 'baseline' }}>
          <Typography variant="subtitle2" sx={{ mr: 1 }}>
            {authorName}
          </Typography>
          <Typography variant="caption" color="text.secondary">
            {formatDate(comment.created_at)}
          </Typography>
        </Box>
        <Typography variant="body2" color="text.secondary">
          {comment.content}
        </Typography>
      </Box>
    </Box>
  );
};

export default CommentItem; 