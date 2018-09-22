import {reduce} from 'lodash'
import {READ_ISSUE_SUCCESSFULL} from '../actions/issues'

export const poems = (state={}, action) => {
  switch (action.type) {
  case READ_ISSUE_SUCCESSFULL:
    // since issues come loaded with poems we need
    // to normalize the issue and load the poems here
    const issue = action.payload
    return mergePoemsById(state, issue.poems)
  default:
    return state
  }
}


export const mergePoemsById = (state, poems) =>
  reduce(
    poems,
    (acc, poem) => ({...acc, [poem.id]: poem}),
    state,
  )
