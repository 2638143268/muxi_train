import { useState, useEffect, useCallback } from 'react';

// 自定义Hook useCounter
function useCounter(initialValue = 0) {
    const [count, setCount] = useState(initialValue);

    // 使用useEffect在数值变化时打印控制台日志
    useEffect(() => {
        console.log(`计数器值已更新为: ${count}`);
    }, [count]);

    // 使用useCallback优化事件处理函数
    const increment = useCallback(() => {
        setCount(prev => prev + 1);
    }, []);

    const decrement = useCallback(() => {
        setCount(prev => prev - 1);
    }, []);

    const reset = useCallback(() => {
        setCount(initialValue);
    }, [initialValue]);

    return { count, increment, decrement, reset };
}

function CounterApp() {
    const { count, increment, decrement, reset } = useCounter(0);

    return (
        <div style={styles.container}>
            <h1 style={styles.title}>计数器</h1>
            <div style={styles.counterDisplay}>{count}</div>
            <div style={styles.buttonGroup}>
                <button style={styles.button} onClick={decrement}>-</button>
                <button style={styles.button} onClick={reset}>重置</button>
                <button style={styles.button} onClick={increment}>+</button>
            </div>
        </div>
    );
}

// CSS样式
const styles = {
    container: {
        display: 'flex',
        flexDirection: 'column',
        alignItems: 'center',
        justifyContent: 'center',
        height: '100vh',
        fontFamily: 'Arial, sans-serif',
    },
    title: {
        fontSize: '2rem',
        marginBottom: '20px',
    },
    counterDisplay: {
        fontSize: '3rem',
        margin: '20px 0',
        padding: '20px 40px',
        border: '2px solid #333',
        borderRadius: '8px',
    },
    buttonGroup: {
        display: 'flex',
        gap: '10px',
    },
    button: {
        padding: '10px 20px',
        fontSize: '1.2rem',
        cursor: 'pointer',
        backgroundColor: '#f0f0f0',
        border: '1px solid #ccc',
        borderRadius: '4px',
        transition: 'background-color 0.3s',
    },
};

export default CounterApp;