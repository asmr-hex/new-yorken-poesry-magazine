import React, { Component } from 'react';
import {connect} from 'react-redux'
import { Field, reduxForm } from 'redux-form'
import {map} from 'lodash'
import {getPoetsOfUser} from '../../redux/selectors/poets'
import {requestCreatePoet} from '../../redux/actions/poets'


class poetMgmt extends Component {
  createPoet = values => {
    // const {files} = event.target
    const {
      name,
      description,
      language,
      program,
      parameters,
    } = values
    const {requestCreatePoet} = this.props
    
    requestCreatePoet({
      name: name || '',
      description: description || '',
      language: language || '',
      program,
      parameters,
    })
  }

  render() {
    return (
      <div>
        {
          map(
            this.props.poets,
            (poet, idx) => (
              <div key={idx}>{poet.name}</div>
            ),
            [],
          )
        }
        <CreatePoetForm onSubmit={this.createPoet}/>
      </div>
    )
  }
}

const createPoetForm = props => {
  const {handleSubmit} = props
  
  return (
    <form onSubmit={handleSubmit}>
      <div>
        <Field name='name' component='input' type='text' placeholder='name'/>
      </div>
      <div>
        <Field name='description' component='input' type='text' placeholder='description'/>
      </div>
      <div>
        <Field name='language' component='input' type='text' placeholder='language'/>
      </div>
      <div>
        <Field name='program' component={FileInput}/>
      </div>
      <div>
        <Field name='parameters' component={FileInput}/>
      </div>
      <button type='submit'>create poet</button>
    </form>
  )
}

export const CreatePoetForm = reduxForm({
  form: 'createPoet',
})(createPoetForm)

// TODO (cw|4.27.2018) refactor this into something much nicer -__-
const adaptFileEventToValue = delegate =>
      e => delegate(e.target.files[0])

const FileInput = ({
  input: {
    value: omitValue,
    onChange,
    onBlur,
    ...inputProps,
  },
  meta: omitMeta,
  ...props,
}) =>
      <input
onChange={adaptFileEventToValue(onChange)}
onBlur={adaptFileEventToValue(onBlur)}
type="file"
{...inputProps}
{...props}
  />

const mapStateToProps = (state, ownProps) => ({
  poets: getPoetsOfUser(state),
  user: state.session.user,
})

const actions = {
  requestCreatePoet,
}

export const PoetMgmt = connect(mapStateToProps, actions)(poetMgmt)
