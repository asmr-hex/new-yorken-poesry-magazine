import {READ_SUPPORTED_LANGUAGES_SUCCESSFUL} from '../actions/languages'

export const languages = (state=[], action) => {
  switch (action.type) {
  case READ_SUPPORTED_LANGUAGES_SUCCESSFUL:
    return action.payload
  default:
    return state
  }
}
