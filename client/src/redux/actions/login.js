import fetch from 'cross-fetch'
import {checkResponse} from './error'


export const LOGIN_REQUESTED = 'LOGIN_REQUESTED'
export const requestLogin = ({username, password}, redirectUponLogin) => dispatch => {
  const payload = {username, password}
  
  dispatch({
    payload,
    type: LOGIN_REQUESTED,
  })

  fetch(
    `/dashboard/login`,
    {
      method: 'POST',
      body: JSON.stringify(payload),
      headers: new Headers({'Content-Type': 'application/json'}),
      credentials: 'same-origin',
    })
    .then(checkResponse)
    .then(
      user => {
        dispatch(loginSuccessful(user))
        redirectUponLogin()
      },
      error => dispatch(loginFailed(error)),
    )
}


export const LOGIN_SUCCESSFUL = 'LOGIN_SUCCESSFUL'
export const loginSuccessful = user => dispatch =>
  dispatch({
    payload: user,
    type: LOGIN_SUCCESSFUL,
  })


export const LOGIN_FAILED = 'LOGIN_FAILED'
export const loginFailed = error => dispatch => {
  error.message = 'login failed: ' + error.message
  
  dispatch({
    error,
    type: LOGIN_FAILED,
  })  
}

