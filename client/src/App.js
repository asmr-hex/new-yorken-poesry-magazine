import React, { Component } from 'react';
import './App.css';

class App extends Component {
  render() {
    return (
      <div className="App">
        <header className="App-header">
          New Yorken Poesry Magazine
        </header>
        <p className="main">
          To get started, edit <code>src/App.js</code> and save to reload.
        </p>
        <footer className="footer">
          {<IssueNumbers/>}
        </footer>
      </div>
    );
  }
}

class IssueNumbers extends Component {
  render() {
    return (
      <div>
        1
      </div>
    )
  }
}

export default App;
