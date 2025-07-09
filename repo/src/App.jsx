import React, { useState, useEffect } from 'react';
import { BrowserRouter as Router, Routes, Route, Link, useNavigate, useParams, useLocation } from 'react-router-dom';
import './App.css';

function App() {
  const [currentUser, setCurrentUser] = useState(
    JSON.parse(localStorage.getItem('user'))
  );

  return (
    <Router>
      <div className="app-container">
        <header className="app-header">
          <h1>贴吧系统</h1>
          <nav>
            <Link to="/">首页</Link>
            {currentUser ? (
              <>
                <Link to="/post/new">发帖</Link>
                {!currentUser.is_admin && <Link to="/my-posts">我的帖子</Link>}
                <span>欢迎, {currentUser.name}{currentUser.is_admin && ' (管理员)'}</span>
                <button onClick={() => {
                  localStorage.removeItem('user');
                  setCurrentUser(null);
                }}>退出</button>
              </>
            ) : (
              <>
                <Link to="/login">登录</Link>
                <Link to="/register">注册</Link>
              </>
            )}
          </nav>
        </header>

        <main className="app-main">
          <Routes>
            <Route path="/" element={<HomePage currentUser={currentUser} />} />
            <Route path="/post/:id" element={<PostDetailPage currentUser={currentUser} />} />
            <Route path="/post/new" element={
              <ProtectedRoute currentUser={currentUser}>
                <NewPostPage />
              </ProtectedRoute>
            } />
            <Route path="/my-posts" element={
              <ProtectedRoute currentUser={currentUser}>
                <MyPostsPage currentUser={currentUser} />
              </ProtectedRoute>
            } />
            <Route path="/login" element={<LoginPage setCurrentUser={setCurrentUser} />} />
            <Route path="/register" element={<RegisterPage setCurrentUser={setCurrentUser} />} />
          </Routes>
        </main>
      </div>
    </Router>
  );
}

// 路由保护组件
function ProtectedRoute({ children, currentUser }) {
  const navigate = useNavigate();

  useEffect(() => {
    if (!currentUser) {
      navigate('/login');
    }
  }, [currentUser, navigate]);

  return currentUser ? children : null;
}

// 首页组件
function HomePage({ currentUser }) {
  const [posts, setPosts] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const navigate = useNavigate();

  useEffect(() => {
    const fetchPosts = async () => {
      try {
        const response = await fetch('http://localhost:8082/repos');

        if (!response.ok) {
          throw new Error('获取帖子失败');
        }

        const data = await response.json();

        // 确保data是数组
        if (!Array.isArray(data)) {
          throw new Error('返回的数据格式不正确');
        }

        setPosts(data);
        setError(null);
      } catch (err) {
        console.error('获取帖子失败:', err);
        setError(err.message);
        setPosts([]); // 确保posts始终是数组
      } finally {
        setLoading(false);
      }
    };

    fetchPosts();
  }, []);

  const handleDeletePost = async (postId) => {
    if (!window.confirm('确定要删除这个帖子吗？')) {
      return;
    }

    try {
      const response = await fetch(`http://localhost:8082/repos/${postId}`, {
        method: 'DELETE',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          user_id: currentUser.id.toString(),
          is_admin: currentUser.is_admin || false
        })
      });

      if (response.ok) {
        // 刷新帖子列表
        setPosts(posts.filter(post => post.id !== postId));
        alert('删除成功');
      } else {
        const error = await response.json();
        alert(error.error || '删除失败');
      }
    } catch (error) {
      console.error('删除帖子失败:', error);
      alert('删除失败，请重试');
    }
  };

  if (loading) return <div className="loading">加载中...</div>;
  if (error) return <div className="error">错误: {error}</div>;

  return (
    <div className="post-list">
      <h2>帖子列表</h2>
      {posts.length === 0 ? (
        <p>暂无帖子</p>
      ) : (
        <ul>
          {posts.map(post => (
            <li key={post.id}>
              <div onClick={() => navigate(`/post/${post.id}`)} style={{ cursor: 'pointer' }}>
                <h3>{post.repo_name}</h3>
                <p className="post-preview">{post.repo_url.substring(0, 100)}...</p>
                <div className="post-meta">
                  <span>作者: {post.user_name || '匿名'}</span>
                  <span>点赞: {post.like_count || 0}</span>
                  <span>评论: {post.comment_count || 0}</span>
                  <span>发布时间: {new Date(post.created_at).toLocaleString()}</span>
                </div>
              </div>
              {currentUser && (currentUser.is_admin || currentUser.id === post.user_id) && (
                <div className="post-actions">
                  <button
                    onClick={(e) => {
                      e.stopPropagation();
                      handleDeletePost(post.id);
                    }}
                    className="delete-btn"
                  >
                    删除帖子
                  </button>
                </div>
              )}
            </li>
          ))}
        </ul>
      )}
    </div>
  );
}

