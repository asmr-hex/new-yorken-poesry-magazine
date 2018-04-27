import React from 'react'
import {render} from 'react-dom'
import {createStore, applyMiddleware} from 'redux'
import {Provider} from 'react-redux'
import thunk from 'redux-thunk'
import {composeWithDevTools} from 'redux-devtools-extension'
import {BrowserRouter} from 'react-router-dom'
import {reducers} from './redux/reducers'
import './index.css'
import {App} from './components/app'


let store = createStore(
  reducers,
  composeWithDevTools(
    applyMiddleware(thunk),
  ),
)

render(
  <Provider store={store}>
    <BrowserRouter>
      <App/>
    </BrowserRouter>
  </Provider>,
  document.getElementById('root')
)
