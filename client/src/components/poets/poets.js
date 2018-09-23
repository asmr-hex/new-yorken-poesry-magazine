import React, { Component } from 'react';
import {Route, Switch} from 'react-router-dom'
import {connect} from 'react-redux'
import {Link} from 'react-router-dom'
import {get, isEmpty, map, values} from 'lodash'
import {requestReadPoets} from '../../redux/actions/poets'


const mapStateToProps = (state, ownProps) => ({
  poets: values(get(state, `poets`, {})),
})

const actions = {
  fetchPoets: requestReadPoets,
}

export class poets extends Component {
  componentDidMount() {
    const {
      poets,
      fetchPoets,
    } = this.props

    fetchPoets()
  }
  
  render() {
    const {
      poets,
    } = this.props
    
    return (
      <div>
        {
          map(
            poets,
            (poet, idx) => (
              <div key={idx}>{poet.name}</div>
            ),
            [],
          )
        }
      </div>
    )
  }
}

export const Poets = connect(mapStateToProps, actions)(poets)
