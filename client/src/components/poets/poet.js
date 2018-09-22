import React, { Component } from 'react';
import {Route, Switch} from 'react-router-dom'
import {connect} from 'react-redux'
import Highlight from 'react-highlight'
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
      <div styles={{textAlign: 'left !important'}}>
        <h3>{code.filename}</h3>
        <Highlight className="python" styles={{textAlign: 'left'}}>
          {code.code}
        </Highlight>
      </div>
    )
  }
}

export const Poet = connect(mapStateToProps, actions)(poet)
