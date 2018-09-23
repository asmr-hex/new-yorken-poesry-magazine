import {reduce} from 'lodash'
import {LOGIN_SUCCESSFUL} from '../actions/login'
import {
  CREATE_POET_SUCCESSFUL,
  READ_POET_SUCCESSFUL,
  READ_POETS_SUCCESSFUL,
} from '../actions/poets'
import {
  READ_ISSUES_SUCCESSFUL,
  READ_ISSUE_SUCCESSFUL,
} from '../actions/issues'


export const poets = (state = {}, action) => {
  switch (action.type) {
  case READ_ISSUES_SUCCESSFUL:
    return mergePoetsFromIssues(state, action.payload)
  case READ_ISSUE_SUCCESSFUL:
    // since issues come loaded with poets (judges + contributors)
    // we need to normalize the issue and load in poets here
    const issue = action.payload
    return mergePoetsById(state, [...issue.committee, ...issue.contributors])
  case READ_POET_SUCCESSFUL:
  case CREATE_POET_SUCCESSFUL:
    return {...state, [action.payload.id]: action.payload}
  case READ_POETS_SUCCESSFUL:
    return mergePoetsById(state, action.payload)
  case LOGIN_SUCCESSFUL:
    // if we are logged in, we want to extract the poets from the
    // user and store them by id in this part of the state tree
    return mergePoetsById(state, action.payload.poets)
  default:
    return state
  }
}

export const mergePoetsFromIssues = (state, issues) =>
  reduce(
    issues,
    (acc, issue) => mergePoetsById(acc, [...issue.committee, ...issue.contributors]),
    state,
  )

export const mergePoetsById = (state, poets) =>
  reduce(
    poets,
    (acc, poet) => ({...acc, [poet.id]: poet}),
    state,
  )
