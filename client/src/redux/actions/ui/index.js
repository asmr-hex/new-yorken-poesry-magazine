export const SIGNIN_SELECTED = 'SIGNIN_SELECTED'
export const selectSignIn = () => dispatch =>
  dispatch({
    type: SIGNIN_SELECTED,
  })

export const TITLE_SHOWN = 'TITLE_SHOWN'
export const showTitle = () => dispatch =>
  dispatch({
    type: TITLE_SHOWN,
  })

export const MENU_SHOWN = 'MENU_SHOWN'
export const showMenu = () => dispatch =>
  dispatch({
    type: MENU_SHOWN,
  })
