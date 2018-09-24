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


export const READ_POETS_REQUESTED = 'READ_POETS_REQUESTED'
export const requestReadPoets = () => dispatch => {
  dispatch({
    type: READ_POETS_REQUESTED,
  })

  // make request for poet
  fetch(
    `/api/v1/poets`,
    {
      method: 'GET',
    })
    .then(checkResponse)
    .then(
      poets => dispatch(readPoetsSuccessful(poets)),
      error => dispatch(readPoetsFailed(error)),
    )  
}

export const READ_POETS_SUCCESSFUL = 'READ_POETS_SUCCESSFUL'
export const readPoetsSuccessful = poets => dispatch =>
  dispatch({
    payload: poets,
    type: READ_POETS_SUCCESSFUL,
  })

export const READ_POETS_FAILED = 'READ_POETS_FAILED'
export const readPoetsFailed = error => dispatch => {
  error.message = 'read poet failed: ' + error.message
  
  dispatch({
    error,
    type: READ_POETS_FAILED,
  })
}



export const READ_POET_REQUESTED = 'READ_POET_REQUESTED'
export const requestReadPoet = poetId => dispatch => {
  const payload = {poetId}

  dispatch({
    payload,
    type: READ_POET_REQUESTED,
  })

  // make request for poet
  fetch(
    `/api/v1/poet/${poetId}`,
    {
      method: 'GET',
    })
    .then(checkResponse)
    .then(
      poet => dispatch(readPoetSuccessful(poet)),
      error => dispatch(readPoetFailed(error)),
    )  
}

export const READ_POET_SUCCESSFUL = 'READ_POET_SUCCESSFUL'
export const readPoetSuccessful = poet => dispatch =>
  dispatch({
    payload: poet,
    type: READ_POET_SUCCESSFUL,
  })

export const READ_POET_FAILED = 'READ_POET_FAILED'
export const readPoetFailed = error => dispatch => {
  error.message = 'read poet failed: ' + error.message
  
  dispatch({
    error,
    type: READ_POET_FAILED,
  })
}


export const READ_POET_CODE_REQUESTED = 'READ_POET_CODE_REQUESTED'
export const requestReadPoetCode = poetId => dispatch => {
  const payload = {poetId}

  dispatch({
    payload,
    type: READ_POET_CODE_REQUESTED,
  })

  // make request for poet code
  fetch(
    `/api/v1/code/${poetId}`,
    {
      method: 'GET',
    })
    .then(checkResponse)
    .then(
      code => dispatch(readPoetCodeSuccessful(poetId, code)),
      error => dispatch(readPoetCodeFailed(error)),
    )
}

export const READ_POET_CODE_SUCCESSFUL = 'READ_POET_CODE_SUCCESSFUL'
export const readPoetCodeSuccessful = (poetId, code) => dispatch =>
  dispatch({
    payload: {poetId, code},
    type: READ_POET_CODE_SUCCESSFUL,
  })

export const READ_POET_CODE_FAILED = 'READ_POET_CODE_FAILED'
export const readPoetCodeFailed = error => dispatch => {
  error.message = 'read poet code failed: ' + error.message
  
  dispatch({
    error,
    type: READ_POET_CODE_FAILED,
  })
}


export const GENERATE_POEM_REQUESTED = 'GENERATE_POEM_REQUESTED'
export const requestGeneratePoem = poetId => dispatch => {
  const payload = {poetId}

  dispatch({
    payload,
    type: GENERATE_POEM_REQUESTED,
  })

  // make request to generate poem
  fetch(
    `/api/v1/poet/${poetId}/write-poem`,
    {
      method: 'GET',
    })
    .then(checkResponse)
    .then(
      poem => dispatch(generatePoemSuccessful(poem)),
      error => dispatch(generatePoemFailed(error)),
    )  
}

export const GENERATE_POEM_SUCCESSFUL = 'GENERATE_POEM_SUCCESSFUL'
export const generatePoemSuccessful = poem => dispatch =>
  dispatch({
    payload: poem,
    type: GENERATE_POEM_SUCCESSFUL,
  })

export const GENERATE_POEM_FAILED = 'GENERATE_POEM_FAILED'
export const generatePoemFailed = error => dispatch => {
  error.message = 'read poet code failed: ' + error.message
  
  dispatch({
    error,
    type: GENERATE_POEM_FAILED,
  })
}
