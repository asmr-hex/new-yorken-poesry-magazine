import {GENERATE_POEM_SUCCESSFUL} from '../actions/poets'


export const generatedPoemsByPoetId = (state={}, action) => {
  switch (action.type) {
  case GENERATE_POEM_SUCCESSFUL:
    return {...state, [action.payload.author.id]: action.payload}
  default:
    return state
  }
}
