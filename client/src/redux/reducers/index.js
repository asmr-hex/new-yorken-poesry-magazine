import {combineReducers} from 'redux'
import { reducer as form } from 'redux-form'
import {ui} from './ui'
import {error} from './error'
import {session} from './session'
import {poets} from './poets'


export const reducers = combineReducers({
  form,
  error,
  poets,
  session,
  ui,
})
