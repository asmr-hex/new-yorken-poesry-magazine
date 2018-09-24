import React, { Component } from 'react';
import {connect} from 'react-redux'
import {Link} from 'react-router-dom'
import {get, map, values} from 'lodash'
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
      fetchPoets,
    } = this.props

    fetchPoets()
  }
  
  render() {
    const {
      poets,
    } = this.props
    
    return (
      <div className='main'>
        <div className='poets-summaries-container'>
          <div className='poets-header'>poets</div>
          {
            map(
              poets,
              (poet, idx) => (
                <PoetRow poet={poet} key={idx}/>
              ),
              [],
            )
          }
        </div>
      </div>
    )
  }
}

export const Poets = connect(mapStateToProps, actions)(poets)

export class PoetRow extends Component {
  render() {
    const {
      id,
      name,
      language,
    } = this.props.poet
    
    return (
      <Link to={`/poet/${id}`}>
        <div className='poet-row-container'>
          <div className='poet-row'>
            <div className='poet-row-name-line'>
              <span className='poet-row-name-item'>{name}</span>
              <span className='poet-row-language-item'>{language}</span>
            </div>
          </div>
        </div>
      </Link>
    )
  }
}
