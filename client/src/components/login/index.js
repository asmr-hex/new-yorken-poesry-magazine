import React, { Component } from 'react';
import {connect} from 'react-redux'
import {LoginForm, SignupForm} from './forms'
import {requestLogin} from '../../redux/actions/login'
import {requestSignup} from '../../redux/actions/signup'


class login extends Component {
  constructor(props) {
    super(props)
    
    this.state = {
      loginFormShown: true,
    }
  }

  // upon login, we want to redirect the route to the dashboard
  redirectUponLogin = () => {
    const {history} = this.props

    history.push('/profile')
  }
  
  login = values => {
    const {password, username} = values
    
    // TODO (cw|4.25.2018) similar validation should happen here
    // to what i am referring to in the comment below.
    
    this.props.requestLogin(
      {
        username: username || '',
        password: password || '',
      },
      this.redirectUponLogin,
    )
  }

  signup = values => {
    const {email, username, password} = values

    // TODO (cw|4.25.2018) invoke some validation.
    // Question: where should validation for outbound request
    // data reside? Certainly not at the component level, but
    // perhaps in the action creators.
    // ps. there *is* server-side validation, but we want to
    // include some before it gets put on the wire for ease-of-use.

    this.props.requestSignup(
      {
        email: email || '',
        username: username || '',
        password: password || '',      
      },
      this.redirectUponLogin,
    )
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
  requestSignup,
}

export const Login = connect(() =>({}), actions)(login)
