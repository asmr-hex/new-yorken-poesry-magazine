import React, { Component } from 'react';
import {connect} from 'react-redux'
import {get} from 'lodash'
import {LoginForm, SignupForm} from './forms'
import {requestLogin} from '../../redux/actions/login'
import {requestSignup} from '../../redux/actions/signup'
import {resetErrorMsg} from '../../redux/actions/error'
import './index.css'


const mapStateToProps = (state, ownProps) => ({
  errors: get(state, `error`, ''),
  pendingVerification: get(state, `session.pendingVerification`, false),
})

const actions = {
  requestLogin,
  requestSignup,
  resetErrorMsg,
}

class login extends Component {
  constructor(props) {
    super(props)
    
    this.state = {
      loginFormShown: true,
    }
  }

  componentDidMount() {
    // if we are reloading this page, reset the error message
    this.props.resetErrorMsg()
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

    this.props.requestSignup({
        email: email || '',
        username: username || '',
        password: password || '',      
      })
  }

  showLoginForm = () => {
    this.setState({
      loginFormShown: true,
    })

    this.props.resetErrorMsg()
  }

  showSignupForm = () => {
    this.setState({
      loginFormShown: false,
    })

    this.props.resetErrorMsg()
  }

  // upon login, we want to redirect the route to the dashboard
  redirectUponLogin = () => {
    const {history} = this.props

    history.push('/profile')
  }

  render() {
    const {loginFormShown} = this.state
    const {pendingVerification} = this.props

    const signupLoginComponent =
      <div className='login-page'>
        <div className='login-container'>
           {
            loginFormShown ?
              <LoginForm onSubmit={this.login}/>
              : <SignupForm onSubmit={this.signup}/> 
           } 
          <div className='login-signup-choose-box-thing'>
            <span className='login-choice-button' onClick={this.showLoginForm}>login</span> / 
            <span className='signup-choice-button' onClick={this.showSignupForm}> signup</span>
          </div>
          <div className='login-signup-error-message'>
                  {this.props.errors}
          </div>
        </div>
      </div>

    const pendingVerificationComponent =
      <div>YO PENDING VERIFICATION FRIEMDB</div>
    
    return (
      <div className='main'>
        {
          pendingVerification ? pendingVerificationComponent : signupLoginComponent
        }
      </div>
    )
  }
}

export const Login = connect(mapStateToProps, actions)(login)
