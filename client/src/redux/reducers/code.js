import {READ_POET_CODE_SUCCESSFUL} from '../actions/poets'


export const codeByPoetId = (state={}, action) => {
  switch (action.type) {
  case READ_POET_CODE_SUCCESSFUL:
    // convert code bytes to utf8 characters
    const code = action.payload.code
    return {...state, [action.payload.poetId]: code}
  default:
    return state
  }
}
