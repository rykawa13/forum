import React, { useState } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { Avatar, Typography, Box, Divider, Button } from '@mui/material';
import { format, parseISO } from 'date-fns';
import { ru } from 'date-fns/locale';
import CommentForm from './CommentForm';
import CommentList from './CommentList';
import { addComment } from '../../store/slices/commentsSlice';

const PostItem = ({ post }) => {
  console.log('PostItem render:', post);

  const [showComments, setShowComments] = useState(false);
  const dispatch = useDispatch();
  const userId = useSelector(state => state.auth.user?.id);

  if (!post || typeof post !== 'object') {
    console.error('Invalid post data:', post);
    return null;
  }

  // Получаем имя автора или используем fallback
  const authorName = post.author?.username || `User ${post.author_id}`;
  const avatarLetter = authorName.charAt(0).toUpperCase();

  // Безопасное форматирование даты
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

  const handleToggleComments = () => {
    setShowComments(!showComments);
  };

  const handleAddComment = async (commentData) => {
    if (!userId) {
      alert('Пожалуйста, войдите в систему, чтобы оставить комментарий');
      return;
    }
    
    dispatch(addComment({
      postId: post.id,
      content: commentData.content,
      userId
    }));
  };

  return (
    <Box sx={{ width: '100%', mb: 3, p: 2, bgcolor: 'background.paper', borderRadius: 1, boxShadow: 1 }}>
      <Box sx={{ display: 'flex', alignItems: 'center', mb: 1 }}>
        <Avatar sx={{ mr: 1 }}>{avatarLetter}</Avatar>
        <Box>
          <Typography variant="subtitle1" fontWeight="bold">
            {authorName}
          </Typography>
          <Typography variant="caption" color="text.secondary">
            {formatDate(post.created_at)}
          </Typography>
        </Box>
      </Box>
      <Typography variant="h6" sx={{ mb: 1 }}>
        {post.title}
      </Typography>
      <Typography variant="body1" sx={{ mb: 2 }}>
        {post.content}
      </Typography>
      
      <Button 
        size="small" 
        onClick={handleToggleComments}
        sx={{ mb: 1 }}
      >
        {showComments ? 'Скрыть комментарии' : 'Показать комментарии'}
      </Button>

      {showComments && (
        <Box sx={{ pl: 2 }}>
          <Divider sx={{ mb: 2 }} />
          <CommentList postId={post.id} />
          <CommentForm postId={post.id} onSubmit={handleAddComment} />
        </Box>
      )}
    </Box>
  );
};

export default PostItem;