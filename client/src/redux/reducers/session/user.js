import {LOGIN_SUCCESSFUL} from '../../actions/login'
import {SIGNUP_SUCCESSFUL} from '../../actions/signup'


export const user = (state = {}, action) => {
  switch (action.type) {
  case SIGNUP_SUCCESSFUL:
  case LOGIN_SUCCESSFUL:
    return action.payload
  default:
    return state
  }
}
