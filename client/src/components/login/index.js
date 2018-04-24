import React, { Component } from 'react';
import {connect} from 'react-redux'
import { Field, reduxForm } from 'redux-form'


class login extends Component {
  submit = values => {
    console.log(values)
  }

  render() {
    return (
      <div>
        <LoginForm onSubmit={this.submit}/>
      </div>
    )
  }
}

const loginForm = props => {
  const { handleSubmit } = props
  return (
    <form onSubmit={handleSubmit}>
        <div>
            <label htmlFor="firstName">First Name</label>
            <Field name="firstName" component="input" type="text" />
        </div>
        <button type="submit">Submit</button>
    </form> 
  )
}

export const LoginForm = reduxForm({
  form: 'login'
})(loginForm)

export const Login = connect(() =>({}), {})(login)
