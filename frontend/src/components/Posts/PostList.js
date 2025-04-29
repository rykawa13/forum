import React, { useEffect, useState } from 'react';
import { fetchPosts } from '../../services/postService';

const PostList = () => {
  const [posts, setPosts] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    const loadPosts = async () => {
      try {
        const data = await fetchPosts();
        setPosts(data);
      } catch (err) {
        setError('Failed to load posts');
      } finally {
        setLoading(false);
      }
    };
    loadPosts();
  }, []);

  if (loading) return <div>Loading posts...</div>;
  if (error) return <div className="error">{error}</div>;

  return (
    <div className="post-list">
      {posts.map(post => (
        <div key={post.id} className="post-item">
          <h4>{post.title}</h4>
          <p>{post.content}</p>
          <div className="post-meta">
            <span>By {post.author}</span>
            <span>{new Date(post.created_at).toLocaleDateString()}</span>
          </div>
        </div>
      ))}
    </div>
  );
};

export default PostList;