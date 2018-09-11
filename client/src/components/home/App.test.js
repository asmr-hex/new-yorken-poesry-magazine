import React from 'react'
import ReactDOM from 'react-dom'
import {createStore, applyMiddleware} from 'redux'
import {Provider} from 'react-redux'
import thunk from 'redux-thunk'
import {reducers} from '../../redux/reducers'
import {Home} from '.'

it('renders without crashing', () => {
  // TODO (cw|3.31.2018) figure out testing suite
  // const div = document.createElement('div')
  // let store = createStore(
  //   reducers,
  //   composeWithDevTools(
  //     applyMiddleware(thunk),
  //   ),
  // )

  // ReactDOM.render(
  //   <Provider store={store}/>
  //     <Home />
  //   </Provider>,
  //   div,
  // )
  // ReactDOM.unmountComponentAtNode(div)
});
