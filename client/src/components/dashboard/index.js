import React, { Component } from 'react';
import {connect} from 'react-redux'
import {Link} from 'react-router-dom'
import {UserMgmt} from './user'
import {PoetMgmt} from './poet'


class dashboard extends Component {
  constructor(props) {
    super(props)

    this.state = {
      userMgmt: true,
    }
  }
  
  manageUserPage = ok => e =>
    this.setState({userMgmt: ok})

  render() {
    return (
      <div>
        <div className='about-tabs'>
          <div className='header-menu-item' onClick={this.manageUserPage(true)}>a</div>
          <div className='header-menu-item' onClick={this.manageUserPage(false)}>p</div>
          <Link to='/' className='header-menu-item'>?</Link>
        </div>
        <div>
          {
            this.state.userMgmt ?
                <UserMgmt/>
              : <PoetMgmt/>
              
          }
        </div>
      </div>
    )
  }
}

export const Dashboard = connect(() =>({}), {})(dashboard)
