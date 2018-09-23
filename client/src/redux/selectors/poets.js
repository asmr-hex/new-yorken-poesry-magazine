import {get, map} from 'lodash'


export const getPoetsOfUser = state =>
  map(
    get(state, `session.user.poets`, []),
    id => get(state, `poets.${id}`, {}),
    [],
  )

export const getJudgesByIssueVolume = (volume, state) =>
  map(
    get(state, `issuesByVolume.${volume}.committee`, []),
    id => get(state, `poets.${id}`, {}),
    [],
  )

export const getContributorsByIssueVolume = (volume, state) =>
  map(
    get(state, `issuesByVolume.${volume}.contributors`, []),
    id => get(state, `poets.${id}`, {}),
    [],
  )

export const getPoetCode = (id, state) =>
  get(state, `codeByPoetId.${id}`, {})
