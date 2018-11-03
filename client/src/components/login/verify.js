import React, { Component } from 'react';
import {connect} from 'react-redux'
import {get} from 'lodash'
import {requestVerification} from '../../redux/actions/verify'
import queryString from 'query-string'


const mapStateToProps = (state, ownProps) => {
  // parse query string
  const {token, email} = queryString.parse(ownProps.location.search)
  
  return {
    token,
    email,
    error: get(state, `error`, null),
  }
}

const actions = {
  requestVerification,
}

class verify extends Component {
  componentDidMount() {
    // immediately send a request to verify given the token
    // and email address in the query parameters of the URL
    this.props.requestVerification(
      {
        token: this.props.token,
        email: this.props.email,
      },
      
      this.redirectUponVerification,
    )
  }

  // upon verification, we want to redirect the route to the dashboard
  redirectUponVerification = () => {
    const {history} = this.props

    history.push('/profile')
  }

  render() {
    const style = {
      display: 'flex',
      flexDirection: 'column',
      justifyContent: 'center',
      textAlign: 'left',
      fontSize: 'x-large',
      color: '#e58de8',
    }

    const msg = this.props.error || 'pls hold on while we verify yr email...'
    
    return (
      <div style={{display: 'flex', justifyContent: 'center'}}>
        <div style={style}>
          <pre>
            {String.raw`
                     ,
                    /|      __
                   / |   ,-~ /
                  Y :|  //  /
                  | jj /( .^
                  >-"~"-v"
                 /       Y
                jo  o    |
               ( ~T~     j
                >._-' _./
               /   "~"  |
              Y     _,  |
             /| ;-"~ _  l
            / l/ ,-"~    \
            \//\/      .- \
             Y        /    Y 
             l       I     !
             ]\      _\    /"\
            (" ~----( ~   Y.  ) row
            `}
          </pre>
          {msg}
        </div>
      </div>
    )
  }
  
}

export const Verify = connect(mapStateToProps, actions)(verify)
