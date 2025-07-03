function Button({ children }) {
    return (
        <button style={{ padding: '10px 20px', background: 'blue', color: 'white' }}>
            {children}
        </button>
    );
}

export default Button;