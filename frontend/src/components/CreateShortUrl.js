import React, { useState } from 'react';
import axios from 'axios';
import './CreateShortUrl.css';

const CreateShortUrl = ({ onBack }) => {
  const [url, setUrl] = useState('');
  const [shortUrl, setShortUrl] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [expiresInDays, setExpiresInDays] = useState('');

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError('');
    setShortUrl('');

    // Basic URL validation
    try {
      new URL(url);
    } catch {
      setError('Please enter a valid URL (including http:// or https://)');
      setLoading(false);
      return;
    }

    try {
      const requestData = { url };
      if (expiresInDays && parseInt(expiresInDays) > 0) {
        requestData.expires_in_days = parseInt(expiresInDays);
      }

      const response = await axios.post('http://localhost:8080/api/shorten', requestData);
      setShortUrl(response.data.short_url);
    } catch (err) {
      console.error('Error occurred:', err);
      
      if (err.response) {
        setError(err.response.data?.error || `Server error: ${err.response.status}`);
      } else if (err.request) {
        setError('Cannot connect to server. Make sure the backend is running on port 8080.');
      } else {
        setError('Something went wrong: ' + err.message);
      }
    } finally {
      setLoading(false);
    }
  };

  const copyToClipboard = async () => {
    try {
      await navigator.clipboard.writeText(shortUrl);
      // You could add a toast notification here
      alert('Short URL copied to clipboard!');
    } catch (err) {
      console.error('Failed to copy: ', err);
    }
  };

  const handleReset = () => {
    setUrl('');
    setShortUrl('');
    setError('');
    setExpiresInDays('');
  };

  return (
    <div className="create-short-url-container">
      <div className="page-header">
        <button className="back-button" onClick={onBack}>
          ← Back to Dashboard
        </button>
        <h1>Create Short URL</h1>
        <p>Transform your long URLs into short, shareable links</p>
      </div>

      <div className="url-shortener-card">
        <form onSubmit={handleSubmit} className="url-form">
          <div className="form-section">
            <label htmlFor="url">Original URL</label>
            <div className="input-group">
              <input
                type="url"
                id="url"
                value={url}
                onChange={(e) => setUrl(e.target.value)}
                placeholder="https://example.com/very-long-url..."
                required
                className="url-input"
              />
            </div>
          </div>

          <div className="form-section">
            <label htmlFor="expires">Expiration (Optional)</label>
            <div className="input-group">
              <input
                type="number"
                id="expires"
                value={expiresInDays}
                onChange={(e) => setExpiresInDays(e.target.value)}
                placeholder="Days until expiration (leave empty for no expiration)"
                min="1"
                max="365"
                className="expires-input"
              />
              <span className="input-suffix">days</span>
            </div>
          </div>

          {error && <div className="error-message">{error}</div>}

          <div className="form-actions">
            <button type="submit" disabled={loading || !url} className="shorten-button">
              {loading ? (
                <>
                  <span className="spinner"></span>
                  Creating...
                </>
              ) : (
                'Create Short URL'
              )}
            </button>
            {(url || shortUrl) && (
              <button type="button" onClick={handleReset} className="reset-button">
                Reset
              </button>
            )}
          </div>
        </form>

        {shortUrl && (
          <div className="result-section">
            <div className="result-header">
              <h3>Your Short URL is Ready!</h3>
              <p>Share this link anywhere</p>
            </div>
            <div className="result-card">
              <div className="url-display">
                <input
                  type="text"
                  value={shortUrl}
                  readOnly
                  className="short-url-input"
                />
                <button onClick={copyToClipboard} className="copy-button">
                  <svg width="20" height="20" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                    <path d="M8 4V16C8 17.1046 8.89543 18 10 18H18C19.1046 18 20 17.1046 20 16V7.24264C20 6.44699 19.6839 5.68393 19.1213 5.12132L16.8787 2.87868C16.3161 2.31607 15.553 2 14.7574 2H10C8.89543 2 8 2.89543 8 4Z" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"/>
                    <path d="M16 18V20C16 21.1046 15.1046 22 14 22H6C4.89543 22 4 21.1046 4 20V9C4 7.89543 4.89543 7 6 7H8" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"/>
                  </svg>
                  Copy
                </button>
              </div>
              <div className="url-info">
                <div className="info-item">
                  <span className="info-label">Original:</span>
                  <span className="info-value">{url}</span>
                </div>
                {expiresInDays && (
                  <div className="info-item">
                    <span className="info-label">Expires in:</span>
                    <span className="info-value">{expiresInDays} days</span>
                  </div>
                )}
              </div>
            </div>
          </div>
        )}
      </div>
    </div>
  );
};

export default CreateShortUrl;
