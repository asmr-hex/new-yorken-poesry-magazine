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

    const pageStyle = {
      display: 'flex',
      flexDirection: 'column',
      width: '100%',
      margin: '2em 25% 5em 25%',
    }
    
    const headerStyle = {
      color: '#ffb2e4',
      fontSize: '4em',
      textShadow: '4px 4px #affbff',
    }
    
    return (
      <div className='main'>
        <div className='poets-summaries-container' style={pageStyle}>
          <div className='poets-header' style={headerStyle}>poets</div>
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

    const contentStyle = {
      fontSize: '2.3em',
      fontWeight: 'bold',
      color: '#19ecff',
      textShadow: '2px 2px #ffb2e4',
    }
    
    return (
      <Link to={`/poet/${id}`}>
        <div className='poet-row-container'>
          <div className='poet-row' style={contentStyle}>
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
