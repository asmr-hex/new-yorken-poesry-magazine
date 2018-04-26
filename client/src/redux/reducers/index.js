import {combineReducers} from 'redux'
import { reducer as form } from 'redux-form'
import {ui} from './ui'
import {error} from './error'
import {session} from './session'


export const reducers = combineReducers({
  form,
  error,
  session,
  ui,
})
