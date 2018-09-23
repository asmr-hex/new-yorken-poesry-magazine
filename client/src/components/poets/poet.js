import React, { Component } from 'react';
import {Route, Switch} from 'react-router-dom'
import {connect} from 'react-redux'
import Highlight from 'react-highlight'
import {get, isEmpty} from 'lodash'
import {getPoetCode} from '../../redux/selectors/poets'
import {
  requestReadPoet,
  requestReadPoetCode,
} from '../../redux/actions/poets'


const mapStateToProps = (state, ownProps) => {
  // either get id from it being passed into the component directly
  // as a prop or get it from the url path params
  const id = ownProps.id || ownProps.match.params.id

  return {
    id,
    code: getPoetCode(id, state),
    poet: get(state, `poets.${id}`, {}),
  }
}

const actions = {
  fetchPoet: requestReadPoet,
  fetchCode: requestReadPoetCode,
}

export class poet extends Component {
  componentDidMount() {
    const {
      id,
      fetchPoet,
      fetchCode,
      poet,
      code,
    } = this.props

    // fetch the code for this poet
    if (isEmpty(code)) {
      fetchCode(id)      
    }

    // fetch poet if we don't already have it
    if (isEmpty(poet)) {
      fetchPoet(id)
    }

  }
  
  render() {
    const {
      code,
      poet,
    } = this.props

    return (
      <div>
        <h3>{code.filename}</h3>
        <p>{poet.name}</p>
        <Highlight className="python">
          {code.code}
        </Highlight>
      </div>
    )
  }
}

export const Poet = connect(mapStateToProps, actions)(poet)
