import React, { Component } from 'react';
import {connect} from 'react-redux'
import {LoginForm, SignupForm} from './forms'
import {requestLogin} from '../../redux/actions/login'


class login extends Component {
  constructor(props) {
    super(props)
    
    this.state = {
      loginFormShown: true,
    }
  }

  login = values => {
    const {password, username} = values

    this.props.requestLogin({
      username: username || '',
      password: password || '',
    })
  }

  signup = values => {
    console.log("signing up")
  }

  showLoginForm = () => {
    this.setState({
      loginFormShown: true,
    })
  }

  showSignupForm = () => {
    this.setState({
      loginFormShown: false,
    })
  }
  
  render() {
    const {loginFormShown} = this.state

    return (
      <div>
        {
          loginFormShown ?
              <LoginForm onSubmit={this.login}/>
            : <SignupForm onSubmit={this.signup}/>
            
        }
        <div>
          <span onClick={this.showLoginForm}>login</span>/
          <span onClick={this.showSignupForm}>signup</span>
        </div>
      </div>
      
    )
  }
}

const actions = {
  requestLogin,
}

export const Login = connect(() =>({}), actions)(login)
