import {combineReducers} from 'redux'
import {ui} from './ui'
import {login} from './login'


export const reducers = combineReducers({
  ui,
  login,
})
