import React, { Component } from 'react';
import {connect} from 'react-redux'
import { Field, reduxForm } from 'redux-form'
import {Link} from 'react-router-dom'
import TextField from 'material-ui/TextField'
import SelectField from 'material-ui/SelectField'
import MenuItem from 'material-ui/MenuItem'
import {get, isEmpty, map} from 'lodash'
import Highlight from 'react-highlight'
import {formatDate} from '../../types/date'
import {getPoetsOfUser} from '../../redux/selectors/poets'
import {
  requestCreatePoet,
  requestDeletePoet,
} from '../../redux/actions/poets'
import {requestReadLanguages} from  '../../redux/actions/languages'
import {resetErrorMsg} from '../../redux/actions/error'


const mapStateToProps = (state, ownProps) => ({
  poets: getPoetsOfUser(state),
  user: state.session.user,
  languages: get(state, `languages`, []),
  errors: get(state, `error`, ''),
  userErrors: get(state, `userError`, '')
})

const actions = {
  requestCreatePoet,
  requestReadLanguages,
  retirePoet: requestDeletePoet,
  resetErrorMsg,
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

    // if we are reloading this page, reset the error message
    this.props.resetErrorMsg()
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
      languages,
      retirePoet,
      user,
    } = this.props
    
    return (
      <div className='profile-poets-container'>
        <div className='profile-poets-list-container'>
          <div className='profile-poets-list-header-container'>
            <div className='profile-poets-list-header'>
              {`poets`}
            </div>
          </div>
          <div className='profile-poets-list'>
            {
              map(
                this.props.poets,
                (poet, idx) => (
                  <PoetSummary poet={poet} retirePoet={retirePoet}key={idx}/>
                ),
                [],
              )
            }
          </div>
        </div>
        <CreatePoetForm onSubmit={this.createPoet} languages={languages} errors={this.props.errors} userErrors={this.props.userErrors}/>
      </div>
    )
  }
}

export class PoetSummary extends Component {
  deletePoet() {
    const {
      poet,
      retirePoet,
    } = this.props
    
    // alert(`are you sure you want to retire ${poet.name}?`)

    retirePoet(poet.id)
  }
  
  render() {
    const {
      poet,
    } = this.props
    
    return (
      <div className='profile-poet-summary-container'>
        <Link className='profile-poet-summary' to={`/poet/${poet.id}`}>
          <div className='profile-poet-name-language'>
            <div className='profile-poet-name'>{poet.name}</div>
            <div className='profile-poet-language'>{poet.language}</div>
          </div>
          <div className='profile-poet-birthday'>{formatDate(poet.birthDate)}</div>
          <div className='profile-poet-description'>{poet.description}</div>
        </Link>
        <div className='profile-poet-delete-container' onClick={() => this.deletePoet()}>
          <div className='profile-poet-delete-button'>x</div>
        </div>
      </div>
    )
  }
}

const renderTextField = multiline => ({input, label, meta: {touched, error}, ...custom}) => (
  <TextField
    hintText={label}
    floatingLabelText={label}
    hintStyle={{color: '#f28cce', fontSize: '1.5em'}}
    inputStyle={{color: '#f28cce', fontSize: '1.5em'}}
    underlineStyle={{borderColor: '#19ecff'}}
    underlineFocusStyle={{borderColor: '#f28cce'}}
    multiLine={multiline}
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
    hintStyle={{color: '#f28cce', fontSize: '1.5em'}}
    inputStyle={{color: '#f28cce', fontSize: '1.5em'}}
    underlineStyle={{borderColor: '#19ecff'}}
    underlineFocusStyle={{borderColor: '#f28cce'}}
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

  onChangeFileName = fileName => text => {
    this.setState({[fileName]: text})
  }
  
  render() {
    const {
      handleSubmit,
      languages,
    } = this.props
    
    return (
      <div className='create-poet-form-container'>
        <div className='create-poet-form-header-container'>
          <div className='create-poet-form-header'>upload a poet</div>
        </div>
        <div className='create-poet-form-and-error'>
          <form className='create-poet-form' onSubmit={handleSubmit}>
          <div>
            <Field name='name' component={renderTextField(false)} type='text' spellcheck='false' placeholder='name'/>
          </div>
          <div styles={{marginTop: '1.5em'}}>
            <Field name='description' component={renderTextField(true)} type='text' placeholder='description'/>
          </div>
          <div className='create-poet-form-language-select'>
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
        <Field className='profile-poet-file-button' id='program' name='program' component={FileInput(this.onChangeFileName('programFileText').bind(this))}/>
            <label htmlFor="program">{this.state.programFileText}</label>
          </div>
          <div className='profile-poet-button'>
        <Field className='profile-poet-file-button' name='parameters' id='parameters' component={FileInput(this.onChangeFileName('parametersFileText').bind(this))}/>
            <label htmlFor="parameters">{this.state.parametersFileText}</label>
            <span style={{padding: '0.6em', marginLeft: '0.8em', fontStyle: 'italic'}}>optional</span>
          </div>
          <button className='profile-poet-submit-button' type='submit'>create poet</button>
        </form>
        <div className='profile-poet-upload-error-message'>
          {this.props.errors}
      </div>
        </div>
        {
        //   <div>
        //   {
        //     this.props.userErrors === '' ?
        //       this.props.userErrors
        //       : <div>
        //       <Highlight>
        //       {this.props.userErrors}
        //     </Highlight>
        //       </div>
        //   }
        // </div>
          
      }
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
const adaptFileEventToValue = (delegate, handler) =>
      e => {
        delegate(e.target.files[0])
        if (e.target.files.length !== 0) {
          handler(e.target.files[0].name)
        }
      }

const FileInput = handler => ({
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
onChange={adaptFileEventToValue(onChange, handler)}
onBlur={adaptFileEventToValue(onChange, handler)}
type="file"
{...inputProps}
{...props}
  />

export const PoetMgmt = connect(mapStateToProps, actions)(poetMgmt)
