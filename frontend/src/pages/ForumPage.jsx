import React, { useEffect } from 'react';
import { useSelector, useDispatch } from 'react-redux';
import { fetchPosts, createNewPost } from '../store/forumSlice';
import PostList from '../components/forum/PostList';
import PostForm from '../components/forum/PostForm';
import { Box, Typography } from '@mui/material';

const ForumPage = () => {
  const dispatch = useDispatch();
  const { posts, loading, error } = useSelector(state => state.forum);
  const { isAuthenticated } = useSelector(state => state.auth);

  useEffect(() => {
    dispatch(fetchPosts());
  }, [dispatch]);

  const handleCreatePost = (postData) => {
    dispatch(createNewPost(postData));
  };

  if (loading) return <Typography>Загрузка постов...</Typography>;
  if (error) return <Typography color="error">Ошибка: {error}</Typography>;

  return (
    <Box sx={{ p: 3 }}>
      <Typography variant="h4" gutterBottom>
        Форум
      </Typography>
      {isAuthenticated && <PostForm onSubmit={handleCreatePost} />}
      <PostList posts={posts} />
    </Box>
  );
};

export default ForumPage;