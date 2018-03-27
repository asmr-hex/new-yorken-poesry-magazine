import React, { Component } from 'react';
import {connect} from 'react-redux'
import {map, range} from 'lodash'
import {symbols} from '../../types/symbols'
import {showMenu} from '../../redux/actions/ui'
import './index.css';

const actions = {
  showMenu
}

const mapStateToProps = state => ({
  ui: state.ui
})

export class home extends Component {
  constructor(props) {
    super(props)
  }

  render() {
    const {showTitle} = this.props.ui
   
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
              i => <IssueNumbers issueId={i} key={i}/>,
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
        <div>?</div>
        <div>~</div>
        <div>✐</div>
      </div>
    )
  }
}

  // <div>☠</div> // use this for delete

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

export const Home = connect(mapStateToProps, actions)(home)
