import {combineReducers} from 'redux'
import {loggedIn} from './loggedIn'
import {user} from './user'


export const session = combineReducers({
  loggedIn,
  user,
})
