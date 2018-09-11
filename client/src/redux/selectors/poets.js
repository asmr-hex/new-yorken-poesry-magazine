import {get, map} from 'lodash'


export const getPoetsOfUser = state =>
  map(
    get(state, `session.user.poets`, []),
    id => get(state, `poets.${id}`, {}),
    [],
  )
