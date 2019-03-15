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
        <div className='profile-body'>
          <UserMgmt user={user}/>
          <PoetMgmt/>
        </div>
      </div>
    )
  }
}

export const Dashboard = connect(mapStateToProps, actions)(dashboard)
