import { Link } from 'react-router-dom';

export default function Products() {
    return (
        <div style={{ padding: '20px' }}>
            <h2>商品列表</h2>
            <ul>
                <li><Link to="/products/1">商品1</Link></li>
                <li><Link to="/products/2">商品2</Link></li>
                <li><Link to="/products/3">商品3</Link></li>
            </ul>
        </div>
    );
}