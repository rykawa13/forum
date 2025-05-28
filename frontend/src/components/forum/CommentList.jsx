import React, { useEffect } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { Box, Typography, CircularProgress } from '@mui/material';
import CommentItem from './CommentItem';
import { fetchComments } from '../../store/slices/commentsSlice';

const CommentList = ({ postId }) => {
  const dispatch = useDispatch();
  const { items, loading, error } = useSelector(state => state.comments);
  const comments = items[postId] || [];

  useEffect(() => {
    if (postId) {
      dispatch(fetchComments(postId));
    }
  }, [dispatch, postId]);

  if (loading) {
    return (
      <Box sx={{ display: 'flex', justifyContent: 'center', my: 2 }}>
        <CircularProgress size={24} />
      </Box>
    );
  }

  if (error) {
    return (
      <Typography color="error" variant="body2" sx={{ my: 2 }}>
        {error}
      </Typography>
    );
  }

  if (!comments.length) {
    return (
      <Typography variant="body2" color="text.secondary" sx={{ my: 2 }}>
        Пока нет комментариев. Будьте первым!
      </Typography>
    );
  }

  return (
    <Box sx={{ mt: 2 }}>
      {comments.map((comment) => (
        <CommentItem key={comment.id} comment={comment} />
      ))}
    </Box>
  );
};

export default CommentList; 