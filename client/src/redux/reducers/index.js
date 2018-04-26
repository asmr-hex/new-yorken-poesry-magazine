import {combineReducers} from 'redux'
import { reducer as form } from 'redux-form'
import {ui} from './ui'
import {error} from './error'


export const reducers = combineReducers({
  ui,
  form,
  error,
})
