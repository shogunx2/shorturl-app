import React, { useState, useEffect } from 'react';
import './App.css';
import Auth from './components/Auth';
import Dashboard from './components/Dashboard';
import CreateShortUrl from './components/CreateShortUrl';
import VisitShortUrl from './components/VisitShortUrl';

function App() {
  const [currentPage, setCurrentPage] = useState('auth');
  const [user, setUser] = useState(null);

  // Check for existing authentication on app load
  useEffect(() => {
    const token = localStorage.getItem('authToken');
    const userId = localStorage.getItem('userId');
    
    if (token && userId) {
      setUser(userId);
      setCurrentPage('dashboard');
    }
  }, []);

  const handleLogin = (userId) => {
    setUser(userId);
    setCurrentPage('dashboard');
  };

  const handleLogout = () => {
    localStorage.removeItem('authToken');
    localStorage.removeItem('userId');
    setUser(null);
    setCurrentPage('auth');
  };

  const handleNavigate = (page) => {
    setCurrentPage(page);
  };

  const renderCurrentPage = () => {
    switch (currentPage) {
      case 'auth':
        return <Auth onLogin={handleLogin} />;
      
      case 'dashboard':
        return (
          <Dashboard 
            user={user} 
            onLogout={handleLogout} 
            onNavigate={handleNavigate} 
          />
        );
      
      case 'create':
        return (
          <CreateShortUrl 
            onBack={() => handleNavigate('dashboard')} 
          />
        );
      
      case 'visit':
        return (
          <VisitShortUrl 
            onBack={() => handleNavigate('dashboard')} 
          />
        );
      
      default:
        return <Auth onLogin={handleLogin} />;
    }
  };

  return (
    <div className="App">
      {renderCurrentPage()}
    </div>
  );
}

export default App;
