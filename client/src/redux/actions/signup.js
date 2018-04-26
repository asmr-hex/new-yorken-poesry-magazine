import fetch from 'cross-fetch'
import {handleErrors} from './login'


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
    .then(handleErrors)
    .then(
      user => dispatch(signupSuccessful(user)),
      err => console.log('login error: ', err),
    )
}


export const SIGNUP_SUCCESSFUL = 'SIGNUP_SUCCESSFUL'
export const signupSuccessful = user => dispatch =>
  dispatch({
    payload: user,
    type: SIGNUP_SUCCESSFUL,
  })


export const SIGNUP_FAILED = 'SIGNUP_FAILED'
export const signupFailed = () => dispatch =>
  dispatch({
    type: SIGNUP_FAILED,
  })
