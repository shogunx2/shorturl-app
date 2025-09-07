import React from 'react';
import './App.css';
import UrlShortener from './components/UrlShortener';

function App() {
  return (
    <div className="App">
      <header className="App-header">
        <h1>URL Shortener</h1>
        <p>Create short, memorable links for any URL</p>
      </header>
      <main className="App-main">
        <UrlShortener />
      </main>
    </div>
  );
}

export default App;
