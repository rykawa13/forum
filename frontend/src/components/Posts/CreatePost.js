import React, { useState } from 'react';
import { createPost } from '../../services/postService';
import { useAuth } from '../../hooks/useAuth';

const CreatePost = () => {
  const [postContent, setPostContent] = useState('');
  const [error, setError] = useState('');
  const { isAuthenticated } = useAuth();

  const handleSubmit = async (e) => {
    e.preventDefault();
    if (!isAuthenticated) {
      setError('You must be logged in to create a post');
      return;
    }
    try {
      await createPost({ content: postContent });
      setPostContent('');
      setError('');
    } catch (err) {
      setError('Error creating post');
    }
  };

  return (
    <div className="create-post">
      <h3>Create New Post</h3>
      <form onSubmit={handleSubmit}>
        <textarea
          value={postContent}
          onChange={(e) => setPostContent(e.target.value)}
          placeholder="Write your post..."
          disabled={!isAuthenticated}
        />
        {error && <div className="error">{error}</div>}
        <button 
          type="submit" 
          disabled={!isAuthenticated}
        >
          Post
        </button>
      </form>
    </div>
  );
};

export default CreatePost;