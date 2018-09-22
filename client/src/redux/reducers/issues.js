import {map} from 'lodash'
import {READ_ISSUE_SUCCESSFULL} from '../actions/issues'


export const issuesByVolume = (state = {}, action) => {
  switch (action.type) {
  case READ_ISSUE_SUCCESSFULL:
    return mergeIssueByVolume(state, action.payload)
  default:
    return state
  }
}

export const mergeIssueByVolume = (state, issue) => {
  // since issues come in with all associations contained
  // within, we want to normalize it so its just ids s.t.
  // we aren't duplicating data (single-source of truth!)
  const normalizedIssue = {
    ...issue,
    committee: map(issue.committee, judge => judge.id, []),
    contributors: map(issue.contributors, poet => poet.id, []),
    poems: map(issue.poems, poem => poem.id, []),
  }

  return {
    ...state,
    [normalizedIssue.volume]: normalizedIssue,
    latest: normalizedIssue.latest ? normalizedIssue : state.latest,
  }
}
