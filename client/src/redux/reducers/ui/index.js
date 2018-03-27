import {SIGNIN_SELECTED} from '../../actions/ui'
import {SIGNIN_UI} from '../../../types/ui'


export const ui = (state = {}, action) => {
  switch (action.type) {
  case SIGNIN_SELECTED:
    return {...state, display: SIGNIN_UI}
  default:
    return state
  }
}
