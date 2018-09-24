import React from 'react';
import { Field, reduxForm } from 'redux-form'
import TextField from 'material-ui/TextField'


const renderTextField = ({input, label, meta: {touched, error}, ...custom}) => (
  <TextField
    hintText={label}
    floatingLabelText={label}
    hintStyle={{color: '#222', fontSize: '1.5em'}}
    inputStyle={{color: '#222', fontSize: '1.5em'}}
    underlineStyle={{borderColor: '#222'}}
    underlineFocusStyle={{borderColor: '#222'}}
    errorText={touched && error}
    {...input}
    {...custom}
    />
)

const loginForm = props => {
  const { handleSubmit } = props
  
  return (
    <form onSubmit={handleSubmit} name='login-form'>
      <div>
        <Field name="username" component={renderTextField} type="text" placeholder="username"/>
      </div>
      <div>
        <Field name="password" component={renderTextField} type="text" placeholder="password"/>
      </div>
      <button className='signin-button' type="submit">login</button>
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
        <Field name="email" component={renderTextField} type="text" placeholder="email address"/>
      </div>
      <div>
        <Field name="username" component={renderTextField} type="text" placeholder="username"/>
      </div>
      <div>
        <Field name="password" component={renderTextField} type="text" placeholder="password"/>
      </div>
      <button className='signin-button' type="submit">signup</button>
    </form> 
  )
}

export const SignupForm = reduxForm({
  form: 'signup'
})(signupForm)
