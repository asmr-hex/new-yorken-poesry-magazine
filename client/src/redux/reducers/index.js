import {combineReducers} from 'redux'
import { reducer as form } from 'redux-form'
import {ui} from './ui'
import {error} from './error'
import {userError} from './userError'
import {session} from './session'
import {poets} from './poets'
import {issuesByVolume} from './issues'
import {poems} from './poems'
import {languages} from './languages'
import {codeByPoetId} from './code'
import {generatedPoemsByPoetId} from './generatedPoems'


export const reducers = combineReducers({
  form,
  error,
  userError, // TODO (cw|11.4.2018) consolidate into error reducer
  poets,
  languages,
  poems,
  issuesByVolume,
  codeByPoetId,
  generatedPoemsByPoetId,
  session,
  ui,
})
