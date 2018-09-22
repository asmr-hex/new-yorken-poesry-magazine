import {get} from 'lodash'

export const getIssueByVolume = (volume, state) =>
  get(state, `issuesByVolume.${volume}`, {})
