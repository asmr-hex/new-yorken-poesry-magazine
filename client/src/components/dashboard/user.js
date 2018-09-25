import React, { Component } from 'react';
import {connect} from 'react-redux'


class userMgmt extends Component {
  render() {
    const {
      user,
    } = this.props
    
    return (
      <div className='profile-user-details-container'>
        <div className='profile-user-details'>
          <span>{user.username}</span>
          <span>{user.email}</span>
        </div>
      </div>
    )
  }
}

export const UserMgmt = connect(() => ({}), {})(userMgmt)
