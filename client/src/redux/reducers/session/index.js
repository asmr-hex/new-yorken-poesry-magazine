import {combineReducers} from 'redux'
import {loggedIn} from './loggedIn'
import {user} from './user'
import {pendingVerification} from './pendingVerification'


export const session = combineReducers({
  loggedIn,
  user,
  pendingVerification,
})
