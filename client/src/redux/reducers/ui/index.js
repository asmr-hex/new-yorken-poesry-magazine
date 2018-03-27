import {
  SIGNIN_SELECTED,
  TITLE_SHOWN,
  MENU_SHOWN,
} from '../../actions/ui'
import {SIGNIN_UI} from '../../../types/ui'


export const ui = (state = {}, action) => {
  switch (action.type) {
  case SIGNIN_SELECTED:
    return {...state, display: SIGNIN_UI}
  case TITLE_SHOWN:
    return {...state, showTitle: true}
  case MENU_SHOWN:
    return {...state, showTitle: false}
  default:
    return state
  }
}
