import { useParams } from 'react-router-dom';

export default function ProductDetail() {
    const { id } = useParams();
    return <h2 style={{ padding: '20px' }}>商品详情 ID: {id}</h2>;
}