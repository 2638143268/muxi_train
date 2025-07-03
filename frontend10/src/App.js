import { BrowserRouter as Router, Routes, Route, Link } from 'react-router-dom';
import Home from './Home';
import About from './About';
import Products from './Products';
import ProductDetail from './ProductDetail';
import Profile from './Profile';
import Login from './Login';
import NotFound from './NotFound';

function App() {
  return (
    <Router>
      <nav style={{ padding: '10px', background: '#f0f0f0' }}>
        <Link to="/" style={{ marginRight: '10px' }}>首页</Link>
        <Link to="/products" style={{ marginRight: '10px' }}>商品</Link>
        <Link to="/profile">用户中心</Link>
      </nav>

      <Routes>
        <Route path="/" element={<Home />} />
        <Route path="/about" element={<About />} />
        <Route path="/products" element={<Products />} />
        <Route path="/products/:id" element={<ProductDetail />} />
        <Route path="/profile/*" element={<Profile />} />
        <Route path="/login" element={<Login />} />
        <Route path="*" element={<NotFound />} />
      </Routes>
    </Router>
  );
}

export default App;