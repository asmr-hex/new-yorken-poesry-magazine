import {READ_POET_CODE_SUCCESSFUL} from '../actions/poets'


export const codeByPoetId = (state={}, action) => {
  switch (action.type) {
  case READ_POET_CODE_SUCCESSFUL:
    return {...state, [action.payload.poetId]: action.payload.code}
  default:
    return state
  }
}
