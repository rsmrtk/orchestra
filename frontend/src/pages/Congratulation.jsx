import { useNavigate } from 'react-router-dom';
import './Auth.css';

function Congratulation() {
  const navigate = useNavigate();

  return (
    <div className="auth-container">
      <div className="auth-card congratulation">
        <h1 className="congratulation-title">Congratulation!</h1>
        <p className="congratulation-message">
          You have successfully logged in to the Orchestra system.
        </p>
        <div className="congratulation-actions">
          <button
            onClick={() => navigate('/login')}
            className="submit-btn"
          >
            Back to Login
          </button>
        </div>
      </div>
    </div>
  );
}

export default Congratulation;
