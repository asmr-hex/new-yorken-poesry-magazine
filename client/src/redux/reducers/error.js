import {
  RESET_ERROR_MESSAGE,
} from '../actions/error'


/**
 * reduce the error message state.
 */
export const error = (state = null, action) => {
  const {type, error} = action

  if (type === RESET_ERROR_MESSAGE) {
    return null
  } else if (error) {
    // print out stack trace in console
    console.error(error)

    // set the error message
    return error.message
  }

  return state
}
