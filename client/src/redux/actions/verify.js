import fetch from 'cross-fetch'
import {checkResponse} from './error'


export const VERIFY_REQUESTED = 'VERIFY_REQUESTED'
export const requestVerification = ({email, token}, redirectUponLogin) => dispatch => {
  const payload = {email, token}

  dispatch({
    payload,
    type: VERIFY_REQUESTED,
  })

  fetch(
    `dashboard/verify?token=${token}&email=${email}`,
    {
      method: 'GET',
    })
    .then(checkResponse)
    .then(
      user => {
        dispatch(verifySuccessful(user))
        redirectUponLogin()
      },
      error => dispatch(verifyFailed(error))
    )
}


export const VERIFY_SUCCESSFUL = 'VERIFY_SUCCESSFUL'
export const verifySuccessful = user => dispatch =>
  dispatch({
    payload: user,
    type: VERIFY_SUCCESSFUL,
  })


export const VERIFY_FAILED = 'VERIFY_FAILED'
export const verifyFailed = error => dispatch => {
  error.text()
    .then(
      msg => {
        error.message = msg
        dispatch({
          error,
          type: VERIFY_FAILED,
        })          
      }
    ) 
}

