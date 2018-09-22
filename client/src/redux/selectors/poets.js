import {get, map} from 'lodash'


export const getPoetsOfUser = state =>
  map(
    get(state, `session.user.poets`, []),
    id => get(state, `poets.${id}`, {}),
    [],
  )

export const getJudgesByIssueVolume = (state, volume) =>
  map(
    get(state, `issuesByVolume.${volume}.committee`, []),
    id => get(state, `poets.${id}`, {}),
    [],
  )

export const getContributorsByIssueVolume = (state, volume) =>
  map(
    get(state, `issuesByVolume.${volume}.contributors`, []),
    id => get(state, `poets.${id}`, {}),
    [],
  )
