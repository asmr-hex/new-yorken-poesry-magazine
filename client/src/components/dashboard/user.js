import React, { Component } from 'react';
import {connect} from 'react-redux'


class userMgmt extends Component {
  render() {
    return (
      <div>
        usr mgmt
      </div>
    )
  }
}

export const UserMgmt = connect(() => ({}), {})(userMgmt)
