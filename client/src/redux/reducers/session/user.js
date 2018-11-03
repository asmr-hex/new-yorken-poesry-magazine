import {filter, map} from 'lodash'
import {LOGIN_SUCCESSFUL} from '../../actions/login'
import {VERIFY_SUCCESSFUL} from '../../actions/verify'
import {
  CREATE_POET_SUCCESSFUL,
  DELETE_POET_SUCCESSFUL,
} from '../../actions/poets'


export const user = (state = {}, action) => {
  switch (action.type) {
  case CREATE_POET_SUCCESSFUL:
    // when a poet is created by our user, we want it to show up here.
    return {...state, poets: [...state.poets, action.payload.id]}
  case DELETE_POET_SUCCESSFUL:
    return {...state, poets: filter(state.poets, id => id !== action.payload.id)}
  case VERIFY_SUCCESSFUL:
  case LOGIN_SUCCESSFUL:
    // both login and signup will return a user which we have to modify slightly.
    // particularly, we want to denormalize how we are storing poets nested in a
    // user. See the #poets reducer for details about how poets are stored.
    return denormalizePoets(action.payload)
  default:
    return state
  }
}

export const denormalizePoets = user => ({
  ...user,
  poets: map(user.poets, poet => poet.id, []),
})
