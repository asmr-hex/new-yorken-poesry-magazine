import {LOGIN_SUCCESSFUL} from '../../actions/login'
import {SIGNUP_SUCCESSFUL} from '../../actions/signup'


export const loggedIn = (state = false, action) => {
  switch (action.type) {
  case SIGNUP_SUCCESSFUL:
  case LOGIN_SUCCESSFUL:
    return true
  default:
    return state
  }
}
