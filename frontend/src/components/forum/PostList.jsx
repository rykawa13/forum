import React from 'react';
import { useSelector } from 'react-redux';
import { List, ListItem, Divider, Typography, Box } from '@mui/material';
import PostItem from './PostItem';

const PostList = ({ posts }) => {
  if (!posts || posts.length === 0) {
    return (
      <Box sx={{ p: 2, textAlign: 'center' }}>
        <Typography variant="body1">Пока нет сообщений на форуме</Typography>
      </Box>
    );
  }

  return (
    <List sx={{ width: '100%', bgcolor: 'background.paper' }}>
      {posts.map((post, index) => (
        <React.Fragment key={post.id}>
          <ListItem alignItems="flex-start">
            <PostItem post={post} />
          </ListItem>
          {index < posts.length - 1 && <Divider component="li" />}
        </React.Fragment>
      ))}
    </List>
  );
};

export default PostList;