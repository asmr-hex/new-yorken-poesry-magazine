import {get, map} from 'lodash'


export const getPoemsByIssueVolume = (state, volume) =>
  map(
    get(state, `issuesByVolume.${volume}.poems`, []),
    id => get(state, `poems.${id}`, {}),
    [],
  )
