import {reduce} from 'lodash'
import {LOGIN_SUCCESSFUL} from '../actions/login'
import {CREATE_POET_SUCCESSFUL} from '../actions/poets'


export const poets = (state = {}, action) => {
  switch (action.type) {
  case CREATE_POET_SUCCESSFUL:
    return {...state, [action.payload.id]: action.payload}
  case LOGIN_SUCCESSFUL:
    // if we are logged in, we want to extract the poets from the
    // user and store them by id in this part of the state tree
    return mergePoetsById(state, action.payload.poets)
  default:
    return state
  }
}

export const mergePoetsById = (state, poets) =>
  reduce(
    poets,
    (acc, poet) => ({...acc, [poet.id]: poet}),
    state,
  )
