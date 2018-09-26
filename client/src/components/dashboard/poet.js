import React, { Component } from 'react';
import {connect} from 'react-redux'
import { Field, reduxForm } from 'redux-form'
import {Link} from 'react-router-dom'
import TextField from 'material-ui/TextField'
import SelectField from 'material-ui/SelectField'
import MenuItem from 'material-ui/MenuItem'
import {get, isEmpty, map} from 'lodash'
import {formatDate} from '../../types/date'
import {getPoetsOfUser} from '../../redux/selectors/poets'
import {requestCreatePoet} from '../../redux/actions/poets'
import {requestReadLanguages} from  '../../redux/actions/languages'


const mapStateToProps = (state, ownProps) => ({
  poets: getPoetsOfUser(state),
  user: state.session.user,
  languages: get(state, `languages`, [])
})

const actions = {
  requestCreatePoet,
  requestReadLanguages,
}

class poetMgmt extends Component {
  componentDidMount() {
    const {
      languages,
      requestReadLanguages,
    } = this.props

    if (isEmpty(languages)) {
      requestReadLanguages()
    }
  }

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
    const {
      poets,
      languages,
    } = this.props
    
    return (
      <div className='profile-poets-container'>
        <div className='profile-poets-list'>
          <div className='profile-poets-list-header'>my poets:</div>
          {
            map(
              this.props.poets,
              (poet, idx) => (
                <PoetSummary poet={poet} key={idx}/>
              ),
              [],
            )
        }
        </div>
        <CreatePoetForm onSubmit={this.createPoet} languages={languages}/>
      </div>
    )
  }
}

export class PoetSummary extends Component {
  render() {
    const {
      poet,
    } = this.props
    
    return (
        <Link className='profile-poet-summary' to={`/poet/${poet.id}`}>
          <div className='profile-poet-name-language'>
            <div className='profile-poet-name'>{poet.name}</div>
            <div className='profile-poet-language'>{poet.language}</div>
          </div>
          <div className='profile-poet-birthday'>{formatDate(poet.birthDate)}</div>
          <div className='profile-poet-description'>{poet.description}</div>
        </Link>
    )
  }
}

const renderTextField = ({input, label, meta: {touched, error}, ...custom}) => (
  <TextField
    hintText={label}
    floatingLabelText={label}
    hintStyle={{color: '#222', fontSize: '1.5em'}}
    inputStyle={{color: '#222', fontSize: '1.5em'}}
    underlineStyle={{borderColor: '#222'}}
    underlineFocusStyle={{borderColor: '#222'}}
    errorText={touched && error}
    {...input}
    {...custom}
    />
)

const renderSelectField = ({
  input,
  label,
  meta: { touched, error },
  children,
  ...custom
}) => (
  <SelectField
    floatingLabelText={label}
    errorText={touched && error}
    {...input}
    onChange={(event, index, value) => input.onChange(value)}
    children={children}
    hintStyle={{color: '#222', fontSize: '1.5em'}}
    inputStyle={{color: '#222', fontSize: '1.5em'}}
    style={{
      textAlign: 'left',
    }}
    {...custom}
    />
)

export class createPoetForm extends Component {
  state = {
    value: '',
    programFileText: 'select poet file...',
    parametersFileText: 'select parameters file...',
  }

  onChangeProgram = event => {
    console.log(event)
    this.setState({ProgramFileText: event.target.value});
  }
  
  render() {
    const {
      handleSubmit,
      languages,
    } = this.props
    
    return (
      <div className='create-poet-form'>
        <form onSubmit={handleSubmit}>
          <div>
            <Field name='name' component={renderTextField} type='text' placeholder='name'/>
          </div>
          <div styles={{marginTop: '1.5em'}}>
            <Field name='description' component={renderTextField} type='text' placeholder='description'/>
          </div>
          <div>
            <Field name='language' component={renderSelectField} label='language'>
              {
                map(
                  languages,
                  (language, idx) => (
                    <MenuItem value={language} primaryText={language} key={idx}/>
                  ),
                  [],
                )
              }
            </Field>
          </div>
          <div className='profile-poet-button'>
        <Field className='profile-poet-file-button' id='program' name='program' component={FileInput}/>
            <label htmlFor="program">{this.state.programFileText}</label>
          </div>
          <div className='profile-poet-button'>
            <Field className='profile-poet-file-button' name='parameters' id='parameters' component={FileInput}/>
            <label htmlFor="parameters">{this.state.parametersFileText}</label>
            <span style={{padding: '0.6em', marginLeft: '0.8em', fontStyle: 'italic'}}>optional</span>
          </div>
          <button className='profile-poet-submit-button' type='submit'>create poet</button>
        </form>
      </div>
    )
  }
}

const validate = values => {
  const errors = {}
  if (!values.name) {
    errors.name = 'required.'
  }

  if (!values.language) {
    errors.language = 'required.'
  }

  if (!values.program) {
    errors.program = 'required.'
  }

  return errors
}

export const CreatePoetForm = reduxForm({
  form: 'createPoet',
  validate,
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

export const PoetMgmt = connect(mapStateToProps, actions)(poetMgmt)
