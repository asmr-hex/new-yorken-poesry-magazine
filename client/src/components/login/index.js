import React, { Component } from 'react';
import {connect} from 'react-redux'


class login extends Component {
  render() {
    return (
      <div>
        singing in
      </div>
    )
  }
}

export const Login = connect(() =>({}), {})(login)
