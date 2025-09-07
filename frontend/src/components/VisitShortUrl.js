import React, { useState } from 'react';
import './VisitShortUrl.css';

const VisitShortUrl = ({ onBack }) => {
  const [shortUrl, setShortUrl] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  const extractShortCode = (url) => {
    // Extract short code from various URL formats
    const patterns = [
      // Full URL with domain: http://localhost:8080/abc123
      /^https?:\/\/[^\/]+\/(.+)$/,
      // Just the short code: abc123
      /^([a-zA-Z0-9]+)$/,
      // URL path: /abc123
      /^\/(.+)$/
    ];

    for (const pattern of patterns) {
      const match = url.match(pattern);
      if (match && match[1]) {
        return match[1];
      }
    }
    return url; // Return as-is if no pattern matches
  };

  const handleSubmit = async (e) => {
    console.log('Handle submit called');
    e.preventDefault();
    setLoading(true);
    setError('');

    if (!shortUrl.trim()) {
      setError('Please enter a short URL or code');
      setLoading(false);
      return;
    }
    console.log('Short URL input:', shortUrl);

    try {
      const shortCode = extractShortCode(shortUrl.trim());
      
      // Make a GET request to the backend to get the original URL
      console.log('Visiting short code:', shortCode);
      const response = await fetch(`http://localhost:8080/${shortCode}`, {
        method: 'GET',
        redirect: 'manual' // Don't follow redirects automatically
      });

      if (response.type === 'opaqueredirect' || response.status === 302 || response.status === 301) {
        // For redirects, we need to construct the redirect URL
        // Since we can't get the Location header in the browser due to CORS,
        // we'll open the short URL directly
        console.log('Redirecting to original URL');
        const fullShortUrl = shortCode.startsWith('http') ? shortCode : `http://localhost:8080/${shortCode}`;
        window.open(fullShortUrl, '_blank');
      } else if (response.ok) {
        // If we get a 200 response, try to extract the URL from the response
        console.log('Response OK, extracting original URL');
        const data = await response.json();
        if (data.original_url) {
          window.open(data.original_url, '_blank');
        } else {
          throw new Error('Invalid response format');
        }
      } else if (response.status === 404) {
        setError('Short URL not found. Please check the URL and try again.');
      } else if (response.status === 410) {
        setError('This short URL has expired.');
      } else {
        setError(`Error: ${response.status} - ${response.statusText}`);
      }
    } catch (err) {
      console.error('Error occurred:', err);
      
      if (err.name === 'TypeError' && err.message.includes('fetch')) {
        setError('Cannot connect to server. Make sure the backend is running on port 8080.');
      } else {
        setError('Something went wrong: ' + err.message);
      }
    } finally {
      setLoading(false);
    }
  };

  const handleInputChange = (e) => {
    setShortUrl(e.target.value);
    setError('');
  };

  return (
    <div className="visit-short-url-container">
      <div className="page-header">
        <button className="back-button" onClick={onBack}>
          ← Back to Dashboard
        </button>
        <h1>Visit Short URL</h1>
        <p>Enter a short URL to be redirected to the original destination</p>
      </div>

      <div className="visit-url-card">
        <form onSubmit={handleSubmit} className="visit-form">
          <div className="form-section">
            <label htmlFor="shortUrl">Short URL or Code</label>
            <div className="input-group">
              <input
                type="text"
                id="shortUrl"
                value={shortUrl}
                onChange={handleInputChange}
                placeholder="Enter short URL (e.g., abc123 or http://localhost:8080/abc123)"
                required
                className="short-url-input"
              />
            </div>
            <div className="input-help">
              You can enter just the short code (e.g., "abc123") or the full URL
            </div>
          </div>

          {error && <div className="error-message">{error}</div>}

          <div className="form-actions">
            <button type="submit" disabled={loading || !shortUrl.trim()} className="visit-button">
              {loading ? (
                <>
                  <span className="spinner"></span>
                  Opening...
                </>
              ) : (
                <>
                  <svg width="20" height="20" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                    <path d="M18 13v6a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h6" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"/>
                    <polyline points="15,3 21,3 21,9" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"/>
                    <line x1="10" y1="14" x2="21" y2="3" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"/>
                  </svg>
                  Visit URL (New Tab)
                </>
              )}
            </button>
          </div>
        </form>

        <div className="info-section">
          <div className="info-card">
            <h3>How it works</h3>
            <div className="info-steps">
              <div className="info-step">
                <div className="step-number">1</div>
                <div className="step-content">
                  <h4>Enter Short URL</h4>
                  <p>Paste the short URL or just the short code</p>
                </div>
              </div>
              <div className="info-step">
                <div className="step-number">2</div>
                <div className="step-content">
                  <h4>Click Visit</h4>
                  <p>We'll fetch the original destination</p>
                </div>
              </div>
              <div className="info-step">
                <div className="step-number">3</div>
                <div className="step-content">
                  <h4>New Tab Opens</h4>
                  <p>The original URL opens in a new browser tab</p>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default VisitShortUrl;
