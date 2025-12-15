import { useState, useEffect } from 'react'
import './App.css'

function App() {
  const [message, setMessage] = useState('')
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState(null)

  const fetchCongratulation = async () => {
    setLoading(true)
    setError(null)
    try {
      const response = await fetch('http://localhost:8383/congratulation')
      if (!response.ok) {
        throw new Error('Failed to fetch')
      }
      const data = await response.json()
      setMessage(data.message)
    } catch (err) {
      setError(err.message)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchCongratulation()
  }, [])

  return (
    <div className="app">
      <h1>Orchestra Deploy Practice</h1>
      {loading && <p>Loading...</p>}
      {error && <p style={{ color: 'red' }}>Error: {error}</p>}
      {message && (
        <div className="congratulation">
          <h2>{message}</h2>
        </div>
      )}
      <button onClick={fetchCongratulation} disabled={loading}>
        Reload Message
      </button>
    </div>
  )
}

export default App
