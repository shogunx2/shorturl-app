import React from 'react';
import './Dashboard.css';

const Dashboard = ({ user, onLogout, onNavigate }) => {
  return (
    <div className="dashboard-container">
      <header className="dashboard-header">
        <div className="header-content">
          <div className="header-left">
            <h1>URL Shortener</h1>
            <p>Welcome back, {user}!</p>
          </div>
          <div className="header-right">
            <button className="logout-button" onClick={onLogout}>
              Sign Out
            </button>
          </div>
        </div>
      </header>

      <main className="dashboard-main">
        <div className="dashboard-grid">
          <div className="feature-card" onClick={() => onNavigate('create')}>
            <div className="card-icon">
              <svg width="48" height="48" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                <path d="M10 6H6C4.89543 6 4 6.89543 4 8V18C4 19.1046 4.89543 20 6 20H16C17.1046 20 18 19.1046 18 18V8C18 6.89543 17.1046 6 16 6H12" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"/>
                <path d="M8 6C8 4.89543 8.89543 4 10 4H12C13.1046 4 14 4.89543 14 6C14 7.10457 13.1046 8 12 8H10C8.89543 8 8 7.10457 8 6Z" stroke="currentColor" strokeWidth="2"/>
                <path d="M12 12H8" stroke="currentColor" strokeWidth="2" strokeLinecap="round"/>
                <path d="M12 16H8" stroke="currentColor" strokeWidth="2" strokeLinecap="round"/>
              </svg>
            </div>
            <div className="card-content">
              <h3>Create Short URL</h3>
              <p>Transform long URLs into short, shareable links that are easy to remember and share.</p>
              <div className="card-action">
                <span>Create Now →</span>
              </div>
            </div>
          </div>

          <div className="feature-card" onClick={() => onNavigate('visit')}>
            <div className="card-icon">
              <svg width="48" height="48" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                <path d="M10 13a5 5 0 0 0 7.54.54l3-3a5 5 0 0 0-7.07-7.07l-1.72 1.71" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"/>
                <path d="M14 11a5 5 0 0 0-7.54-.54l-3 3a5 5 0 0 0 7.07 7.07l1.71-1.71" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"/>
              </svg>
            </div>
            <div className="card-content">
              <h3>Visit Short URL</h3>
              <p>Enter a short URL to be redirected to the original destination in a new browser tab.</p>
              <div className="card-action">
                <span>Visit Now →</span>
              </div>
            </div>
          </div>
        </div>

        <div className="stats-section">
          <div className="stats-card">
            <h4>Quick Stats</h4>
            <div className="stats-grid">
              <div className="stat-item">
                <span className="stat-number">∞</span>
                <span className="stat-label">URLs Shortened</span>
              </div>
              <div className="stat-item">
                <span className="stat-number">∞</span>
                <span className="stat-label">Clicks Tracked</span>
              </div>
              <div className="stat-item">
                <span className="stat-number">24/7</span>
                <span className="stat-label">Availability</span>
              </div>
            </div>
          </div>
        </div>
      </main>
    </div>
  );
};

export default Dashboard;
