import fetch from 'cross-fetch'
import {checkResponse} from './error'


export const SIGNUP_REQUESTED = 'SIGNUP_REQUESTED'
export const requestSignup = ({email, username, password}) => dispatch => {
  const payload = {email, username, password}

  dispatch({
    payload,
    type: SIGNUP_REQUESTED,
  })

  fetch(
    `/dashboard/signup`,
    {
      method: 'POST',
      body: JSON.stringify(payload),
      headers: new Headers({'Content-Type': 'application/json'}),
    })
    .then(checkResponse)
    .then(
      user => dispatch(signupSuccessful(user)),
      error => dispatch(signupFailed(error)),
    )
}


export const SIGNUP_SUCCESSFUL = 'SIGNUP_SUCCESSFUL'
export const signupSuccessful = user => dispatch =>
  dispatch({
    payload: user,
    type: SIGNUP_SUCCESSFUL,
  })


export const SIGNUP_FAILED = 'SIGNUP_FAILED'
export const signupFailed = error => dispatch => {
  error.message = 'signup failed: ' + error.message

  dispatch({
    error,
    type: SIGNUP_FAILED,
  })  
}
