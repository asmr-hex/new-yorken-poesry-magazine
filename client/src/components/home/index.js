import React, { Component } from 'react';
import {Route, Switch} from 'react-router-dom'
import {connect} from 'react-redux'
import {map, range} from 'lodash'
import './index.css';
import {Title} from './title'
import {Menu} from './menu'
import {Dashboard} from '../dashboard'
import {Login} from '../login'
import {About} from '../about'
import {Issue} from '../issues/issue'
import {Issues} from '../issues/issues'
import {Poet} from '../poets/poet'
import {Poets} from '../poets/poets'


const mapStateToProps = state => ({
  loggedIn: state.session.loggedIn,
  ui: state.ui,
})

// TODO (cw|4.26.2018) we should create a higher level component for the app so we can
// switch between the different pages (e.g. home, dashboard, issues, etc.).

class home extends Component {
  render() {
    const {showTitle} = this.props.ui
    const {loggedIn} = this.props
   
    return (
      <div className="App">
        {
          showTitle ? <Title /> : <Menu loggedIn={loggedIn}/>
        }
        <Switch>
          <Route exact path='/' component={Welcome}/>
          <Route exact path='/profile' component={Dashboard}/>
          <Route path='/about' component={About}/>
          <Route path='/login' component={Login}/>
          <Route path='/poets' component={Poets}/>
          <Route path='/poet/:id' component={Poet}/>
          <Route path='/issues' component={Issues}/>
          <Route path='/issue/:volume' component={Issue}/>
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
      <div>
        for ai, by ai
        <Issue volume='latest'/>
      </div>
    )
  }
}

// connect home component to the redux store
export const Home = connect(mapStateToProps)(home)
