import React, { Component } from 'react';
import {connect} from 'react-redux'
import {LoginForm, SignupForm} from './forms'


class login extends Component {
  constructor(props) {
    super(props)
    
    this.state = {
      loginFormShown: true,
    }
  }

  submit = values => {
    console.log(values)
  }

  showLoginForm() {
    this.setState({
      loginFormShown: true,
    })
  }

  showSignupForm() {
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
              <LoginForm onSubmit={this.submit}/>
            : <SignupForm onSubmit={this.submit}/>
            
        }
        <div>
          <span onClick={() => this.showLoginForm()}>login</span>/
          <span onClick={() => this.showSignupForm()}>signup</span>
        </div>
      </div>
      
    )
  }
}

export const Login = connect(() =>({}), {})(login)
