import {CREATE_POET_FAILED_USER_ERROR} from '../actions/poets'
import { RESET_ERROR_MESSAGE } from '../actions/error'


export const userError = (state=null, action) => {
  switch (action.type) {
  case RESET_ERROR_MESSAGE:
    return null
  case CREATE_POET_FAILED_USER_ERROR:
    return action.error.message
  default:
    return state
  }
}
