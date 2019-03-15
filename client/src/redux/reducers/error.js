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

  } else if (type === '@@redux-form/SET_SUBMIT_FAILED') {
    // TODO (cw|10.12.2018) this is kinda janky...it would be nice to use
    // some constant from the redux-form library instead of a magic string
    // i found from one of their unexported constants -__- so maybe i'm just
    // approaching this the wrong way...
    // this is a redux-form error
    
    return 'u forgot to fill out  some stuff buddy'
    
  } else if (error) {
    // this is an error from the server
    
    // set the error message
    return error.message
    
  }

  return state
}
