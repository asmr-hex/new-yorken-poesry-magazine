import React from 'react'
import {render} from 'react-dom'
import {createStore} from 'redux'
import {Provider} from 'react-redux'
import {BrowserRouter} from 'react-router-dom'
import {reducers} from './redux/reducers'
import './index.css'
import {Home} from './components/home'


let store = createStore(reducers)

render(
  <Provider store={store}>
    <BrowserRouter>
      <Home/>
    </BrowserRouter>
  </Provider>,
  document.getElementById('root')
)
