import { Routes, Route, Link } from 'react-router-dom';

function Orders() {
    return <h3 style={{ padding: '20px' }}>我的订单页面</h3>;
}

function Settings() {
    return <h3 style={{ padding: '20px' }}>设置页面</h3>;
}

export default function Profile() {
    return (
        <div style={{ padding: '20px' }}>
            <h2>用户中心</h2>
            <nav>
                <Link to="orders" style={{ marginRight: '10px' }}>我的订单</Link>
                <Link to="settings">设置</Link>
            </nav>

            <Routes>
                <Route path="orders" element={<Orders />} />
                <Route path="settings" element={<Settings />} />
                <Route index element={<div>请选择子页面</div>} />
            </Routes>
        </div>
    );
}