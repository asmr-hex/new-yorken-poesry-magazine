import fetch from 'cross-fetch'


export const LOGIN_REQUESTED = 'LOGIN_REQUESTED'
export const requestLogin = ({username, password}) => dispatch => {
  const payload = {username, password}
  
  dispatch({
    payload,
    type: LOGIN_REQUESTED,
  })

  fetch(
    `/api/dashboard/login`,
    {
      method: 'POST',
      body: JSON.stringify(payload),
      headers: new Headers({'Content-Type': 'application/json'}),
    })
    .then(handleErrors)
    .then(
      response => console.log(response.json()),
      err => console.log('login error: ', err),
    )
    .then(json => dispatch(loginSuccessful()))  
}

// TODO (cw|4.24.2018) wrap fetch in a status checker?
const handleErrors = response => {
  if (response.ok) {
    return response
  }
  
  // TODO (cw|4.24.2018) create custom error classes
  throw new Error(response.statusText)
}

      

export const LOGIN_SUCCESSFUL = 'LOGIN_SUCCESSFUL'
export const loginSuccessful = () => dispatch =>
  dispatch({
    type: LOGIN_SUCCESSFUL,
  })


export const LOGIN_FAILED = 'LOGIN_FAILED'
export const loginFailed = () => dispatch =>
  dispatch({
    type: LOGIN_FAILED,
  })
