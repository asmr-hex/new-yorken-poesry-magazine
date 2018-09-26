import React, { Component } from 'react';
import {connect} from 'react-redux'
import {get} from 'lodash'
import {UserMgmt} from './user'
import {PoetMgmt} from './poet'
import './index.css'


const mapStateToProps = (state, ownProps) => ({
  user: get(state, `session.user`, {}),
})

const actions = {}

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
    const {
      user,
    } = this.props
    
    return (
      <div className='main'>
        <div className='profile-container'>
          <div className='profile-tabs'>
            <div
              className={this.state.userMgmt ? 'profile-tab-selected' : 'profile-tab'}
              onClick={this.manageUserPage(true)}>
              profile
            </div>
            <div
              className={this.state.userMgmt ? 'profile-tab' : 'profile-tab-selected'}
              onClick={this.manageUserPage(false)}>
              poets
            </div>
          </div>
          <div className='profile-body'>
            {
              this.state.userMgmt ?
              <UserMgmt user={user}/>
                : <PoetMgmt/>
                
              }
          </div>
        </div>
      </div>
    )
  }
}

export const Dashboard = connect(mapStateToProps, actions)(dashboard)
