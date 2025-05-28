import React, { useMemo } from 'react';
import { Box, Typography, CircularProgress } from '@mui/material';
import PostItem from './PostItem';
import PostForm from './PostForm';

const PostList = ({ posts = [], loading, error, onCreatePost }) => {
  console.log('PostList render:', { posts, loading, error });
  
  // Мемоизируем проверку массива и первого поста
  const { isValidArray, postsCount, firstPost } = useMemo(() => ({
    isValidArray: Array.isArray(posts),
    postsCount: Array.isArray(posts) ? posts.length : 0,
    firstPost: Array.isArray(posts) && posts.length > 0 ? posts[0] : null
  }), [posts]);

  console.log('Posts validation:', {
    isValidArray,
    postsCount,
    firstPost
  });

  if (loading) {
    return (
      <Box sx={{ display: 'flex', justifyContent: 'center', mt: 4 }}>
        <CircularProgress />
      </Box>
    );
  }

  if (error) {
    return (
      <Typography color="error" sx={{ mt: 2 }}>
        {error}
      </Typography>
    );
  }

  // Проверяем, что posts определен и является массивом
  const validPosts = isValidArray ? posts : [];

  return (
    <Box>
      <PostForm onSubmit={onCreatePost} />
      {validPosts.length === 0 ? (
        <Typography variant="body1" sx={{ mt: 2 }}>
          Пока нет ни одного поста. Будьте первым!
        </Typography>
      ) : (
        <Box sx={{ mt: 2 }}>
          {validPosts.map((post) => {
            console.log('Rendering post:', post);
            return post ? (
              <PostItem key={post.id} post={post} />
            ) : null;
          })}
        </Box>
      )}
    </Box>
  );
};

export default React.memo(PostList);