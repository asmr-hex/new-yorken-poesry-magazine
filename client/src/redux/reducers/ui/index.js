import {
  SIGNIN_SELECTED,
  TITLE_SHOWN,
  MENU_SHOWN,
} from '../../actions/ui'
import {
  SIGNIN_UI,
  MINIMAL_UI,
} from '../../../types/ui'

const defaultUiState = {
  display: MINIMAL_UI,
  showTitle: true,
}

export const ui = (state = defaultUiState, action) => {
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
