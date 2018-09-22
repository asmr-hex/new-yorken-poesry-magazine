import React, { Component } from 'react';
import {Route, Switch} from 'react-router-dom'
import {connect} from 'react-redux'
import {getPoetCode} from '../../redux/selectors/poets'
import {requestReadPoetCode} from '../../redux/actions/poets'


const mapStateToProps = (state, ownProps) => {
  // either get id from it being passed into the component directly
  // as a prop or get it from the url path params
  const id = ownProps.id || ownProps.match.params.id

  return {
    id,
    code: getPoetCode(id, state),
  }
}

const actions = {
  fetchCode: requestReadPoetCode,
}

export class poet extends Component {
  componentDidMount() {
    const {id, fetchCode} = this.props

    // fetch the code for this poet
    fetchCode(id)
  }
  
  render() {
    const {
      code
    } = this.props
    
    return (
      <div>
        <h3>{code.filename}</h3>
        <p>{code.code}</p>
      </div>
    )
  }
}

export const Poet = connect(mapStateToProps, actions)(poet)
