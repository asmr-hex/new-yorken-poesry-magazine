import fetch from 'cross-fetch'
import {checkResponse} from './error'


export const READ_SUPPORTED_LANGAUGES_REQUESTED = 'READ_SUPPORTED_LANGUAGES_REQUESTED'
export const requestReadLanguages = () => dispatch => {
  dispatch({
    type: READ_SUPPORTED_LANGAUGES_REQUESTED,
  })

  fetch(
    `/api/v1/supported-languages`,
    {
      method: 'GET',
    })
    .then(checkResponse)
    .then(
      languages => dispatch(readLanguagesSuccessful(languages)),
      error => dispatch(readLanguagesFailed(error)),
    )
}


export const READ_SUPPORTED_LANGUAGES_SUCCESSFUL = 'READ_SUPPORTED_LANGUAGES_SUCCESSFUL'
export const readLanguagesSuccessful = languages => dispatch =>
  dispatch({
    payload: languages,
    type: READ_SUPPORTED_LANGUAGES_SUCCESSFUL,
  })


export const READ_SUPPOERTED_LANGUAGES_FAILED = 'READ_SUPPOERTED_LANGUAGES_FAILED'
export const readLanguagesFailed = error => dispatch => {
  error.message = 'read supported languages failed: ' + error.message
  
  dispatch({
    error,
    type: READ_SUPPOERTED_LANGUAGES_FAILED,
  })  
}

