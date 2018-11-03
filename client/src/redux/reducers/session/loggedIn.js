import {LOGIN_SUCCESSFUL} from '../../actions/login'
import {VERIFY_SUCCESSFUL} from '../../actions/verify'


export const loggedIn = (state = false, action) => {
  switch (action.type) {
  case VERIFY_SUCCESSFUL:
  case LOGIN_SUCCESSFUL:
    return true
  default:
    return state
  }
}
