import React from 'react';
import { Field, reduxForm } from 'redux-form'


const loginForm = props => {
  const { handleSubmit } = props
  return (
    <form onSubmit={handleSubmit}>
      <div>
        <Field name="username" component="input" type="text" placeholder="username"/>
      </div>
      <div>
        <Field name="password" component="input" type="text" placeholder="password"/>
      </div>
      <button type="submit">login</button>
    </form> 
  )
}

export const LoginForm = reduxForm({
  form: 'login'
})(loginForm)


// Signup Form

const signupForm = props => {
  const { handleSubmit } = props
  return (
    <form onSubmit={handleSubmit}>
      <div>
        <Field name="email" component="input" type="text" placeholder="email address"/>
      </div>
      <div>
        <Field name="username" component="input" type="text" placeholder="username"/>
      </div>
      <div>
        <Field name="password" component="input" type="text" placeholder="password"/>
      </div>
      <button type="submit">signup</button>
    </form> 
  )
}

export const SignupForm = reduxForm({
  form: 'signup'
})(signupForm)