// 我的帖子页面
function MyPostsPage({ currentUser }) {
  const [posts, setPosts] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const navigate = useNavigate();

  useEffect(() => {
    const fetchMyPosts = async () => {
      try {
        const response = await fetch(`http://localhost:8082/repos?user_id=${currentUser.id}`);

        if (!response.ok) {
          throw new Error('获取我的帖子失败');
        }

        const data = await response.json();
        setPosts(Array.isArray(data) ? data : []);
        setError(null);
      } catch (err) {
        console.error('获取我的帖子失败:', err);
        setError(err.message);
        setPosts([]);
      } finally {
        setLoading(false);
      }
    };

    if (currentUser && currentUser.id) {
      fetchMyPosts();
    }
  }, [currentUser]);

  const handleDeletePost = async (postId) => {
    if (!window.confirm('确定要删除这个帖子吗？')) {
      return;
    }

    try {
      const response = await fetch(`http://localhost:8082/repos/${postId}`, {
        method: 'DELETE',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          user_id: currentUser.id.toString(),
          is_admin: false
        })
      });

      if (response.ok) {
        setPosts(posts.filter(post => post.id !== postId));
        alert('删除成功');
      } else {
        const error = await response.json();
        alert(error.error || '删除失败');
      }
    } catch (error) {
      console.error('删除帖子失败:', error);
      alert('删除失败，请重试');
    }
  };

  if (loading) return <div className="loading">加载中...</div>;
  if (error) return <div className="error">错误: {error}</div>;

  return (
    <div className="post-list">
      <h2>我的帖子</h2>
      {posts.length === 0 ? (
        <p>您还没有发布任何帖子</p>
      ) : (
        <ul>
          {posts.map(post => (
            <li key={post.id}>
              <div onClick={() => navigate(`/post/${post.id}`)} style={{ cursor: 'pointer' }}>
                <h3>{post.repo_name}</h3>
                <p className="post-preview">{post.repo_url.substring(0, 100)}...</p>
                <div className="post-meta">
                  <span>点赞: {post.like_count || 0}</span>
                  <span>评论: {post.comment_count || 0}</span>
                  <span>发布时间: {new Date(post.created_at).toLocaleString()}</span>
                </div>
              </div>
              <div className="post-actions">
                <button
                  onClick={(e) => {
                    e.stopPropagation();
                    handleDeletePost(post.id);
                  }}
                  className="delete-btn"
                >
                  删除帖子
                </button>
              </div>
            </li>
          ))}
        </ul>
      )}
    </div>
  );
}

