import fetch from 'cross-fetch'
import {reduce} from 'lodash'
import {checkResponse} from './error'

// TODO: stopping point-- when i submit the create poet form, it modifies the
// url for somereason with the path parameters.. OH OH
// THIS MIGHT ACTUALLY BE A RESULT OF THE WRITTEN RESPONSE FROM THE BACKEND??
// actually idk... :/ (i'll debug after the interview)

export const CREATE_POET_REQUESTED = 'CREATE_POET_REQUESTED'
export const requestCreatePoet = ({name, description, language, program, parameters}) => dispatch => {
  const payload = {name, description, language}

  dispatch({
    payload,
    type: CREATE_POET_REQUESTED,
  })

  // create FormData and add key value data
  let formData = new FormData()
  reduce(
    payload,
    (_, v, k) => {
      console.log(k, v)
      formData.append(k, v) 
    },
    {},
  )

  // add files to FormData
  // TODO (cw|4.27.2018) we are following the conventions described in
  // server/core/handlers#CreatePoet. However, they seem somewhat arbitrary.
  // it would be nice to not have to have this "[]" after src.
  formData.append("src[]", program, "program")
  formData.append("src[]", parameters, "parameters")
  
  fetch(
    `/dashboard/poet`,
    {
      method: 'POST',
      body: formData,
    })
    .then(checkResponse)
    .then(
      poet => dispatch(createPoetSuccessful(poet)),
      error => dispatch(createPoetFailed(error)),
    )
}


export const CREATE_POET_SUCCESSFUL = 'CREATE_POET_SUCCESSFUL'
export const createPoetSuccessful = poet => dispatch =>
  dispatch({
    payload: poet,
    type: CREATE_POET_SUCCESSFUL,
  })


export const CREATE_POET_FAILED = 'CREATE_POET_FAILED'
export const createPoetFailed = error => dispatch => {
  error.message = 'create poet failed: ' + error.message
  
  dispatch({
    error,
    type: CREATE_POET_FAILED,
  })  
}
