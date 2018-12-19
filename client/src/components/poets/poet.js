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
import '../app/highlight.css'
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

    const pageStyle = {
      display: 'flex',
      flexDirection: 'column',
      width: '100%',
      margin: '5em 25% 5em 25%',
    }
    
    const headerStyle = {
      color: '#ffb2e4',
      fontSize: '4em',
      textShadow: '4px 4px #affbff',
    }

    const subheaderStyle = {
      color: '#75b0ff',
      textShadow: '4px 4px #affbff',
    }

    return (
      <div className='main'>
        <div className='poet-container' style={pageStyle}>
          <div className='poet-header'>
            <div className='poet-header-name' style={headerStyle}>{poet.name}</div>
            <div className='poet-subheader' style={subheaderStyle}>
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
            <PoetOverview poet={poet} writePoem={writePoem} generatedPoem={generatedPoem}/>           
            <PoetCode {...{ poet, code}}/>
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
        <table className='poet-overview-details'>
          <tr>
            <td className='poet-overview-detail'>birthday   </td>
            <td className='poet-overview-detail-value'>
              {formatDate(poet.birthDate)}
            </td>
          </tr>
          <tr>
            <td className='poet-overview-detail'>published works   </td>
            <td className='poet-overview-detail-value'>
              {'-'}
            </td>
          </tr>
          <tr>
            <td className='poet-overview-detail'>volumes curated   </td>
            <td className='poet-overview-detail-value'>
              {'-'}
            </td>
          </tr>
        </table>
        
        <span className='poet-overview-description'>
          {poet.description}
        </span>
        <div style={{display: 'flex', justifyContent: 'center', width: '100%'}}>
          <div className='poet-overview-generate-poem-button'
               onClick={() => writePoem(poet.id)}
            >
            generate a poem
          </div>
        </div>
        {
          isEmpty(generatedPoem) ? null : <Poem poem={generatedPoem}/>
        }
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
        <div className='poet-overview-header'>code</div>
        <Highlight className="python poet-body-code">
          {code.code}
        </Highlight>
      </div>
    )
  }
}
