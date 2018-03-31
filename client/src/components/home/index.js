import React, { Component } from 'react';
import {connect} from 'react-redux'
import {map, range} from 'lodash'
import {symbols} from '../../types/symbols'
import './index.css';
import {Title} from './title'
import {Menu} from './menu'


const mapStateToProps = state => ({
  ui: state.ui
})

export class home extends Component {
  render() {
    const {showTitle} = this.props.ui
   
    return (
      <div className="App">
        {
          showTitle ?
             <Title /> :
             <Menu />
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

// connect home component to the redux store
export const Home = connect(mapStateToProps)(home)
