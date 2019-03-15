
/**
 * dispatch this action when we want to clear the global error message.
 */
export const RESET_ERROR_MESSAGE = 'RESET_ERROR_MESSAGE'
export const resetErrorMsg = () => dispatch =>
  dispatch({
    type: RESET_ERROR_MESSAGE,
  })

// TODO (cw|4.24.2018) wrap fetch in a status checker?
export const checkResponse = response => {
  if (response.ok) {
      // decode response body from json
      return response.json()
  }

  // TODO (cw|4.24.2018) create custom error classes?
  // TODO (cw|4.27.2018) eventually get the response.text() to work.
  // rn it returns another promise.
  throw response// new Error(`${response.statusText}: ${response.text()}`)
}