// 帖子详情页
function PostDetailPage({ currentUser }) {
  const { id } = useParams();
  const [post, setPost] = useState(null);
  const [comments, setComments] = useState([]); // 初始化为空数组而不是null
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [commentContent, setCommentContent] = useState('');
  const [isLiked, setIsLiked] = useState(false);
  const navigate = useNavigate();

  useEffect(() => {
    const fetchData = async () => {
      try {
        setLoading(true);
        setError(null);

        // 获取帖子详情
        const postRes = await fetch(`http://localhost:8082/repos/${id}`);
        if (postRes.ok) {
          const postData = await postRes.json();
          setPost(postData);
        } else {
          throw new Error('帖子不存在或获取失败');
        }

        // 获取评论列表
        const commentsRes = await fetch(`http://localhost:8082/comments?repo_id=${id}`);
        if (commentsRes.ok) {
          const commentsData = await commentsRes.json();
          // 确保comments始终是数组
          setComments(Array.isArray(commentsData) ? commentsData : []);
        } else {
          console.warn('获取评论失败，使用空数组');
          setComments([]);
        }
      } catch (error) {
        console.error('获取数据失败:', error);
        setError(error.message);
        setComments([]); // 确保comments始终是数组
      } finally {
        setLoading(false);
      }
    };

    if (id) {
      fetchData();
    }
  }, [id]);

  const handleLike = async () => {
    if (!currentUser) {
      navigate('/login', { state: { from: `/post/${id}` } });
      return;
    }

    try {
      const response = await fetch('http://localhost:8082/like', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          user_id: currentUser.id.toString(),
          repo_id: id
        })
      });

      if (response.ok) {
        // 刷新帖子数据
        const res = await fetch(`http://localhost:8082/repos/${id}`);
        if (res.ok) {
          const updatedPost = await res.json();
          setPost(updatedPost);
        }
        setIsLiked(!isLiked);
      }
    } catch (error) {
      console.error('点赞失败:', error);
      alert('点赞失败，请重试');
    }
  };

  // 提交评论
  const handleCommentSubmit = async (e) => {
    e.preventDefault();
    if (!currentUser) {
      navigate('/login', { state: { from: `/post/${id}` } });
      return;
    }

    if (!commentContent.trim()) {
      alert('评论内容不能为空');
      return;
    }

    try {
      const response = await fetch('http://localhost:8082/comment', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          user_id: currentUser.id.toString(),
          repo_id: id,
          content: commentContent
        })
      });

      if (response.ok) {
        // 刷新评论数据
        const res = await fetch(`http://localhost:8082/comments?repo_id=${id}`);
        if (res.ok) {
          const updatedComments = await res.json();
          setComments(Array.isArray(updatedComments) ? updatedComments : []);
        }
        setCommentContent('');
      } else {
        throw new Error('评论失败');
      }
    } catch (error) {
      console.error('评论失败:', error);
      alert('评论失败，请重试');
    }
  };

  const handleDeletePost = async () => {
    if (!window.confirm('确定要删除这个帖子吗？')) {
      return;
    }

    try {
      const response = await fetch(`http://localhost:8082/repos/${id}`, {
        method: 'DELETE',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          user_id: currentUser.id.toString(),
          is_admin: currentUser.is_admin || false
        })
      });

      if (response.ok) {
        alert('删除成功');
        navigate('/');
      } else {
        const error = await response.json();
        alert(error.error || '删除失败');
      }
    } catch (error) {
      console.error('删除帖子失败:', error);
      alert('删除失败，请重试');
    }
  };

  const handleDeleteComment = async (commentId) => {
    if (!window.confirm('确定要删除这条评论吗？')) {
      return;
    }

    try {
      const response = await fetch(`http://localhost:8082/comments/${commentId}`, {
        method: 'DELETE',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          user_id: currentUser.id.toString(),
          is_admin: currentUser.is_admin || false
        })
      });

      if (response.ok) {
        setComments(comments.filter(comment => comment.id !== commentId));
        alert('删除成功');
      } else {
        const error = await response.json();
        alert(error.error || '删除失败');
      }
    } catch (error) {
      console.error('删除评论失败:', error);
      alert('删除失败，请重试');
    }
  };

  if (loading) return <div className="loading">加载中...</div>;
  if (error) return <div className="error">错误: {error}</div>;
  if (!post) return <div className="error">帖子不存在</div>;

  return (
    <div className="post-detail">
      <h2>{post.repo_name || '无标题'}</h2>
      <div className="post-author">作者: {post.user_name || '匿名'}</div>
      <div className="post-time">发布时间: {post.created_at ? new Date(post.created_at).toLocaleString() : '未知'}</div>
      <div className="post-content">{post.repo_url || '无内容'}</div>

      <div className="post-actions">
        <button onClick={handleLike} className={isLiked ? 'liked' : ''}>
          {isLiked ? '已点赞' : '点赞'} ({post.like_count || 0})
        </button>
        {currentUser && (currentUser.is_admin || currentUser.id === post.user_id) && (
          <button onClick={handleDeletePost} className="delete-btn">
            删除帖子
          </button>
        )}
      </div>

      <div className="comment-section">
        <h3>评论 ({Array.isArray(comments) ? comments.length : 0})</h3>

        {currentUser && (
          <form onSubmit={handleCommentSubmit} className="comment-form">
            <textarea
              value={commentContent}
              onChange={(e) => setCommentContent(e.target.value)}
              placeholder="写下你的评论..."
              required
            />
            <button type="submit">发表评论</button>
          </form>
        )}

        <ul className="comment-list">
          {Array.isArray(comments) && comments.length > 0 ? (
            comments.map(comment => (
              <li key={comment.id} className="comment-item">
                <div className="comment-content">{comment.content || '无内容'}</div>
                <div className="comment-meta">
                  <span>用户: {comment.user_name || '匿名'}</span>
                  <span>{comment.created_at ? new Date(comment.created_at).toLocaleString() : '未知时间'}</span>
                  {currentUser && (currentUser.is_admin || currentUser.id === comment.user_id) && (
                    <button
                      onClick={() => handleDeleteComment(comment.id)}
                      className="delete-btn small"
                    >
                      删除
                    </button>
                  )}
                </div>
              </li>
            ))
          ) : (
            <li className="no-comments">暂无评论</li>
          )}
        </ul>
      </div>
    </div>
  );
}

