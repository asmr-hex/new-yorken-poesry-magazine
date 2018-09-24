import React, { Component } from 'react';
import {connect} from 'react-redux'
import Highlight from 'react-highlight'
import {get, isEmpty} from 'lodash'
import {formatDate} from '../../types/date'
import {getPoetCode} from '../../redux/selectors/poets'
import {
  requestReadPoet,
  requestReadPoetCode,
  requestGeneratePoem,
} from '../../redux/actions/poets'
import './index.css'
import {Poem} from '../poems/poem'


const mapStateToProps = (state, ownProps) => {
  // either get id from it being passed into the component directly
  // as a prop or get it from the url path params
  const id = ownProps.id || ownProps.match.params.id

  return {
    id,
    code: getPoetCode(id, state),
    poet: get(state, `poets.${id}`, {}),
    generatedPoem: get(state, `generatedPoemsByPoetId.${id}`, {}),
  }
}

const actions = {
  fetchPoet: requestReadPoet,
  fetchCode: requestReadPoetCode,
  writePoem: requestGeneratePoem,
}

export class poet extends Component {
  constructor(props) {
    super(props)

    this.state = {
      view: 'overview' // can also be 'code'
    }
  }
  
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

  chooseView(view) {
    this.setState({view})
  }
  
  render() {
    const {
      code,
      poet,
      writePoem,
      generatedPoem,
    } = this.props
    
    return (
      <div className='main'>
        <div className='poet-container'>
          <div className='poet-header'>
            <div className='poet-header-name'>{poet.name}</div>
            <div className='poet-subheader'>
              <span className='poet-subheader-language'>{poet.language}</span>
              <span className='poet-subheader-designer-text'>
                designed by
                {/* <Link to={`/user/${poet.designer}`} className='text-link'> */}
                  <span className='poet-subheader-designer'>{get(poet, `designer.username`, '')}</span>
                  {/* </Link> */}
              </span>

            </div>
          </div>
          <div className='poet-body'>
            <div className='poet-body-menu'>
              <span className={this.state.view === 'overview' ?
                                   'poet-body-menu-item-selected'
                                   : 'poet-body-menu-item'}
                    onClick={() => this.chooseView('overview')}
                    >
                overview
              </span>
              <span className={this.state.view === 'code' ?
                                   'poet-body-menu-item-selected'
                                   : 'poet-body-menu-item'}
                    onClick={() => this.chooseView('code')}
                    >
                code
              </span>
            </div>
            <div className='poet-body-content'>
              {
                this.state.view === 'overview' ?
                <PoetOverview poet={poet} writePoem={writePoem} generatedPoem={generatedPoem}/>
                  :<PoetCode poet={poet} code={code}/>
              }
            </div>
          </div>
        </div>
      </div>
    )
  }
}

export const Poet = connect(mapStateToProps, actions)(poet)

export class PoetOverview extends Component {
  render() {
    const {
      poet,
      writePoem,
      generatedPoem,
    } = this.props

    return (
      <div className='poet-overview'>
        <div className='poet-overview-details'>
          <span className='poet-overview-details-birthday'>
            birthday:   {formatDate(poet.birthDate)}
          </span>
          <span className='poet-overview-details-description'>
            {poet.description}
          </span>
        </div>
        <div className='poet-overview-generate-poem'>
          <div className='poet-overview-generate-poem-button'
               onClick={() => writePoem(poet.id)}
            >
            generate a poem
          </div>
          {
            isEmpty(generatedPoem) ? null : <Poem poem={generatedPoem}/>
          }
        </div>
      </div>
    )
  }
}

export class PoetCode extends Component {
  render() {
    const {
      code
    } = this.props
    
    return (
      <div className='poet-code'>
        <Highlight className="python poet-body-code">
          {code.code}
        </Highlight>
      </div>
    )
  }
}
