import {reduce} from 'lodash'
import {
  READ_ISSUES_SUCCESSFUL,
  READ_ISSUE_SUCCESSFUL,
} from '../actions/issues'

export const poems = (state={}, action) => {
  switch (action.type) {
  case READ_ISSUES_SUCCESSFUL:
    return mergePoemsFromIssues(state, action.payload)
  case READ_ISSUE_SUCCESSFUL:
    // since issues come loaded with poems we need
    // to normalize the issue and load the poems here
    const issue = action.payload
    return mergePoemsById(state, issue.poems)
  default:
    return state
  }
}


export const mergePoemsFromIssues = (state, issues) =>
  reduce(
    issues,
    (acc, issue) => mergePoemsById(acc, issue.poems),
    state,
  )

export const mergePoemsById = (state, poems) =>
  reduce(
    poems,
    (acc, poem) => ({...acc, [poem.id]: poem}),
    state,
  )
