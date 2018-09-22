import {reduce} from 'lodash'
import {LOGIN_SUCCESSFUL} from '../actions/login'
import {CREATE_POET_SUCCESSFUL} from '../actions/poets'
import {READ_ISSUE_SUCCESSFULL} from '../actions/issues'


export const poets = (state = {}, action) => {
  switch (action.type) {
  case READ_ISSUE_SUCCESSFULL:
    // since issues come loaded with poets (judges + contributors)
    // we need to normalize the issue and load in poets here
    const issue = action.payload
    return mergePoetsById(state, [...issue.committee, ...issue.contributors])
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
