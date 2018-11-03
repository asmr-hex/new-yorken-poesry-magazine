import {SIGNUP_SUCCESSFUL} from '../../actions/signup'


export const pendingVerification = (state = false, action) => {
  switch (action.type) {
  case SIGNUP_SUCCESSFUL:
    return true
  default:
    return state
  }
}
