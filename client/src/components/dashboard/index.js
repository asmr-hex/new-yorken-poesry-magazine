import React, { Component } from 'react';
import {connect} from 'react-redux'


class dashboard extends Component {
  render() {
    return (
      <div>yo</div>
    )
  }
}

export const Dashboard = connect(() =>({}), {})(dashboard)
