import { createContext, useContext, useState, useEffect } from 'react';

// 创建ThemeContext
const ThemeContext = createContext();

// 自定义Hook useTheme
function useTheme() {
    const [theme, setTheme] = useState('light');

    const toggleTheme = () => {
        setTheme(prevTheme => prevTheme === 'light' ? 'dark' : 'light');
    };

    // 使用useEffect同步body背景色变化
    useEffect(() => {
        document.body.style.backgroundColor = theme === 'light' ? '#fff' : '#333';
        document.body.style.color = theme === 'light' ? '#333' : '#fff';
    }, [theme]);

    return { theme, toggleTheme };
}

// 主题卡片组件
function ThemeCard() {
    const { theme } = useContext(ThemeContext);

    const cardStyle = {
        padding: '20px',
        margin: '20px 0',
        borderRadius: '8px',
        backgroundColor: theme === 'light' ? '#f0f0f0' : '#444',
        color: theme === 'light' ? '#333' : '#fff',
        boxShadow: '0 2px 4px rgba(0,0,0,0.1)',
        textAlign: 'center',
    };

    return (
        <div style={cardStyle}>
            <h2>当前主题: {theme === 'light' ? '明亮模式' : '暗黑模式'}</h2>
        </div>
    );
}

// 主题切换按钮组件
function ThemeButton() {
    const { toggleTheme } = useContext(ThemeContext);

    const buttonStyle = {
        padding: '10px 20px',
        fontSize: '1rem',
        cursor: 'pointer',
        backgroundColor: '#4CAF50',
        color: 'white',
        border: 'none',
        borderRadius: '4px',
    };

    return (
        <button style={buttonStyle} onClick={toggleTheme}>
            切换主题
        </button>
    );
}

function ThemeApp() {
    const theme = useTheme();

    const appStyle = {
        minHeight: '100vh',
        padding: '20px',
        fontFamily: 'Arial, sans-serif',
        backgroundColor: theme.theme === 'light' ? '#fff' : '#333',
        color: theme.theme === 'light' ? '#333' : '#fff',
    };

    return (
        <ThemeContext.Provider value={theme}>
            <div style={appStyle}>
                <h1 style={{ textAlign: 'center' }}>主题切换应用</h1>
                <div style={{ maxWidth: '500px', margin: '0 auto' }}>
                    <ThemeCard />
                    <ThemeButton />
                </div>
            </div>
        </ThemeContext.Provider>
    );
}

export default ThemeApp;