// 发帖页面
function NewPostPage() {
  const [repoName, setRepoName] = useState('');
  const [repoUrl, setRepoUrl] = useState('');
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();

  const handleSubmit = async (e) => {
    e.preventDefault();
    if (!repoName.trim() || !repoUrl.trim()) {
      alert('标题和内容不能为空');
      return;
    }

    setLoading(true);
    try {
      const user = JSON.parse(localStorage.getItem('user'));
      const response = await fetch('http://localhost:8082/repo', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          user_id: user.id.toString(),
          repo_name: repoName,
          repo_url: repoUrl
        })
      });

      if (!response.ok) {
        throw new Error('发帖失败');
      }

      const result = await response.json();
      // 确保返回的数据结构正确
      if (result.repo && result.repo.id) {
        navigate(`/post/${result.repo.id}`);
      } else {
        // 如果返回的数据结构不符合预期，跳转到首页
        navigate('/');
      }
    } catch (error) {
      console.error('发帖失败:', error);
      alert('发帖失败，请重试');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="new-post">
      <h2>发表新帖</h2>
      <form onSubmit={handleSubmit}>
        <div className="form-group">
          <label>标题</label>
          <input
            type="text"
            value={repoName}
            onChange={(e) => setRepoName(e.target.value)}
            placeholder="请输入帖子标题"
            required
          />
        </div>
        <div className="form-group">
          <label>内容</label>
          <textarea
            value={repoUrl}
            onChange={(e) => setRepoUrl(e.target.value)}
            placeholder="请输入帖子内容"
            rows="10"
            required
          />
        </div>
        <button type="submit" disabled={loading}>
          {loading ? '发布中...' : '发布'}
        </button>
      </form>
    </div>
  );
}

// 注册页面
function RegisterPage({ setCurrentUser }) {
  const [formData, setFormData] = useState({
    name: '',
    username: '',
    password: '',
    confirmPassword: ''
  });
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();

  const handleChange = (e) => {
    setFormData({
      ...formData,
      [e.target.name]: e.target.value
    });
  };

  const handleSubmit = async (e) => {
    e.preventDefault();

    if (!formData.name.trim() || !formData.username.trim() || !formData.password.trim()) {
      alert('所有字段都不能为空');
      return;
    }

    if (formData.password !== formData.confirmPassword) {
      alert('两次输入的密码不一致');
      return;
    }

    setLoading(true);
    try {
      const response = await fetch('http://localhost:8082/register', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          name: formData.name,
          username: formData.username,
          password: formData.password
        })
      });

      if (!response.ok) {
        throw new Error('注册失败');
      }

      const result = await response.json();
      localStorage.setItem('user', JSON.stringify(result.user));
      setCurrentUser(result.user);
      navigate('/');
    } catch (error) {
      console.error('注册失败:', error);
      alert('注册失败，用户名可能已存在');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="register-page">
      <h2>注册</h2>
      <form onSubmit={handleSubmit}>
        <div className="form-group">
          <label>姓名</label>
          <input
            type="text"
            name="name"
            value={formData.name}
            onChange={handleChange}
            placeholder="请输入您的姓名"
            required
          />
        </div>
        <div className="form-group">
          <label>用户名</label>
          <input
            type="text"
            name="username"
            value={formData.username}
            onChange={handleChange}
            placeholder="请输入用户名"
            required
          />
        </div>
        <div className="form-group">
          <label>密码</label>
          <input
            type="password"
            name="password"
            value={formData.password}
            onChange={handleChange}
            placeholder="请输入密码"
            required
          />
        </div>
        <div className="form-group">
          <label>确认密码</label>
          <input
            type="password"
            name="confirmPassword"
            value={formData.confirmPassword}
            onChange={handleChange}
            placeholder="请再次输入密码"
            required
          />
        </div>
        <button type="submit" disabled={loading}>
          {loading ? '注册中...' : '注册'}
        </button>
      </form>
      <p>
        已有账号？<Link to="/login">立即登录</Link>
      </p>
    </div>
  );
}

// 登录页面
function LoginPage({ setCurrentUser }) {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();
  const location = useLocation();

  const from = location.state?.from?.pathname || '/';

  const handleSubmit = async (e) => {
    e.preventDefault();
    if (!username.trim() || !password.trim()) {
      alert('用户名和密码不能为空');
      return;
    }

    setLoading(true);
    try {
      const response = await fetch('http://localhost:8082/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ username, password })
      });

      if (!response.ok) {
        throw new Error('登录失败');
      }

      const result = await response.json();
      localStorage.setItem('user', JSON.stringify(result.user));
      setCurrentUser(result.user);
      navigate(from, { replace: true });
    } catch (error) {
      console.error('登录失败:', error);
      alert('用户名或密码错误');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="login-page">
      <h2>登录</h2>
      <form onSubmit={handleSubmit}>
        <div className="form-group">
          <label>用户名</label>
          <input
            type="text"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
            placeholder="请输入用户名"
            required
          />
        </div>
        <div className="form-group">
          <label>密码</label>
          <input
            type="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            placeholder="请输入密码"
            required
          />
        </div>
        <button type="submit" disabled={loading}>
          {loading ? '登录中...' : '登录'}
        </button>
      </form>
      <p>
        没有账号？<Link to="/register">立即注册</Link>
      </p>
    </div>
  );
}

export default App;

