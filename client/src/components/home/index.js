import React, { Component } from 'react';
import {Route, Switch} from 'react-router-dom'
import {connect} from 'react-redux'
import {map, range} from 'lodash'
import './index.css';
import {Title} from './title'
import {Menu} from './menu'
import {Login} from '../login'
import {About} from '../about'
import {Issue} from '../issues/issue'
import {Issues} from '../issues/issues'
import {Poet} from '../poets/poet'


const mapStateToProps = state => ({
  loggedIn: state.session.loggedIn,
  ui: state.ui,
})

// TODO (cw|4.26.2018) we should create a higher level component for the app so we can
// switch between the different pages (e.g. home, dashboard, issues, etc.).

class home extends Component {
  render() {
    const {showTitle} = this.props.ui
   
    return (
      <div className="App">
        {
          showTitle ? <Title /> : <Menu />
        }
        <Switch>
          <Route exact path='/' component={Welcome}/>
          <Route path='/about' component={About}/>
          <Route path='/login' component={Login}/>
          <Route path='/poet/:id' component={Poet}/>
          <Route path='/issues' component={Issues}/>
        </Switch>
        {
            // <footer className="footer">
            //     {
            //       map(
            //         range(8),
            //         i => <IssueNumbers issueId={i} key={i}/>,
            //       )
            //     }
            // </footer>
        }
      </div>
    );
  }
}

class Welcome extends Component {
  render() {
    return (
      <div className="main">
        for ai, by ai
        <Issue volume='latest'/>
      </div>
    )
  }
}

// connect home component to the redux store
export const Home = connect(mapStateToProps)(home)
