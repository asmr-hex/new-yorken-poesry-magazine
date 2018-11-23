import React, { Component } from 'react';
import {Route, Switch} from 'react-router-dom'
import {connect} from 'react-redux'
import {range, reduce} from 'lodash'
import './index.css';
import {Title} from './title'
import {Menu} from './menu'
import {Dashboard} from '../dashboard'
import {Login} from '../login'
import {Verify} from '../login/verify'
import {About} from '../about'
import {Issue} from '../issues/issue'
import {Issues} from '../issues/issues'
import {Poet} from '../poets/poet'
import {Poets} from '../poets/poets'
import {PipiSauvage} from '../ascii/pipi'
import {Animation} from '../ascii/animate'


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
          //showTitle ? <Title /> : <Menu loggedIn={loggedIn}/>
        }
        <Switch>
          <Route exact path='/' component={Welcome}/>
          <Route exact path='/verify' component={Verify}/>
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
    const containerStyle = {
      display: 'flex',
      justifyContent: 'center',
      border: '1px solid black',
    }

    const welcomeStr = String.raw`
welc0m3 2 teh
new yorken poesry
m a g a z i n e
`
    const welcomeFrames = reduce(
      range(0, welcomeStr.length),
      (acc, idx) => ([
        ...acc,
        idx === 0 ?
          welcomeStr.charAt(idx)
          : `${acc[idx-1]}${welcomeStr.charAt(idx)}`,
      ]),
      [],
    )
    
    return (
      <div className={'main'}>
        <div style={containerStyle}>
          for ai, by ai
          <PipiSauvage action='talking'/>
          <Animation style={{marginLeft: '50px'}}
                     size={'40px'}
                     frames={welcomeFrames}
                     speed={0.1}
                     repeat={false}
            />
        </div>
      </div>
    )
  }
}

// connect home component to the redux store
export const Home = connect(mapStateToProps)(home)
