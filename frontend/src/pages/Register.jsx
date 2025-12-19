import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import './Auth.css';

const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8282';

function Register() {
  const [firstName, setFirstName] = useState('');
  const [lastName, setLastName] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');
    setLoading(true);

    try {
      const response = await fetch(`${API_URL}/auth/register`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          first_name: firstName,
          last_name: lastName,
        }),
      });

      const data = await response.json();

      if (data.success) {
        // Успішна реєстрація - перенаправляємо на логін
        alert(data.message);
        navigate('/login');
      } else {
        // Клієнт вже існує - показуємо повідомлення і можливо перенаправляємо
        setError(data.message);
        if (data.redirect === '/login') {
          setTimeout(() => navigate('/login'), 2000);
        }
      }
    } catch (err) {
      setError('Failed to connect to server. Please try again.');
      console.error('Register error:', err);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="auth-container">
      <div className="auth-card">
        <h1>Register</h1>
        <form onSubmit={handleSubmit}>
          <div className="form-group">
            <label htmlFor="firstName">First Name</label>
            <input
              type="text"
              id="firstName"
              value={firstName}
              onChange={(e) => setFirstName(e.target.value)}
              placeholder="Enter your first name"
              required
            />
          </div>
          <div className="form-group">
            <label htmlFor="lastName">Last Name</label>
            <input
              type="text"
              id="lastName"
              value={lastName}
              onChange={(e) => setLastName(e.target.value)}
              placeholder="Enter your last name"
              required
            />
          </div>
          {error && <div className="error-message">{error}</div>}
          <button type="submit" className="submit-btn" disabled={loading}>
            {loading ? 'Loading...' : 'Register'}
          </button>
        </form>
        <p className="auth-link">
          Already have an account? <a href="/login">Login here</a>
        </p>
      </div>
    </div>
  );
}

export default Register;
