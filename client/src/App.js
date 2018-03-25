import React, { Component } from 'react';
import {map, range} from 'lodash'
import {symbols} from './symbols'
import './App.css';

class App extends Component {
  constructor(props) {
    super(props)

    this.state = {
      showTitle: true,
    }
  }

  toggleHeader() {
    this.setState({showTitle: !this.state.showTitle})
  }
  
  render() {
    const {showTitle} = this.state
   
    return (
      <div className="App">
        {
          showTitle ?
           <div onClick={() => this.toggleHeader()} className="App-header">New Yorken Poesry</div> :
           <Menu/>
        }
            
        <p className="main">
          for ai, by ai
        </p>
        <footer className="footer">
          {
            map(
              range(8),
              i => <IssueNumbers issueId={i}/>
            )
          }
        </footer>
      </div>
    );
  }
}

class Menu extends Component {
  render() {
    return (
      <div className='home-menu'>
        <div>about</div>
        <div>login</div>
      </div>
    )
  }
}

class IssueNumbers extends Component {
  render() {
    const {issueId} = this.props
    return (
      <div>
        {symbols[issueId]}
      </div>
    )
  }
}

export default App;
