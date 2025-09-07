import React, { useState } from 'react';
import axios from 'axios';
import './UrlShortener.css';

const UrlShortener = () => {
  const [url, setUrl] = useState('');
  const [shortUrl, setShortUrl] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError('');
    setShortUrl('');

    try {
      console.log('Sending request to:', '/api/shorten');
      console.log('URL to shorten:', url);
      
      const response = await axios.post('http://localhost:8080/api/shorten', {
        url: url
      });

      console.log('Response received:', response.data);
      setShortUrl(response.data.short_url);
    } catch (err) {
      console.error('Error occurred:', err);
      console.error('Error response:', err.response);
      
      if (err.response) {
        // Server responded with error status
        setError(err.response.data?.error || `Server error: ${err.response.status}`);
      } else if (err.request) {
        // Request was made but no response received
        setError('Cannot connect to server. Make sure the backend is running on port 8080.');
      } else {
        // Something else happened
        setError('Something went wrong: ' + err.message);
      }
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="url-shortener">
      <form onSubmit={handleSubmit} className="url-form">
        <div className="input-group">
          <input
            type="url"
            value={url}
            onChange={(e) => setUrl(e.target.value)}
            placeholder="Enter your long URL here..."
            required
            className="url-input"
          />
          <button 
            type="submit" 
            disabled={loading}
            className="submit-btn"
            onClick={(e) => {
              console.log('Button clicked!');
              console.log('Current URL value:', url);
              console.log('Loading state:', loading);
            }}
          >
            {loading ? 'Shortening...' : 'Shorten URL'}
          </button>
        </div>
      </form>

      {error && (
        <div className="error-message">
          {error}
        </div>
      )}

      {shortUrl && (
        <div className="result">
          <h3>Your short URL:</h3>
          <div className="short-url-container">
            <input
              type="text"
              value={shortUrl}
              readOnly
              className="short-url-input"
            />
          </div>
        </div>
      )}
    </div>
  );
};

export default UrlShortener;
