import React, { useEffect, useMemo, useCallback } from 'react';
import { useSelector, useDispatch } from 'react-redux';
import { fetchPosts, createNewPost } from '../store/forumSlice';
import PostList from '../components/forum/PostList';
import { Box, Typography } from '@mui/material';

// Создаем селекторы вне компонента для оптимизации
const selectForumPosts = state => state.forum.posts;
const selectForumLoading = state => state.forum.loading;
const selectForumError = state => state.forum.error;
const selectAuth = state => state.auth;

const ForumPage = () => {
  const dispatch = useDispatch();
  
  // Используем отдельные селекторы для каждого поля
  const posts = useSelector(selectForumPosts);
  const loading = useSelector(selectForumLoading);
  const error = useSelector(selectForumError);
  const { isAuthenticated, user } = useSelector(selectAuth);

  useEffect(() => {
    console.log('ForumPage: Fetching posts');
    dispatch(fetchPosts());
  }, [dispatch]);

  // Добавляем отладочный вывод при изменении данных
  useEffect(() => {
    console.log('ForumPage: Posts state changed', {
      postsCount: posts?.length || 0,
      loading,
      error,
      firstPost: posts?.[0]
    });
  }, [posts, loading, error]);

  const handleCreatePost = useCallback((postData) => {
    if (!isAuthenticated) {
      alert('Для создания поста необходимо войти в систему');
      return;
    }
    dispatch(createNewPost({
      ...postData,
      author_id: user.id
    }));
  }, [dispatch, isAuthenticated, user]);

  // Мемоизируем пропсы для PostList
  const postListProps = useMemo(() => ({
    posts: posts || [],
    loading,
    error,
    onCreatePost: handleCreatePost
  }), [posts, loading, error, handleCreatePost]);

  return (
    <Box sx={{ p: 3 }}>
      <Typography variant="h4" gutterBottom>
        Форум
      </Typography>
      <PostList {...postListProps} />
    </Box>
  );
};

export default ForumPage;